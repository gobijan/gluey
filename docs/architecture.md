# Gluey Architecture

## Overview

Gluey is a web framework for Go that brings Rails-like productivity through code generation from a declarative DSL.
This document describes the internal architecture and design decisions.

## Core Philosophy

1. **DSL → Structure → Implementation**: Define structure in DSL, generate boilerplate, developers implement business logic
2. **Convention Over Configuration**: Smart defaults with escape hatches
3. **Zero Magic**: All generated code is readable, debuggable Go
4. **Type Safety**: Compile-time checking throughout

## Architecture Layers

```
┌─────────────────────────────────────────┐
│           User DSL File                  │
│         (design/app.go)                  │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│          DSL Functions                   │
│           (dsl/*.go)                     │
│   WebApp(), Resource(), Type(), etc.     │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│        Expression Tree                   │
│          (expr/*.go)                     │
│   AppExpr, ResourceExpr, FormExpr, etc.  │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│       Evaluation Pipeline                │
│          (eval/*.go)                     │
│   Execute → Prepare → Validate → Finalize│
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│        Code Generation                   │
│        (codegen/*.go)                    │
│   Types, Controllers, Router, Views      │
└─────────────────┬───────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│        Generated Code                    │
│         (gen/webapp/*)                   │
│   Ready for developer implementation     │
└─────────────────────────────────────────┘
```

## Package Structure

### `/dsl` - Domain Specific Language

The public API that users interact with. Each file provides DSL functions for specific aspects:

- `app.go` - WebApp() function, top-level app definition
- `resource.go` - Resource() for RESTful resources
- `page.go` - Page() for non-resource pages
- `types.go` - Type(), Attribute() for form definitions
- `validation.go` - Required(), MaxLength(), Format() validators
- `layout.go` - Layout management functions
- `middleware.go` - Middleware composition
- `auth.go` - Authentication/authorization DSL

Example usage:

```go
var _ = gluey.WebApp("myapp", func() {
    Resource("posts")
})
```

### `/expr` - Expression Types

Internal representation of the DSL. Each expression type corresponds to a DSL concept:

- `app.go` - AppExpr struct, root of expression tree
- `resource.go` - ResourceExpr for resources
- `page.go` - PageExpr for pages
- `form.go` - FormExpr for form types
- `attribute.go` - AttributeExpr for fields
- `root.go` - Global Root variable holding the expression tree
- `types.go` - Type system definitions

Expression types implement interfaces from `/eval`:

```go
type AppExpr struct {
    Name        string
    Description string
    Resources   []*ResourceExpr
    Pages       []*PageExpr
    Forms       []*FormExpr
}

func (a *AppExpr) EvalName() string { return a.Name }
func (a *AppExpr) Prepare() { /* ... */ }
func (a *AppExpr) Validate() error { /* ... */ }
```

### `/eval` - Evaluation Engine

The runtime that processes the DSL:

- `eval.go` - RunDSL() orchestrates the pipeline
- `expression.go` - Core interfaces (Expression, Preparer, Validator, Finalizer)
- `context.go` - Evaluation context and state
- `error.go` - Error collection and reporting

Pipeline phases:

1. **Execute**: Run DSL functions, build expression tree
2. **Prepare**: Process expressions (e.g., resolve references)
3. **Validate**: Check for errors and inconsistencies
4. **Finalize**: Last preparations before code generation

### `/codegen` - Code Generation

Transforms expressions into Go code:

- `generator.go` - Main generation orchestration
- `types.go` - Generate form structs with validation
- `controllers.go` - Generate controller interfaces
- `router.go` - Generate HTTP router setup
- `views.go` - Generate HTML templates
- `helpers.go` - Generate BaseController and helpers
- `templates/` - Go templates for code generation

Generation strategy:

```go
func Generate(app *expr.AppExpr) error {
    generateTypes(app.Forms)
    generateControllers(app.Resources)
    generateRouter(app)
    generateViews(app)
    generateHelpers(app)
}
```

### `/runtime` - Runtime Support Library

Code that generated apps depend on:

- `controller.go` - BaseController with common helpers
- `middleware.go` - Built-in middleware implementations
- `flash.go` - Flash message handling
- `pagination.go` - Pagination utilities
- `binding.go` - Form binding from requests
- `validation.go` - Runtime validation execution

Generated code imports this package:

```go
import "gluey.dev/gluey/runtime"

type PostsController struct {
    *runtime.BaseController
}
```

