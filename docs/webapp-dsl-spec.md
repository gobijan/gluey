# Goa WebApp DSL Specification

## Vision
Extend Goa to become a full-stack web framework that brings Rails-like productivity to Go. For developers coming from Ruby on Rails who want type safety and compiled performance without manually creating boilerplate files, Goa WebApp generates the entire application structure from a declarative DSL.

## Target Audience
- Rails developers seeking Go's performance and type safety
- Go developers wanting Rails-like conventions and productivity
- Teams building web applications alongside APIs
- Developers who prefer design-first development

## Core Principles

1. **Structure, Not Logic**: DSL defines what exists (routes, forms, resources), not how it behaves
2. **Convention Over Configuration**: Rails-like conventions with Go idioms
3. **Generate Everything Once**: No manual file creation - DSL generates entire structure
4. **Type Safety Throughout**: Compile-time checking for forms, routes, and data
5. **Zero Magic**: All generated code is readable, debuggable Go
6. **Progressive Enhancement**: Start minimal, add complexity only when needed

## DSL Structure

### Top-Level Definition

```go
// Separate from API - no pollution between web and API concerns
var _ = WebApp("portal", func() {
    Description("Customer portal application")
    
    // Global configuration (optional)
    Sessions(func() {
        Store("cookie")  // or "redis", "memory"
    })
    
    Assets(func() {
        Path("/static")
    })
    
    Layouts(func() {
        Default("application")  // Sets default layout
        Layout("admin")         // Additional layout
        Layout("marketing")     // Another layout option
    })
    
    // Global middleware stack
    Use("RequestID", "Logger", "Recover", "CSRF")
})
```

### Resource-Based Forms Philosophy

**Everything is a resource.** Forms are always defined within the context of the resource that handles them. This creates a consistent, RESTful approach where:

- Login forms belong to a `sessions` resource
- Search forms belong to a `searches` resource  
- Contact forms belong to a `contacts` resource
- User registration belongs to the `users` resource

No separate form types exist at the application level - all forms are resource forms.

### Resources (RESTful Controllers)

```go
var _ = WebApp("portal", func() {
    // Standard CRUD resource with forms
    Resource("posts", func() {
        // Define forms within the resource
        Form("PostForm", func() {
            Attribute("title", String, Required(), MaxLength(200))
            Attribute("content", String, Required())
            Attribute("published", Boolean)
        })
        
        // Use the same form for create and update
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
                Param("page", Int, Default(1))
                Param("per_page", Int, Default(20))
            })
        })
    })
    
    // Resource with different forms for different actions
    Resource("users", func() {
        // Signup form for creating users
        Form("SignupForm", func() {
            Attribute("name", String, Required())
            Attribute("email", String, Required(), Format(FormatEmail))
            Attribute("password", String, Required(), MinLength(8))
            Attribute("password_confirmation", String, Required())
        })
        
        // Profile form for editing users
        Form("ProfileForm", func() {
            Attribute("name", String)
            Attribute("email", String, Format(FormatEmail))
            Attribute("bio", String, MaxLength(500))
        })
        
        Create(func() {
            UseForm("SignupForm")
        })
        Update(func() {
            UseForm("ProfileForm")
        })
        
        // Authentication requirements
        Auth("authenticated").Except("new", "create")
        Auth("owner").Only("edit", "update", "destroy")
    })
    
    // Singular resource for session (login/logout)
    Resource("session", func() {
        Singular()  // Makes routes singular (/session not /sessions)
        
        Form("LoginForm", func() {
            Attribute("email", String, Required(), Format(FormatEmail))
            Attribute("password", String, Required())
            Attribute("remember_me", Boolean)
        })
        
        Actions("new", "create", "destroy")  // Only login/logout actions
        
        Create(func() {
            UseForm("LoginForm")
        })
    })
    
    // Search as a resource
    Resource("searches", func() {
        Form("SearchForm", func() {
            Attribute("q", String, Required())
            Attribute("category", String)
            Attribute("from_date", String, Format(FormatDate))
            Attribute("to_date", String, Format(FormatDate))
        })
        
        Actions("new", "create")  // Search form and results
        
        Create(func() {
            UseForm("SearchForm")
        })
    })
    
    // Nested resources
    Resource("users", func() {
        Resource("posts", func() {
            // Generates /users/{user_id}/posts/... routes
            BelongsTo("user")
        })
    })
})
```

