# Testable Rules

Rules ensuring safe sandbox environments for agent experimentation.

## ax-testable-sandboxable

**Severity:** info

Operations should indicate sandbox safety with `x-ax-sandboxable`.

### Why

Agents need to experiment with APIs safely. Knowing which operations can be called in test/sandbox environments enables:

- Safe exploration during development
- Automated testing without side effects
- Confidence in agent behavior

### Bad

```yaml
post:
  operationId: createPayment
```

### Good

```yaml
post:
  operationId: createPayment
  x-ax-sandboxable: true
  x-ax-retryable: false
```

### Guidelines

| Operation Type | Sandboxable |
|----------------|-------------|
| Read operations | Yes |
| Test/mock endpoints | Yes |
| Reversible operations | Yes |
| Financial transactions | Depends on sandbox mode |
| External integrations | Usually no |
| Irreversible mutations | Usually no |

## ax-testable-example-values

**Severity:** info

Properties should have example values.

### Why

Examples help agents understand:

- Expected data formats
- Reasonable default values
- Valid input ranges

### Bad

```yaml
properties:
  email:
    type: string
    format: email
  created_at:
    type: string
    format: date-time
```

### Good

```yaml
properties:
  email:
    type: string
    format: email
    example: "user@example.com"
  created_at:
    type: string
    format: date-time
    example: "2024-01-15T09:30:00Z"
```

### Schema-Level Examples

For complex objects, use schema-level examples:

```yaml
User:
  type: object
  properties:
    id:
      type: string
    name:
      type: string
    email:
      type: string
  example:
    id: "usr_123"
    name: "Jane Doe"
    email: "jane@example.com"
```
