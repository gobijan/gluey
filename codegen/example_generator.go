package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobijan/gluey/expr"
)

// ExampleGenerator generates example implementations.
// This is used by the 'example' command and only creates files that don't exist.
type ExampleGenerator struct {
	app       *expr.AppExpr
	OutputDir string // Base output directory (defaults to ".")
}

// NewExampleGenerator creates a new example generator.
func NewExampleGenerator(app *expr.AppExpr) *ExampleGenerator {
	return &ExampleGenerator{
		app:       app,
		OutputDir: ".",
	}
}

// Generate generates example implementations.
func (g *ExampleGenerator) Generate() error {
	// Create app directories if they don't exist
	dirs := []string{
		filepath.Join(g.OutputDir, "app/controllers"),
		filepath.Join(g.OutputDir, "app/views/layouts"),
		filepath.Join(g.OutputDir, "app/views/shared"),
	}

	// Add view directories for each resource
	for _, resource := range g.app.Resources {
		dirs = append(dirs, filepath.Join(g.OutputDir, "app/views", resource.Name))
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate base controller (if doesn't exist)
	if err := g.generateBaseController(); err != nil {
		return err
	}

	// Generate example controllers (if don't exist)
	for _, resource := range g.app.Resources {
		if err := g.generateResourceController(resource); err != nil {
			return err
		}
	}

	// Generate pages controller (if doesn't exist)
	if len(g.app.Pages) > 0 {
		if err := g.generatePagesController(); err != nil {
			return err
		}
	}

	// Generate views (if don't exist)
	if err := g.generateViews(); err != nil {
		return err
	}

	fmt.Println("✅ Example files generated in app/")
	fmt.Println("\nCreated:")
	fmt.Println("  - app/controllers/ - Example controller implementations")
	fmt.Println("  - app/views/ - HTML templates")
	fmt.Println("\nThese files are yours to modify. They won't be overwritten.")

	return nil
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// generateBaseController generates the base controller if it doesn't exist.
func (g *ExampleGenerator) generateBaseController() error {
	filename := filepath.Join(g.OutputDir, "app/controllers/base.go")
	if fileExists(filename) {
		fmt.Printf("  Skipping %s (already exists)\n", filename)
		return nil
	}

	content := `package controllers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// BaseController provides common functionality for all controllers.
type BaseController struct {
	templates *template.Template
}

// NewBaseController creates a new base controller.
func NewBaseController() *BaseController {
	// Load templates with custom functions
	funcMap := template.FuncMap{
		"title": strings.Title,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}
	
	// Load all template files
	tmpl := template.New("").Funcs(funcMap)
	
	// Load layout templates
	layoutFiles, _ := filepath.Glob("app/views/layouts/*.html")
	for _, file := range layoutFiles {
		content, err := os.ReadFile(file)
		if err == nil {
			name := filepath.Base(file)
			template.Must(tmpl.New(name).Parse(string(content)))
		}
	}
	
	// Load shared templates  
	sharedFiles, _ := filepath.Glob("app/views/shared/*.html")
	for _, file := range sharedFiles {
		content, err := os.ReadFile(file)
		if err == nil {
			name := filepath.Base(file)
			template.Must(tmpl.New(name).Parse(string(content)))
		}
	}
	
	// Load view templates
	viewFiles, _ := filepath.Glob("app/views/*/*.html")
	for _, file := range viewFiles {
		if !strings.Contains(file, "/layouts/") && !strings.Contains(file, "/shared/") {
			content, err := os.ReadFile(file)
			if err == nil {
				// Use relative path as template name (e.g., "posts/index.html")
				name := strings.TrimPrefix(file, "app/views/")
				template.Must(tmpl.New(name).Parse(string(content)))
			}
		}
	}
	
	return &BaseController{
		templates: tmpl,
	}
}

// Render renders a template with the given data.
func (c *BaseController) Render(w http.ResponseWriter, view string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	
	// Add common data
	data["AppName"] = "` + ToTitle(g.app.Name) + `"
	
	// Render the view content first
	var contentBuf strings.Builder
	if c.templates != nil {
		err := c.templates.ExecuteTemplate(&contentBuf, view, data)
		if err != nil {
			// If the specific view doesn't exist, just use empty content
			contentBuf.Reset()
		}
	}
	
	// Add the rendered content to data
	data["Content"] = template.HTML(contentBuf.String())
	
	// Execute the layout with the rendered content
	if c.templates != nil {
		err := c.templates.ExecuteTemplate(w, "application.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Templates not loaded", http.StatusInternalServerError)
	}
}

// Redirect redirects to the given URL.
func (c *BaseController) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// Flash sets a flash message cookie.
func (c *BaseController) Flash(w http.ResponseWriter, typ, message string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "flash_" + typ,
		Value: message,
		Path:  "/",
	})
}
`

	fmt.Printf("  Creating %s\n", filename)
	return os.WriteFile(filename, []byte(content), 0644)
}

// generateResourceController generates an example controller for a resource.
func (g *ExampleGenerator) generateResourceController(resource *expr.ResourceExpr) error {
	filename := filepath.Join(g.OutputDir, fmt.Sprintf("app/controllers/%s.go", resource.Name))
	if fileExists(filename) {
		fmt.Printf("  Skipping %s (already exists)\n", filename)
		return nil
	}

	singular := toSingular(resource.Name)
	controllerType := resource.Name + "Controller"

	content := fmt.Sprintf(`package controllers

import (
	"net/http"
	"%s/gen/interfaces"
)

// %s handles requests for %s resources.
type %s struct {
	BaseController
}

// New%s creates a new %s controller.
func New%s() interfaces.%sController {
	return &%s{
		BaseController: *NewBaseController(),
	}
}

// Index displays a list of %s
func (c *%s) Index(w http.ResponseWriter, r *http.Request) {
	// TODO: Fetch %s from database
	%s := []map[string]interface{}{
		{"ID": 1, "Name": "Sample %s 1"},
		{"ID": 2, "Name": "Sample %s 2"},
	}
	
	c.Render(w, "%s/index", map[string]interface{}{
		"Title": "%s",
		"%s": %s,
	})
}

// Show displays a single %s
func (c *%s) Show(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	// TODO: Fetch %s from database
	%s := map[string]interface{}{
		"ID": id,
		"Name": "Sample %s",
	}
	
	c.Render(w, "%s/show", map[string]interface{}{
		"Title": "%s Details",
		"%s": %s,
	})
}

// New displays the form for creating a new %s
func (c *%s) New(w http.ResponseWriter, r *http.Request) {
	c.Render(w, "%s/new", map[string]interface{}{
		"Title": "New %s",
	})
}

// Create handles the creation of a new %s
func (c *%s) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse form, validate, and save to database
	
	c.Flash(w, "success", "%s created successfully!")
	c.Redirect(w, r, "/%s")
}

// Edit displays the form for editing a %s
func (c *%s) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	// TODO: Fetch %s from database
	%s := map[string]interface{}{
		"ID": id,
		"Name": "Sample %s",
	}
	
	c.Render(w, "%s/edit", map[string]interface{}{
		"Title": "Edit %s",
		"%s": %s,
	})
}

// Update handles updating a %s
func (c *%s) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	// TODO: Parse form, validate, and update in database
	
	c.Flash(w, "success", "%s updated successfully!")
	c.Redirect(w, r, "/%s/"+id)
}

// Destroy handles deleting a %s
func (c *%s) Destroy(w http.ResponseWriter, r *http.Request) {
	// TODO: Delete from database
	
	c.Flash(w, "success", "%s deleted successfully!")
	c.Redirect(w, r, "/%s")
}
`,
		g.app.Name,
		controllerType, resource.Name,
		controllerType,
		ToTitle(resource.Name), resource.Name,
		ToTitle(resource.Name), ToTitle(resource.Name),
		controllerType,
		resource.Name,
		controllerType,
		resource.Name,
		resource.Name,
		ToTitle(singular),
		ToTitle(singular),
		resource.Name,
		ToTitle(resource.Name),
		ToTitle(resource.Name), resource.Name,
		resource.Name,
		controllerType,
		singular,
		singular,
		ToTitle(singular),
		resource.Name,
		ToTitle(singular),
		ToTitle(singular), singular,
		resource.Name,
		controllerType,
		resource.Name,
		ToTitle(singular),
		resource.Name,
		controllerType,
		ToTitle(singular),
		resource.Name,
		resource.Name,
		controllerType,
		singular,
		singular,
		ToTitle(singular),
		resource.Name,
		ToTitle(singular),
		ToTitle(singular), singular,
		resource.Name,
		controllerType,
		ToTitle(singular),
		resource.Name,
		resource.Name,
		controllerType,
		ToTitle(singular),
		resource.Name,
	)

	fmt.Printf("  Creating %s\n", filename)
	return os.WriteFile(filename, []byte(content), 0644)
}

// generatePagesController generates an example pages controller.
func (g *ExampleGenerator) generatePagesController() error {
	filename := filepath.Join(g.OutputDir, "app/controllers/pages.go")
	if fileExists(filename) {
		fmt.Printf("  Skipping %s (already exists)\n", filename)
		return nil
	}

	content := fmt.Sprintf(`package controllers

import (
	"net/http"
	"%s/gen/interfaces"
)

// pagesController handles static page requests.
type pagesController struct {
	BaseController
}

// NewPagesController creates a new pages controller.
func NewPagesController() interfaces.PagesController {
	return &pagesController{
		BaseController: *NewBaseController(),
	}
}
`, g.app.Name)

	// Add method for each page
	for _, page := range g.app.Pages {
		for _, route := range page.Routes {
			methodName := toTitle(page.Name)
			if route.Method != "GET" {
				methodName += toTitle(strings.ToLower(route.Method))
			}

			content += fmt.Sprintf(`
// %s handles %s %s
func (c *pagesController) %s(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement %s page
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>%s Page</h1><p>TODO: Implement this page</p>"))
}
`, methodName, route.Method, route.Path, methodName, page.Name, ToTitle(page.Name))
		}
	}

	fmt.Printf("  Creating %s\n", filename)
	return os.WriteFile(filename, []byte(content), 0644)
}

// generateViews generates all view templates.
func (g *ExampleGenerator) generateViews() error {
	// Generate layout
	if err := g.generateLayout(); err != nil {
		return err
	}

	// Generate shared partials
	if err := g.generateSharedViews(); err != nil {
		return err
	}

	// Generate resource views
	for _, resource := range g.app.Resources {
		if err := g.generateResourceViews(resource); err != nil {
			return err
		}
	}

	return nil
}

// generateLayout generates the main layout template.
func (g *ExampleGenerator) generateLayout() error {
	filename := filepath.Join(g.OutputDir, "app/views/layouts/application.html")
	if fileExists(filename) {
		fmt.Printf("  Skipping %s (already exists)\n", filename)
		return nil
	}

	viewGen := NewViewsGenerator(g.app)
	content, _ := viewGen.GenerateLayout()

	fmt.Printf("  Creating %s\n", filename)
	return os.WriteFile(filename, []byte(content), 0644)
}

// generateSharedViews generates shared view partials.
func (g *ExampleGenerator) generateSharedViews() error {
	viewGen := NewViewsGenerator(g.app)

	// Generate errors partial
	filename := filepath.Join(g.OutputDir, "app/views/shared/_errors.html")
	if !fileExists(filename) {
		content := viewGen.GenerateErrors()
		fmt.Printf("  Creating %s\n", filename)
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Generate flash partial
	filename = filepath.Join(g.OutputDir, "app/views/shared/_flash.html")
	if !fileExists(filename) {
		content := viewGen.GenerateFlash()
		fmt.Printf("  Creating %s\n", filename)
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// generateResourceViews generates views for a resource.
func (g *ExampleGenerator) generateResourceViews(resource *expr.ResourceExpr) error {
	viewGen := NewViewsGenerator(g.app)
	views, err := viewGen.GenerateResourceViews(resource)
	if err != nil {
		return err
	}

	for name, content := range views {
		filename := filepath.Join(g.OutputDir, "app/views", resource.Name, name)
		if fileExists(filename) {
			fmt.Printf("  Skipping %s (already exists)\n", filename)
			continue
		}

		// Update template definition
		content = strings.Replace(content,
			`{{define "content"}}`,
			fmt.Sprintf(`{{define "%s/%s"}}`, resource.Name, strings.TrimSuffix(name, ".html")),
			1)

		fmt.Printf("  Creating %s\n", filename)
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// toSingular converts a plural resource name to singular.
func toSingular(plural string) string {
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
