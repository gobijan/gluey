# Blog Example

A complete blog application demonstrating Gluey's features.

## Features

- 📝 **Blog posts** - Full CRUD operations with drafts and publishing
- 💬 **Comments** - Nested comments on posts
- 👤 **User management** - Registration, login, profiles
- 🔍 **Search & filter** - Find posts by title, content, category
- 📄 **Pagination** - Efficient listing of posts
- 🔒 **Authentication** - Protect actions with auth requirements
- ✅ **Form validation** - Type-safe forms with validation rules

## Quick Start

```bash
# From the examples/blog directory

# 1. Initialize the module (if not already done)
go mod init blog
go mod edit -replace github.com/gobijan/gluey=../..

# 2. Generate interfaces and types
gluey gen design

# 3. Generate example implementations
gluey example design

# 4. Run the application
go run main.go
```

**Note:** The generated code (`gen/` and `app/` folders) is not committed to the repository. You must run the generation commands above to create the application.

Visit http://localhost:8000

## Project Structure

```
blog/
├── design/
│   └── app.go           # DSL design file
├── gen/                 # Generated code (DO NOT EDIT)
│   ├── interfaces/      # Controller interfaces
│   ├── types/          # Form types and validations
│   └── http/           # HTTP transport layer
├── app/                # Your implementations
│   ├── controllers/    # Controller implementations
│   └── views/         # HTML templates
└── main.go            # Application entry point
```

## Design Overview

The blog is defined using Gluey's DSL in `design/app.go`:

### Resources

- **Posts** - Articles with title, content, category, and status
- **Comments** - Nested under posts with author information
- **Users** - User profiles with authentication

### Pages

- **Home** - Landing page at `/`
- **About** - Static about page
- **Contact** - Contact form with GET/POST routes
- **Login/Register** - Authentication pages

### Forms

All forms include validation rules:
- `PostForm` - Create/edit posts with markdown content
- `CommentForm` - Add comments with email validation
- `LoginForm` - User authentication
- `RegisterForm` - New user registration with password requirements
- `ContactForm` - Contact form with message length limits

## Customization

After generating the code, you can customize:

1. **Controllers** (`app/controllers/`) - Add your business logic
2. **Views** (`app/views/`) - Customize HTML templates
3. **Database** - Integrate your preferred database/ORM
4. **Middleware** - Add authentication, logging, etc.

## Example Controller Implementation

```go
// app/controllers/posts.go
type PostsController struct {
    DB *YourDatabase
}

func (c *PostsController) Index(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters for search/filter
    search := r.URL.Query().Get("search")
    category := r.URL.Query().Get("category")
    page := r.URL.Query().Get("page")
    
    // Fetch posts from database
    posts := c.DB.GetPosts(search, category, page)
    
    // Render template
    interfaces.Render(w, "posts/index", map[string]any{
        "Posts": posts,
        "Search": search,
        "Category": category,
    })
}
```

## Authentication Flow

The example includes auth requirements:
- Posts can be viewed by anyone
- Creating/editing posts requires authentication
- Users can only edit their own profiles

Implement your auth middleware and integrate with the generated routes.

## Database Integration

The generated code is database-agnostic. Choose your preferred approach:

```go
// Option 1: Standard library
import "database/sql"
db, _ := sql.Open("postgres", dsn)

// Option 2: GORM
import "gorm.io/gorm"
db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// Option 3: sqlx
import "github.com/jmoiron/sqlx"
db, _ := sqlx.Connect("postgres", dsn)
```

## Testing

```bash
# Run tests
go test ./...

# With coverage
go test -cover ./...
```

## Next Steps

1. Implement database models
2. Add authentication middleware
3. Customize templates with your design
4. Deploy to your preferred hosting platform

## License

MIT