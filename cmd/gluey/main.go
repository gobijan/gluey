package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const Version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "gen", "generate":
		runGenerate()
	case "example":
		runExampleCommand()
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("Error: 'new' command requires a project name")
			fmt.Println("Usage: gluey new <project-name>")
			os.Exit(1)
		}
		runNew(os.Args[2])
	case "version", "-v", "--version":
		fmt.Printf("gluey version %s\n", Version)
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
                Options: --local  Use local gluey source (for development)
  gen           Generate interfaces and contracts from DSL (alias: generate)
  example       Generate example implementation (only creates new files)
  version       Show version information
  help          Show this help message

Examples:
  gluey new myapp       # Create a new project called 'myapp'
  gluey new myapp --local  # Create project using local gluey source
  gluey gen            # Generate interfaces from design/app.go
  gluey example        # Generate example controllers and views
  gluey version        # Show version

For more information, visit: https://gluey.dev`)
}

func runGenerate() {
	if err := runGenerateImpl(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runExampleCommand() {
	if err := runExample(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runNew(projectName string) {
	// Check if we're in local development mode
	localMode := false
	glueyPath := ""

	// Check if --local flag is present or if we're in the gluey source directory
	for i, arg := range os.Args {
		if arg == "--local" {
			localMode = true
			// Remove the flag from args
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			break
		}
	}

	// Auto-detect if we're in gluey source directory
	if !localMode {
		if cwd, err := os.Getwd(); err == nil {
			if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
				if data, err := os.ReadFile(filepath.Join(cwd, "go.mod")); err == nil {
					if strings.Contains(string(data), "module github.com/gobijan/gluey") {
						fmt.Println("📍 Detected gluey source directory. Using local mode.")
						localMode = true
						glueyPath = cwd
					}
				}
			}
		}
	}

	// If local mode and no path found, try to find it
	if localMode && glueyPath == "" {
		if cwd, err := os.Getwd(); err == nil {
			glueyPath = cwd
		}
	}

	fmt.Printf("🚀 Creating new Gluey project: %s\n", projectName)
	if localMode {
		fmt.Println("   Using local development mode")
	}

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

import . "github.com/gobijan/gluey/dsl"

var _ = WebApp("%s", func() {
	Description("My %s application")
	
	// Posts resource with CRUD forms
	Resource("posts", func() {
		// Define a form for posts
		Form("PostForm", func() {
			Attribute("title", String, Required(), MaxLength(200))
			Attribute("content", String, Required(), MinLength(10))
			Attribute("published", Boolean)
		})
		
		// Use the same form for both create and update
		Create(func() {
			UseForm("PostForm")
		})
		Update(func() {
			UseForm("PostForm")
		})
		
		// Add search and pagination to index
		Index(func() {
			Params(func() {
				Param("search", String)
				Param("page", Int)
				Param("per_page", Int)
			})
		})
	})
	
	// Users resource with different forms for signup and profile
	Resource("users", func() {
		// Signup form for creating users
		Form("SignupForm", func() {
			Attribute("name", String, Required())
			Attribute("email", String, Required(), Format(FormatEmail))
			Attribute("password", String, Required(), MinLength(8))
			Attribute("password_confirmation", String, Required())
		})
		
		// Profile form for updating users
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
		
		// Only allow certain actions
		Actions("index", "show", "new", "create", "edit", "update")
	})
	
	// Session resource for authentication (singular)
	Resource("session", func() {
		Singular() // Makes routes singular (/session not /sessions)
		
		Form("LoginForm", func() {
			Attribute("email", String, Required(), Format(FormatEmail))
			Attribute("password", String, Required(), MinLength(8))
			Attribute("remember_me", Boolean)
		})
		
		Actions("new", "create", "destroy") // Only login/logout actions
		
		Create(func() {
			UseForm("LoginForm")
		})
	})
	
	// Static pages
	Page("home", "/")
	Page("about", "/about")
})
`, projectName, projectName)

	designFile := filepath.Join(projectName, "design", "app.go")
	if err := os.WriteFile(designFile, []byte(designContent), 0644); err != nil {
		fmt.Printf("Error creating design file: %v\n", err)
		os.Exit(1)
	}

	// Create go.mod
	var goModContent string
	if localMode {
		// Use replace directive for local development
		absProjectPath, _ := filepath.Abs(projectName)
		relativePath, _ := filepath.Rel(absProjectPath, glueyPath)
		goModContent = fmt.Sprintf(`module %s

go 1.21

require github.com/gobijan/gluey v0.0.0

replace github.com/gobijan/gluey => %s
`, projectName, relativePath)
	} else {
		goModContent = fmt.Sprintf(`module %s

go 1.21

require github.com/gobijan/gluey v%s
`, projectName, Version)
	}

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
	genhttp "` + projectName + `/gen/http"
)

func main() {
	// Initialize controllers
	ctrls := genhttp.Controllers{
		// Initialize your controllers here
		// Example:
		// Posts: controllers.NewPosts(),
		// Users: controllers.NewUsers(),
		// Session: controllers.NewSession(),
		// Pages: controllers.NewPagesController(),
	}
	
	// Setup routes
	mux := http.NewServeMux()
	genhttp.MountRoutes(mux, ctrls)
	
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

1. Generate interfaces from DSL:
   `+"```bash"+`
   gluey gen
   `+"```"+`

2. Generate example implementation:
   `+"```bash"+`
   gluey example
   `+"```"+`

3. Customize controllers in app/controllers/

4. Run the application:
   `+"```bash"+`
   go run main.go
   `+"```"+`

5. Visit http://localhost:8000

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
	if localMode {
		fmt.Println("  go mod tidy      # Download dependencies")
	}
	fmt.Println("  gluey gen        # Generate interfaces from DSL")
	fmt.Println("  gluey example    # Generate example implementations")
	fmt.Println("  go run main.go   # Run your application")
	if localMode {
		fmt.Println("\n📝 Note: Using local gluey source for development")
	}
	fmt.Println("\nHappy coding! 🎉")
}
