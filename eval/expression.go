package eval

// Expression is the base interface for all DSL expressions.
type Expression interface {
	// EvalName returns the name of the expression for debugging.
	EvalName() string
}

// Root is the top-level expression that serves as the entry point.
type Root interface {
	Expression
	// WalkSets walks through expression sets for evaluation.
	WalkSets(walker SetWalker)
	// Packages returns the import paths for error reporting.
	Packages() []string
}

// Source expressions contain DSL functions to be executed.
type Source interface {
	Expression
	// DSL returns the DSL function to be executed.
	DSL() func()
}

// Preparer expressions require preparation after DSL execution.
type Preparer interface {
	Expression
	// Prepare is called after DSL execution but before validation.
	Prepare()
}

// Validator expressions can be validated for correctness.
type Validator interface {
	Expression
	// Validate returns an error if the expression is invalid.
	Validate() error
}

// Finalizer expressions require finalization after validation.
type Finalizer interface {
	Expression
	// Finalize is called after validation, before code generation.
	Finalize()
}

// SetWalker is a function that processes expression sets.
type SetWalker func(set ExpressionSet)

// ExpressionSet is a collection of expressions to be processed together.
type ExpressionSet []Expression

// TopExpr is a marker interface for top-level expressions.
type TopExpr interface {
	Expression
	top()
}
