package runtime

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// TemplateEngine handles template rendering.
type TemplateEngine struct {
	templates *template.Template
	funcMap   template.FuncMap
}

// NewTemplateEngine creates a new template engine.
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		funcMap: DefaultFuncMap(),
	}
}

// LoadTemplates loads templates from a directory.
func (e *TemplateEngine) LoadTemplates(viewsPath string) error {
	// Create template with function map
	tmpl := template.New("").Funcs(e.funcMap)

	// Load all templates
	patterns := []string{
		filepath.Join(viewsPath, "*.html"),
		filepath.Join(viewsPath, "*/*.html"),
		filepath.Join(viewsPath, "*/*/*.html"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		for _, file := range matches {
			// Get relative path for template name
			name, err := filepath.Rel(viewsPath, file)
			if err != nil {
				return err
			}

			// Read and parse template
			_, err = tmpl.New(name).ParseFiles(file)
			if err != nil {
				return err
			}
		}
	}

	e.templates = tmpl
	return nil
}

// Render renders a template with data.
func (e *TemplateEngine) Render(w io.Writer, name string, data any) error {
	return e.templates.ExecuteTemplate(w, name, data)
}

// DefaultFuncMap returns the default template functions.
func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		// String helpers
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,

		// Path helpers
		"link_to":  linkTo,
		"path_for": pathFor,

		// Form helpers
		"form_for":   formFor,
		"text_field": textField,
		"submit":     submitButton,

		// Formatting
		"truncate":  truncate,
		"pluralize": pluralize,

		// Safety
		"safe": safe,
		"raw":  raw,
	}
}

// Template helper functions

func linkTo(text, path string) template.HTML {
	return template.HTML(`<a href="` + template.HTMLEscapeString(path) + `">` + template.HTMLEscapeString(text) + `</a>`)
}

func pathFor(resource string, args ...any) string {
	// Simple path helper
	if len(args) == 0 {
		return "/" + resource
	}
	// Assume first arg is ID for show/edit paths
	return fmt.Sprintf("/%s/%v", resource, args[0])
}

func formFor(resource string, args ...any) template.HTML {
	method := "POST"
	action := "/" + resource

	if len(args) > 0 {
		// Assume it's an edit form
		action = fmt.Sprintf("/%s/%v", resource, args[0])
		method = "POST" // Will add _method=PUT hidden field
	}

	return template.HTML(fmt.Sprintf(`<form method="%s" action="%s">`, method, action))
}

func textField(name, value string, attrs ...string) template.HTML {
	extra := strings.Join(attrs, " ")
	return template.HTML(fmt.Sprintf(
		`<input type="text" name="%s" value="%s" %s>`,
		template.HTMLEscapeString(name),
		template.HTMLEscapeString(value),
		extra,
	))
}

func submitButton(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<button type="submit">%s</button>`,
		template.HTMLEscapeString(text),
	))
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

func safe(s string) template.HTML {
	return template.HTML(s)
}

func raw(s string) template.HTML {
	return template.HTML(s)
}
