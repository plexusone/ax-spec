# Introspectable Rules

Rules ensuring machine-readable capabilities and metadata.

## ax-introspectable-operation-id

**Severity:** error

Operations must have `operationId`.

### Why

Operation IDs are essential for:

- Code generation (function names)
- Agent tool selection
- Logging and tracing
- Documentation cross-references

### Bad

```yaml
paths:
  /users:
    get:
      summary: List users
```

### Good

```yaml
paths:
  /users:
    get:
      operationId: listUsers
      summary: List users
```

## ax-introspectable-summary

**Severity:** error

Operations must have `summary`.

### Why

Summaries help agents understand what an operation does. They're used for:

- Tool descriptions in agent prompts
- Documentation generation
- API discovery

### Bad

```yaml
paths:
  /users:
    get:
      operationId: listUsers
```

### Good

```yaml
paths:
  /users:
    get:
      operationId: listUsers
      summary: List all users with optional filtering
```

## ax-introspectable-capabilities

**Severity:** info

Operations should declare `x-ax-capabilities`.

### Why

Capabilities enable agents to discover operations by intent rather than by name. An agent looking for "how to create a payment" can find operations with `create_payment` capability.

### Bad

```yaml
paths:
  /payments:
    post:
      operationId: createPayment
      summary: Create a payment
```

### Good

```yaml
paths:
  /payments:
    post:
      operationId: createPayment
      summary: Create a payment
      x-ax-capabilities:
        - create_payment
        - transfer_funds
```
