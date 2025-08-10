package eval

import (
	"fmt"
	"strings"
)

// Context holds the evaluation context.
var Context = &evalContext{}

// evalContext manages the state during DSL evaluation.
type evalContext struct {
	root    Root       // The root expression
	current Expression // Currently evaluating expression
	Errors  error      // Accumulated errors
}

// recordError records an error in the context.
func (c *evalContext) recordError(err error) {
	if err == nil {
		return
	}

	if c.Errors == nil {
		c.Errors = err
	} else {
		// Accumulate multiple errors
		c.Errors = &multiError{
			errors: []error{c.Errors, err},
		}
	}
}

// Reset clears the evaluation context.
func (c *evalContext) Reset() {
	c.root = nil
	c.current = nil
	c.Errors = nil
}

// multiError represents multiple errors.
type multiError struct {
	errors []error
}

// Error implements the error interface.
func (m *multiError) Error() string {
	var msgs []string
	for _, err := range m.errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("%d errors:\n%s", len(m.errors), strings.Join(msgs, "\n"))
}
