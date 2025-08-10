package codegen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gluey.dev/gluey/codegen"
	"gluey.dev/gluey/expr"
)

func TestInterfaceGenerator(t *testing.T) {
	// Create a test app
	app := &expr.AppExpr{
		Name:        "testapp",
		Description: "Test application",
		Resources: []*expr.ResourceExpr{
			{
				Name:    "posts",
				Actions: []string{"index", "show", "new", "create", "edit", "update", "destroy"},
			},
			{
				Name:    "users",
				Actions: []string{"index", "show"},
			},
		},
		Pages: []*expr.PageExpr{
			{
				Name: "home",
				Routes: []expr.RouteExpr{
					{Method: "GET", Path: "/"},
				},
			},
			{
				Name: "about",
				Routes: []expr.RouteExpr{
					{Method: "GET", Path: "/about"},
				},
			},
		},
		Forms: []*expr.FormExpr{
			{
				Name: "LoginForm",
				Attributes: []*expr.AttributeExpr{
					{Name: "email", Type: expr.String},
					{Name: "password", Type: expr.String},
				},
			},
		},
	}

	// Create temp directory for output
	tmpDir := t.TempDir()

	// Generate interfaces
	gen := codegen.NewInterfaceGenerator(app, tmpDir)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Check that files were created
	expectedFiles := []string{
		"interfaces/posts_controller.go",
		"interfaces/users_controller.go",
		"interfaces/pages_controller.go",
		"types/forms.go",
		"http/router.go",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(tmpDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", file)
		}
	}

	// Check content of posts controller interface
	postsPath := filepath.Join(tmpDir, "interfaces/posts_controller.go")
	content, err := os.ReadFile(postsPath)
	if err != nil {
		t.Fatalf("Failed to read posts controller: %v", err)
	}

	postsContent := string(content)
	if !strings.Contains(postsContent, "type PostsController interface") {
		t.Error("Posts controller should define PostsController interface")
	}

	// Check that all actions are present
	actions := []string{"Index", "Show", "New", "Create", "Edit", "Update", "Destroy"}
	for _, action := range actions {
		if !strings.Contains(postsContent, action+"(w http.ResponseWriter, r *http.Request)") {
			t.Errorf("Posts controller should have %s method", action)
		}
	}

	// Check router content
	routerPath := filepath.Join(tmpDir, "http/router.go")
	routerContent, err := os.ReadFile(routerPath)
	if err != nil {
		t.Fatalf("Failed to read router: %v", err)
	}

	routerStr := string(routerContent)
	if !strings.Contains(routerStr, "type Controllers struct") {
		t.Error("Router should define Controllers struct")
	}
	if !strings.Contains(routerStr, "func MountRoutes(mux *http.ServeMux, c Controllers)") {
		t.Error("Router should define MountRoutes function")
	}
}

func TestExampleGenerator(t *testing.T) {
	// Create a test app
	app := &expr.AppExpr{
		Name: "testapp",
		Resources: []*expr.ResourceExpr{
			{
				Name:    "posts",
				Actions: []string{"index", "show", "new", "create", "edit", "update", "destroy"},
			},
		},
		Pages: []*expr.PageExpr{
			{
				Name: "home",
				Routes: []expr.RouteExpr{
					{Method: "GET", Path: "/"},
				},
			},
		},
	}

	// Create temp directory and change to it
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Generate examples
	gen := codegen.NewExampleGenerator(app)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Check that files were created
	expectedFiles := []string{
		"app/controllers/base.go",
		"app/controllers/posts.go",
		"app/controllers/pages.go",
		"app/views/layouts/application.html",
		"app/views/shared/_errors.html",
		"app/views/shared/_flash.html",
		"app/views/posts/index.html",
		"app/views/posts/show.html",
		"app/views/posts/new.html",
		"app/views/posts/edit.html",
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", file)
		}
	}

	// Test overwrite protection - run again
	err = gen.Generate()
	if err != nil {
		t.Fatalf("Second Generate() failed: %v", err)
	}

	// Files should still exist but not be overwritten
	// We can check by creating a marker file and seeing if it persists
	markerPath := "app/controllers/marker.txt"
	os.WriteFile(markerPath, []byte("test"), 0644)

	err = gen.Generate()
	if err != nil {
		t.Fatalf("Third Generate() failed: %v", err)
	}

	// Marker should still exist
	if _, err := os.Stat(markerPath); os.IsNotExist(err) {
		t.Error("Marker file was deleted, directory was recreated")
	}
}

func TestTypesGenerator(t *testing.T) {
	app := &expr.AppExpr{
		Name: "testapp",
		Forms: []*expr.FormExpr{
			{
				Name: "LoginForm",
				Attributes: []*expr.AttributeExpr{
					{
						Name: "email",
						Type: expr.String,
						Validations: []expr.Validation{
							&expr.RequiredValidation{},
							&expr.FormatValidation{Format: expr.FormatEmail},
						},
					},
					{
						Name: "password",
						Type: expr.String,
						Validations: []expr.Validation{
							&expr.RequiredValidation{},
							&expr.MinLengthValidation{Min: 8},
						},
					},
					{
						Name: "remember_me",
						Type: expr.Boolean,
					},
				},
			},
		},
	}

	gen := codegen.NewTypesGenerator(app)
	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Check that form struct is generated
	if !strings.Contains(content, "type LoginForm struct") {
		t.Error("Should generate LoginForm struct")
	}

	// Check fields
	if !strings.Contains(content, "Email string") {
		t.Error("Should have Email field")
	}
	if !strings.Contains(content, "Password string") {
		t.Error("Should have Password field")
	}
	if !strings.Contains(content, "RememberMe bool") {
		t.Error("Should have RememberMe field")
	}

	// Check validation tags
	if !strings.Contains(content, `validate:"required,email"`) {
		t.Error("Email should have required and email validation tags")
	}
	if !strings.Contains(content, `validate:"required,min=8"`) {
		t.Error("Password should have required and min=8 validation tags")
	}

	// Check Validate method
	if !strings.Contains(content, "func (f *LoginForm) Validate() error") {
		t.Error("Should generate Validate method")
	}
}

func TestViewsGenerator(t *testing.T) {
	app := &expr.AppExpr{
		Name: "testapp",
		Resources: []*expr.ResourceExpr{
			{
				Name:    "posts",
				Actions: []string{"index", "show", "new", "edit"},
			},
		},
	}

	gen := codegen.NewViewsGenerator(app)

	// Test layout generation
	layout, err := gen.GenerateLayout()
	if err != nil {
		t.Fatalf("GenerateLayout() failed: %v", err)
	}

	if !strings.Contains(layout, "<!DOCTYPE html>") {
		t.Error("Layout should contain DOCTYPE")
	}
	if !strings.Contains(layout, "{{.Title}}") {
		t.Error("Layout should contain title placeholder")
	}

	// Test resource views generation
	views, err := gen.GenerateResourceViews(app.Resources[0])
	if err != nil {
		t.Fatalf("GenerateResourceViews() failed: %v", err)
	}

	expectedViews := []string{"index.html", "show.html", "new.html", "edit.html"}
	for _, viewName := range expectedViews {
		if _, ok := views[viewName]; !ok {
			t.Errorf("Expected view %s was not generated", viewName)
		}
	}

	// Check index view content
	indexView := views["index.html"]
	if !strings.Contains(indexView, "posts-index") {
		t.Error("Index view should have posts-index class")
	}
	if !strings.Contains(indexView, "{{range .Posts}}") {
		t.Error("Index view should iterate over posts")  
	}
}