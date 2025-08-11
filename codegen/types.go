package codegen

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gobijan/gluey/expr"
)

// TypesGenerator generates form types.
type TypesGenerator struct {
	app     *expr.AppExpr
	version string
	command string
}

// NewTypesGenerator creates a new types generator.
func NewTypesGenerator(app *expr.AppExpr) *TypesGenerator {
	return &TypesGenerator{
		app:     app,
		version: "0.1.0",
		command: "gluey gen design",
	}
}

// SetVersion sets the version for generated headers.
func (g *TypesGenerator) SetVersion(version string) {
	g.version = version
}

// SetCommand sets the command for generated headers.
func (g *TypesGenerator) SetCommand(command string) {
	g.command = command
}

// Generate generates all form types.
func (g *TypesGenerator) Generate() (string, error) {
	var buf bytes.Buffer

	// Header MUST come first, before package declaration
	description := "form types and validation"
	buf.WriteString(GenerateHeader(description, g.version, g.command))

	// Write package header
	buf.WriteString("package types\n\n")

	// Check if we need imports
	needsImports := len(g.app.Forms) > 0

	// Also check resource forms
	if !needsImports {
		for _, resource := range g.app.Resources {
			if len(resource.Forms) > 0 {
				needsImports = true
				break
			}
		}
	}

	// Add imports if needed
	if needsImports {
		buf.WriteString("import (\n")
		buf.WriteString("\t\"github.com/gobijan/gluey/runtime\"\n")
		buf.WriteString(")\n\n")
	}

	// Generate app-level form types (legacy support)
	for _, form := range g.app.Forms {
		code, err := g.generateForm(form)
		if err != nil {
			return "", err
		}
		buf.WriteString(code)
		buf.WriteString("\n")
	}

	// Generate resource-level forms
	for _, resource := range g.app.Resources {
		// Generate forms defined within the resource
		for _, form := range resource.Forms {
			code, err := g.generateForm(form)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
			buf.WriteString("\n")
		}

		// Generate query parameter types for actions with params
		if config, ok := resource.ActionConfigs["index"]; ok && len(config.Params) > 0 {
			typeName := ToCamelCase(resource.Name) + "IndexParams"
			code := g.generateParamsType(typeName, config.Params)
			buf.WriteString(code)
			buf.WriteString("\n")
		}

		// Check if we need to generate default forms
		newFormName := resource.NewFormName()
		needsNewForm := true
		editFormName := resource.EditFormName()
		needsEditForm := true

		// Check if forms exist in resource.Forms
		for name := range resource.Forms {
			if name == newFormName {
				needsNewForm = false
			}
			if name == editFormName {
				needsEditForm = false
			}
		}

		// Also check app-level forms (legacy)
		if g.app.Form(newFormName) != nil {
			needsNewForm = false
		}
		if g.app.Form(editFormName) != nil {
			needsEditForm = false
		}

		// Generate default forms if needed
		if needsNewForm {
			code := g.generateDefaultNewForm(resource)
			buf.WriteString(code)
			buf.WriteString("\n")
		}

		if needsEditForm {
			code := g.generateDefaultEditForm(resource)
			buf.WriteString(code)
			buf.WriteString("\n")
		}
	}

	return buf.String(), nil
}

// generateForm generates a single form type.
func (g *TypesGenerator) generateForm(form *expr.FormExpr) (string, error) {
	var buf bytes.Buffer

	// Generate struct
	buf.WriteString(fmt.Sprintf("// %s represents form data.\n", form.Name))
	buf.WriteString(fmt.Sprintf("type %s struct {\n", form.Name))

	for _, attr := range form.Attributes {
		fieldCode := g.generateField(attr)
		buf.WriteString(fieldCode)
	}

	buf.WriteString("}\n\n")

	// Generate Validate method
	buf.WriteString(fmt.Sprintf("// Validate validates %s.\n", form.Name))
	buf.WriteString(fmt.Sprintf("func (f *%s) Validate() error {\n", form.Name))
	buf.WriteString("\tv := runtime.NewValidator()\n\n")

	for _, attr := range form.Attributes {
		validationCode := g.generateValidation(attr)
		if validationCode != "" {
			buf.WriteString(validationCode)
		}
	}

	buf.WriteString("\n\tif !v.Valid() {\n")
	buf.WriteString("\t\treturn v.Errors()\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\treturn nil\n")
	buf.WriteString("}\n")

	return buf.String(), nil
}

// generateField generates a struct field.
func (g *TypesGenerator) generateField(attr *expr.AttributeExpr) string {
	fieldName := g.toGoName(attr.Name)
	fieldType := g.goType(attr.Type)
	tags := g.generateTags(attr)

	return fmt.Sprintf("\t%s %s %s\n", fieldName, fieldType, tags)
}

