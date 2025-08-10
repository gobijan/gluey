# Gluey

The web framework that glues Go together. Rails-like productivity with Go's performance and type safety.

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
go install gluey.dev/gluey/cmd/gluey@latest
```

### Create a New Project

```bash
gluey new myapp
cd myapp
```

### Define Your Application

Edit `design/app.go`:

```go
package design

import . "gluey.dev/gluey/dsl"

var _ = WebApp("myapp", func() {
    Description("My awesome web application")
    
    // Define resources (generates all RESTful routes)
    Resource("posts")
    Resource("users")
    
    // Define custom pages
    Page("home", "/")
    Page("about", "/about")
    
    // Define form types with validation
    Type("LoginForm", func() {
        Attribute("email", String, Required(), Format(FormatEmail))
        Attribute("password", String, Required(), MinLength(8))
        Attribute("remember_me", Boolean)
    })
})
```

### Generate Code

```bash
gluey gen
```

This generates:
- Controller interfaces
- Router setup
- Form types with validation
- HTML templates
- Base controller with helpers

### Implement Your Business Logic

```go
// app/controllers/posts.go
type PostsController struct {
    *gen.BaseController
    DB *sql.DB
}

func (c *PostsController) Index(w http.ResponseWriter, r *http.Request) {
    posts := c.fetchPosts()
    c.Render(w, "posts/index", posts)
}

func (c *PostsController) Create(w http.ResponseWriter, r *http.Request) {
    var form gen.NewPostForm
    if err := c.Bind(r, &form); err != nil {
        c.Render(w, "posts/new", map[string]any{"Errors": err})
        return
    }
    
    // Your business logic here
    post := c.createPost(form)
    
    c.Flash(r, "success", "Post created!")
    c.Redirect(w, r, "/posts/" + post.ID)
}
```

### Run Your Application

```bash
go run main.go
```

Visit http://localhost:8080

## Documentation

- [Architecture](docs/architecture.md) - Internal design and structure
- [DSL Specification](docs/webapp-dsl-spec.md) - Complete DSL reference

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

Gluey is in early development. Core features are being implemented:

- ✅ DSL foundation
- ✅ Expression system
- ✅ Basic CLI
- 🚧 Code generation
- 🚧 Runtime support
- 📝 Documentation

## Contributing

Contributions are welcome! Please read the architecture documentation first.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Author

Created by [@gobijan](https://github.com/gobijan)

---

**From the creator:** Coming from Ruby on Rails, I wanted the same productivity in Go without sacrificing type safety or performance. Gluey generates the tedious boilerplate while letting you write the business logic in pure, idiomatic Go
