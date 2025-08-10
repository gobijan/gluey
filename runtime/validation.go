package runtime

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

// Error returns the error message.
func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

// Error returns all error messages.
func (v ValidationErrors) Error() string {
	var msgs []string
	for _, err := range v {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// HasErrors returns true if there are validation errors.
func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

// Validator provides validation functions.
type Validator struct {
	errors ValidationErrors
}

// NewValidator creates a new validator.
func NewValidator() *Validator {
	return &Validator{
		errors: make(ValidationErrors, 0),
	}
}

// Required validates that a field is not empty.
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "is required",
		})
	}
	return v
}

// Email validates an email address.
func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}
	
	_, err := mail.ParseAddress(value)
	if err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "must be a valid email address",
		})
	}
	return v
}

// URL validates a URL.
func (v *Validator) URL(field, value string) *Validator {
	if value == "" {
		return v
	}
	
	u, err := url.Parse(value)
	if err != nil || u.Scheme == "" || u.Host == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "must be a valid URL",
		})
	}
	return v
}

// MinLength validates minimum string length.
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at least %d characters", min),
		})
	}
	return v
}

// MaxLength validates maximum string length.
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at most %d characters", max),
		})
	}
	return v
}

// Pattern validates against a regex pattern.
func (v *Validator) Pattern(field, value, pattern string) *Validator {
	if value == "" {
		return v
	}
	
	matched, err := regexp.MatchString(pattern, value)
	if err != nil || !matched {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "format is invalid",
		})
	}
	return v
}

// Match validates that two fields match (e.g., password confirmation).
func (v *Validator) Match(field1, value1, field2, value2 string) *Validator {
	if value1 != value2 {
		v.errors = append(v.errors, ValidationError{
			Field:   field2,
			Message: fmt.Sprintf("must match %s", field1),
		})
	}
	return v
}

// Errors returns the validation errors.
func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

// Valid returns true if there are no validation errors.
func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

// ValidateStruct validates a struct based on its tags.
// This is a simple implementation - production apps should use
// a proper validation library like go-playground/validator.
func ValidateStruct(s any) ValidationErrors {
	// TODO: Implement struct tag-based validation
	// For now, return empty errors
	return ValidationErrors{}
}