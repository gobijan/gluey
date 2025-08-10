package eval_test

import (
	"errors"
	"testing"

	"gluey.dev/gluey/eval"
)

// TestExpression is a test implementation of Expression
type TestExpression struct {
	name      string
	prepared  bool
	validated bool
	finalized bool
	dslFunc   func()
	validErr  error
}

func (t *TestExpression) EvalName() string {
	return t.name
}

func (t *TestExpression) Prepare() {
	t.prepared = true
}

func (t *TestExpression) Validate() error {
	t.validated = true
	return t.validErr
}

func (t *TestExpression) Finalize() {
	t.finalized = true
}

func (t *TestExpression) DSL() func() {
	return t.dslFunc
}

func (t *TestExpression) WalkSets(walk eval.SetWalker) {}

func (t *TestExpression) Packages() []string {
	return []string{}
}

func TestRunDSL(t *testing.T) {
	// Reset context before each test
	eval.Context.Reset()

	tests := []struct {
		name    string
		root    eval.Root
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil root returns error",
			root:    nil,
			wantErr: true,
			errMsg:  "no root expression found",
		},
		{
			name: "successful execution",
			root: &TestExpression{
				name:    "test",
				dslFunc: func() {},
			},
			wantErr: false,
		},
		{
			name: "validation error propagates",
			root: &TestExpression{
				name:     "test",
				dslFunc:  func() {},
				validErr: errors.New("validation failed"),
			},
			wantErr: true,
			errMsg:  "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval.SetRoot(tt.root)
			err := eval.RunDSL()

			if (err != nil) != tt.wantErr {
				t.Errorf("RunDSL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("RunDSL() error = %v, want %v", err.Error(), tt.errMsg)
			}

			// Check phases were executed in order
			if expr, ok := tt.root.(*TestExpression); ok && !tt.wantErr {
				if !expr.prepared {
					t.Error("Prepare() was not called")
				}
				if !expr.validated {
					t.Error("Validate() was not called")
				}
				if !expr.finalized {
					t.Error("Finalize() was not called")
				}
			}
		})
	}
}

func TestExecute(t *testing.T) {
	eval.Context.Reset()

	executed := false
	dsl := func() {
		executed = true
	}

	expr := &TestExpression{name: "test"}

	result := eval.Execute(dsl, expr)

	if !result {
		t.Error("Execute() returned false for successful execution")
	}

	if !executed {
		t.Error("DSL function was not executed")
	}

	if eval.Current() != nil {
		t.Error("Current expression was not reset after execution")
	}
}

func TestReportError(t *testing.T) {
	eval.Context.Reset()

	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	eval.ReportError(err1)
	if eval.Context.Errors == nil {
		t.Error("First error was not recorded")
	}

	eval.ReportError(err2)
	if eval.Context.Errors == nil {
		t.Error("Second error was not recorded")
	}

	// Check that nil errors are ignored
	eval.Context.Reset()
	eval.ReportError(nil)
	if eval.Context.Errors != nil {
		t.Error("Nil error should not be recorded")
	}
}

func TestIncompatibleDSL(t *testing.T) {
	eval.Context.Reset()

	expr := &TestExpression{name: "test"}
	eval.Execute(func() {}, expr)

	// Inside the DSL execution, report incompatible DSL
	eval.Execute(func() {
		eval.IncompatibleDSL()
	}, expr)

	if eval.Context.Errors == nil {
		t.Error("IncompatibleDSL should record an error")
	}
}

func TestInvalidArgError(t *testing.T) {
	eval.Context.Reset()

	eval.InvalidArgError("string", 123)

	if eval.Context.Errors == nil {
		t.Error("InvalidArgError should record an error")
	}

	errStr := eval.Context.Errors.Error()
	if errStr != "invalid argument: expected string, got int" {
		t.Errorf("InvalidArgError message = %v, want 'invalid argument: expected string, got int'", errStr)
	}
}

func TestTooManyArgError(t *testing.T) {
	eval.Context.Reset()

	eval.TooManyArgError()

	if eval.Context.Errors == nil {
		t.Error("TooManyArgError should record an error")
	}

	errStr := eval.Context.Errors.Error()
	if errStr != "too many arguments provided" {
		t.Errorf("TooManyArgError message = %v, want 'too many arguments provided'", errStr)
	}
}