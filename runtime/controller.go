package runtime

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

// BaseController provides common functionality for all controllers.
type BaseController struct {
	templates *template.Template
	viewsPath string
}

// NewBaseController creates a new base controller.
func NewBaseController(viewsPath string) *BaseController {
	return &BaseController{
		viewsPath: viewsPath,
	}
}

// LoadTemplates loads all templates from the views directory.
func (c *BaseController) LoadTemplates() error {
	if c.viewsPath == "" {
		c.viewsPath = "gen/webapp/views"
	}

	pattern := filepath.Join(c.viewsPath, "**/*.html")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		// Try with single level glob
		pattern = filepath.Join(c.viewsPath, "*/*.html")
		tmpl, err = template.ParseGlob(pattern)
		if err != nil {
			return fmt.Errorf("failed to load templates: %w", err)
		}
	}

	c.templates = tmpl
	return nil
}

// Render renders an HTML template with the given data.
func (c *BaseController) Render(w http.ResponseWriter, templateName string, data any) error {
	if c.templates == nil {
		if err := c.LoadTemplates(); err != nil {
			http.Error(w, "Templates not loaded", http.StatusInternalServerError)
			return err
		}
	}

	// Ensure template name has .html extension
	if !strings.HasSuffix(templateName, ".html") {
		templateName = templateName + ".html"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Wrap data with flash messages if available
	viewData := map[string]any{
		"Data":  data,
		"Flash": c.getFlash(w, nil),
	}

	if err := c.templates.ExecuteTemplate(w, templateName, viewData); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	return nil
}

// JSON sends a JSON response.
func (c *BaseController) JSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// Redirect sends an HTTP redirect response.
func (c *BaseController) Redirect(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusSeeOther)
}

// RedirectBack redirects to the referrer or fallback path.
func (c *BaseController) RedirectBack(w http.ResponseWriter, r *http.Request, fallback string) {
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = fallback
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// Bind binds form data to a struct.
func (c *BaseController) Bind(r *http.Request, dest any) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	// Simple form binding - in production, use a proper form binding library
	// This is a basic implementation for the MVP
	// TODO: Implement proper form binding with struct tags

	return nil
}

// Param gets a URL parameter value.
func (c *BaseController) Param(r *http.Request, name string) string {
	return r.PathValue(name)
}

// Params gets all URL parameters.
func (c *BaseController) Params(r *http.Request) map[string]string {
	params := make(map[string]string)

	// Get query parameters
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// In Go 1.22+, we can use r.PathValue for route params
	// For now, return query params
	return params
}

// Flash sets a flash message.
func (c *BaseController) Flash(w http.ResponseWriter, r *http.Request, level, message string) {
	cookie := &http.Cookie{
		Name:     "flash_" + level,
		Value:    message,
		Path:     "/",
		MaxAge:   0, // Session cookie
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

// getFlash retrieves and clears flash messages.
func (c *BaseController) getFlash(w http.ResponseWriter, r *http.Request) map[string]string {
	if r == nil {
		return make(map[string]string)
	}

	flash := make(map[string]string)
	levels := []string{"success", "error", "warning", "info"}

	for _, level := range levels {
		cookie, err := r.Cookie("flash_" + level)
		if err == nil && cookie.Value != "" {
			flash[level] = cookie.Value

			// Clear the flash message
			clearCookie := &http.Cookie{
				Name:     "flash_" + level,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			}
			http.SetCookie(w, clearCookie)
		}
	}

	return flash
}

// Error sends an error response.
func (c *BaseController) Error(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

// NotFound sends a 404 response.
func (c *BaseController) NotFound(w http.ResponseWriter) {
	http.NotFound(w, nil)
}

// CurrentUser is a placeholder for getting the current user.
// Applications should override this based on their auth system.
func (c *BaseController) CurrentUser(r *http.Request) any {
	// Placeholder - implement based on your auth system
	return nil
}
