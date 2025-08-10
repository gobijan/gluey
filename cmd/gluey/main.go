package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "gen", "generate":
		runGenerate()
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("Error: 'new' command requires a project name")
			fmt.Println("Usage: gluey new <project-name>")
			os.Exit(1)
		}
		runNew(os.Args[2])
	case "version", "-v", "--version":
		fmt.Printf("gluey version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Gluey - Rails-like web framework for Go

Usage:
  gluey <command> [arguments]

Commands:
  new <name>    Create a new Gluey project
  gen           Generate code from DSL (alias: generate)
  version       Show version information
  help          Show this help message

Examples:
  gluey new myapp       # Create a new project called 'myapp'
  gluey gen            # Generate code from design/app.go
  gluey version        # Show version

For more information, visit: https://gluey.dev`)
}

func runGenerate() {
	if err := runGenerateImpl(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runNew(projectName string) {
	fmt.Printf("🚀 Creating new Gluey project: %s\n", projectName)

	// Create project directory
	if err := os.Mkdir(projectName, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Create subdirectories
	dirs := []string{
		filepath.Join(projectName, "design"),
		filepath.Join(projectName, "app"),
		filepath.Join(projectName, "app", "controllers"),
		filepath.Join(projectName, "public"),
		filepath.Join(projectName, "public", "css"),
		filepath.Join(projectName, "public", "js"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Create design/app.go
	designContent := fmt.Sprintf(`package design

import . "gluey.dev/gluey/dsl"

var _ = WebApp("%s", func() {
	Description("My %s application")
	
	// Define your resources
	Resource("posts")
	Resource("users")
	
	// Define pages
	Page("home", "/")
	Page("about", "/about")
	
	// Define forms
	Type("LoginForm", func() {
		Attribute("email", String, Required(), Format(FormatEmail))
		Attribute("password", String, Required(), MinLength(8))
		Attribute("remember_me", Boolean)
	})
})
`, projectName, projectName)

	designFile := filepath.Join(projectName, "design", "app.go")
	if err := os.WriteFile(designFile, []byte(designContent), 0644); err != nil {
		fmt.Printf("Error creating design file: %v\n", err)
		os.Exit(1)
	}

	// Create go.mod
	goModContent := fmt.Sprintf(`module %s

go 1.21

require gluey.dev/gluey v%s
`, projectName, version)

	goModFile := filepath.Join(projectName, "go.mod")
	if err := os.WriteFile(goModFile, []byte(goModContent), 0644); err != nil {
		fmt.Printf("Error creating go.mod: %v\n", err)
		os.Exit(1)
	}

	// Create main.go
	mainContent := `package main

import (
	"fmt"
	"log"
	"net/http"
	
	"` + projectName + `/app/controllers"
	"` + projectName + `/gen/webapp/` + projectName + `"
)

func main() {
	// Initialize controllers
	controllers := &controllers.Controllers{
		// Initialize your controllers here
	}
	
	// Setup routes
	mux := http.NewServeMux()
	webapp.MountRoutes(mux, controllers)
	
	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	
	// Start server
	fmt.Println("🚀 Server starting on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
`

	mainFile := filepath.Join(projectName, "main.go")
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		fmt.Printf("Error creating main.go: %v\n", err)
		os.Exit(1)
	}

	// Create README
	readmeContent := fmt.Sprintf(`# %s

A Gluey web application.

## Getting Started

1. Generate code from DSL:
   `+"```bash"+`
   gluey gen
   `+"```"+`

2. Implement your controllers in app/controllers/

3. Run the application:
   `+"```bash"+`
   go run main.go
   `+"```"+`

4. Visit http://localhost:8000

## Project Structure

- design/     - DSL definitions
- gen/        - Generated code (do not edit)
- app/        - Your application code
- public/     - Static assets

## Learn More

Visit https://gluey.dev for documentation.
`, projectName)

	readmeFile := filepath.Join(projectName, "README.md")
	if err := os.WriteFile(readmeFile, []byte(readmeContent), 0644); err != nil {
		fmt.Printf("Error creating README: %v\n", err)
		os.Exit(1)
	}

	// Create .gitignore
	gitignoreContent := `# Generated files
/gen/

# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.work
go.work.sum

# IDE
.idea/
.vscode/
*.swp
*.swo
*~
.DS_Store
`

	gitignoreFile := filepath.Join(projectName, ".gitignore")
	if err := os.WriteFile(gitignoreFile, []byte(gitignoreContent), 0644); err != nil {
		fmt.Printf("Error creating .gitignore: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Project '%s' created successfully!\n\n", projectName)
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  gluey gen        # Generate code from DSL")
	fmt.Println("  go run main.go   # Run your application")
	fmt.Println("\nHappy coding! 🎉")
}
