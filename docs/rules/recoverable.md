# Recoverable Rules

Rules ensuring structured errors with actionable suggestions.

## ax-recoverable-error-schema

**Severity:** error

Error responses must have a defined schema.

### Why

Without schemas, agents receive unstructured error responses they cannot parse programmatically. This prevents automated error recovery.

### Bad

```yaml
responses:
  '400':
    description: Bad request
```

### Good

```yaml
responses:
  '400':
    description: Bad request
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Error'
```

## ax-recoverable-error-structure

**Severity:** warn

Error responses should include `x-ax-error-code`.

### Why

Human-readable error messages are insufficient for agents. Machine-readable error codes enable:

- Programmatic error handling
- Automatic retry decisions
- Self-healing workflows

### Bad

```yaml
responses:
  '400':
    description: Invalid payment amount
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Error'
```

### Good

```yaml
responses:
  '400':
    description: Invalid payment amount
    x-ax-error-code: INVALID_PAYMENT_AMOUNT
    x-ax-error-suggestion: Ensure amount is positive and less than account limit
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/Error'
```

## ax-recoverable-retryable

**Severity:** info

Mutating operations should indicate retry semantics with `x-ax-retryable`.

### Why

Agents need to know whether failed operations can be safely retried to avoid duplicate mutations.

### Bad

```yaml
post:
  operationId: createPayment
```

### Good

```yaml
post:
  operationId: createPayment
  x-ax-retryable: false
```

## ax-recoverable-idempotent

**Severity:** info

PUT/DELETE operations should indicate idempotency with `x-ax-idempotent`.

### Why

Idempotent operations can be safely retried. Agents need this information for reliable error recovery.

### Bad

```yaml
put:
  operationId: updateUser
```

### Good

```yaml
put:
  operationId: updateUser
  x-ax-idempotent: true
  x-ax-retryable: true
```
