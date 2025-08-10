# Todo Example

A minimal todo list application showcasing Gluey's simplicity.

## Features

- ✅ Create, read, update, delete todos
- 🏷️ Priority levels (low, medium, high)
- 📅 Due dates
- 🔍 Filter by status and priority
- 📝 Optional descriptions

## Quick Start

```bash
# From the examples/todo directory
go mod init todo
go mod edit -replace gluey.dev/gluey=../../

# Generate the application
gluey gen design
gluey example design

# Run
go run main.go
```

Visit http://localhost:8000

## The Entire Design

The complete application is defined in just 35 lines of DSL:

```go
WebApp("todo", func() {
    Description("A simple todo list application")
    
    Resource("todos", func() {
        Index(func() {
            Filterable("status", "priority")
        })
    })
    
    Type("TodoForm", func() {
        Attribute("title", String, Required(), MinLength(1))
        Attribute("priority", String, Enum("low", "medium", "high"))
        Attribute("due_date", String, Format(FormatDate))
        Attribute("completed", Boolean, Default(false))
    })
})
```

## What Gets Generated

From this simple DSL, Gluey generates:

### Routes
- `GET /todos` - List all todos
- `GET /todos/new` - New todo form
- `POST /todos` - Create todo
- `GET /todos/{id}` - View todo
- `GET /todos/{id}/edit` - Edit form
- `PUT /todos/{id}` - Update todo
- `DELETE /todos/{id}` - Delete todo

### Code Structure
```
gen/
├── interfaces/
│   └── todos.go         # TodosController interface
├── types/
│   └── forms.go         # TodoForm with validation
└── http/
    └── server.go        # HTTP routing and setup

app/
├── controllers/
│   └── todos.go         # Your implementation
└── views/
    └── todos/
        ├── index.html   # List view
        ├── show.html    # Detail view
        ├── new.html     # Create form
        └── edit.html    # Edit form
```

## Simple Implementation

```go
// app/controllers/todos.go
type TodosController struct {
    todos []Todo  // In-memory storage for demo
}

func (c *TodosController) Index(w http.ResponseWriter, r *http.Request) {
    interfaces.Render(w, "todos/index", map[string]any{
        "Todos": c.todos,
    })
}

func (c *TodosController) Create(w http.ResponseWriter, r *http.Request) {
    var form gen.TodoForm
    if err := interfaces.Bind(r, &form); err != nil {
        interfaces.Render(w, "todos/new", map[string]any{"Errors": err})
        return
    }
    
    // Add to todos
    c.todos = append(c.todos, Todo{
        ID:       len(c.todos) + 1,
        Title:    form.Title,
        Priority: form.Priority,
        DueDate:  form.DueDate,
    })
    
    http.Redirect(w, r, "/todos", http.StatusSeeOther)
}
```

## Add Persistence

Replace in-memory storage with a database:

```go
// Using database/sql
type TodosController struct {
    db *sql.DB
}

func (c *TodosController) Index(w http.ResponseWriter, r *http.Request) {
    rows, _ := c.db.Query("SELECT * FROM todos ORDER BY created_at DESC")
    // ... fetch and render
}
```

## Deployment

The generated application is a standard Go binary:

```bash
# Build
go build -o todo

# Run in production
./todo

# Or use Docker
docker build -t todo .
docker run -p 8000:8000 todo
```

## License

MIT