package eval

import (
	"errors"
	"fmt"
)

// RunDSL executes the DSL evaluation pipeline.
// It runs through four phases: Execute, Prepare, Validate, and Finalize.
func RunDSL() error {
	root := CurrentRoot()
	if root == nil {
		return errors.New("no root expression found")
	}

	// Phase 1: Execute DSL functions
	if src, ok := root.(Source); ok {
		if dsl := src.DSL(); dsl != nil {
			Execute(dsl, root)
		}
	}
	root.WalkSets(executeSet)

	// Check for errors after execution
	if Context.Errors != nil {
		return Context.Errors
	}

	// Phase 2: Prepare expressions
	if prep, ok := root.(Preparer); ok {
		prep.Prepare()
	}
	root.WalkSets(prepareSet)

	// Phase 3: Validate expressions
	if val, ok := root.(Validator); ok {
		if err := val.Validate(); err != nil {
			return err
		}
	}
	root.WalkSets(validateSet)

	// Phase 4: Finalize expressions
	if fin, ok := root.(Finalizer); ok {
		fin.Finalize()
	}
	root.WalkSets(finalizeSet)

	return Context.Errors
}

// executeSet executes DSL functions in the expression set.
func executeSet(set ExpressionSet) {
	for _, expr := range set {
		if src, ok := expr.(Source); ok {
			if dsl := src.DSL(); dsl != nil {
				Execute(dsl, expr)
			}
		}
	}
}

// prepareSet prepares expressions in the set.
func prepareSet(set ExpressionSet) {
	for _, expr := range set {
		if prep, ok := expr.(Preparer); ok {
			prep.Prepare()
		}
	}
}

// validateSet validates expressions in the set.
func validateSet(set ExpressionSet) {
	for _, expr := range set {
		if val, ok := expr.(Validator); ok {
			if err := val.Validate(); err != nil {
				ReportError(err)
			}
		}
	}
}

// finalizeSet finalizes expressions in the set.
func finalizeSet(set ExpressionSet) {
	for _, expr := range set {
		if fin, ok := expr.(Finalizer); ok {
			fin.Finalize()
		}
	}
}

// Execute runs a DSL function in the context of an expression.
func Execute(dsl func(), expr Expression) bool {
	if dsl == nil {
		return true
	}

	// Set current expression in context
	oldCurrent := Context.current
	Context.current = expr
	defer func() {
		Context.current = oldCurrent
	}()

	// Execute the DSL
	dsl()

	return Context.Errors == nil
}

// Current returns the current expression being evaluated.
func Current() Expression {
	return Context.current
}

// CurrentRoot returns the root expression.
func CurrentRoot() Root {
	return Context.root
}

// SetRoot sets the root expression.
func SetRoot(root Root) {
	Context.root = root
}

// ReportError reports an error during evaluation.
func ReportError(err error) {
	if err == nil {
		return
	}
	Context.recordError(err)
}

// IncompatibleDSL reports that a DSL function was called in the wrong context.
func IncompatibleDSL() {
	ReportError(fmt.Errorf("incompatible DSL function called in %T context", Current()))
}

// InvalidArgError reports an invalid argument error.
func InvalidArgError(expected string, actual interface{}) {
	ReportError(fmt.Errorf("invalid argument: expected %s, got %T", expected, actual))
}

// TooManyArgError reports too many arguments were provided.
func TooManyArgError() {
	ReportError(errors.New("too many arguments provided"))
}