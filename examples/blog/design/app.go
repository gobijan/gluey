package design

import . "github.com/gobijan/gluey/dsl"

var _ = WebApp("blog", func() {
	Description("A simple blog application demonstrating Gluey features")

	// Blog posts resource with full CRUD operations
	Resource("posts", func() {
		// Configure index action
		Index(func() {
			Paginate(10)
			Searchable("title", "content", "author")
			Filterable("status", "category")
		})

		// Only authenticated users can create/edit/delete
		Auth("authenticated").Except("index", "show")
	})

	// Comments nested under posts
	Resource("posts", func() {
		Resource("comments", func() {
			BelongsTo("post")
			Actions("index", "create", "destroy")
			Auth("authenticated").Only("create", "destroy")
		})
	})

	// User management
	Resource("users", func() {
		Actions("show", "new", "create", "edit", "update")
		Auth("authenticated").Only("edit", "update")
		Auth("self").Only("edit", "update") // Can only edit own profile
	})

	// Static pages
	Page("home", "/")
	Page("about", "/about")
	Page("contact", func() {
		Route("GET", "/contact")
		Route("POST", "/contact")
	})

	// Authentication pages
	Page("login", func() {
		Route("GET", "/login")
		Route("POST", "/login")
	})
	Page("logout", "/logout")
	Page("register", func() {
		Route("GET", "/register")
		Route("POST", "/register")
	})

	// Form definitions with validation
	Type("PostForm", func() {
		Attribute("title", String, func() {
			Required()
			MinLength(3)
			MaxLength(200)
			Description("Post title")
		})
		Attribute("content", String, func() {
			Required()
			MinLength(10)
			Description("Post content in markdown")
		})
		Attribute("category", String, func() {
			Enum("tech", "lifestyle", "travel", "food", "other")
			Default("other")
		})
		Attribute("status", String, func() {
			Enum("draft", "published", "archived")
			Default("draft")
		})
		Attribute("tags", ArrayOf(String), func() {
			Description("Tags for categorization")
		})
		Attribute("published_at", String, func() {
			Format(FormatDateTime)
			Description("Publication date and time")
		})
	})

	Type("CommentForm", func() {
		Attribute("author", String, func() {
			Required()
			MinLength(2)
			MaxLength(100)
		})
		Attribute("email", String, func() {
			Required()
			Format(FormatEmail)
		})
		Attribute("content", String, func() {
			Required()
			MinLength(5)
			MaxLength(1000)
		})
	})

	Type("LoginForm", func() {
		Attribute("email", String, func() {
			Required()
			Format(FormatEmail)
		})
		Attribute("password", String, func() {
			Required()
			MinLength(8)
		})
		Attribute("remember_me", Boolean, func() {
			Default(false)
		})
	})

	Type("RegisterForm", func() {
		Attribute("name", String, func() {
			Required()
			MinLength(2)
			MaxLength(100)
		})
		Attribute("email", String, func() {
			Required()
			Format(FormatEmail)
		})
		Attribute("password", String, func() {
			Required()
			MinLength(8)
			Pattern("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).*$")
			Description("Must contain uppercase, lowercase, and number")
		})
		Attribute("password_confirmation", String, func() {
			Required()
			Description("Must match password")
		})
		Attribute("terms_accepted", Boolean, func() {
			Required()
			Validation(func() {
				// Custom validation - must be true
			})
		})
	})

	Type("ContactForm", func() {
		Attribute("name", String, Required())
		Attribute("email", String, func() {
			Required()
			Format(FormatEmail)
		})
		Attribute("subject", String, func() {
			Required()
			MinLength(5)
			MaxLength(100)
		})
		Attribute("message", String, func() {
			Required()
			MinLength(20)
			MaxLength(5000)
		})
	})
})
