package dsl

import (
	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
)

// Resource defines a RESTful resource.
//
// Resource must appear in a WebApp expression.
//
// Resource generates all 7 RESTful routes by default:
// - GET    /resources       (index)
// - GET    /resources/new   (new)
// - POST   /resources       (create)
// - GET    /resources/{id}  (show)
// - GET    /resources/{id}/edit (edit)
// - PUT    /resources/{id}  (update)
// - DELETE /resources/{id}  (destroy)
//
// Example:
//
//	WebApp("myapp", func() {
//	    Resource("posts")
//	})
//
// With customization:
//
//	WebApp("myapp", func() {
//	    Resource("posts", func() {
//	        Actions("index", "show", "edit", "update")
//	        Auth("user")
//	        Paginate(20)
//	    })
//	})
func Resource(name string, fn ...func()) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		// Check if we're inside another resource (nested)
		if parent, ok := eval.Current().(*expr.ResourceExpr); ok {
			nested := &expr.ResourceExpr{
				Name:   name,
				Parent: parent,
			}
			if len(fn) > 0 {
				nested.DSLFunc = fn[0]
				// Don't execute here - let RunDSL handle it
			}
			// Add to parent app
			if app := expr.Root; app != nil {
				app.Resources = append(app.Resources, nested)
			}
			return
		}
		eval.IncompatibleDSL()
		return
	}

	resource := &expr.ResourceExpr{
		Name: name,
	}

	if len(fn) > 0 {
		resource.DSLFunc = fn[0]
		// Don't execute here - let RunDSL handle it
	}

	app.Resources = append(app.Resources, resource)
}

// Actions specifies which RESTful actions to generate.
//
// Actions must appear in a Resource expression.
//
// Example:
//
//	Resource("posts", func() {
//	    Actions("index", "show", "edit", "update")
//	})
func Actions(actions ...string) {
	resource, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	resource.Actions = actions
}

// Auth specifies authentication requirements for the resource.
//
// Auth must appear in a Resource expression.
//
// Example:
//
//	Resource("posts", func() {
//	    Auth("authenticated").Except("index", "show")
//	    Auth("admin").Only("destroy")
//	})
func Auth(requirement string) *authBuilder {
	resource, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return nil
	}

	if resource.AuthRequirements == nil {
		resource.AuthRequirements = make(map[string][]string)
	}

	return &authBuilder{
		resource:     resource,
		requirement: requirement,
	}
}

// authBuilder helps build authentication requirements.
type authBuilder struct {
	resource    *expr.ResourceExpr
	requirement string
}

// Except applies the requirement to all actions except the specified ones.
func (a *authBuilder) Except(actions ...string) *authBuilder {
	if a == nil || a.resource == nil {
		return a
	}

	// Apply to all actions except specified
	allActions := []string{"index", "show", "new", "create", "edit", "update", "destroy"}
	excluded := make(map[string]bool)
	for _, action := range actions {
		excluded[action] = true
	}

	for _, action := range allActions {
		if !excluded[action] {
			a.resource.AuthRequirements[action] = append(
				a.resource.AuthRequirements[action], 
				a.requirement,
			)
		}
	}

	return a
}

// Only applies the requirement only to the specified actions.
func (a *authBuilder) Only(actions ...string) *authBuilder {
	if a == nil || a.resource == nil {
		return a
	}

	for _, action := range actions {
		a.resource.AuthRequirements[action] = append(
			a.resource.AuthRequirements[action],
			a.requirement,
		)
	}

	return a
}

// BelongsTo specifies that this resource belongs to a parent resource.
//
// BelongsTo must appear in a nested Resource expression.
//
// Example:
//
//	Resource("users", func() {
//	    Resource("posts", func() {
//	        BelongsTo("user")
//	    })
//	})
func BelongsTo(parent string) {
	_, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	// Parent relationship is already set during Resource() call
	// This is mainly for clarity and potential future use
}

// Index configures the index action.
//
// Index must appear in a Resource expression.
//
// Example:
//
//	Resource("posts", func() {
//	    Index(func() {
//	        Paginate(20)
//	        Searchable("title", "content")
//	        Filterable("status", "category")
//	    })
//	})
func Index(fn func()) {
	res, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if fn != nil {
		eval.Execute(fn, res)
	}
}

// Paginate sets pagination for an action.
//
// Paginate must appear in a Resource or action configuration.
//
// Example:
//
//	Index(func() {
//	    Paginate(20)
//	})
func Paginate(perPage int) {
	resource, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if resource.Pagination == nil {
		resource.Pagination = make(map[string]int)
	}
	resource.Pagination["index"] = perPage
}

// Searchable marks fields as searchable.
//
// Searchable must appear in a Resource or action configuration.
//
// Example:
//
//	Index(func() {
//	    Searchable("title", "content", "author")
//	})
func Searchable(fields ...string) {
	resource, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if resource.SearchableFields == nil {
		resource.SearchableFields = make(map[string][]string)
	}
	resource.SearchableFields["index"] = fields
}

// Filterable marks fields as filterable.
//
// Filterable must appear in a Resource or action configuration.
//
// Example:
//
//	Index(func() {
//	    Filterable("status", "category", "author")
//	})
func Filterable(fields ...string) {
	resource, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if resource.FilterableFields == nil {
		resource.FilterableFields = make(map[string][]string)
	}
	resource.FilterableFields["index"] = fields
}

// Update configures the update action.
//
// Update must appear in a Resource expression.
//
// Example:
//
//	Resource("posts", func() {
//	    Update(func() {
//	        Form("EditPostForm")
//	    })
//	})
func Update(fn func()) {
	res, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if fn != nil {
		eval.Execute(fn, res)
	}
}

// Form specifies a custom form for an action.
//
// Form must appear in an action configuration.
//
// Example:
//
//	Update(func() {
//	    Form("CustomEditForm")
//	})
func Form(name string) {
	resource, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if resource.CustomForms == nil {
		resource.CustomForms = make(map[string]string)
	}
	// Infer the action from context (would need more context tracking)
	// For now, set for common actions
	resource.CustomForms["new"] = name
	resource.CustomForms["create"] = name
	resource.CustomForms["edit"] = name
	resource.CustomForms["update"] = name
}