# x-ax-retryable

Indicates whether an operation is safe to retry on failure.

## Purpose

Agents need to know whether failed operations can be safely retried. This prevents:

- Duplicate mutations (double-charging, duplicate records)
- Data corruption
- Wasted resources on non-retryable failures

## Schema

```yaml
x-ax-retryable:
  type: boolean
  default: false
```

## Example

```yaml
paths:
  /users:
    get:
      operationId: listUsers
      x-ax-retryable: true  # Safe - read operation
    post:
      operationId: createUser
      x-ax-retryable: false  # Unsafe - may create duplicate

  /payments/{id}:
    get:
      operationId: getPayment
      x-ax-retryable: true  # Safe - read operation
    put:
      operationId: updatePayment
      x-ax-retryable: true  # Safe - idempotent update
      x-ax-idempotent: true
```

## Generated Code

```go
// RetryPolicy defines retry behavior for an operation
type RetryPolicy struct {
    Retryable    bool
    MaxRetries   int
    BackoffMs    int
    BackoffScale float64
}

// RetryPolicies maps operations to their retry behavior
var RetryPolicies = map[string]RetryPolicy{
    "listUsers":     {true, 3, 1000, 2.0},
    "createUser":    {false, 0, 0, 0},
    "getPayment":    {true, 3, 1000, 2.0},
    "updatePayment": {true, 3, 1000, 2.0},
}
```

## Decision Matrix

| HTTP Method | Idempotent | Retryable |
|-------------|------------|-----------|
| GET | Yes | Yes |
| HEAD | Yes | Yes |
| OPTIONS | Yes | Yes |
| PUT | Yes* | Yes* |
| DELETE | Yes* | Yes* |
| POST | No | No |
| PATCH | No | No |

*Depends on implementation — use `x-ax-idempotent` to override

## Related Extensions

- `x-ax-idempotent` — Whether repeated calls produce the same result
- `x-ax-error-code` — Error codes may indicate retry eligibility (e.g., `RATE_LIMIT_EXCEEDED`)

## Best Practices

1. **Default to false** — Safe by default
2. **Mark all reads as retryable** — GET, HEAD, OPTIONS
3. **Consider idempotency** — PUT/DELETE may be retryable if idempotent
4. **Never retry non-idempotent mutations** — POST, most PATCH
