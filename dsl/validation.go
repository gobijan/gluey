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


