package expr_test

import (
	"testing"

	"github.com/gobijan/gluey/expr"
)

func TestAppExpr(t *testing.T) {
	app := &expr.AppExpr{
		Name:        "testapp",
		Description: "Test application",
	}

	// Test EvalName
	if app.EvalName() != "testapp" {
		t.Errorf("EvalName() = %v, want %v", app.EvalName(), "testapp")
	}

	// Test adding resources
	resource := &expr.ResourceExpr{Name: "posts"}
	app.Resources = append(app.Resources, resource)

	if len(app.Resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(app.Resources))
	}

	// Test Resource() method
	found := app.Resource("posts")
	if found != resource {
		t.Error("Resource() did not return the correct resource")
	}

	notFound := app.Resource("users")
	if notFound != nil {
		t.Error("Resource() should return nil for non-existent resource")
	}

	// Test adding pages
	page := &expr.PageExpr{Name: "home"}
	app.Pages = append(app.Pages, page)

	// Check page exists
	if len(app.Pages) != 1 {
		t.Error("Page was not added")
	}
	if app.Pages[0].Name != "home" {
		t.Error("Page name incorrect")
	}

	// Test adding forms
	form := &expr.FormExpr{Name: "LoginForm"}
	app.Forms = append(app.Forms, form)

	found3 := app.Form("LoginForm")
	if found3 != form {
		t.Error("Form() did not return the correct form")
	}
}

func TestResourceExpr(t *testing.T) {
	resource := &expr.ResourceExpr{
		Name:    "posts",
		Actions: []string{"index", "show", "create"},
	}

	// Test EvalName
	if resource.EvalName() != "posts" {
		t.Errorf("EvalName() = %v, want %v", resource.EvalName(), "posts")
	}

	// Test HasAction
	if !resource.HasAction("index") {
		t.Error("HasAction('index') should return true")
	}

	if resource.HasAction("delete") {
		t.Error("HasAction('delete') should return false")
	}

	// Test form naming conventions
	if resource.NewFormName() != "NewPostsForm" {
		t.Errorf("NewFormName() = %v, want %v", resource.NewFormName(), "NewPostsForm")
	}

	if resource.EditFormName() != "EditPostsForm" {
		t.Errorf("EditFormName() = %v, want %v", resource.EditFormName(), "EditPostsForm")
	}

	// Test Validate
	err := resource.Validate()
	if err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}

	// Test validation with empty name
	emptyResource := &expr.ResourceExpr{Name: ""}
	err = emptyResource.Validate()
	if err == nil {
		t.Error("Validate() should return error for empty name")
	}
}

func TestPageExpr(t *testing.T) {
	page := &expr.PageExpr{
		Name: "home",
		Routes: []expr.RouteExpr{
			{Method: "GET", Path: "/"},
			{Method: "POST", Path: "/contact"},
		},
	}

	// Test EvalName
	if page.EvalName() != "home" {
		t.Errorf("EvalName() = %v, want %v", page.EvalName(), "home")
	}

	// Test validation
	err := page.Validate()
	if err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}

	// Test validation with empty name
	emptyPage := &expr.PageExpr{Name: ""}
	err = emptyPage.Validate()
	if err == nil {
		t.Error("Validate() should return error for empty name")
	}

	// Test validation with no routes
	noRoutePage := &expr.PageExpr{Name: "about"}
	err = noRoutePage.Validate()
	if err == nil {
		t.Error("Validate() should return error for page with no routes")
	}
}

func TestFormExpr(t *testing.T) {
	form := &expr.FormExpr{
		Name: "LoginForm",
		Attributes: []*expr.AttributeExpr{
			{
				Name: "email",
				Type: expr.String,
			},
			{
				Name: "password",
				Type: expr.String,
			},
		},
	}

	// Test EvalName
	if form.EvalName() != "LoginForm" {
		t.Errorf("EvalName() = %v, want %v", form.EvalName(), "LoginForm")
	}

	// Test Attribute lookup
	attr := form.Attribute("email")
	if attr == nil {
		t.Error("Attribute('email') should return the email attribute")
	}

	notFound := form.Attribute("username")
	if notFound != nil {
		t.Error("Attribute('username') should return nil")
	}

	// Test validation
	err := form.Validate()
	if err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}

	// Test validation with empty name
	emptyForm := &expr.FormExpr{Name: ""}
	err = emptyForm.Validate()
	if err == nil {
		t.Error("Validate() should return error for empty name")
	}
}

func TestAttributeExpr(t *testing.T) {
	attr := &expr.AttributeExpr{
		Name:        "email",
		Type:        expr.String,
		Description: "User email",
		Validations: []expr.Validation{
			&expr.RequiredValidation{},
			&expr.FormatValidation{Format: expr.FormatEmail},
		},
	}

	// Test EvalName
	if attr.EvalName() != "email" {
		t.Errorf("EvalName() = %v, want %v", attr.EvalName(), "email")
	}

	// Test IsRequired
	if !attr.IsRequired() {
		t.Error("IsRequired() should return true")
	}

	// Test Format
	format, ok := attr.Format()
	if !ok || format != expr.FormatEmail {
		t.Error("Format() should return FormatEmail")
	}

	// Test validation
	err := attr.Validate()
	if err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}

	// Test validation with empty name
	emptyAttr := &expr.AttributeExpr{Name: ""}
	err = emptyAttr.Validate()
	if err == nil {
		t.Error("Validate() should return error for empty name")
	}
}

func TestPrimitiveTypes(t *testing.T) {
	tests := []struct {
		typ      expr.DataType
		expected string
	}{
		{expr.Boolean, "boolean"},
		{expr.Int, "int"},
		{expr.Int32, "int32"},
		{expr.Int64, "int64"},
		{expr.Float32, "float32"},
		{expr.Float64, "float64"},
		{expr.String, "string"},
		{expr.Bytes, "bytes"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.typ.Name() != tt.expected {
				t.Errorf("Name() = %v, want %v", tt.typ.Name(), tt.expected)
			}
		})
	}
}

func TestArrayType(t *testing.T) {
	arr := &expr.ArrayType{
		ElemType: expr.String,
	}

	if arr.Name() != "array<string>" {
		t.Errorf("Name() = %v, want %v", arr.Name(), "array<string>")
	}
}

func TestMapType(t *testing.T) {
	m := &expr.MapType{
		KeyType:  expr.String,
		ElemType: expr.Int,
	}

	if m.Name() != "map<string,int>" {
		t.Errorf("Name() = %v, want %v", m.Name(), "map<string,int>")
	}
}

func TestReset(t *testing.T) {
	// Set a root
	expr.Root = &expr.AppExpr{Name: "test"}

	// Reset should clear it
	expr.Reset()

	if expr.Root != nil {
		t.Error("Reset() should set Root to nil")
	}
}
