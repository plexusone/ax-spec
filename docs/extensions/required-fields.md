# x-ax-required-fields

Explicit list of fields required for deterministic operation execution.

## Purpose

While OpenAPI schemas define required fields, agents benefit from a consolidated list at the operation level that includes:

- Request body fields
- Path parameters
- Query parameters
- Header requirements

This enables pre-flight validation before making API calls.

## Schema

```yaml
x-ax-required-fields:
  type: array
  items:
    type: string
```

## Example

```yaml
paths:
  /payments:
    post:
      operationId: createPayment
      x-ax-required-fields:
        - amount
        - currency
        - recipient_id
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - amount
                - currency
                - recipient_id
              properties:
                amount:
                  type: number
                currency:
                  type: string
                recipient_id:
                  type: string
```

## Generated Code

```go
// RequiredFieldsForOperation returns the required fields for an operation.
var RequiredFieldsForOperation = map[string][]string{
    "createPayment": {"amount", "currency", "recipient_id"},
}
```

## Best Practices

1. **Include all required inputs** — Not just request body, but also required path/query params
2. **Use consistent naming** — Match the field names in your schema
3. **Order by importance** — Put the most important fields first
