package codegen

import (
	"fmt"
	"strings"

	"gluey.dev/gluey/expr"
)

// ViewsGenerator generates HTML templates.
type ViewsGenerator struct {
	app *expr.AppExpr
}

// NewViewsGenerator creates a new views generator.
func NewViewsGenerator(app *expr.AppExpr) *ViewsGenerator {
	return &ViewsGenerator{app: app}
}

// GenerateLayout generates the main layout template.
func (g *ViewsGenerator) GenerateLayout() (string, error) {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - ` + g.app.Name + `</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: system-ui, -apple-system, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        header { background: #2c3e50; color: white; padding: 1rem 0; margin-bottom: 2rem; }
        header h1 { font-size: 1.5rem; }
        nav { margin-top: 1rem; }
        nav a { color: white; text-decoration: none; margin-right: 1rem; }
        nav a:hover { text-decoration: underline; }
        .flash { padding: 10px; margin-bottom: 20px; border-radius: 4px; }
        .flash.success { background: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .flash.error { background: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
        .flash.info { background: #d1ecf1; color: #0c5460; border: 1px solid #bee5eb; }
        .flash.warning { background: #fff3cd; color: #856404; border: 1px solid #ffeaa7; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background: #f5f5f5; }
        .btn { display: inline-block; padding: 8px 16px; background: #3498db; color: white; text-decoration: none; border-radius: 4px; border: none; cursor: pointer; }
        .btn:hover { background: #2980b9; }
        .btn.danger { background: #e74c3c; }
        .btn.danger:hover { background: #c0392b; }
        form { margin: 20px 0; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input[type="text"], input[type="email"], input[type="password"], textarea, select { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
        textarea { min-height: 100px; resize: vertical; }
        .actions { margin-top: 20px; }
        .actions a { margin-right: 10px; }
    </style>
</head>
<body>
    <header>
        <div class="container">
            <h1>` + strings.Title(g.app.Name) + `</h1>
            <nav>
                <a href="/">Home</a>
                {{range $name := .Resources}}
                <a href="/{{$name}}">{{$name | title}}</a>
                {{end}}
            </nav>
        </div>
    </header>
    
    <div class="container">
        {{template "shared/_flash.html" .Flash}}
        
        {{template "content" .Data}}
    </div>
</body>
</html>
`, nil
}

// GenerateErrors generates the errors partial.
func (g *ViewsGenerator) GenerateErrors() string {
	return `{{if .Errors}}
<div class="errors">
    <h3>Please correct the following errors:</h3>
    <ul>
        {{range .Errors}}
        <li>{{.Field}}: {{.Message}}</li>
        {{end}}
    </ul>
</div>
{{end}}`
}

// GenerateFlash generates the flash messages partial.
func (g *ViewsGenerator) GenerateFlash() string {
	return `{{if .success}}
<div class="flash success">{{.success}}</div>
{{end}}
{{if .error}}
<div class="flash error">{{.error}}</div>
{{end}}
{{if .warning}}
<div class="flash warning">{{.warning}}</div>
{{end}}
{{if .info}}
<div class="flash info">{{.info}}</div>
{{end}}`
}

// GenerateResourceViews generates all views for a resource.
func (g *ViewsGenerator) GenerateResourceViews(resource *expr.ResourceExpr) (map[string]string, error) {
	views := make(map[string]string)
	
	if resource.HasAction("index") {
		views["index.html"] = g.generateIndexView(resource)
	}
	
	if resource.HasAction("show") {
		views["show.html"] = g.generateShowView(resource)
	}
	
	if resource.HasAction("new") {
		views["new.html"] = g.generateNewView(resource)
	}
	
	if resource.HasAction("edit") {
		views["edit.html"] = g.generateEditView(resource)
	}
	
	return views, nil
}

// generateIndexView generates the index view for a resource.
func (g *ViewsGenerator) generateIndexView(resource *expr.ResourceExpr) string {
	singular := g.toSingular(resource.Name)
	
	return fmt.Sprintf(`{{define "content"}}
<div class="%s-index">
    <h1>%s</h1>
    
    <div class="actions">
        <a href="/%s/new" class="btn">New %s</a>
    </div>
    
    {{if .%s}}
    <table>
        <thead>
            <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            {{range .%s}}
            <tr>
                <td>{{.ID}}</td>
                <td>{{.Name}}</td>
                <td>
                    <a href="/%s/{{.ID}}">View</a>
                    <a href="/%s/{{.ID}}/edit">Edit</a>
                    <form method="post" action="/%s/{{.ID}}/delete" style="display:inline">
                        <button type="submit" onclick="return confirm('Are you sure?')" class="btn danger">Delete</button>
                    </form>
                </td>
            </tr>
            {{end}}
        </tbody>
    </table>
    {{else}}
    <p>No %s found.</p>
    {{end}}
</div>
{{end}}`,
		resource.Name,
		strings.Title(resource.Name),
		resource.Name,
		strings.Title(singular),
		strings.Title(resource.Name),
		strings.Title(resource.Name),
		resource.Name,
		resource.Name,
		resource.Name,
		resource.Name,
	)
}

// generateShowView generates the show view for a resource.
func (g *ViewsGenerator) generateShowView(resource *expr.ResourceExpr) string {
	singular := g.toSingular(resource.Name)
	
	return fmt.Sprintf(`{{define "content"}}
<div class="%s-show">
    <h1>%s Details</h1>
    
    {{with .%s}}
    <dl>
        <dt>ID:</dt>
        <dd>{{.ID}}</dd>
        
        <dt>Name:</dt>
        <dd>{{.Name}}</dd>
        
        <!-- Add more fields as needed -->
    </dl>
    
    <div class="actions">
        <a href="/%s/{{.ID}}/edit" class="btn">Edit</a>
        <a href="/%s">Back to List</a>
        
        <form method="post" action="/%s/{{.ID}}/delete" style="display:inline">
            <button type="submit" onclick="return confirm('Are you sure?')" class="btn danger">Delete</button>
        </form>
    </div>
    {{else}}
    <p>%s not found.</p>
    {{end}}
</div>
{{end}}`,
		singular,
		strings.Title(singular),
		strings.Title(singular),
		resource.Name,
		resource.Name,
		resource.Name,
		strings.Title(singular),
	)
}

// generateNewView generates the new view for a resource.
func (g *ViewsGenerator) generateNewView(resource *expr.ResourceExpr) string {
	singular := g.toSingular(resource.Name)
	formName := resource.NewFormName()
	
	return fmt.Sprintf(`{{define "content"}}
<div class="%s-new">
    <h1>New %s</h1>
    
    {{template "shared/_errors.html" .}}
    
    <form method="post" action="/%s">
        <div class="form-group">
            <label for="name">Name</label>
            <input type="text" id="name" name="name" value="{{.Form.Name}}" required>
        </div>
        
        <!-- Add more form fields based on your %s struct -->
        
        <div class="actions">
            <button type="submit" class="btn">Create %s</button>
            <a href="/%s">Cancel</a>
        </div>
    </form>
</div>
{{end}}`,
		singular,
		strings.Title(singular),
		resource.Name,
		formName,
		strings.Title(singular),
		resource.Name,
	)
}

// generateEditView generates the edit view for a resource.
func (g *ViewsGenerator) generateEditView(resource *expr.ResourceExpr) string {
	singular := g.toSingular(resource.Name)
	formName := resource.EditFormName()
	
	return fmt.Sprintf(`{{define "content"}}
<div class="%s-edit">
    <h1>Edit %s</h1>
    
    {{template "shared/_errors.html" .}}
    
    <form method="post" action="/%s/{{.%s.ID}}">
        <div class="form-group">
            <label for="name">Name</label>
            <input type="text" id="name" name="name" value="{{.%s.Name}}">
        </div>
        
        <!-- Add more form fields based on your %s struct -->
        
        <div class="actions">
            <button type="submit" class="btn">Update %s</button>
            <a href="/%s/{{.%s.ID}}">Cancel</a>
        </div>
    </form>
</div>
{{end}}`,
		singular,
		strings.Title(singular),
		resource.Name,
		strings.Title(singular),
		strings.Title(singular),
		formName,
		strings.Title(singular),
		resource.Name,
		strings.Title(singular),
	)
}

// toSingular converts a plural resource name to singular.
func (g *ViewsGenerator) toSingular(plural string) string {
	if strings.HasSuffix(plural, "ies") {
		return plural[:len(plural)-3] + "y"
	}
	if strings.HasSuffix(plural, "es") {
		return plural[:len(plural)-2]
	}
	if strings.HasSuffix(plural, "s") {
		return plural[:len(plural)-1]
	}
	return plural
}