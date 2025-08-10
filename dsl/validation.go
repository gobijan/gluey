package dsl

import "gluey.dev/gluey/expr"

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
func Required() expr.Validation {
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
	return &expr.FormatValidation{Format: format}
}