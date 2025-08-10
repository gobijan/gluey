package expr

// ResourceExpr represents a RESTful resource.
type ResourceExpr struct {
	// Name is the resource name (e.g., "posts").
	Name string
	// DSLFunc contains the DSL function.
	DSLFunc func()
	// Actions to generate (default: all 7 RESTful actions).
	Actions []string
	// Parent resource for nested resources.
	Parent *ResourceExpr
	// Auth requirements.
	AuthRequirements map[string][]string // action -> requirements
	// Pagination settings.
	Pagination map[string]int // action -> per page
	// Searchable fields.
	SearchableFields map[string][]string // action -> fields
	// Filterable fields.
	FilterableFields map[string][]string // action -> fields
	// Custom form for actions.
	CustomForms map[string]string // action -> form name
	// Layout override.
	Layout string
}

// EvalName returns the name of the resource.
func (r *ResourceExpr) EvalName() string {
	return r.Name
}

// DSL returns the DSL function.
func (r *ResourceExpr) DSL() func() {
	return r.DSLFunc
}

// Prepare prepares the resource expression.
func (r *ResourceExpr) Prepare() {
	// Set default actions if not specified
	if len(r.Actions) == 0 {
		r.Actions = []string{"index", "show", "new", "create", "edit", "update", "destroy"}
	}

	// Initialize maps if needed
	if r.AuthRequirements == nil {
		r.AuthRequirements = make(map[string][]string)
	}
	if r.Pagination == nil {
		r.Pagination = make(map[string]int)
	}
	if r.SearchableFields == nil {
		r.SearchableFields = make(map[string][]string)
	}
	if r.FilterableFields == nil {
		r.FilterableFields = make(map[string][]string)
	}
	if r.CustomForms == nil {
		r.CustomForms = make(map[string]string)
	}
}

// Validate validates the resource expression.
func (r *ResourceExpr) Validate() error {
	if r.Name == "" {
		return &ValidationError{Message: "resource name cannot be empty"}
	}

	// Validate action names
	validActions := map[string]bool{
		"index": true, "show": true, "new": true, "create": true,
		"edit": true, "update": true, "destroy": true,
	}

	for _, action := range r.Actions {
		if !validActions[action] {
			return &ValidationError{
				Message: "invalid action: " + action,
			}
		}
	}

	return nil
}

// HasAction returns true if the resource has the specified action.
func (r *ResourceExpr) HasAction(action string) bool {
	for _, a := range r.Actions {
		if a == action {
			return true
		}
	}
	return false
}

// NewFormName returns the form name for the new/create actions.
func (r *ResourceExpr) NewFormName() string {
	if form, ok := r.CustomForms["new"]; ok {
		return form
	}
	if form, ok := r.CustomForms["create"]; ok {
		return form
	}
	// Convention: New{Resource}Form
	return "New" + capitalize(r.Name) + "Form"
}

// EditFormName returns the form name for the edit/update actions.
func (r *ResourceExpr) EditFormName() string {
	if form, ok := r.CustomForms["edit"]; ok {
		return form
	}
	if form, ok := r.CustomForms["update"]; ok {
		return form
	}
	// Convention: Edit{Resource}Form
	return "Edit" + capitalize(r.Name) + "Form"
}

// capitalize capitalizes the first letter of a string.
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return string(s[0]-32) + s[1:]
}