### Static Pages (Non-Resource)

```go
var _ = WebApp("portal", func() {
    // Only truly static content uses Page
    Page("home", "/")
    Page("about", "/about")
    Page("terms", "/terms")
    Page("privacy", "/privacy")
    
    // Pages can have layouts and auth
    Page("dashboard", func() {
        Route("GET", "/dashboard")
        Layout("admin")
        Auth("authenticated")
    })
    
    // Note: Contact forms should be a resource instead:
    // Resource("contacts", func() { ... })
})
```

## Conventions

### Form Naming Conventions

Forms are defined within resources and can be named anything, but these conventions help:

| Resource | Conventional Form Names | Usage |
|----------|------------------------|-------|
| `posts` | `PostForm` | Single form for both create and update |
| `posts` | `CreatePostForm`, `UpdatePostForm` | Different forms for create vs update |
| `users` | `SignupForm`, `ProfileForm` | Semantic names for user actions |
| `session` | `LoginForm` | Form for session creation |
| `searches` | `SearchForm` | Form for search parameters |

If no forms are defined, the generator creates placeholder forms following the pattern:
- `New{Resource}Form` for create actions
- `Edit{Resource}Form` for update actions

### RESTful Routes

For standard resources like `Resource("posts")`:

| HTTP Method | Path | Controller Method | View Template |
|-------------|------|------------------|---------------|
| GET | /posts | Index() | posts/index.html |
| GET | /posts/new | New() | posts/new.html |
| POST | /posts | Create() | - (typically redirects) |
| GET | /posts/{id} | Show() | posts/show.html |
| GET | /posts/{id}/edit | Edit() | posts/edit.html |
| PUT/PATCH | /posts/{id} | Update() | - (typically redirects) |
| DELETE | /posts/{id} | Destroy() | - (typically redirects) |

For singular resources with `Singular()` like `Resource("session")`:

| HTTP Method | Path | Controller Method | Purpose |
|-------------|------|------------------|---------|
| GET | /session/new | New() | Login page |
| POST | /session | Create() | Process login |
| DELETE | /session | Destroy() | Logout |

### Generated Directory Structure

Everything is generated - no manual file creation:

```
gen/webapp/{app_name}/
├── types.go              # All form structs with validation
├── router.go             # Route mounting and middleware
├── controllers/
│   ├── base.go          # BaseController with all helpers
│   ├── posts.go         # PostsController interface
│   ├── users.go         # UsersController interface
│   └── pages.go         # PagesController interface
├── views/
│   ├── layouts/
│   │   ├── application.html
│   │   ├── admin.html
│   │   └── marketing.html
│   ├── posts/
│   │   ├── index.html
│   │   ├── show.html
│   │   ├── new.html
│   │   └── edit.html
│   ├── users/
│   │   └── ...
│   └── shared/
│       ├── _errors.html
│       ├── _flash.html
│       └── _pagination.html
└── assets/
    ├── css/
    │   └── application.css
    └── js/
        └── application.js
```

## What Gets Generated

### 1. Form Types with Validation

```go
// gen/webapp/portal/types.go

// From Resource("posts") with Form("PostForm")
type PostForm struct {
    Title     string `form:"title" json:"title" validate:"required,max=200"`
    Content   string `form:"content" json:"content" validate:"required"`
    Published bool   `form:"published" json:"published"`
}

func (f *PostForm) Validate() error {
    v := runtime.NewValidator()
    v.Required("title", f.Title)
    v.MaxLength("title", f.Title, 200)
    v.Required("content", f.Content)
    
    if !v.Valid() {
        return v.Errors()
    }
    return nil
}

// From Resource("users") with different forms
type SignupForm struct {
    Name                 string `form:"name" json:"name" validate:"required"`
    Email                string `form:"email" json:"email" validate:"required,email"`
    Password             string `form:"password" json:"password" validate:"required,min=8"`
    PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" validate:"required"`
}

type ProfileForm struct {
    Name  string `form:"name" json:"name,omitempty"`
    Email string `form:"email" json:"email,omitempty" validate:"omitempty,email"`
    Bio   string `form:"bio" json:"bio,omitempty" validate:"omitempty,max=500"`
}

