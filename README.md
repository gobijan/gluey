# Gluey (Weekend fun POC - Do not use!)

**The web framework that glues Go together.** Rails-like productivity with Go's performance and type safety.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-00ADD8)](https://go.dev/)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![CI Status](https://github.com/gobijan/gluey/actions/workflows/ci.yml/badge.svg)](https://github.com/gobijan/gluey/actions)

## Features

- 🚀 **Rails-like conventions** - RESTful resources, MVC pattern
- 🔧 **DSL-driven development** - Define your app structure declaratively  
- ⚡ **Code generation** - Generate controllers, routes, forms, and views
- 🛡️ **Type safety** - Compile-time checking throughout
- 📦 **Zero magic** - All generated code is readable Go
- 🎯 **Progressive enhancement** - Start simple, add complexity as needed

## Quick Start

### Installation

```bash
# Install the latest version
go install github.com/gobijan/gluey/cmd/gluey@latest

# Or clone and build from source
git clone https://github.com/gobijan/gluey
cd gluey
make install
```

### Create Your First App

#### 1. Generate a new Gluey project

```bash
gluey new blogapp
cd blogapp
```

This creates a complete project structure with example DSL showing best practices.

#### 2. Or create from scratch

Create `design/app.go`:

```go
package design

import . "github.com/gobijan/gluey/dsl"

var _ = WebApp("blogapp", func() {
    Description("A simple blog application")
    
    // Posts resource with CRUD forms
    Resource("posts", func() {
        // Define form inline with the resource
        Form("PostForm", func() {
            Attribute("title", String, Required(), MaxLength(200))
            Attribute("content", String, Required(), MinLength(10))
            Attribute("published", Boolean)
        })
        
        // Use the same form for create and update
        Create(func() {
            UseForm("PostForm")
        })
        Update(func() {
            UseForm("PostForm")
        })
        
        // Add search and pagination
        Index(func() {
            Params(func() {
                Param("search", String)
                Param("page", Int)
            })
        })
    })
    
    // Users with different forms for signup vs profile
    Resource("users", func() {
        Form("SignupForm", func() {
            Attribute("email", String, Required(), Format(FormatEmail))
            Attribute("password", String, Required(), MinLength(8))
        })
        
        Form("ProfileForm", func() {
            Attribute("name", String)
            Attribute("bio", String, MaxLength(500))
        })
        
        Create(func() {
            UseForm("SignupForm")
        })
        Update(func() {
            UseForm("ProfileForm")
        })
    })
    
    // Static pages
    Page("home", "/")
    Page("about", "/about")
})
```

#### 3. Generate the application structure

```bash
# Generate interfaces, types, and routes from DSL
gluey gen

# Generate example controller implementations and views
gluey example
```

This creates:
- `gen/` - Generated interfaces, types, forms with validation, and HTTP router
- `app/controllers/` - Controller implementations (yours to customize)
- `app/views/` - HTML templates (yours to customize)
- `main.go` - Server entry point

#### 4. Implement your business logic

The generated forms and controllers are ready to use:

```go
package controllers

import (
    "net/http"
    "blogapp/gen/interfaces"
    "blogapp/gen/types"
)

type PostsController struct {
    BaseController
    // Add your DB, services, etc.
}

func (c *PostsController) Create(w http.ResponseWriter, r *http.Request) {
    // Parse form into generated struct
    r.ParseForm()
    var form types.PostForm
    form.Title = r.FormValue("title")
    form.Content = r.FormValue("content")
    form.Published = r.FormValue("published") == "true"
    
    // Use generated validation
    if err := form.Validate(); err != nil {
        c.Render(w, "posts/new", map[string]any{
            "Form":   form,
            "Errors": err,
        })
        return
    }
    
    // Save to database using form fields
    // post := savePost(form.Title, form.Content, form.Published)
    
    c.Flash(w, "success", "Post created!")
    c.Redirect(w, r, "/posts")
}

func (c *PostsController) Index(w http.ResponseWriter, r *http.Request) {
    // Use generated query params struct
    var params types.PostsIndexParams
    params.Search = r.URL.Query().Get("search")
    params.Page = parseIntOrDefault(r.URL.Query().Get("page"), 1)
    
    // Fetch with search/pagination
    posts := fetchPosts(params.Search, params.Page)
    
    c.Render(w, "posts/index", map[string]any{
        "Posts": posts,
        "Params": params,
    })
}
```

