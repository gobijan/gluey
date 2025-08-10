# Gluey Examples

This directory contains example applications demonstrating Gluey's features and best practices.

## Available Examples

### 1. [Todo](./todo) - Minimal Example
A simple todo list showing the basics:
- Single resource with CRUD operations
- Form validation
- Filtering
- Clean, minimal code

**Perfect for:** Getting started, understanding the basics

### 2. [Blog](./blog) - Full-Featured Example  
A complete blog application demonstrating:
- Multiple resources (posts, comments, users)
- Nested resources
- Authentication requirements
- Search and pagination
- Complex forms with validation
- Static pages

**Perfect for:** Learning advanced features, real-world patterns

## Running Examples

Each example can be run independently:

```bash
cd examples/[example-name]
go mod init [example-name]
go mod edit -replace github.com/gobijan/gluey=../../
gluey gen design
gluey example design
go run main.go
```

## Creating Your Own Example

1. Create a new directory under `examples/`
2. Create `design/app.go` with your DSL
3. Follow the pattern in existing examples
4. Add a README explaining your example

## Learning Path

1. **Start with Todo** - Understand the basic flow
2. **Explore Blog** - See advanced features
3. **Build your own** - Apply what you've learned

## Common Patterns

### Resource with Limited Actions
```go
Resource("items", func() {
    Actions("index", "show", "create")  // Only these actions
})
```

### Nested Resources
```go
Resource("posts", func() {
    Resource("comments", func() {
        BelongsTo("post")
    })
})
```

### Authentication
```go
Resource("admin_panel", func() {
    Auth("admin").All()  // Require admin for all actions
})
```

### Custom Forms
```go
Type("CustomForm", func() {
    Attribute("field", String, Required(), MinLength(5))
    Attribute("email", String, Format(FormatEmail))
})
```

## Contributing Examples

We welcome new examples! Good examples:
- Demonstrate specific use cases
- Include clear documentation
- Show best practices
- Are self-contained

Submit a PR with your example in its own directory.

## Questions?

- Check the [main documentation](../docs/)
- Open an [issue](https://github.com/gobijan/gluey/issues)
- Join the [discussion](https://github.com/gobijan/gluey/discussions)