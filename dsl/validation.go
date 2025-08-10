package dsl

import (
	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
)

// Format constants
const (
	FormatEmail    = expr.FormatEmail
	FormatDate     = expr.FormatDate
	FormatDateTime = expr.FormatDateTime
	FormatURL      = expr.FormatURL
	FormatUUID     = expr.FormatUUID
)

// Required marks an attribute as required.
//
// Required must appear in an Attribute expression.
//
// Example:
//
//	Attribute("email", String, Required())
//	// Or in nested form:
//	Attribute("email", String, func() {
//	    Required()
//	})
func Required() expr.Validation {
	// Check if we're in nested context
	if attr, ok := eval.Current().(*expr.AttributeExpr); ok {
		v := &expr.RequiredValidation{}
		attr.Validations = append(attr.Validations, v)
		return v
	}
	return &expr.RequiredValidation{}
}

// MaxLength sets the maximum length for a string attribute.
//
// MaxLength must appear in an Attribute expression.
//
// Example:
//
//	Attribute("title", String, MaxLength(200))
func MaxLength(max int) expr.Validation {
	// Check if we're in nested context
	if attr, ok := eval.Current().(*expr.AttributeExpr); ok {
		v := &expr.MaxLengthValidation{Max: max}
		attr.Validations = append(attr.Validations, v)
		return v
	}
	return &expr.MaxLengthValidation{Max: max}
}

// MinLength sets the minimum length for a string attribute.
//
// MinLength must appear in an Attribute expression.
//
// Example:
//
//	Attribute("password", String, MinLength(8))
func MinLength(min int) expr.Validation {
	// Check if we're in nested context
	if attr, ok := eval.Current().(*expr.AttributeExpr); ok {
		v := &expr.MinLengthValidation{Min: min}
		attr.Validations = append(attr.Validations, v)
		return v
	}
	return &expr.MinLengthValidation{Min: min}
}

// Format sets the format for a string attribute.
//
// Format must appear in an Attribute expression.
//
// Example:
//
//	Attribute("email", String, Format(FormatEmail))
//	Attribute("website", String, Format(FormatURL))
func Format(format string) expr.Validation {
	// Check if we're in nested context
	if attr, ok := eval.Current().(*expr.AttributeExpr); ok {
		v := &expr.FormatValidation{Format: format}
		attr.Validations = append(attr.Validations, v)
		return v
	}
	return &expr.FormatValidation{Format: format}
}


// Pattern sets a regex pattern for validation.
//
// Pattern must appear in an Attribute expression.
//
// Example:
//
//	Attribute("username", String, Pattern("^[a-zA-Z0-9_]+$"))
func Pattern(pattern string) expr.Validation {
	// Check if we're in nested context
	if attr, ok := eval.Current().(*expr.AttributeExpr); ok {
		v := &expr.PatternValidation{Pattern: pattern}
		attr.Validations = append(attr.Validations, v)
		return v
	}
	return &expr.PatternValidation{Pattern: pattern}
}

// Enum specifies allowed values for an attribute.
//
// Enum must appear in an Attribute expression.
//
// Example:
//
//	Attribute("status", String, Enum("draft", "published", "archived"))
func Enum(values ...string) expr.Validation {
	// Check if we're in nested context
	if attr, ok := eval.Current().(*expr.AttributeExpr); ok {
		v := &expr.EnumValidation{Values: values}
		attr.Validations = append(attr.Validations, v)
		return v
	}
	return &expr.EnumValidation{Values: values}
}

// Min sets the minimum value for numeric attributes.
//
// Min must appear in an Attribute expression.
//
// Example:
//
//	Attribute("age", Int, Min(0))
func Min(min int) expr.Validation {
	// Check if we're in nested context
	if attr, ok := eval.Current().(*expr.AttributeExpr); ok {
		v := &expr.MinValidation{Min: min}
		attr.Validations = append(attr.Validations, v)
		return v
	}
	return &expr.MinValidation{Min: min}
}

// Max sets the maximum value for numeric attributes.
//
// Max must appear in an Attribute expression.
//
// Example:
//
//	Attribute("quantity", Int, Max(100))
func Max(max int) expr.Validation {
	// Check if we're in nested context
	if attr, ok := eval.Current().(*expr.AttributeExpr); ok {
		v := &expr.MaxValidation{Max: max}
		attr.Validations = append(attr.Validations, v)
		return v
	}
	return &expr.MaxValidation{Max: max}
}

// Validation sets a custom validation function.
//
// Validation must appear in an Attribute expression.
//
// Example:
//
//	Attribute("password_confirmation", String, func() {
//	    Validation(func() {
//	        // Custom validation logic
//	    })
//	})
func Validation(fn func()) {
	// This is a placeholder for custom validation logic
	// In a real implementation, this would register the validation
}


