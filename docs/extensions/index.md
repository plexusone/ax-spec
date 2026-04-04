# x-ax-* Extensions

AX Spec defines OpenAPI extensions that provide machine-readable metadata for AI agents.

## Overview

| Extension | Type | Description |
|-----------|------|-------------|
| [`x-ax-required-fields`](required-fields.md) | `string[]` | Explicit list of fields required for operation |
| [`x-ax-capabilities`](capabilities.md) | `string[]` | Semantic capabilities for agent discovery |
| [`x-ax-error-code`](error-code.md) | `string` | Machine-readable error code |
| `x-ax-error-suggestion` | `string` | Actionable fix suggestion |
| [`x-ax-retryable`](retryable.md) | `boolean` | Whether operation is safe to retry |
| `x-ax-idempotent` | `boolean` | Whether repeated calls are safe |
| `x-ax-sandboxable` | `boolean` | Whether safe to call in test environment |
| `x-ax-category` | `string` | Functional grouping for discovery |
| `x-ax-cost-estimate` | `object` | Resource consumption estimate |

## Example

```yaml
paths:
  /payments:
    post:
      operationId: createPayment
      summary: Create a new payment
      x-ax-required-fields:
        - amount
        - currency
        - recipient_id
      x-ax-capabilities:
        - create_payment
        - transfer_funds
      x-ax-retryable: false
      x-ax-sandboxable: true
      responses:
        '400':
          x-ax-error-code: INVALID_PAYMENT_REQUEST
          x-ax-error-suggestion: Check amount is positive and currency is valid ISO 4217 code
```

## Schema

Extensions are validated against the [x-ax-extensions.schema.json](https://github.com/plexusone/ax-spec/blob/main/schemas/x-ax-extensions.schema.json) schema.
