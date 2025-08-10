package expr

// PageExpr represents a non-resource page.
type PageExpr struct {
	// Name is the page name.
	Name string
	// DSLFunc contains the DSL function.
	DSLFunc func()
	// Routes for the page.
	Routes []RouteExpr
	// Form used on the page.
	FormName string
	// Layout override.
	Layout string
	// Auth requirements.
	AuthRequirements []string
}

// RouteExpr represents an HTTP route.
type RouteExpr struct {
	// Method is the HTTP method (GET, POST, etc.).
	Method string
	// Path is the URL path.
	Path string
}

// EvalName returns the name of the page.
func (p *PageExpr) EvalName() string {
	return p.Name
}

// DSL returns the DSL function.
func (p *PageExpr) DSL() func() {
	return p.DSLFunc
}

// Prepare prepares the page expression.
func (p *PageExpr) Prepare() {
	// If no routes specified but path given, assume GET
	if len(p.Routes) == 0 && p.Name != "" {
		// Simple pages might just have a name that becomes the path
		p.Routes = []RouteExpr{
			{Method: "GET", Path: "/" + p.Name},
		}
	}
}

// Validate validates the page expression.
func (p *PageExpr) Validate() error {
	if p.Name == "" {
		return &ValidationError{Message: "page name cannot be empty"}
	}
	
	if len(p.Routes) == 0 {
		return &ValidationError{Message: "page must have at least one route"}
	}
	
	// Validate HTTP methods
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true,
		"DELETE": true, "HEAD": true, "OPTIONS": true,
	}
	
	for _, route := range p.Routes {
		if !validMethods[route.Method] {
			return &ValidationError{
				Message: "invalid HTTP method: " + route.Method,
			}
		}
		if route.Path == "" {
			return &ValidationError{
				Message: "route path cannot be empty",
			}
		}
	}
	
	return nil
}