// generateTags generates struct tags for a field.
func (g *TypesGenerator) generateTags(attr *expr.AttributeExpr) string {
	var tags []string

	// Form tag
	tags = append(tags, fmt.Sprintf(`form:"%s"`, attr.Name))

	// JSON tag
	jsonTag := attr.Name
	if !attr.IsRequired() {
		jsonTag += ",omitempty"
	}
	tags = append(tags, fmt.Sprintf(`json:"%s"`, jsonTag))

	// Validation tags
	var validations []string
	if attr.IsRequired() {
		validations = append(validations, "required")
	}

	if max, ok := attr.MaxLength(); ok {
		validations = append(validations, fmt.Sprintf("max=%d", max))
	}

	if min, ok := attr.MinLength(); ok {
		validations = append(validations, fmt.Sprintf("min=%d", min))
	}

	if format, ok := attr.Format(); ok {
		switch format {
		case expr.FormatEmail:
			validations = append(validations, "email")
		case expr.FormatURL:
			validations = append(validations, "url")
		}
	}

	if len(validations) > 0 {
		tags = append(tags, fmt.Sprintf(`validate:"%s"`, strings.Join(validations, ",")))
	}

	return fmt.Sprintf("`%s`", strings.Join(tags, " "))
}

// generateValidation generates validation code for an attribute.
func (g *TypesGenerator) generateValidation(attr *expr.AttributeExpr) string {
	var buf bytes.Buffer
	fieldName := g.toGoName(attr.Name)

	if attr.IsRequired() {
		buf.WriteString(fmt.Sprintf("\tv.Required(\"%s\", f.%s)\n", attr.Name, fieldName))
	}

	if format, ok := attr.Format(); ok {
		switch format {
		case expr.FormatEmail:
			buf.WriteString(fmt.Sprintf("\tv.Email(\"%s\", f.%s)\n", attr.Name, fieldName))
		case expr.FormatURL:
			buf.WriteString(fmt.Sprintf("\tv.URL(\"%s\", f.%s)\n", attr.Name, fieldName))
		}
	}

	if min, ok := attr.MinLength(); ok {
		buf.WriteString(fmt.Sprintf("\tv.MinLength(\"%s\", f.%s, %d)\n", attr.Name, fieldName, min))
	}

	if max, ok := attr.MaxLength(); ok {
		buf.WriteString(fmt.Sprintf("\tv.MaxLength(\"%s\", f.%s, %d)\n", attr.Name, fieldName, max))
	}

	return buf.String()
}

// generateDefaultNewForm generates a default form for creating resources.
func (g *TypesGenerator) generateDefaultNewForm(resource *expr.ResourceExpr) string {
	name := resource.NewFormName()
	return fmt.Sprintf(`// %s represents form data for creating %s.
type %s struct {
	// Add your fields here
	// Example:
	// Title   string `+"`"+`form:"title" json:"title" validate:"required"`+"`"+`
	// Content string `+"`"+`form:"content" json:"content"`+"`"+`
}

// Validate validates %s.
func (f *%s) Validate() error {
	return nil
}
`, name, resource.Name, name, name, name)
}

// generateDefaultEditForm generates a default form for editing resources.
func (g *TypesGenerator) generateDefaultEditForm(resource *expr.ResourceExpr) string {
	name := resource.EditFormName()
	return fmt.Sprintf(`// %s represents form data for editing %s.
type %s struct {
	// Add your fields here
	// Example:
	// Title   string `+"`"+`form:"title" json:"title,omitempty"`+"`"+`
	// Content string `+"`"+`form:"content" json:"content,omitempty"`+"`"+`
}

// Validate validates %s.
func (f *%s) Validate() error {
	return nil
}
`, name, resource.Name, name, name, name)
}

// goType converts an expression type to a Go type.
func (g *TypesGenerator) goType(dataType expr.DataType) string {
	if dataType == nil {
		return "string"
	}

	switch dataType {
	case expr.Boolean:
		return "bool"
	case expr.Int:
		return "int"
	case expr.Int32:
		return "int32"
	case expr.Int64:
		return "int64"
	case expr.Float32:
		return "float32"
	case expr.Float64:
		return "float64"
	case expr.String:
		return "string"
	case expr.Bytes:
		return "[]byte"
	}

	// Handle array types
	if arrayType, ok := dataType.(*expr.ArrayType); ok {
		elemType := g.goType(arrayType.ElemType)
		return "[]" + elemType
	}

	// Handle map types
	if mapType, ok := dataType.(*expr.MapType); ok {
		keyType := g.goType(mapType.KeyType)
		valueType := g.goType(mapType.ElemType)
		return fmt.Sprintf("map[%s]%s", keyType, valueType)
	}

	return "any"
}

// generateParamsType generates a query parameters type.
func (g *TypesGenerator) generateParamsType(name string, params []*expr.ParamExpr) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("// %s represents query parameters.\n", name))
	buf.WriteString(fmt.Sprintf("type %s struct {\n", name))

	for _, param := range params {
		fieldName := g.toGoName(param.Name)
		fieldType := g.goType(param.Type)

		// Build tags
		tags := fmt.Sprintf("`form:\"%s\" json:\"%s,omitempty\"`", param.Name, param.Name)

		buf.WriteString(fmt.Sprintf("\t%s %s %s\n", fieldName, fieldType, tags))
	}

	buf.WriteString("}\n")

	return buf.String()
}

// toGoName converts a field name to a Go field name.
func (g *TypesGenerator) toGoName(name string) string {
	return ToCamelCase(name)
}
