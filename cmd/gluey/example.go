package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// runExample executes the example generation.
func runExample() error {
	fmt.Println("🎨 Generating example implementation...")

	// Check if design/app.go exists
	designFile := filepath.Join("design", "app.go")
	if _, err := os.Stat(designFile); os.IsNotExist(err) {
		return fmt.Errorf("%s not found - make sure you're in a Gluey project directory", designFile)
	}

	// Read go.mod to get the module name
	goModContent, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	modFile, err := modfile.Parse("go.mod", goModContent, nil)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod: %w", err)
	}
	moduleName := modFile.Module.Mod.Path

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create a temporary directory for the generator
	tmpDir, err := os.MkdirTemp("", "gluey-example-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary main.go that imports and executes the design
	mainContent := fmt.Sprintf(`package main

import (
	"fmt"
	"log"
	"os"
	
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
	
	// Get output directory from environment
	outDir := os.Getenv("GLUEY_OUTPUT")
	if outDir == "" {
		outDir = "."
	}
	
	// Generate examples
	gen := codegen.NewExampleGenerator(expr.Root)
	gen.OutputDir = outDir
	if err := gen.Generate(); err != nil {
		log.Fatal("Example generation failed:", err)
	}
	
	fmt.Println("✅ Example generation complete!")
	fmt.Println("\nGenerated files:")
	fmt.Println("  app/controllers/ - Controller implementations")
	fmt.Println("  app/views/       - HTML templates")
	fmt.Println("  main.go         - Server entry point")
}
`, moduleName)

	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write generator main.go: %w", err)
	}

	// Create a go.mod for the temp directory
	tmpGoModContent := fmt.Sprintf(`module gluey-example-generator

go %s

require (
	%s v0.0.0
	gluey.dev/gluey v0.0.0
)

replace %s => %s
replace gluey.dev/gluey => %s
`, 
		strings.TrimPrefix(modFile.Go.Version, "go"),
		moduleName,
		moduleName, cwd,
		getGlueyPathForExample())
	
	tmpGoMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(tmpGoMod, []byte(tmpGoModContent), 0644); err != nil {
		return fmt.Errorf("failed to write temp go.mod: %w", err)
	}

	// Run go mod tidy in the temp directory to fetch dependencies
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = tmpDir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w\nOutput: %s", err, output)
	}

	// Run the generator
	cmd := exec.Command("go", "run", mainFile)
	cmd.Dir = tmpDir // Run in the temp directory
	cmd.Env = append(os.Environ(), "GLUEY_OUTPUT="+cwd) // Pass output directory
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return fmt.Errorf("example generation failed: %w", err)
	}
	if len(output) > 0 {
		fmt.Println(string(output))
	}

	return nil
}

// getGlueyPathForExample returns the path to the Gluey module for example generation.
func getGlueyPathForExample() string {
	// First, check if we're in the Gluey repo itself
	if _, err := os.Stat(filepath.Join(".", "go.mod")); err == nil {
		if content, err := os.ReadFile("go.mod"); err == nil {
			if strings.Contains(string(content), "module gluey.dev/gluey") {
				cwd, _ := os.Getwd()
				return cwd
			}
		}
	}
	
	// Check parent directories (for examples)
	for i := 1; i <= 3; i++ {
		parentPath := strings.Repeat("../", i)
		goModPath := filepath.Join(parentPath, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			if content, err := os.ReadFile(goModPath); err == nil {
				if strings.Contains(string(content), "module gluey.dev/gluey") {
					abs, _ := filepath.Abs(parentPath)
					return abs
				}
			}
		}
	}
	
	// Default: assume it's available as a module
	return ""
}