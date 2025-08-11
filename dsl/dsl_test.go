package dsl_test

import (
	"testing"

	"github.com/gobijan/gluey/dsl"
	"github.com/gobijan/gluey/eval"
	"github.com/gobijan/gluey/expr"
)

func TestWebApp(t *testing.T) {
	// Reset before test
	expr.Reset()
	eval.Context.Reset()

	// Create a web app
	app := dsl.WebApp("testapp", func() {
		// This should be executed during eval.RunDSL()
	})

	if app == nil {
		t.Fatal("WebApp() returned nil")
	}

	if app.Name != "testapp" {
		t.Errorf("App name = %v, want %v", app.Name, "testapp")
	}

	if expr.Root != app {
		t.Error("WebApp() did not set expr.Root")
	}

	if eval.CurrentRoot() != app {
		t.Error("WebApp() did not set eval root")
	}
}

func TestDescription(t *testing.T) {
	expr.Reset()
	eval.Context.Reset()

	executed := false
	app := dsl.WebApp("testapp", func() {
		dsl.Description("Test application")
		executed = true
	})

	// Execute the DSL
	err := eval.RunDSL()
	if err != nil {
		t.Fatalf("RunDSL() failed: %v", err)
	}

	if !executed {
		t.Error("WebApp DSL function was not executed")
	}

	if app.Description != "Test application" {
		t.Errorf("Description = %v, want %v", app.Description, "Test application")
	}
}

func TestResource(t *testing.T) {
	expr.Reset()
	eval.Context.Reset()

	dsl.WebApp("testapp", func() {
		dsl.Resource("posts")
		dsl.Resource("users", func() {
			// With configuration function
		})
	})

	err := eval.RunDSL()
	if err != nil {
		t.Fatalf("RunDSL() failed: %v", err)
	}

	app := expr.Root

	if len(app.Resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(app.Resources))
	}

	posts := app.Resource("posts")
	if posts == nil {
		t.Error("posts resource not found")
	}

	users := app.Resource("users")
	if users == nil {
		t.Error("users resource not found")
	}
}

func TestPage(t *testing.T) {
	expr.Reset()
	eval.Context.Reset()

	dsl.WebApp("testapp", func() {
		dsl.Page("home", "/")
		dsl.Page("about", "/about")
	})

	err := eval.RunDSL()
	if err != nil {
		t.Fatalf("RunDSL() failed: %v", err)
	}

	app := expr.Root

	if len(app.Pages) != 2 {
		t.Errorf("Expected 2 pages, got %d", len(app.Pages))
	}

	// Find pages by name
	var home, about *expr.PageExpr
	for _, p := range app.Pages {
		switch p.Name {
		case "home":
			home = p
		case "about":
			about = p
		}
	}

	if home == nil {
		t.Error("home page not found")
	} else if len(home.Routes) != 1 || home.Routes[0].Path != "/" {
		t.Error("home page route incorrect")
	}

	if about == nil {
		t.Error("about page not found")
	} else if len(about.Routes) != 1 || about.Routes[0].Path != "/about" {
		t.Error("about page route incorrect")
	}
}

func TestType(t *testing.T) {
	expr.Reset()
	eval.Context.Reset()

	dsl.WebApp("testapp", func() {
		dsl.Type("LoginForm", func() {
			dsl.Attribute("email", dsl.String, dsl.Required())
			dsl.Attribute("password", dsl.String, dsl.Required(), dsl.MinLength(8))
			dsl.Attribute("remember_me", dsl.Boolean)
		})
	})

	err := eval.RunDSL()
	if err != nil {
		t.Fatalf("RunDSL() failed: %v", err)
	}

	app := expr.Root

	if len(app.Forms) != 1 {
		t.Fatalf("Expected 1 form, got %d", len(app.Forms))
	}

	form := app.Form("LoginForm")
	if form == nil {
		t.Fatal("LoginForm not found")
	}

	if len(form.Attributes) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(form.Attributes))
	}

	// Check email attribute
	email := form.Attribute("email")
	if email == nil {
		t.Error("email attribute not found")
	}
	if !email.IsRequired() {
		t.Error("email should be required")
	}

	// Check password attribute
	password := form.Attribute("password")
	if password == nil {
		t.Error("password attribute not found")
	}
	if !password.IsRequired() {
		t.Error("password should be required")
	}
	min, ok := password.MinLength()
	if !ok || min != 8 {
		t.Error("password should have min length of 8")
	}

	// Check remember_me attribute
	rememberMe := form.Attribute("remember_me")
	if rememberMe == nil {
		t.Error("remember_me attribute not found")
	} else if rememberMe.Type != dsl.Boolean {
		t.Error("remember_me should be boolean type")
	}
}

func TestValidations(t *testing.T) {
	// Test Required
	req := dsl.Required()
	if _, ok := req.(*expr.RequiredValidation); !ok {
		t.Error("Required() should return RequiredValidation")
	}

	// Test MinLength
	min := dsl.MinLength(5)
	if minVal, ok := min.(*expr.MinLengthValidation); !ok || minVal.Min != 5 {
		t.Error("MinLength(5) should return MinLengthValidation with Min=5")
	}

	// Test MaxLength
	max := dsl.MaxLength(100)
	if maxVal, ok := max.(*expr.MaxLengthValidation); !ok || maxVal.Max != 100 {
		t.Error("MaxLength(100) should return MaxLengthValidation with Max=100")
	}

	// Test Format
	email := dsl.Format(dsl.FormatEmail)
	if fmtVal, ok := email.(*expr.FormatValidation); !ok || fmtVal.Format != expr.FormatEmail {
		t.Error("Format(FormatEmail) should return FormatValidation with Format=FormatEmail")
	}
}

func TestArrayOf(t *testing.T) {
	arr := dsl.ArrayOf(dsl.String)
	if arr.ElemType != dsl.String {
		t.Error("ArrayOf(String) should have String element type")
	}
}

func TestMapOf(t *testing.T) {
	m := dsl.MapOf(dsl.String, dsl.Int)
	if m.KeyType != dsl.String {
		t.Error("MapOf key type should be String")
	}
	if m.ElemType != dsl.Int {
		t.Error("MapOf element type should be Int")
	}
}

func TestIncompatibleDSLContext(t *testing.T) {
	expr.Reset()
	eval.Context.Reset()

	// Try to call Resource outside of WebApp
	dsl.Resource("posts")

	// This should have recorded an error
	if eval.Context.Errors == nil {
		t.Error("Resource outside WebApp should record an error")
	}
}

func TestNestedResource(t *testing.T) {
	expr.Reset()
	eval.Context.Reset()

	dsl.WebApp("testapp", func() {
		dsl.Resource("posts", func() {
			dsl.Actions("index", "show")
			dsl.Resource("comments") // Nested resource
		})
	})

	err := eval.RunDSL()
	if err != nil {
		t.Fatalf("RunDSL() failed: %v", err)
	}

	app := expr.Root

	// Should have 2 resources (posts and comments)
	if len(app.Resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(app.Resources))
	}

	posts := app.Resource("posts")
	if posts == nil {
		t.Fatal("posts resource not found")
	}

	if len(posts.Actions) != 2 {
		t.Errorf("Expected 2 actions for posts, got %d", len(posts.Actions))
	}

	comments := app.Resource("comments")
	if comments == nil {
		t.Fatal("comments resource not found")
	}

	if comments.Parent != posts {
		t.Error("comments should have posts as parent")
	}
}
