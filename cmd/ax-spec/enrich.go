package main

import (
	"fmt"
	"os"

	"github.com/plexusone/ax-spec/pkg/enrich"
	"github.com/spf13/cobra"
)

var enrichCmd = &cobra.Command{
	Use:   "enrich <openapi-spec>",
	Short: "Add x-ax-* extensions to OpenAPI spec",
	Long: `Enrich an OpenAPI specification with x-ax-* extensions for Agent Experience.

This command analyzes your OpenAPI spec and adds:
  - x-ax-capabilities: Derived from operationId (e.g., createTrace -> create_trace)
  - x-ax-required-fields: Copied from schema required arrays
  - x-ax-retryable: true for GET, false for mutations (configurable)
  - x-ax-sandboxable: true by default (configurable)

With --discover flag, makes actual API calls to discover:
  - x-ax-error-code: From real error responses
  - Retry behavior validation

Examples:
  # Infer extensions from spec structure
  ax-spec enrich openapi.yaml -o openapi-ax.yaml

  # Also discover via API calls
  ax-spec enrich openapi.yaml -o openapi-ax.yaml \
    --discover \
    --api-base https://api.comet.com/opik/api \
    --api-token $OPIK_API_KEY

  # Preview changes without writing
  ax-spec enrich openapi.yaml --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runEnrich,
}

var (
	enrichOutput   string
	enrichDiscover bool
	enrichAPIBase  string
	enrichAPIToken string
	enrichDryRun   bool
)

func init() {
	enrichCmd.Flags().StringVarP(&enrichOutput, "output", "o", "", "Output file (default: stdout)")
	enrichCmd.Flags().BoolVar(&enrichDiscover, "discover", false, "Make API calls to discover error codes")
	enrichCmd.Flags().StringVar(&enrichAPIBase, "api-base", "", "API base URL for discovery")
	enrichCmd.Flags().StringVar(&enrichAPIToken, "api-token", "", "API token for discovery (or use env var)")
	enrichCmd.Flags().BoolVar(&enrichDryRun, "dry-run", false, "Preview changes without writing")
}

func runEnrich(cmd *cobra.Command, args []string) error {
	specPath := args[0]

	// Load spec
	spec, err := enrich.LoadSpec(specPath)
	if err != nil {
		return fmt.Errorf("failed to load spec: %w", err)
	}

	// Create enricher
	opts := enrich.Options{
		InferCapabilities:   true,
		InferRequiredFields: true,
		InferRetryable:      true,
		InferSandboxable:    true,
	}

	if enrichDiscover {
		if enrichAPIBase == "" {
			return fmt.Errorf("--api-base required when using --discover")
		}
		token := enrichAPIToken
		if token == "" {
			token = os.Getenv("API_TOKEN")
		}
		opts.Discover = true
		opts.APIBase = enrichAPIBase
		opts.APIToken = token
	}

	enricher := enrich.New(opts)

	// Enrich spec
	enrichedSpec, report, err := enricher.Enrich(spec)
	if err != nil {
		return fmt.Errorf("enrichment failed: %w", err)
	}

	// Print report
	fmt.Fprintf(os.Stderr, "Enrichment report:\n")
	fmt.Fprintf(os.Stderr, "  Operations processed: %d\n", report.OperationsProcessed)
	fmt.Fprintf(os.Stderr, "  Capabilities added:   %d\n", report.CapabilitiesAdded)
	fmt.Fprintf(os.Stderr, "  Required fields added: %d\n", report.RequiredFieldsAdded)
	fmt.Fprintf(os.Stderr, "  Retryable flags added: %d\n", report.RetryableFlagsAdded)
	if opts.Discover {
		fmt.Fprintf(os.Stderr, "  Error codes discovered: %d\n", report.ErrorCodesDiscovered)
	}

	if enrichDryRun {
		fmt.Fprintf(os.Stderr, "\n[dry-run] Would write to: %s\n", enrichOutput)
		return nil
	}

	// Write output
	output := os.Stdout
	if enrichOutput != "" {
		f, err := os.Create(enrichOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		output = f
	}

	if err := enrich.WriteSpec(enrichedSpec, output); err != nil {
		return fmt.Errorf("failed to write spec: %w", err)
	}

	if enrichOutput != "" {
		fmt.Fprintf(os.Stderr, "\nWrote enriched spec to: %s\n", enrichOutput)
	}

	return nil
}
