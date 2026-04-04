# x-ax-error-code

Machine-readable error code for agent error classification.

## Purpose

Human-readable error messages are insufficient for agents. Error codes enable:

- Programmatic error handling
- Automatic retry decisions
- Self-healing workflows
- Error categorization

## Schema

```yaml
x-ax-error-code:
  type: string
  pattern: "^[A-Z][A-Z0-9_]*$"
```

## Example

```yaml
paths:
  /payments:
    post:
      operationId: createPayment
      responses:
        '400':
          description: Invalid request
          x-ax-error-code: INVALID_PAYMENT_REQUEST
          x-ax-error-suggestion: Check amount is positive and currency is valid ISO 4217 code
        '402':
          description: Insufficient funds
          x-ax-error-code: INSUFFICIENT_FUNDS
          x-ax-error-suggestion: Add funds to the account or use a different payment method
        '429':
          description: Rate limited
          x-ax-error-code: RATE_LIMIT_EXCEEDED
          x-ax-error-suggestion: Wait and retry with exponential backoff
```

## Generated Code

```go
// Error codes
const (
    ErrInvalidPaymentRequest = "INVALID_PAYMENT_REQUEST"
    ErrInsufficientFunds     = "INSUFFICIENT_FUNDS"
    ErrRateLimitExceeded     = "RATE_LIMIT_EXCEEDED"
)

// ErrorMetadata provides details for error handling
var ErrorMetadata = map[string]struct {
    HTTPStatus int
    Retryable  bool
    Category   string
}{
    ErrInvalidPaymentRequest: {400, false, "validation"},
    ErrInsufficientFunds:     {402, false, "business"},
    ErrRateLimitExceeded:     {429, true, "rate_limit"},
}
```

## Error Categories

| Category | Description | Retryable |
|----------|-------------|-----------|
| `validation` | Invalid input | No |
| `authentication` | Auth failure | No |
| `authorization` | Permission denied | No |
| `not_found` | Resource missing | No |
| `business` | Business rule violation | No |
| `rate_limit` | Too many requests | Yes (with backoff) |
| `server` | Internal error | Yes |
| `unavailable` | Service down | Yes |

## Best Practices

1. **Use SCREAMING_SNAKE_CASE** — Consistent with the pattern requirement
2. **Be specific** — `INVALID_DATE_RANGE` is better than `INVALID_INPUT`
3. **Include suggestions** — Pair with `x-ax-error-suggestion`
4. **Map to HTTP status** — Error codes should align with response codes
