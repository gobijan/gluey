package expr

// ActionConfig holds configuration for a resource action.
type ActionConfig struct {
	// Action name (for identification)
	Action string
	// FormName is the name of the form to use for this action
	FormName string
	// Params holds query parameter definitions for index/show actions
	Params []*ParamExpr
}

// EvalName returns the name of the action config.
func (a *ActionConfig) EvalName() string {
	return a.Action
}

// Validate validates the action config.
func (a *ActionConfig) Validate() error {
	return nil
}

// Prepare prepares the action config.
func (a *ActionConfig) Prepare() {
	// Nothing to prepare yet
}

// ParamExpr represents a query parameter.
type ParamExpr struct {
	// Name is the parameter name
	Name string
	// Type is the parameter type
	Type DataType
	// Default value
	Default interface{}
	// Max value (for numeric types)
	Max interface{}
	// Description
	Description string
}

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
	// Forms defined within this resource
	Forms map[string]*FormExpr
	// Whether this is a singular resource (e.g., session vs sessions)
	Singular bool
	// Action configurations
	ActionConfigs map[string]*ActionConfig
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
	if r.Forms == nil {
		r.Forms = make(map[string]*FormExpr)
	}
	if r.ActionConfigs == nil {
		r.ActionConfigs = make(map[string]*ActionConfig)
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
	// Check action configs first
	if config, ok := r.ActionConfigs["create"]; ok && config.FormName != "" {
		return config.FormName
	}
	if config, ok := r.ActionConfigs["new"]; ok && config.FormName != "" {
		return config.FormName
	}
	// Then check custom forms (legacy)
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
	// Check action configs first
	if config, ok := r.ActionConfigs["update"]; ok && config.FormName != "" {
		return config.FormName
	}
	if config, ok := r.ActionConfigs["edit"]; ok && config.FormName != "" {
		return config.FormName
	}
	// Then check custom forms (legacy)
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