#### 5. Run your application

```bash
go run main.go
```

Visit http://localhost:8000 and you'll see your blog!

## Full Example

Here's a complete example showing the modern resource-centric approach:

```go
package design

import . "github.com/gobijan/gluey/dsl"

var _ = WebApp("myapp", func() {
    Description("Complete web application with authentication")
    
    // Posts with full CRUD and search
    Resource("posts", func() {
        // Single form for both create and update
        Form("PostForm", func() {
            Attribute("title", String, Required(), MaxLength(200))
            Attribute("content", String, Required(), MinLength(10))
            Attribute("category", String, Required())
            Attribute("tags", ArrayOf(String))
            Attribute("published", Boolean)
        })
        
        Create(func() {
            UseForm("PostForm")
        })
        Update(func() {
            UseForm("PostForm")
        })
        
        // Query parameters for index
        Index(func() {
            Params(func() {
                Param("search", String)
                Param("category", String)
                Param("page", Int)
                Param("per_page", Int)
            })
        })
    })
    
    // Users with different forms for different actions
    Resource("users", func() {
        // Signup form for registration
        Form("SignupForm", func() {
            Attribute("name", String, Required())
            Attribute("email", String, Required(), Format(FormatEmail))
            Attribute("password", String, Required(), MinLength(8))
            Attribute("password_confirmation", String, Required())
            Attribute("terms_accepted", Boolean, Required())
        })
        
        // Profile form for updates
        Form("ProfileForm", func() {
            Attribute("name", String)
            Attribute("email", String, Format(FormatEmail))
            Attribute("bio", String, MaxLength(500))
            Attribute("website", String, Format(FormatURL))
        })
        
        Create(func() {
            UseForm("SignupForm")
        })
        Update(func() {
            UseForm("ProfileForm")
        })
        
        // Limit available actions
        Actions("index", "show", "new", "create", "edit", "update")
    })
    
    // Session for authentication (singular resource)
    Resource("session", func() {
        Singular() // Makes routes like /session instead of /sessions
        
        Form("LoginForm", func() {
            Attribute("email", String, Required(), Format(FormatEmail))
            Attribute("password", String, Required())
            Attribute("remember_me", Boolean)
        })
        
        // Only login/logout actions
        Actions("new", "create", "destroy")
        
        Create(func() {
            UseForm("LoginForm")
        })
    })
    
    // Password reset as a resource
    Resource("password_resets", func() {
        Form("RequestResetForm", func() {
            Attribute("email", String, Required(), Format(FormatEmail))
        })
        
        Form("ResetPasswordForm", func() {
            Attribute("password", String, Required(), MinLength(8))
            Attribute("password_confirmation", String, Required())
            Attribute("token", String, Required())
        })
        
        Actions("new", "create", "edit", "update")
        
        Create(func() {
            UseForm("RequestResetForm")
        })
        Update(func() {
            UseForm("ResetPasswordForm")
        })
    })
    
    // Search as a resource (for complex search forms)
    Resource("searches", func() {
        Form("AdvancedSearchForm", func() {
            Attribute("query", String, Required())
            Attribute("type", String, Enum("posts", "users", "all"))
            Attribute("date_from", String, Format(FormatDate))
            Attribute("date_to", String, Format(FormatDate))
            Attribute("sort_by", String, Enum("relevance", "date", "popularity"))
        })
        
        Actions("new", "create") // Only search form and results
        
        Create(func() {
            UseForm("AdvancedSearchForm")
        })
    })
    
    // Static pages
    Page("home", "/")
    Page("about", "/about")
    Page("terms", "/terms")
    Page("privacy", "/privacy")
    })
```

## Key Concepts

### Everything is a Resource

In Gluey, forms belong to resources. This creates a consistent, RESTful approach:

- **Authentication**: `session` resource with login form
- **User Registration**: `users` resource with signup form  
- **Search**: `searches` resource with search form
- **Password Reset**: `password_resets` resource with reset forms

### Resource-Level Forms

Forms are defined within resources, keeping related code together:

