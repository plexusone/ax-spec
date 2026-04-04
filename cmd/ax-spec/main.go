// ax-spec is a CLI tool for Agent Experience (AX) compliance.
//
// It provides three main commands:
//   - lint: Check OpenAPI specs against AX rules
//   - enrich: Add x-ax-* extensions to OpenAPI specs
//   - gen: Generate Go code from x-ax-* extensions
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ax-spec",
	Short: "Agent Experience (AX) specification tools",
	Long: `ax-spec provides tools for building agent-friendly APIs.

Commands:
  lint    - Check OpenAPI specs against AX rules (wraps vacuum)
  enrich  - Add x-ax-* extensions to OpenAPI specs
  gen     - Generate Go code from x-ax-* extensions

The AX workflow:
  1. ax-spec lint openapi.yaml           # Check compliance
  2. ax-spec enrich openapi.yaml         # Add x-ax-* extensions
  3. ax-spec gen openapi-ax.yaml         # Generate Go code

Learn more: https://github.com/plexusone/ax-spec`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(lintCmd)
	rootCmd.AddCommand(enrichCmd)
	rootCmd.AddCommand(genCmd)
}
