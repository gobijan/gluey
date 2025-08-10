package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
)

// runGenerateImpl executes the interface generation.
func runGenerateImpl() error {
	fmt.Println("🔨 Generating interfaces and contracts from DSL...")

	// Check if design/app.go exists
	designFile := filepath.Join("design", "app.go")
	if _, err := os.Stat(designFile); os.IsNotExist(err) {
		return fmt.Errorf("%s not found - make sure you're in a Gluey project directory", designFile)
	}

	// Get current directory name for module
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	moduleName := filepath.Base(cwd)

	// Create a temporary directory for the generator
	tmpDir, err := ioutil.TempDir("", "gluey-gen-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary main.go that imports and executes the design
	mainContent := fmt.Sprintf(`package main

import (
	"fmt"
	"log"
	
	_ "%s/design"
	"gluey.dev/gluey/codegen"
	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
)

func main() {
	// Execute the DSL
	if err := eval.RunDSL(); err != nil {
		log.Fatal("DSL execution failed:", err)
	}
	
	// Check we have an app
	if expr.Root == nil {
		log.Fatal("No WebApp found in design")
	}
	
	// Generate interfaces only
	gen := codegen.NewInterfaceGenerator(expr.Root, "gen")
	if err := gen.Generate(); err != nil {
		log.Fatal("Interface generation failed:", err)
	}
	
	fmt.Println("✅ Interface generation complete!")
	fmt.Println("\nGenerated files in gen/")
}
`, moduleName)

	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write generator main.go: %w", err)
	}

	// Copy go.mod to temp directory
	goModContent, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}
	
	tmpGoMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(tmpGoMod, goModContent, 0644); err != nil {
		return fmt.Errorf("failed to write temp go.mod: %w", err)
	}

	// Copy go.sum if it exists
	if _, err := os.Stat("go.sum"); err == nil {
		goSumContent, _ := os.ReadFile("go.sum")
		tmpGoSum := filepath.Join(tmpDir, "go.sum")
		os.WriteFile(tmpGoSum, goSumContent, 0644)
	}

	// Run the generator
	cmd := exec.Command("go", "run", mainFile)
	cmd.Dir = cwd // Run in the project directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		// Try alternative approach - run directly without temp file
		return runGenerateDirect()
	}

	return nil
}

// runGenerateDirect runs generation directly in-process.
func runGenerateDirect() error {
	// Reset the expression root
	expr.Reset()
	eval.Context.Reset()

	// Import and run the design package
	designPkg, err := build.Import("./design", ".", 0)
	if err != nil {
		return fmt.Errorf("failed to import design package: %w", err)
	}

	// This is a simplified approach - in production, we'd need to properly
	// compile and run the design package. For now, we'll provide instructions.
	
	fmt.Println("\n⚠️  Alternative generation method:")
	fmt.Println("\nTo generate code, create a file called 'generate.go' with:")
	fmt.Println("\n```go")
	fmt.Printf(`// +build ignore

package main

import (
	_ "%s/design"
	"gluey.dev/gluey/codegen"
	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
	"log"
)

func main() {
	if err := eval.RunDSL(); err != nil {
		log.Fatal(err)
	}
	
	gen := codegen.NewInterfaceGenerator(expr.Root, "gen")
	if err := gen.Generate(); err != nil {
		log.Fatal(err)
	}
}
`, filepath.Base(designPkg.Dir))
	fmt.Println("```")
	fmt.Println("\nThen run: go run generate.go")
	
	return fmt.Errorf("automatic generation not available - see instructions above")
}