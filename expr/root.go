package expr

// Root is the global expression root that holds the web application.
var Root *AppExpr

// Reset clears the root expression.
func Reset() {
	Root = nil
}