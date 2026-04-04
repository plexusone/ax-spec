package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/plexusone/ax-spec/pkg/gen"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen <openapi-spec>",
	Short: "Generate Go code from x-ax-* extensions",
	Long: `Generate Go code from x-ax-* extensions in an OpenAPI specification.

This command reads x-ax-* extensions and generates:
  - errors.go: Error code constants from x-ax-error-code
  - retry.go: Retry policies from x-ax-retryable
  - capabilities.go: Capability mappings from x-ax-capabilities
  - validation.go: Required field validators from x-ax-required-fields

The generated code is designed to complement ogen-generated SDKs,
providing agent-friendly utilities for error handling and retry logic.

Examples:
  # Generate Go code
  ax-spec gen openapi-ax.yaml -o pkg/ax --package ax

  # Generate specific files only
  ax-spec gen openapi-ax.yaml -o pkg/ax --only errors,retry

  # Preview what would be generated
  ax-spec gen openapi-ax.yaml --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runGen,
}

var (
	genOutput  string
	genPackage string
	genOnly    string
	genDryRun  bool
)

func init() {
	genCmd.Flags().StringVarP(&genOutput, "output", "o", ".", "Output directory")
	genCmd.Flags().StringVarP(&genPackage, "package", "p", "ax", "Go package name")
	genCmd.Flags().StringVar(&genOnly, "only", "", "Generate only specific files (comma-separated: errors,retry,capabilities)")
	genCmd.Flags().BoolVar(&genDryRun, "dry-run", false, "Preview what would be generated")
}

func runGen(cmd *cobra.Command, args []string) error {
	specPath := args[0]

	// Load spec and extract x-ax-* extensions
	spec, err := gen.LoadSpec(specPath)
	if err != nil {
		return fmt.Errorf("failed to load spec: %w", err)
	}

	// Extract AX metadata
	meta, err := gen.ExtractAXMetadata(spec)
	if err != nil {
		return fmt.Errorf("failed to extract AX metadata: %w", err)
	}

	// Report what was found
	fmt.Fprintf(os.Stderr, "AX metadata extracted:\n")
	fmt.Fprintf(os.Stderr, "  Error codes:    %d\n", len(meta.ErrorCodes))
	fmt.Fprintf(os.Stderr, "  Operations:     %d\n", len(meta.Operations))
	fmt.Fprintf(os.Stderr, "  Capabilities:   %d\n", len(meta.Capabilities))
	fmt.Fprintf(os.Stderr, "  Required fields: %d operations\n", len(meta.RequiredFields))

	if genDryRun {
		fmt.Fprintf(os.Stderr, "\n[dry-run] Would generate:\n")
		fmt.Fprintf(os.Stderr, "  %s/errors.go\n", genOutput)
		fmt.Fprintf(os.Stderr, "  %s/retry.go\n", genOutput)
		fmt.Fprintf(os.Stderr, "  %s/capabilities.go\n", genOutput)
		fmt.Fprintf(os.Stderr, "  %s/validation.go\n", genOutput)
		return nil
	}

	// Create output directory
	if err := os.MkdirAll(genOutput, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate files
	generator := gen.New(genPackage, meta)

	files := map[string]func() ([]byte, error){
		"errors.go":       generator.GenerateErrors,
		"retry.go":        generator.GenerateRetry,
		"capabilities.go": generator.GenerateCapabilities,
		"validation.go":   generator.GenerateValidation,
	}

	for filename, genFunc := range files {
		content, err := genFunc()
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}

		path := filepath.Join(genOutput, filename)
		if err := os.WriteFile(path, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
		fmt.Fprintf(os.Stderr, "  Generated: %s\n", path)
	}

	fmt.Fprintf(os.Stderr, "\nGeneration complete. Import with:\n")
	fmt.Fprintf(os.Stderr, "  import \"your-module/%s\"\n", genOutput)

	return nil
}
