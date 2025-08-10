package codegen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobijan/gluey/expr"
)

// Generator is the main code generator.
type Generator struct {
	app        *expr.AppExpr
	outputPath string
}

// NewGenerator creates a new generator.
func NewGenerator(app *expr.AppExpr, outputPath string) *Generator {
	if outputPath == "" {
		outputPath = "gen/webapp"
	}
	return &Generator{
		app:        app,
		outputPath: outputPath,
	}
}

// Generate generates all code for the application.
func (g *Generator) Generate() error {
	if g.app == nil {
		return fmt.Errorf("no application expression found")
	}

	// Create output directory structure
	if err := g.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate types (forms)
	if err := g.generateTypes(); err != nil {
		return fmt.Errorf("failed to generate types: %w", err)
	}

	// Generate controllers
	if err := g.generateControllers(); err != nil {
		return fmt.Errorf("failed to generate controllers: %w", err)
	}

	// Generate router
	if err := g.generateRouter(); err != nil {
		return fmt.Errorf("failed to generate router: %w", err)
	}

	// Generate views
	if err := g.generateViews(); err != nil {
		return fmt.Errorf("failed to generate views: %w", err)
	}

	return nil
}

// createDirectories creates the output directory structure.
func (g *Generator) createDirectories() error {
	dirs := []string{
		g.outputPath,
		filepath.Join(g.outputPath, g.app.Name),
		filepath.Join(g.outputPath, g.app.Name, "controllers"),
		filepath.Join(g.outputPath, g.app.Name, "views"),
		filepath.Join(g.outputPath, g.app.Name, "views", "layouts"),
		filepath.Join(g.outputPath, g.app.Name, "views", "shared"),
	}

	// Add directories for each resource
	for _, resource := range g.app.Resources {
		dir := filepath.Join(g.outputPath, g.app.Name, "views", resource.Name)
		dirs = append(dirs, dir)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// generateTypes generates form types.
func (g *Generator) generateTypes() error {
	gen := NewTypesGenerator(g.app)
	code, err := gen.Generate()
	if err != nil {
		return err
	}

	outputFile := filepath.Join(g.outputPath, g.app.Name, "types.go")
	return os.WriteFile(outputFile, []byte(code), 0644)
}

// generateControllers generates controller interfaces.
func (g *Generator) generateControllers() error {
	gen := NewControllersGenerator(g.app)

	// Generate base controller
	baseCode, err := gen.GenerateBase()
	if err != nil {
		return err
	}

	baseFile := filepath.Join(g.outputPath, g.app.Name, "controllers", "base.go")
	if err := os.WriteFile(baseFile, []byte(baseCode), 0644); err != nil {
		return err
	}

	// Generate resource controllers
	for _, resource := range g.app.Resources {
		code, err := gen.GenerateResource(resource)
		if err != nil {
			return err
		}

		file := filepath.Join(g.outputPath, g.app.Name, "controllers", resource.Name+".go")
		if err := os.WriteFile(file, []byte(code), 0644); err != nil {
			return err
		}
	}

	// Generate page controllers if any
	if len(g.app.Pages) > 0 {
		pagesCode, err := gen.GeneratePages(g.app.Pages)
		if err != nil {
			return err
		}

		pagesFile := filepath.Join(g.outputPath, g.app.Name, "controllers", "pages.go")
		if err := os.WriteFile(pagesFile, []byte(pagesCode), 0644); err != nil {
			return err
		}
	}

	return nil
}

// generateRouter generates the router.
func (g *Generator) generateRouter() error {
	gen := NewRouterGenerator(g.app)
	code, err := gen.Generate()
	if err != nil {
		return err
	}

	outputFile := filepath.Join(g.outputPath, g.app.Name, "router.go")
	return os.WriteFile(outputFile, []byte(code), 0644)
}

// generateViews generates HTML templates.
func (g *Generator) generateViews() error {
	gen := NewViewsGenerator(g.app)

	// Generate layout
	layoutCode, err := gen.GenerateLayout()
	if err != nil {
		return err
	}

	layoutFile := filepath.Join(g.outputPath, g.app.Name, "views", "layouts", "application.html")
	if err := os.WriteFile(layoutFile, []byte(layoutCode), 0644); err != nil {
		return err
	}

	// Generate shared partials
	sharedFiles := map[string]string{
		"_errors.html": gen.GenerateErrors(),
		"_flash.html":  gen.GenerateFlash(),
	}

	for name, content := range sharedFiles {
		file := filepath.Join(g.outputPath, g.app.Name, "views", "shared", name)
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Generate resource views
	for _, resource := range g.app.Resources {
		views, err := gen.GenerateResourceViews(resource)
		if err != nil {
			return err
		}

		for name, content := range views {
			file := filepath.Join(g.outputPath, g.app.Name, "views", resource.Name, name)
			if err := os.WriteFile(file, []byte(content), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
