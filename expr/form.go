package expr

// FormExpr represents a form type.
type FormExpr struct {
	// Name is the form name.
	Name string
	// DSLFunc contains the DSL function.
	DSLFunc func()
	// Attributes are the form fields.
	Attributes []*AttributeExpr
}

// EvalName returns the name of the form.
func (f *FormExpr) EvalName() string {
	return f.Name
}

// DSL returns the DSL function.
func (f *FormExpr) DSL() func() {
	return f.DSLFunc
}

// Prepare prepares the form expression.
func (f *FormExpr) Prepare() {
	// Prepare attributes
	for _, attr := range f.Attributes {
		attr.Prepare()
	}
}

// Validate validates the form expression.
func (f *FormExpr) Validate() error {
	if f.Name == "" {
		return &ValidationError{Message: "form name cannot be empty"}
	}
	
	// Validate attributes
	for _, attr := range f.Attributes {
		if err := attr.Validate(); err != nil {
			return err
		}
	}
	
	return nil
}

// Attribute returns an attribute by name.
func (f *FormExpr) Attribute(name string) *AttributeExpr {
	for _, attr := range f.Attributes {
		if attr.Name == name {
			return attr
		}
	}
	return nil
}