package enrich

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DiscoveryResult contains discovered error information.
type DiscoveryResult struct {
	OperationID string
	StatusCode  int
	ErrorCode   string
	ErrorType   string
	Message     string
}

// Discoverer makes API calls to discover error codes.
type Discoverer struct {
	client  *http.Client
	baseURL string
	token   string
	results []DiscoveryResult
}

// NewDiscoverer creates a new API discoverer.
func NewDiscoverer(baseURL, token string) *Discoverer {
	return &Discoverer{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		results: []DiscoveryResult{},
	}
}

// DiscoverErrors attempts to trigger errors from the API to discover error codes.
func (d *Discoverer) DiscoverErrors(spec *Spec) ([]DiscoveryResult, error) {
	paths, ok := spec.raw["paths"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("no paths found in spec")
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

			operationID, _ := operation["operationId"].(string)
			if operationID == "" {
				continue
			}

			// Try to trigger validation errors
			result := d.probeEndpoint(method, pathKey, operation)
			if result != nil {
				result.OperationID = operationID
				d.results = append(d.results, *result)
			}
		}
	}

	return d.results, nil
}

// probeEndpoint makes a request to discover error response format.
func (d *Discoverer) probeEndpoint(method, path string, operation map[string]any) *DiscoveryResult {
	// Build URL - replace path parameters with invalid values
	url := d.baseURL + replacePlaceholders(path, "invalid-id-12345")

	var body io.Reader
	if method == "post" || method == "put" || method == "patch" {
		// Send empty or invalid body to trigger validation error
		body = bytes.NewReader([]byte("{}"))
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, body)
	if err != nil {
		return nil
	}

	// Add auth header
	if d.token != "" {
		req.Header.Set("xi-api-key", d.token) // ElevenLabs style
		req.Header.Set("Authorization", "Bearer "+d.token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Only interested in error responses
	if resp.StatusCode < 400 {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return parseErrorResponse(resp.StatusCode, respBody)
}

// replacePlaceholders replaces {param} with a test value.
func replacePlaceholders(path, value string) string {
	result := path
	for {
		start := strings.Index(result, "{")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		result = result[:start] + value + result[start+end+1:]
	}
	return result
}

// parseErrorResponse extracts error code from various error formats.
func parseErrorResponse(statusCode int, body []byte) *DiscoveryResult {
	result := &DiscoveryResult{
		StatusCode: statusCode,
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		// Not JSON, can't parse
		return result
	}

	// Try common error code field names
	codeFields := []string{
		"error_code", "errorCode", "code",
		"error.code", "error.type",
		"status.error_code",
	}

	for _, field := range codeFields {
		if code := extractField(data, field); code != "" {
			result.ErrorCode = code
			break
		}
	}

	// Try common error type field names
	typeFields := []string{
		"error_type", "errorType", "type",
		"error.type", "detail.type",
	}

	for _, field := range typeFields {
		if errType := extractField(data, field); errType != "" {
			result.ErrorType = errType
			break
		}
	}

	// Try common message field names
	msgFields := []string{
		"message", "error", "detail", "error_message",
		"error.message", "detail.message", "detail.msg",
	}

	for _, field := range msgFields {
		if msg := extractField(data, field); msg != "" {
			result.Message = msg
			break
		}
	}

	// ElevenLabs specific: detail.status field
	if status := extractField(data, "detail.status"); status != "" {
		result.ErrorCode = status
	}

	return result
}

// extractField extracts a potentially nested field from a map.
func extractField(data map[string]any, field string) string {
	parts := strings.Split(field, ".")
	current := any(data)

	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return ""
		}
		current = m[part]
	}

	switch v := current.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%v", v)
	default:
		return ""
	}
}

// ApplyDiscoveredErrors adds discovered error codes to the spec.
func ApplyDiscoveredErrors(spec *Spec, results []DiscoveryResult) int {
	count := 0

	paths, ok := spec.raw["paths"].(map[string]any)
	if !ok {
		return 0
	}

	// Build lookup by operationId
	errorsByOp := make(map[string]DiscoveryResult)
	for _, r := range results {
		if r.ErrorCode != "" || r.ErrorType != "" {
			errorsByOp[r.OperationID] = r
		}
	}

	for _, pathValue := range paths {
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

			operationID, _ := operation["operationId"].(string)
			result, found := errorsByOp[operationID]
			if !found {
				continue
			}

			// Add to responses
			responses, ok := operation["responses"].(map[string]any)
			if !ok {
				responses = make(map[string]any)
				operation["responses"] = responses
			}

			statusKey := fmt.Sprintf("%d", result.StatusCode)
			resp, ok := responses[statusKey].(map[string]any)
			if !ok {
				resp = make(map[string]any)
				responses[statusKey] = resp
			}

			code := result.ErrorCode
			if code == "" {
				code = result.ErrorType
			}
			if code != "" {
				resp["x-ax-error-code"] = strings.ToUpper(strings.ReplaceAll(code, " ", "_"))
				count++
			}

			if result.Message != "" {
				resp["x-ax-error-suggestion"] = result.Message
			}
		}
	}

	return count
}