### `/cmd/gluey` - CLI Tool

Command-line interface for users:

- `main.go` - CLI entry point
- `gen.go` - Generate command implementation
- `new.go` - Scaffold new project

Commands:

```bash
gluey new myapp      # Create new project
gluey gen            # Generate code from DSL
gluey version        # Show version
```

## Code Generation Strategy

### Conventions

1. **Form Naming**: `New{Resource}Form` for create, `Edit{Resource}Form` for update
2. **Routes**: RESTful routes automatically generated
3. **Templates**: One template per action in `views/{resource}/`
4. **Controllers**: Interface generated, implementation by developer

### Generated Structure

```
gen/webapp/{app_name}/
├── types.go              # Form structs
├── router.go             # Route definitions
├── controllers/
│   ├── base.go          # BaseController
│   ├── posts.go         # PostsController interface
│   └── pages.go         # PagesController interface
└── views/
    ├── layouts/
    │   └── application.html
    └── posts/
        ├── index.html
        ├── show.html
        ├── new.html
        └── edit.html
```

### Type Generation

Forms become structs with validation tags:

```go
// DSL
Type("NewPostForm", func() {
    Attribute("title", String, Required())
})

// Generated
type NewPostForm struct {
    Title string `form:"title" validate:"required"`
}
```

### Controller Generation

Resources become interfaces:

```go
// DSL
Resource("posts")

// Generated
type PostsController interface {
    Index(w http.ResponseWriter, r *http.Request)
    Show(w http.ResponseWriter, r *http.Request)
    // ... other actions
}
```

## Design Decisions

### Why Not Generate Implementation?

We generate structure but not business logic because:

1. Business logic is unique to each application
2. Developers need control over database queries, external services, etc.
3. Generated business logic becomes technical debt
4. Clear boundary between framework and application code

### Why Separate from API Frameworks?

Gluey focuses exclusively on web applications because:

1. Web apps have different concerns than APIs (views, forms, sessions)
2. Simpler mental model without transport abstraction
3. Can integrate with any API framework (Goa, Gin, Echo)
4. Optimized for Rails-like developer experience

### Why DSL Over Configuration?

DSL provides:

1. Type safety at design time
2. Better IDE support with Go tooling
3. Single source of truth
4. Composability through functions
5. Compile-time validation

### Why Generate Templates?

Starting templates provide:

1. Consistent structure across projects
2. Best practices built-in
3. Immediate working application
4. Examples for customization

## Extension Points

### Custom Middleware

Developers can add middleware in the DSL:

```go
Use("CustomAuth", "RateLimiter")
```

### Custom Generators

Future support for plugins:

```go
gluey gen --generator=custom
```

### Integration with ORMs

Generated code is ORM-agnostic:

```go
func (c *PostsController) Index(w http.ResponseWriter, r *http.Request) {
    // Use any ORM or raw SQL
    posts := c.DB.Find(&Post{})
}
```

## Performance Considerations

1. **Zero Runtime Reflection**: All type information at compile time
2. **No Middleware Overhead**: Generated code calls handlers directly
3. **Template Caching**: Templates compiled once at startup
4. **Static Type Checking**: Errors caught at compile time

## Security Considerations

1. **CSRF Protection**: Built into form generation
2. **SQL Injection Prevention**: Developers use parameterized queries
3. **XSS Protection**: Template auto-escaping
4. **Authentication Hooks**: Auth() DSL for access control

## Testing Strategy

### Unit Tests

- Test DSL functions create correct expressions
- Test expression validation catches errors
- Test code generation produces expected output

### Integration Tests

- Generate full applications
- Compile generated code
- Test generated routes work correctly

### Golden Files

- Store expected generated code
- Compare against actual generation
- Detect regressions

## Future Enhancements

### Phase 1 (Current)

- Basic DSL and generation
- RESTful resources
- Form handling
- Templates

### Phase 2

- WebSocket support
- Background jobs
- Email templates
- I18n

### Phase 3

- Visual designer
- Live reload
- Database migrations
- Admin interface

## Comparison with Similar Tools

### vs Ruby on Rails

- **Similar**: Conventions, RESTful resources, MVC pattern
- **Different**: Compiled, type-safe, no runtime magic

### vs Buffalo

- **Similar**: Full-stack Go framework
- **Different**: DSL-driven, generate don't abstract

### vs Goa

- **Similar**: DSL approach, code generation
- **Different**: Web-focused, simpler architecture

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](../LICENSE) for details.