```go
Resource("posts", func() {
    // Form defined inside the resource
    Form("PostForm", func() {
        Attribute("title", String, Required())
        Attribute("content", String, Required())
    })
    
    // Bind form to actions
    Create(func() {
        UseForm("PostForm")
    })
    Update(func() {
        UseForm("PostForm")  // Reuse same form
    })
})
```

### Different Forms for Different Actions

When create and update need different fields:

```go
Resource("users", func() {
    Form("SignupForm", func() {
        Attribute("email", String, Required())
        Attribute("password", String, Required())
        Attribute("password_confirmation", String, Required())
    })
    
    Form("ProfileForm", func() {
        Attribute("name", String)
        Attribute("bio", String)
        // No password fields for profile updates
    })
    
    Create(func() {
        UseForm("SignupForm")
    })
    Update(func() {
        UseForm("ProfileForm")
    })
})
```

### Singular Resources

For resources that don't have multiple instances (like session):

```go
Resource("session", func() {
    Singular()  // Routes: /session/new, /session, DELETE /session
    Actions("new", "create", "destroy")  // Only login/logout
})
```

### Query Parameters

Define typed query parameters for index/search actions:

```go
Index(func() {
    Params(func() {
        Param("search", String)
        Param("page", Int)
        Param("per_page", Int)
    })
})
```

This generates a typed struct with all parameters.

## Documentation

- [Getting Started Guide](docs/getting-started.md) - Step-by-step tutorial
- [DSL Reference](docs/webapp-dsl-spec.md) - Complete DSL documentation
- [Architecture](docs/architecture.md) - Internal design and structure
- [Examples](examples/) - Sample applications

## For Rails Developers

Gluey brings familiar concepts to Go:

| Rails | Gluey |
|-------|-------|
| `rails new` | `gluey new` |
| `rails generate scaffold` | `Resource()` in DSL |
| `app/controllers/` | Generated interfaces, you implement |
| `app/views/` | Generated templates |
| `config/routes.rb` | DSL generates router |
| ActiveRecord | Your choice (database/sql, GORM, etc.) |

## Project Status

Gluey is in active development. Current status:

- ✅ DSL foundation
- ✅ Expression system with 4-phase evaluation
- ✅ CLI with gen/example commands
- ✅ Code generation for interfaces and implementations
- ✅ Test suite with CI/CD
- 🚧 Advanced features (middleware, sessions, auth)
- 🚧 Database integrations
- 📝 Documentation and examples

## Development

```bash
# Clone the repository
git clone https://github.com/gobijan/gluey
cd gluey

# Run tests
make test

# Run linter
make lint

# Build CLI
make build

# Create example app
make example-app

# See all available commands
make help
```

## Contributing

Contributions are welcome! Please:

1. Read the [architecture documentation](docs/architecture.md)
2. Check existing issues and PRs
3. Write tests for new features
4. Follow Go conventions and keep code simple

## License

MIT License - see [LICENSE](LICENSE) for details.

## Philosophy

Gluey follows these principles:

1. **Convention over configuration** - Sensible defaults that just work
2. **Explicit is better than magic** - All generated code is readable
3. **Type safety first** - Leverage Go's type system fully
4. **Progressive disclosure** - Start simple, add complexity as needed
5. **Own your code** - Generated code is yours to modify

## Comparison

| Feature | Gluey | Gin/Echo | Buffalo | Rails |
|---------|-------|----------|---------|-------|
| DSL-based design | ✅ | ❌ | ❌ | ❌ |
| Code generation | ✅ | ❌ | ✅ | ✅ |
| Type safety | ✅ | ✅ | ✅ | ❌ |
| RESTful conventions | ✅ | Manual | ✅ | ✅ |
| Form validation | ✅ | Manual | ✅ | ✅ |
| Zero magic | ✅ | ✅ | ❌ | ❌ |

## Support

- 🐛 [Report bugs](https://github.com/gobijan/gluey/issues)
- 💡 [Request features](https://github.com/gobijan/gluey/discussions)
- 💬 [Join discussions](https://github.com/gobijan/gluey/discussions)

## Author

Created by [@gobijan](https://github.com/gobijan)

---

**Why Gluey?** Coming from Rails, I wanted the same productivity in Go without sacrificing type safety or performance. Gluey generates the tedious boilerplate while letting you write business logic in pure, idiomatic Go. The name reflects its purpose: gluing together the best of both worlds - Rails' conventions and Go's simplicity.
