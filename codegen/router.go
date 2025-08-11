package codegen

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gobijan/gluey/expr"
)

// RouterGenerator generates the router setup.
type RouterGenerator struct {
	app     *expr.AppExpr
	version string
	command string
}

// NewRouterGenerator creates a new router generator.
func NewRouterGenerator(app *expr.AppExpr) *RouterGenerator {
	return &RouterGenerator{
		app:     app,
		version: "0.1.0",
		command: "gluey gen design",
	}
}

// SetVersion sets the version for generated headers.
func (g *RouterGenerator) SetVersion(version string) {
	g.version = version
}

// SetCommand sets the command for generated headers.
func (g *RouterGenerator) SetCommand(command string) {
	g.command = command
}

// Generate generates the router setup code.
func (g *RouterGenerator) Generate() (string, error) {
	var buf bytes.Buffer

	// Header MUST come first, before package declaration
	description := "HTTP router setup"
	buf.WriteString(GenerateHeader(description, g.version, g.command))

	buf.WriteString(fmt.Sprintf("package %s\n\n", g.app.Name))
	buf.WriteString("import (\n")
	buf.WriteString("\t\"net/http\"\n")
	buf.WriteString(fmt.Sprintf("\t\"%s/app/controllers\"\n", g.app.Name))
	buf.WriteString(")\n\n")

	// Generate Controllers interface
	buf.WriteString("// Controllers holds all controller implementations.\n")
	buf.WriteString("type Controllers struct {\n")

	for _, resource := range g.app.Resources {
		controllerName := g.toControllerName(resource.Name)
		buf.WriteString(fmt.Sprintf("\t%s controllers.%s\n",
			ToTitle(resource.Name), controllerName))
	}

	if len(g.app.Pages) > 0 {
		buf.WriteString("\tPages controllers.PagesController\n")
	}

	buf.WriteString("}\n\n")

	// Generate MountRoutes function
	buf.WriteString("// MountRoutes mounts all routes on the given mux.\n")
	buf.WriteString("func MountRoutes(mux *http.ServeMux, c Controllers) {\n")

	// Add middleware comment
	if len(g.app.Middleware) > 0 {
		buf.WriteString("\t// TODO: Apply middleware stack: ")
		buf.WriteString(strings.Join(g.app.Middleware, ", "))
		buf.WriteString("\n\n")
	}

	// Mount resource routes
	for _, resource := range g.app.Resources {
		g.generateResourceRoutes(&buf, resource)
		buf.WriteString("\n")
	}

	// Mount page routes
	for _, page := range g.app.Pages {
		g.generatePageRoutes(&buf, page)
	}

	// Mount static files
	if g.app.AssetsPath != "" {
		buf.WriteString("\t// Static files\n")
		buf.WriteString(fmt.Sprintf("\tmux.Handle(\"%s/\", http.StripPrefix(\"%s/\", http.FileServer(http.Dir(\"public\"))))\n",
			g.app.AssetsPath, g.app.AssetsPath))
	}

	buf.WriteString("}\n")

	return buf.String(), nil
}

// generateResourceRoutes generates routes for a resource.
func (g *RouterGenerator) generateResourceRoutes(buf *bytes.Buffer, resource *expr.ResourceExpr) {
	controllerVar := "c." + ToTitle(resource.Name)
	basePath := "/" + resource.Name

	// Handle nested resources
	if resource.Parent != nil {
		basePath = "/" + resource.Parent.Name + "/{" + g.toSingular(resource.Parent.Name) + "_id}/" + resource.Name
	}

	fmt.Fprintf(buf, "\t// %s routes\n", ToTitle(resource.Name))

	for _, action := range resource.Actions {
		method, path := g.getRouteForAction(action, basePath, resource.Name)
		handler := fmt.Sprintf("%s.%s", controllerVar, ToTitle(action))

		// Add auth comment if required
		if auths, ok := resource.AuthRequirements[action]; ok && len(auths) > 0 {
			fmt.Fprintf(buf, "\t// Requires: %s\n", strings.Join(auths, ", "))
		}

		fmt.Fprintf(buf, "\tmux.HandleFunc(\"%s %s\", %s)\n", method, path, handler)
	}
}

// generatePageRoutes generates routes for a page.
func (g *RouterGenerator) generatePageRoutes(buf *bytes.Buffer, page *expr.PageExpr) {
	for _, route := range page.Routes {
		methodName := g.toPageMethodName(page.Name, route.Method)
		handler := fmt.Sprintf("c.Pages.%s", methodName)

		// Add auth comment if required
		if len(page.AuthRequirements) > 0 {
			fmt.Fprintf(buf, "\t// Requires: %s\n", strings.Join(page.AuthRequirements, ", "))
		}

		fmt.Fprintf(buf, "\tmux.HandleFunc(\"%s %s\", %s)\n", route.Method, route.Path, handler)
	}
}

// getRouteForAction returns the HTTP method and path for a RESTful action.
func (g *RouterGenerator) getRouteForAction(action, basePath, resourceName string) (string, string) {
	switch action {
	case "index":
		return "GET", basePath
	case "show":
		return "GET", basePath + "/{id}"
	case "new":
		return "GET", basePath + "/new"
	case "create":
		return "POST", basePath
	case "edit":
		return "GET", basePath + "/{id}/edit"
	case "update":
		// Using POST with _method override for browser compatibility
		// Real apps would handle both PUT and PATCH
		return "POST", basePath + "/{id}"
	case "destroy":
		// Using POST with _method override for browser compatibility
		// Real apps would handle DELETE
		return "POST", basePath + "/{id}/delete"
	default:
		// Custom action
		return "POST", basePath + "/{id}/" + action
	}
}

// toControllerName converts a resource name to a controller name.
func (g *RouterGenerator) toControllerName(resourceName string) string {
	name := ToTitle(resourceName)
	if !strings.HasSuffix(name, "s") {
		name += "s"
	}
	return name + "Controller"
}

// toSingular converts a plural resource name to singular.
func (g *RouterGenerator) toSingular(plural string) string {
	return ToSingular(plural)
}

// toPageMethodName converts a page name and method to a method name.
func (g *RouterGenerator) toPageMethodName(pageName, method string) string {
	name := ToTitle(pageName)
	if method != "GET" {
		name += ToTitle(strings.ToLower(method))
	}
	return name
}
