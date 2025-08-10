package dsl

import (
	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
)

// Primitive type constants
var (
	Boolean = expr.Boolean
	Int     = expr.Int
	Int32   = expr.Int32
	Int64   = expr.Int64
	Float32 = expr.Float32
	Float64 = expr.Float64
	String  = expr.String
	Bytes   = expr.Bytes
)

// Type defines a form type.
//
// Type must appear in a WebApp expression.
//
// Example:
//
//	WebApp("myapp", func() {
//	    Type("LoginForm", func() {
//	        Attribute("email", String, Required(), Format(FormatEmail))
//	        Attribute("password", String, Required(), MinLength(8))
//	        Attribute("remember_me", Boolean)
//	    })
//	})
func Type(name string, fn func()) {
	app, ok := eval.Current().(*expr.AppExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	form := &expr.FormExpr{
		Name:    name,
		DSLFunc: fn,
	}

	if fn != nil {
		eval.Execute(fn, form)
	}

	app.Forms = append(app.Forms, form)
}

// Attribute defines a field in a form.
//
// Attribute must appear in a Type expression.
//
// Example:
//
//	Type("PostForm", func() {
//	    Attribute("title", String, Required(), MaxLength(200))
//	    Attribute("content", String, Required())
//	    Attribute("published", Boolean)
//	})
func Attribute(name string, args ...interface{}) {
	form, ok := eval.Current().(*expr.FormExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	attr := &expr.AttributeExpr{
		Name: name,
	}

	// Parse arguments - type, validations, and description
	for _, arg := range args {
		switch v := arg.(type) {
		case expr.DataType:
			attr.Type = v
		case *expr.PrimitiveType:
			attr.Type = v
		case *expr.ArrayType:
			attr.Type = v
		case *expr.MapType:
			attr.Type = v
		case expr.Validation:
			attr.Validations = append(attr.Validations, v)
		case string:
			// It's a description
			attr.Description = v
		case func():
			// DSL function for nested configuration
			eval.Execute(v, attr)
		default:
			eval.InvalidArgError("type, validation, or description", arg)
		}
	}

	form.Attributes = append(form.Attributes, attr)
}

// ArrayOf creates an array type.
//
// Example:
//
//	Type("PostForm", func() {
//	    Attribute("tags", ArrayOf(String))
//	})
func ArrayOf(elemType expr.DataType) *expr.ArrayType {
	return &expr.ArrayType{
		ElemType: elemType,
	}
}

// MapOf creates a map type.
//
// Example:
//
//	Type("ConfigForm", func() {
//	    Attribute("settings", MapOf(String, String))
//	})
func MapOf(keyType, elemType expr.DataType) *expr.MapType {
	return &expr.MapType{
		KeyType:  keyType,
		ElemType: elemType,
	}
}