// From Resource("session") 
type LoginForm struct {
    Email      string `form:"email" json:"email" validate:"required,email"`
    Password   string `form:"password" json:"password" validate:"required"`
    RememberMe bool   `form:"remember_me" json:"remember_me"`
}

// Query parameters from Index()
type PostsIndexParams struct {
    Search  string `form:"search" json:"search,omitempty"`
    Page    int    `form:"page" json:"page,omitempty"`
    PerPage int    `form:"per_page" json:"per_page,omitempty"`
}
```

### 2. Controller Interfaces

```go
// gen/webapp/portal/controllers/posts.go

type PostsController interface {
    Index(w http.ResponseWriter, r *http.Request)
    Show(w http.ResponseWriter, r *http.Request)
    New(w http.ResponseWriter, r *http.Request)
    Create(w http.ResponseWriter, r *http.Request)
    Edit(w http.ResponseWriter, r *http.Request)
    Update(w http.ResponseWriter, r *http.Request)
    Destroy(w http.ResponseWriter, r *http.Request)
}
```

### 3. Base Controller with Helpers

```go
// gen/webapp/portal/controllers/base.go

type BaseController struct {
    templates *template.Template
    // ... other common fields
}

// Rails-like helpers
func (c *BaseController) Render(w http.ResponseWriter, template string, data any) error
func (c *BaseController) Redirect(w http.ResponseWriter, r *http.Request, path string)
func (c *BaseController) RedirectBack(w http.ResponseWriter, r *http.Request)
func (c *BaseController) Flash(r *http.Request, level, message string)
func (c *BaseController) CurrentUser(r *http.Request) *User
func (c *BaseController) Params(r *http.Request) map[string]string
func (c *BaseController) Param(r *http.Request, name string) string

// Form helpers
func (c *BaseController) Bind(r *http.Request, v any) error
func (c *BaseController) FormFor(model any) *FormBuilder

// Pagination helpers (if Paginate() was in DSL)
func (c *BaseController) Paginate(items any, perPage int) *Pagination

// Search/Filter helpers (if specified in DSL)
func (c *BaseController) Search(r *http.Request, fields ...string) string
func (c *BaseController) Filter(r *http.Request, field string) string
```

### 4. Router with Middleware

```go
// gen/webapp/portal/router.go

func MountRoutes(mux *http.ServeMux, controllers Controllers) {
    // Global middleware
    stack := middleware.Chain(
        middleware.RequestID,
        middleware.Logger,
        middleware.Recover,
        middleware.CSRF,
    )
    
    // Posts resource
    posts := controllers.Posts
    mux.Handle("GET /posts", stack.Then(posts.Index))
    mux.Handle("GET /posts/new", stack.Then(posts.New))
    mux.Handle("POST /posts", stack.Then(posts.Create))
    mux.Handle("GET /posts/{id}", stack.Then(posts.Show))
    mux.Handle("GET /posts/{id}/edit", stack.Then(posts.Edit))
    mux.Handle("PUT /posts/{id}", stack.Then(posts.Update))
    mux.Handle("PATCH /posts/{id}", stack.Then(posts.Update))
    mux.Handle("DELETE /posts/{id}", stack.Then(posts.Destroy))
    
    // Nested resources
    mux.Handle("GET /users/{user_id}/posts", stack.Then(userPosts.Index))
    // ... etc
}
```

### 5. Starter Templates (Rails-like)

```html
<!-- gen/webapp/portal/views/layouts/application.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}} - Portal</title>
    <link rel="stylesheet" href="/static/css/application.css">
    {{block "head" .}}{{end}}
</head>
<body>
    {{template "shared/_flash" .}}
    
    <main>
        {{block "content" .}}{{end}}
    </main>
    
    <script src="/static/js/application.js"></script>
    {{block "scripts" .}}{{end}}
</body>
</html>
```

```html
<!-- gen/webapp/portal/views/posts/index.html -->
{{template "layout" .}}

