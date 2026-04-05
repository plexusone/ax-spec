// Package enrich provides functionality to add x-ax-* extensions to OpenAPI specs.
package enrich

import (
	"io"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Options configures the enrichment process.
type Options struct {
	InferCapabilities   bool
	InferRequiredFields bool
	InferRetryable      bool
	InferSandboxable    bool
	Discover            bool
	APIBase             string
	APIToken            string
}

// Report contains statistics from the enrichment process.
type Report struct {
	OperationsProcessed   int
	CapabilitiesAdded     int
	RequiredFieldsAdded   int
	RetryableFlagsAdded   int
	SandboxableFlagsAdded int
	ErrorCodesDiscovered  int
}

// Enricher adds x-ax-* extensions to OpenAPI specs.
type Enricher struct {
	opts Options
}

// New creates a new Enricher with the given options.
func New(opts Options) *Enricher {
	return &Enricher{opts: opts}
}

// Spec represents a parsed OpenAPI specification.
type Spec struct {
	raw  map[string]any
	node *yaml.Node
}

// LoadSpec loads an OpenAPI spec from a file.
func LoadSpec(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, err
	}

	return &Spec{raw: raw, node: &node}, nil
}

// WriteSpec writes a spec to the given writer.
func WriteSpec(spec *Spec, w io.Writer) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	return enc.Encode(spec.raw)
}

// Enrich adds x-ax-* extensions to the spec.
func (e *Enricher) Enrich(spec *Spec) (*Spec, *Report, error) {
	report := &Report{}

	paths, ok := spec.raw["paths"].(map[string]any)
	if !ok {
		return spec, report, nil
	}

	for pathKey, pathValue := range paths {
		pathItem, ok := pathValue.(map[string]any)
		if !ok {
			continue
		}

		for method, opValue := range pathItem {
			if !isHTTPMethod(method) {
				continue
			}

			operation, ok := opValue.(map[string]any)
			if !ok {
				continue
			}

			report.OperationsProcessed++

			// Infer x-ax-capabilities from operationId
			if e.opts.InferCapabilities {
				if operationID, ok := operation["operationId"].(string); ok {
					if _, exists := operation["x-ax-capabilities"]; !exists {
						caps := inferCapabilities(operationID, method, pathKey)
						operation["x-ax-capabilities"] = caps
						report.CapabilitiesAdded++
					}
				}
			}

			// Infer x-ax-required-fields from requestBody schema
			if e.opts.InferRequiredFields && (method == "post" || method == "put" || method == "patch") {
				if _, exists := operation["x-ax-required-fields"]; !exists {
					if reqFields := inferRequiredFields(operation, spec.raw); len(reqFields) > 0 {
						operation["x-ax-required-fields"] = reqFields
						report.RequiredFieldsAdded++
					}
				}
			}

			// Infer x-ax-retryable
			if e.opts.InferRetryable {
				if _, exists := operation["x-ax-retryable"]; !exists {
					retryable := inferRetryable(method)
					operation["x-ax-retryable"] = retryable
					report.RetryableFlagsAdded++
				}
			}

			// Infer x-ax-sandboxable
			if e.opts.InferSandboxable {
				if _, exists := operation["x-ax-sandboxable"]; !exists {
					operation["x-ax-sandboxable"] = true
					report.SandboxableFlagsAdded++
				}
			}
		}
	}

	// API discovery for error codes
	if e.opts.Discover && e.opts.APIBase != "" {
		discovered := e.discoverErrorCodes(spec, e.opts.APIBase, e.opts.APIToken)
		report.ErrorCodesDiscovered = discovered
	}

	return spec, report, nil
}

func isHTTPMethod(method string) bool {
	switch method {
	case "get", "post", "put", "patch", "delete", "head", "options":
		return true
	}
	return false
}

// inferCapabilities derives x-ax-capabilities from operationId.
// Example: "createTrace" -> ["create_trace"]
func inferCapabilities(operationID, _, _ string) []string {
	// Convert camelCase to snake_case
	cap := camelToSnake(operationID)

	// Add method-based capability if not already implied
	caps := []string{cap}

	// Add semantic capabilities based on common patterns
	switch {
	case strings.HasPrefix(operationID, "create") || strings.HasPrefix(operationID, "add"):
		caps = append(caps, "write")
	case strings.HasPrefix(operationID, "get") || strings.HasPrefix(operationID, "list") || strings.HasPrefix(operationID, "find"):
		caps = append(caps, "read")
	case strings.HasPrefix(operationID, "update") || strings.HasPrefix(operationID, "patch"):
		caps = append(caps, "write")
	case strings.HasPrefix(operationID, "delete") || strings.HasPrefix(operationID, "remove"):
		caps = append(caps, "write", "delete")
	}

	return caps
}

var camelRegex = regexp.MustCompile("([a-z0-9])([A-Z])")

func camelToSnake(s string) string {
	snake := camelRegex.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

// inferRequiredFields extracts required fields from requestBody schema.
func inferRequiredFields(operation map[string]any, spec map[string]any) []string {
	requestBody, ok := operation["requestBody"].(map[string]any)
	if !ok {
		return nil
	}

	content, ok := requestBody["content"].(map[string]any)
	if !ok {
		return nil
	}

	// Try application/json first
	jsonContent, ok := content["application/json"].(map[string]any)
	if !ok {
		return nil
	}

	schema, ok := jsonContent["schema"].(map[string]any)
	if !ok {
		return nil
	}

	// Handle $ref
	if ref, ok := schema["$ref"].(string); ok {
		schema = resolveRef(ref, spec)
		if schema == nil {
			return nil
		}
	}

	required, ok := schema["required"].([]any)
	if !ok {
		return nil
	}

	result := make([]string, 0, len(required))
	for _, r := range required {
		if s, ok := r.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func resolveRef(ref string, spec map[string]any) map[string]any {
	// Parse "#/components/schemas/MySchema"
	if !strings.HasPrefix(ref, "#/") {
		return nil
	}

	parts := strings.Split(ref[2:], "/")
	var current any = spec

	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = m[part]
	}

	result, _ := current.(map[string]any)
	return result
}

// inferRetryable determines if an operation should be retryable.
func inferRetryable(method string) bool {
	switch method {
	case "get", "head", "options":
		return true // Read operations are generally safe to retry
	default:
		return false // Mutations need explicit declaration
	}
}

// discoverErrorCodes makes API calls to discover actual error codes.
func (e *Enricher) discoverErrorCodes(spec *Spec, apiBase, apiToken string) int {
	discoverer := NewDiscoverer(apiBase, apiToken)
	results, err := discoverer.DiscoverErrors(spec)
	if err != nil {
		return 0
	}

	return ApplyDiscoveredErrors(spec, results)
}
