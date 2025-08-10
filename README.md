# Gluey

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

#### 1. Create a new directory and initialize Go module

```bash
mkdir blogapp && cd blogapp
go mod init blogapp
```

#### 2. Create your design file

Create `design/app.go`:

```go
package design

import . "github.com/gobijan/gluey/dsl"

var _ = WebApp("blogapp", func() {
    Description("A simple blog application")
    
    // Define a posts resource with all RESTful actions
    Resource("posts", func() {
        // Generates: index, show, new, create, edit, update, destroy
    })
    
    // Define a users resource with limited actions
    Resource("users", func() {
        Actions("index", "show", "new", "create")
    })
    
    // Define custom pages
    Page("home", "/")
    Page("about", "/about")
    
    // Define form types with validation
    Type("PostForm", func() {
        Attribute("title", String, Required(), MinLength(3), MaxLength(100))
        Attribute("content", String, Required(), MinLength(10))
        Attribute("published", Boolean)
    })
})
```

#### 3. Generate the application structure

```bash
# Generate interfaces and types
gluey gen design

# Generate example implementations
gluey example design
```

This creates:
- `gen/` - Generated interfaces, types, and HTTP setup
- `app/controllers/` - Controller implementations (yours to customize)
- `app/views/` - HTML templates (yours to customize)
- `main.go` - Server entry point

#### 4. Implement your business logic

The generated controllers in `app/controllers/posts.go` are ready to customize:

```go
```go
package controllers

import (
    "net/http"
    "blogapp/gen/interfaces"
)

type PostsController struct {
    // Add your dependencies here (DB, services, etc.)
}

func (c *PostsController) Index(w http.ResponseWriter, r *http.Request) {
    // Fetch posts from your database
    posts := []map[string]any{
        {"ID": "1", "Title": "Hello World", "Content": "My first post!"},
    }
    
    // Render the index template
    interfaces.Render(w, "posts/index", map[string]any{
        "Posts": posts,
    })
}

func (c *PostsController) Create(w http.ResponseWriter, r *http.Request) {
    var form gen.PostForm
    if err := interfaces.Bind(r, &form); err != nil {
        interfaces.Render(w, "posts/new", map[string]any{
            "Errors": err,
            "Form":   form,
        })
        return
    }
    
    // Save to database
    // post := savePost(form)
    
    http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
```
```

#### 5. Run your application

```bash
go run main.go
```

Visit http://localhost:8000 and you'll see your blog!

## Full Example

Here's a more complete example showing advanced features:

```go
package design

import . "github.com/gobijan/gluey/dsl"

var _ = WebApp("myapp", func() {
    Description("E-commerce platform")
    
    // Products resource with search and pagination
    Resource("products", func() {
        Index(func() {
            Paginate(20)
            Searchable("name", "description")
            Filterable("category", "price_range")
        })
        
        // Custom forms for different actions
        Update(func() {
            Form("EditProductForm")
        })
    })
    
    // Nested resources
    Resource("users", func() {
        Resource("orders", func() {
            BelongsTo("user")
            Actions("index", "show", "create")
        })
        
        // Authentication requirements
        Auth("authenticated").Except("new", "create")
        Auth("admin").Only("destroy")
    })
    
    // Custom pages with multiple routes
    Page("checkout", func() {
        Route("GET", "/checkout")
        Route("POST", "/checkout")
        Form("CheckoutForm")
    })
    
    // Complex form with nested attributes
    Type("CheckoutForm", func() {
        Attribute("shipping_address", func() {
            Attribute("street", String, Required())
            Attribute("city", String, Required())
            Attribute("zip", String, Pattern("^\\d{5}$"))
        })
        Attribute("payment_method", String, Enum("credit_card", "paypal"))
        Attribute("items", ArrayOf(func() {
            Attribute("product_id", String, Required())
            Attribute("quantity", Int, Min(1))
        }))
    })
})
```

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