{{define "content"}}
<div class="posts-index">
    <h1>Posts</h1>
    
    <a href="/posts/new" class="btn btn-primary">New Post</a>
    
    {{if .Search}}
    <form method="get" class="search-form">
        <input type="text" name="q" value="{{.SearchQuery}}" placeholder="Search...">
        <button type="submit">Search</button>
    </form>
    {{end}}
    
    <div class="posts">
        {{range .Posts}}
        <article class="post">
            <h2><a href="/posts/{{.ID}}">{{.Title}}</a></h2>
            <div class="post-meta">
                <span class="category">{{.Category}}</span>
                <time>{{.CreatedAt}}</time>
            </div>
            <div class="post-actions">
                <a href="/posts/{{.ID}}/edit">Edit</a>
                <form method="post" action="/posts/{{.ID}}" style="display:inline">
                    <input type="hidden" name="_method" value="DELETE">
                    <button type="submit" onclick="return confirm('Are you sure?')">Delete</button>
                </form>
            </div>
        </article>
        {{end}}
    </div>
    
    {{if .Pagination}}
        {{template "shared/_pagination" .Pagination}}
    {{end}}
</div>
{{end}}
```

```html
<!-- gen/webapp/portal/views/posts/new.html -->
{{template "layout" .}}

{{define "content"}}
<div class="posts-new">
    <h1>New Post</h1>
    
    <form method="post" action="/posts">
        {{template "shared/_errors" .Errors}}
        
        <div class="form-group">
            <label for="title">Title</label>
            <input type="text" id="title" name="title" value="{{.Form.Title}}" required>
        </div>
        
        <div class="form-group">
            <label for="content">Content</label>
            <textarea id="content" name="content" required>{{.Form.Content}}</textarea>
        </div>
        
        <div class="form-group">
            <label for="category_id">Category</label>
            <select id="category_id" name="category_id" required>
                <option value="">Select a category</option>
                {{range .Categories}}
                <option value="{{.ID}}" {{if eq .ID $.Form.CategoryID}}selected{{end}}>
                    {{.Name}}
                </option>
                {{end}}
            </select>
        </div>
        
        <button type="submit" class="btn btn-primary">Create Post</button>
        <a href="/posts" class="btn btn-secondary">Cancel</a>
    </form>
</div>
{{end}}
```

## Developer Implementation

The developer only writes business logic - all structure is generated:

```go
// app/controllers/posts.go - Only file developer needs to create

package controllers

import (
    gen "myapp/gen/webapp/portal/controllers"
)

type PostsController struct {
    *gen.BaseController  // Inherit all helpers
    DB *sql.DB
}

// Implement the interface - structure already defined
func (c *PostsController) Index(w http.ResponseWriter, r *http.Request) {
    // Search if provided (helper available because DSL had Searchable)
    query := c.Search(r, "title", "content")
    
    // Developer writes their query
    posts := c.queryPosts(query)
    
    // Pagination (helper available because DSL had Paginate)
    paged := c.Paginate(posts, 20)
    
    // Render using generated template
    c.Render(w, "posts/index", map[string]any{
        "Posts":      paged.Items,
        "Pagination": paged,
        "SearchQuery": query,
    })
}

func (c *PostsController) Create(w http.ResponseWriter, r *http.Request) {
    // Parse form data into the generated struct
    r.ParseForm()
    var form gen.PostForm
    form.Title = r.FormValue("title")
    form.Content = r.FormValue("content")
    form.Published = r.FormValue("published") == "true"
    
    // Validate using generated validation
    if err := form.Validate(); err != nil {
        // Re-render with errors
        c.Render(w, "posts/new", map[string]any{
            "Form":   form,
            "Errors": err,
        })
        return
    }
    
    // Business logic - developer's responsibility
    post, err := c.createPost(form)
    if err != nil {
        c.Flash(r, "error", "Could not create post")
        c.Render(w, "posts/new", map[string]any{"Form": form})
        return
    }
    
    // Success - developer controls the flow
    c.Flash(r, "success", "Post created successfully!")
    c.Redirect(w, r, fmt.Sprintf("/posts/%d", post.ID))
}

func (c *PostsController) Update(w http.ResponseWriter, r *http.Request) {
    id := c.Param(r, "id")
    
    // Parse and validate form
    r.ParseForm()
    var form gen.PostForm
    form.Title = r.FormValue("title")
    form.Content = r.FormValue("content")
    form.Published = r.FormValue("published") == "true"
    
    if err := form.Validate(); err != nil {
        c.Render(w, "posts/edit", map[string]any{
            "Form": form,
            "Errors": err,
            "PostID": id,
        })
        return
    }
    
    // Update the post
    if err := c.updatePost(id, &form); err != nil {
        c.Flash(r, "error", "Could not update post")
        c.Render(w, "posts/edit", map[string]any{"Form": form})
        return
    }
    
    c.Flash(r, "info", "Post updated")
    c.Redirect(w, r, fmt.Sprintf("/posts/%s", id))
}

