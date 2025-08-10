package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
)

// runExample executes the example generation.
func runExample() error {
	fmt.Println("🎨 Generating example implementation...")

	// Check if design/app.go exists
	designFile := filepath.Join("design", "app.go")
	if _, err := os.Stat(designFile); os.IsNotExist(err) {
		return fmt.Errorf("%s not found - make sure you're in a Gluey project directory", designFile)
	}

	// Reset the expression root
	expr.Reset()
	eval.Context.Reset()

	// TODO: This is a simplified approach - in production, we'd compile and run the design
	// For now, provide instructions
	fmt.Println("\n⚠️  To generate examples, create a file called 'generate_examples.go' with:")
	fmt.Println("\n```go")
	fmt.Println(`// +build ignore

package main

import (
	"log"
	
	_ "YOUR_MODULE/design"
	"gluey.dev/gluey/codegen"
	"gluey.dev/gluey/eval"
	"gluey.dev/gluey/expr"
)

func main() {
	if err := eval.RunDSL(); err != nil {
		log.Fatal(err)
	}
	
	gen := codegen.NewExampleGenerator(expr.Root)
	if err := gen.Generate(); err != nil {
		log.Fatal(err)
	}
}`)
	fmt.Println("```")
	fmt.Println("\nThen run: go run generate_examples.go")
	
	return fmt.Errorf("automatic example generation not yet implemented - see instructions above")
}

// runExampleDirect runs example generation directly in-process.
func runExampleDirect() error {
	// This will be implemented when we have proper design package loading
	return fmt.Errorf("direct example generation not yet implemented")
}