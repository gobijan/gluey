package dsl

import (
	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
)

// Page defines a non-resource page.
//
// Page must appear in a WebApp expression.
//
// Example:
//
//	WebApp("myapp", func() {
//	    Page("home", "/")
//	    Page("about", "/about")
//	})
//
// With configuration:
//
//	WebApp("myapp", func() {
//	    Page("contact", func() {
//	        Route("GET", "/contact")
//	        Route("POST", "/contact")
//	        Form("ContactForm")
//	    })
//	})
func Page(name string, args ...interface{}) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	page := &expr.PageExpr{
		Name: name,
	}

	// Parse arguments - can be a path string or a DSL function
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			// It's a path
			page.Routes = append(page.Routes, expr.RouteExpr{
				Method: "GET",
				Path:   v,
			})
		case func():
			// It's a DSL function
			page.DSLFunc = v
			eval.Execute(v, page)
		default:
			eval.InvalidArgError("path string or DSL function", arg)
		}
	}

	app.Pages = append(app.Pages, page)
}

// Route adds a route to a page.
//
// Route must appear in a Page expression.
//
// Example:
//
//	Page("contact", func() {
//	    Route("GET", "/contact")
//	    Route("POST", "/contact")
//	})
func Route(method, path string) {
	page, ok := eval.Current().(*expr.PageExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	page.Routes = append(page.Routes, expr.RouteExpr{
		Method: method,
		Path:   path,
	})
}