// Private methods - developer's business logic
func (c *PostsController) queryPosts(search string) []*Post {
    // Implementation
}

func (c *PostsController) createPost(form gen.PostForm) (*Post, error) {
    // Implementation - use form.Title, form.Content, etc.
}

func (c *PostsController) updatePost(id string, form *gen.PostForm) error {
    // Implementation
}
```

## Benefits for Rails Developers

### What You Get From Rails
- ✅ Convention over configuration
- ✅ RESTful resources
- ✅ Form objects with validation
- ✅ Layouts and templates
- ✅ Flash messages
- ✅ Strong params (via form binding)
- ✅ Before filters (via middleware)
- ✅ Nested resources
- ✅ Partial templates
- ✅ Asset pipeline structure

### What You Gain With Go + Goa
- ✅ **Type Safety**: Compile-time checking everywhere
- ✅ **Performance**: 10-50x faster than Rails
- ✅ **Single Binary**: Deploy one file
- ✅ **No Magic**: All code is readable and debuggable
- ✅ **Explicit Flow**: You control the business logic
- ✅ **Design-First**: Structure defined upfront

### What's Different (By Design)
- **No ORM**: You write SQL or use your preferred Go database library
- **No Migrations in DSL**: Use golang-migrate or similar
- **No Asset Compilation**: Use esbuild, webpack, or Go embed
- **Explicit Business Logic**: You write the controller implementation

## Migration Path for Rails Developers

### Rails Concept → Goa WebApp Equivalent

| Rails | Goa WebApp |
|-------|------------|
| `rails new myapp` | Define WebApp DSL, run `goa gen` |
| `rails g scaffold Post` | `Resource("posts")` in DSL |
| `app/controllers/` | Generated interfaces, you implement |
| `app/models/` | Your choice (structs + database/sql, GORM, etc.) |
| `app/views/` | Generated templates (customizable) |
| `config/routes.rb` | DSL generates router.go |
| `db/migrations/` | Use golang-migrate (separate) |
| ActiveRecord | Your choice of DB library |
| `before_action` | Middleware in DSL |
| `respond_to` | Check Accept header in controller |
| `form_for` | Generated form types with validation |
| `link_to`, `form_tag` | Template helpers (or use templ) |

## Getting Started

```go
// design/webapp.go
package design

import . "goa.design/goa/v3/dsl"

var _ = WebApp("myapp", func() {
    Description("My Rails-like Go application")
    
    // Everything is a resource
    Resource("posts", func() {
        Form("PostForm", func() {
            Attribute("title", String, Required())
            Attribute("content", String, Required())
        })
        Create(func() {
            UseForm("PostForm")
        })
        Update(func() {
            UseForm("PostForm")
        })
    })
    
    Resource("users", func() {
        Form("SignupForm", func() {
            Attribute("email", String, Required(), Format(FormatEmail))
            Attribute("password", String, Required(), MinLength(8))
        })
        Actions("new", "create")  // Only signup actions
        Create(func() {
            UseForm("SignupForm")
        })
    })
    
    Resource("session", func() {
        Singular()  // Login/logout
        Form("LoginForm", func() {
            Attribute("email", String, Required(), Format(FormatEmail))
            Attribute("password", String, Required())
        })
        Actions("new", "create", "destroy")
        Create(func() {
            UseForm("LoginForm")
        })
    })
    
    // Only truly static pages use Page()
    Page("home", "/")
    Page("about", "/about")
})
```

Run generation:
```bash
goa gen myapp/design
```

This generates your entire application structure. Then implement your controllers:

```go
// app/controllers/posts.go
type PostsController struct {
    *gen.BaseController
    // Add your dependencies
}

// Implement the generated interface...
```

## Summary

Goa WebApp brings Rails' productivity to Go by:
1. **Generating all boilerplate** from a simple DSL
2. **Providing Rails-like conventions** with Go idioms
3. **Maintaining type safety** throughout
4. **Keeping business logic separate** from structure
5. **Making the switch from Rails natural** for developers

The framework handles the tedious parts (routing, form binding, validation, templates) while you focus on your business logic with the full power and performance of Go.