package dsl

import (
	"github.com/gobijan/gluey/eval"
	"github.com/gobijan/gluey/expr"
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
		resource:    resource,
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

	// Create or get action config
	if res.ActionConfigs == nil {
		res.ActionConfigs = make(map[string]*expr.ActionConfig)
	}
	if res.ActionConfigs["index"] == nil {
		res.ActionConfigs["index"] = &expr.ActionConfig{Action: "index"}
	}

	if fn != nil {
		// Execute in the context of the action config
		eval.Execute(fn, res.ActionConfigs["index"])
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

// Create configures the create action.
//
// Create must appear in a Resource expression.
//
// Example:
//
//	Resource("posts", func() {
//	    Create(func() {
//	        UseForm("PostForm")
//	    })
//	})
func Create(fn func()) {
	res, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	// Create or get action config
	if res.ActionConfigs == nil {
		res.ActionConfigs = make(map[string]*expr.ActionConfig)
	}
	if res.ActionConfigs["create"] == nil {
		res.ActionConfigs["create"] = &expr.ActionConfig{Action: "create"}
	}

	if fn != nil {
		// Execute in the context of the action config
		eval.Execute(fn, res.ActionConfigs["create"])
	}
}

// Update configures the update action.
//
// Update must appear in a Resource expression.
//
// Example:
//
//	Resource("posts", func() {
//	    Update(func() {
//	        UseForm("EditPostForm")
//	    })
//	})
func Update(fn func()) {
	res, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	// Create or get action config
	if res.ActionConfigs == nil {
		res.ActionConfigs = make(map[string]*expr.ActionConfig)
	}
	if res.ActionConfigs["update"] == nil {
		res.ActionConfigs["update"] = &expr.ActionConfig{Action: "update"}
	}

	if fn != nil {
		// Execute in the context of the action config
		eval.Execute(fn, res.ActionConfigs["update"])
	}
}

// Form defines a form within a resource.
//
// Form must appear in a Resource expression.
//
// Example:
//
//	Resource("posts", func() {
//	    Form("PostForm", func() {
//	        Attribute("title", String, Required())
//	        Attribute("content", String, Required())
//	    })
//	})
func Form(name string, fn func()) {
	resource, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if resource.Forms == nil {
		resource.Forms = make(map[string]*expr.FormExpr)
	}

	form := &expr.FormExpr{
		Name: name,
	}

	if fn != nil {
		eval.Execute(fn, form)
	}

	resource.Forms[name] = form
}

// UseForm binds a form to the current action.
//
// UseForm must appear in an action configuration (Create, Update, etc).
//
// Example:
//
//	Create(func() {
//	    UseForm("PostForm")
//	})
func UseForm(name string) {
	// Try to get action config from context
	if config, ok := eval.Current().(*expr.ActionConfig); ok {
		config.FormName = name
		return
	}

	// Fallback to resource level (legacy)
	resource, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if resource.CustomForms == nil {
		resource.CustomForms = make(map[string]string)
	}
	// Set for common actions
	resource.CustomForms["new"] = name
	resource.CustomForms["create"] = name
	resource.CustomForms["edit"] = name
	resource.CustomForms["update"] = name
}

// Singular marks a resource as singular (e.g., session vs sessions).
//
// Singular must appear in a Resource expression.
//
// Example:
//
//	Resource("session", func() {
//	    Singular()
//	})
func Singular() {
	resource, ok := eval.Current().(*expr.ResourceExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	resource.Singular = true
}

// Params defines query parameters for an action.
//
// Params must appear in an action configuration.
//
// Example:
//
//	Index(func() {
//	    Params(func() {
//	        Param("search", String)
//	        Param("page", Int, Default(1))
//	    })
//	})
func Params(fn func()) {
	config, ok := eval.Current().(*expr.ActionConfig)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if fn != nil {
		eval.Execute(fn, config)
	}
}

// Param defines a query parameter.
//
// Param must appear in a Params expression.
//
// Example:
//
//	Param("page", Int, Default(1), Max(100))
func Param(name string, dataType interface{}, fns ...func()) {
	config, ok := eval.Current().(*expr.ActionConfig)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	param := &expr.ParamExpr{
		Name: name,
		Type: dataType.(expr.DataType),
	}

	// Process functions like Default, Max, etc.
	for _, fn := range fns {
		fn() // These would set context on param
	}

	if config.Params == nil {
		config.Params = make([]*expr.ParamExpr, 0)
	}
	config.Params = append(config.Params, param)
}
