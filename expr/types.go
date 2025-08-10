package expr

import "fmt"

// Primitive types
var (
	Boolean = &PrimitiveType{name: "boolean", kind: BooleanKind}
	Int     = &PrimitiveType{name: "int", kind: IntKind}
	Int32   = &PrimitiveType{name: "int32", kind: IntKind}
	Int64   = &PrimitiveType{name: "int64", kind: IntKind}
	Float32 = &PrimitiveType{name: "float32", kind: FloatKind}
	Float64 = &PrimitiveType{name: "float64", kind: FloatKind}
	String  = &PrimitiveType{name: "string", kind: StringKind}
	Bytes   = &PrimitiveType{name: "bytes", kind: BytesKind}
)

// PrimitiveType represents a primitive data type.
type PrimitiveType struct {
	name string
	kind TypeKind
}

// Name returns the type name.
func (p *PrimitiveType) Name() string {
	return p.name
}

// Kind returns the type kind.
func (p *PrimitiveType) Kind() TypeKind {
	return p.kind
}

// ArrayType represents an array type.
type ArrayType struct {
	ElemType DataType
}

// Name returns the type name.
func (a *ArrayType) Name() string {
	return fmt.Sprintf("array<%s>", a.ElemType.Name())
}

// Kind returns the type kind.
func (a *ArrayType) Kind() TypeKind {
	return ArrayKind
}

// MapType represents a map type.
type MapType struct {
	KeyType  DataType
	ElemType DataType
}

// Name returns the type name.
func (m *MapType) Name() string {
	return fmt.Sprintf("map<%s,%s>", m.KeyType.Name(), m.ElemType.Name())
}

// Kind returns the type kind.
func (m *MapType) Kind() TypeKind {
	return MapKind
}

// LayoutExpr represents a layout definition.
type LayoutExpr struct {
	// Name is the layout name.
	Name string
	// Template path.
	Template string
	// DSLFunc contains the DSL function.
	DSLFunc func()
}

// EvalName returns the name of the layout.
func (l *LayoutExpr) EvalName() string {
	return l.Name
}

// ValidationError represents a validation error.
type ValidationError struct {
	Message string
}

// Error returns the error message.
func (v *ValidationError) Error() string {
	return v.Message
}

// RequiredValidation validates that a value is not empty.
type RequiredValidation struct{}

// Name returns the validation name.
func (r *RequiredValidation) Name() string {
	return "required"
}

// Validate checks if the value is not empty.
func (r *RequiredValidation) Validate(value interface{}) error {
	if value == nil {
		return &ValidationError{Message: "value is required"}
	}
	if s, ok := value.(string); ok && s == "" {
		return &ValidationError{Message: "value is required"}
	}
	return nil
}

// MaxLengthValidation validates maximum length.
type MaxLengthValidation struct {
	Max int
}

// Name returns the validation name.
func (m *MaxLengthValidation) Name() string {
	return "max_length"
}

// Validate checks if the value exceeds maximum length.
func (m *MaxLengthValidation) Validate(value interface{}) error {
	if s, ok := value.(string); ok {
		if len(s) > m.Max {
			return &ValidationError{
				Message: fmt.Sprintf("value exceeds maximum length of %d", m.Max),
			}
		}
	}
	return nil
}

// MinLengthValidation validates minimum length.
type MinLengthValidation struct {
	Min int
}

// Name returns the validation name.
func (m *MinLengthValidation) Name() string {
	return "min_length"
}

// Validate checks if the value meets minimum length.
func (m *MinLengthValidation) Validate(value interface{}) error {
	if s, ok := value.(string); ok {
		if len(s) < m.Min {
			return &ValidationError{
				Message: fmt.Sprintf("value must be at least %d characters", m.Min),
			}
		}
	}
	return nil
}

// FormatValidation validates a value format.
type FormatValidation struct {
	Format string
}

// Name returns the validation name.
func (f *FormatValidation) Name() string {
	return "format"
}

// Validate checks if the value matches the format.
func (f *FormatValidation) Validate(value interface{}) error {
	// TODO: Implement format validation (email, date, etc.)
	return nil
}

// Common formats
const (
	FormatEmail    = "email"
	FormatDate     = "date"
	FormatDateTime = "datetime"
	FormatURL      = "url"
	FormatUUID     = "uuid"
)