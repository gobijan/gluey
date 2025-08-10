package expr

// AttributeExpr represents a form field or type attribute.
type AttributeExpr struct {
	// Name is the attribute name.
	Name string
	// Type is the attribute type.
	Type DataType
	// Description is optional documentation.
	Description string
	// Validations applied to the attribute.
	Validations []Validation
	// DefaultValue if any.
	DefaultValue interface{}
	// Meta contains additional metadata.
	Meta map[string]interface{}
}

// DataType represents the type of an attribute.
type DataType interface {
	// Name returns the type name.
	Name() string
	// Kind returns the type kind.
	Kind() TypeKind
}

// TypeKind represents the kind of type.
type TypeKind int

const (
	// BooleanKind is a boolean type.
	BooleanKind TypeKind = iota + 1
	// IntKind is an integer type.
	IntKind
	// FloatKind is a floating point type.
	FloatKind
	// StringKind is a string type.
	StringKind
	// BytesKind is a bytes type.
	BytesKind
	// ArrayKind is an array type.
	ArrayKind
	// ObjectKind is an object type.
	ObjectKind
	// MapKind is a map type.
	MapKind
)

// Validation represents a validation rule.
type Validation interface {
	// Name returns the validation name.
	Name() string
	// Validate checks if a value is valid.
	Validate(value interface{}) error
}

// EvalName returns the name of the attribute.
func (a *AttributeExpr) EvalName() string {
	return a.Name
}

// Prepare prepares the attribute expression.
func (a *AttributeExpr) Prepare() {
	// Set default type if not specified
	if a.Type == nil {
		a.Type = String
	}

	// Initialize meta if needed
	if a.Meta == nil {
		a.Meta = make(map[string]interface{})
	}
}

// Validate validates the attribute expression.
func (a *AttributeExpr) Validate() error {
	if a.Name == "" {
		return &ValidationError{Message: "attribute name cannot be empty"}
	}

	if a.Type == nil {
		return &ValidationError{Message: "attribute type cannot be nil"}
	}

	return nil
}

// IsRequired returns true if the attribute is required.
func (a *AttributeExpr) IsRequired() bool {
	for _, v := range a.Validations {
		if v.Name() == "required" {
			return true
		}
	}
	return false
}

// MaxLength returns the maximum length validation if any.
func (a *AttributeExpr) MaxLength() (int, bool) {
	for _, v := range a.Validations {
		if ml, ok := v.(*MaxLengthValidation); ok {
			return ml.Max, true
		}
	}
	return 0, false
}

// MinLength returns the minimum length validation if any.
func (a *AttributeExpr) MinLength() (int, bool) {
	for _, v := range a.Validations {
		if ml, ok := v.(*MinLengthValidation); ok {
			return ml.Min, true
		}
	}
	return 0, false
}

// Format returns the format validation if any.
func (a *AttributeExpr) Format() (string, bool) {
	for _, v := range a.Validations {
		if f, ok := v.(*FormatValidation); ok {
			return f.Format, true
		}
	}
	return "", false
}
