package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint <openapi-spec>",
	Short: "Lint OpenAPI spec against AX rules",
	Long: `Lint an OpenAPI specification against AX rules using vacuum.

This command wraps vacuum with the AX ruleset to check for:
  - Deterministic schemas (explicit types, no additionalProperties)
  - Introspectable operations (operationId, summary)
  - Explicit parameters (required flags)
  - Testable schemas (example values)

Levels:
  L1 - Structured: Basic OpenAPI validity
  L2 - Deterministic: Strict schemas
  L3 - Agent-Ready: Full x-ax-* compliance

Examples:
  ax-spec lint openapi.yaml
  ax-spec lint openapi.yaml --level L2
  ax-spec lint openapi.yaml --ruleset custom-rules.yaml`,
	Args: cobra.ExactArgs(1),
	RunE: runLint,
}

var (
	lintLevel   string
	lintRuleset string
)

func init() {
	lintCmd.Flags().StringVarP(&lintLevel, "level", "l", "L3", "Compliance level: L1, L2, or L3")
	lintCmd.Flags().StringVarP(&lintRuleset, "ruleset", "r", "", "Custom ruleset path (overrides --level)")
}

func runLint(cmd *cobra.Command, args []string) error {
	specPath := args[0]

	// Check if vacuum is installed
	if _, err := exec.LookPath("vacuum"); err != nil {
		return fmt.Errorf("vacuum not found. Install with: npm install -g @quobix/vacuum")
	}

	// Determine ruleset
	ruleset := lintRuleset
	if ruleset == "" {
		// Use embedded ruleset based on level
		// For now, use the main ruleset
		// TODO: Support level-specific profiles
		ruleset = getRulesetPath()
	}

	// Run vacuum
	vacuumArgs := []string{"lint", "--ruleset", ruleset, specPath}
	vacuumCmd := exec.Command("vacuum", vacuumArgs...)
	vacuumCmd.Stdout = os.Stdout
	vacuumCmd.Stderr = os.Stderr

	return vacuumCmd.Run()
}

func getRulesetPath() string {
	// Try to find the ruleset relative to the binary or use a known path
	// For development, use the local path
	paths := []string{
		"rules/ax-openapi.json",
		"../rules/ax-openapi.json",
		"../../rules/ax-openapi.json",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// Fallback to GitHub raw URL
	return "https://raw.githubusercontent.com/grokify/ax-spec/main/rules/ax-openapi.json"
}
