package dsl

import (
	"fmt"
	
	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
)

// WebApp defines a web application.
// It is the top-level DSL function that creates the root expression.
//
// Example:
//
//	var _ = WebApp("myapp", func() {
//	    Description("My application")
//	    Resource("posts")
//	    Page("home", "/")
//	})
func WebApp(name string, fn func()) *expr.AppExpr {
	// Create the app expression
	app := &expr.AppExpr{
		Name:    name,
		DSLFunc: fn,
	}

	// Set as root
	expr.Root = app
	eval.SetRoot(app)

	// Don't execute the DSL here - let RunDSL do it
	// The DSL will be executed during eval.RunDSL()

	return app
}

// Description sets the application description.
//
// Description must appear in a WebApp expression.
//
// Example:
//
//	WebApp("myapp", func() {
//	    Description("My awesome web application")
//	})
func Description(desc string) {
	switch e := eval.Current().(type) {
	case *expr.AppExpr:
		e.Description = desc
	case *expr.ResourceExpr:
		// Resources can have descriptions too (for future use)
	case *expr.PageExpr:
		// Pages can have descriptions too (for future use)
	case *expr.AttributeExpr:
		// Handle attribute description
		e.Description = desc
	default:
		eval.IncompatibleDSL()
	}
}

// Use adds middleware to the application stack.
//
// Use must appear in a WebApp expression.
//
// Example:
//
//	WebApp("myapp", func() {
//	    Use("RequestID", "Logger", "Recover")
//	})
func Use(middleware ...string) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	app.Middleware = append(app.Middleware, middleware...)
}

// Sessions configures session storage.
//
// Sessions must appear in a WebApp expression.
//
// Example:
//
//	WebApp("myapp", func() {
//	    Sessions(func() {
//	        Store("cookie")
//	    })
//	})
func Sessions(fn func()) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if fn != nil {
		eval.Execute(fn, app)
	}
}

// Store sets the session store type.
//
// Store must appear in a Sessions expression.
//
// Example:
//
//	Sessions(func() {
//	    Store("cookie")  // or "redis", "memory"
//	})
func Store(store string) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	app.SessionStore = store
}

// Assets configures static asset serving.
//
// Assets must appear in a WebApp expression.
//
// Example:
//
//	WebApp("myapp", func() {
//	    Assets(func() {
//	        Path("/static")
//	    })
//	})
func Assets(fn func()) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if fn != nil {
		eval.Execute(fn, app)
	}
}

// Path sets the assets path.
//
// Path must appear in an Assets expression.
//
// Example:
//
//	Assets(func() {
//	    Path("/static")
//	})
func Path(path string) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	app.AssetsPath = path
}

// Layouts configures the application layouts.
//
// Layouts must appear in a WebApp expression.
//
// Example:
//
//	WebApp("myapp", func() {
//	    Layouts(func() {
//	        Default("application")
//	        Layout("admin")
//	        Layout("marketing")
//	    })
//	})
func Layouts(fn func()) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if fn != nil {
		eval.Execute(fn, app)
	}
}

// Default sets the default layout.
//
// Default must appear in a Layouts expression.
//
// Example:
//
//	Layouts(func() {
//	    Default("application")
//	})
func Default(value interface{}) {
	switch e := eval.Current().(type) {
	case *expr.AppExpr:
		// Default layout
		if name, ok := value.(string); ok {
			e.DefaultLayout = name
		}
	case *expr.AttributeExpr:
		// Default value for attribute - store in description for now
		// TODO: Add proper Default field to AttributeExpr
		if e.Description != "" {
			e.Description += " "
		}
		e.Description += fmt.Sprintf("(default: %v)", value)
	default:
		eval.IncompatibleDSL()
	}
}

// Layout defines a layout.
//
// Layout must appear in a Layouts or WebApp expression.
//
// Example:
//
//	Layouts(func() {
//	    Layout("admin")
//	})
func Layout(name string, fn ...func()) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	layout := &expr.LayoutExpr{
		Name:     name,
		Template: "layouts/" + name + ".html",
	}

	if len(fn) > 0 {
		layout.DSLFunc = fn[0]
		eval.Execute(fn[0], layout)
	}

	app.Layouts = append(app.Layouts, layout)
}