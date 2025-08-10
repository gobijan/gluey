package expr

import (
	"gluey.dev/gluey/eval"
)

// AppExpr represents a web application.
type AppExpr struct {
	// Name is the application name.
	Name string
	// Description is the optional description.
	Description string
	// DSLFunc contains the DSL function.
	DSLFunc func()
	// Resources defined in the app.
	Resources []*ResourceExpr
	// Pages defined in the app.
	Pages []*PageExpr
	// Forms defined in the app.
	Forms []*FormExpr
	// Layouts defined in the app.
	Layouts []*LayoutExpr
	// DefaultLayout is the default layout name.
	DefaultLayout string
	// Middleware stack.
	Middleware []string
	// Session configuration.
	SessionStore string
	// Assets path.
	AssetsPath string
}

// EvalName returns the name of the application.
func (a *AppExpr) EvalName() string {
	return a.Name
}

// DSL returns the DSL function.
func (a *AppExpr) DSL() func() {
	return a.DSLFunc
}

// WalkSets walks through the expression sets.
func (a *AppExpr) WalkSets(walker eval.SetWalker) {
	// Walk resources
	for _, r := range a.Resources {
		walker(eval.ExpressionSet{r})
	}
	// Walk pages
	for _, p := range a.Pages {
		walker(eval.ExpressionSet{p})
	}
	// Walk forms
	for _, f := range a.Forms {
		walker(eval.ExpressionSet{f})
	}
}

// Packages returns the import paths.
func (a *AppExpr) Packages() []string {
	return []string{"gluey.dev/gluey/dsl"}
}

// Prepare prepares the application expression.
func (a *AppExpr) Prepare() {
	// Set default layout if not specified
	if a.DefaultLayout == "" && len(a.Layouts) > 0 {
		a.DefaultLayout = "application"
	}
	
	// Set default assets path
	if a.AssetsPath == "" {
		a.AssetsPath = "/static"
	}
}

// Validate validates the application expression.
func (a *AppExpr) Validate() error {
	if a.Name == "" {
		return eval.Context.Errors
	}
	return nil
}

// Resource returns a resource by name.
func (a *AppExpr) Resource(name string) *ResourceExpr {
	for _, r := range a.Resources {
		if r.Name == name {
			return r
		}
	}
	return nil
}

// Form returns a form by name.
func (a *AppExpr) Form(name string) *FormExpr {
	for _, f := range a.Forms {
		if f.Name == name {
			return f
		}
	}
	return nil
}