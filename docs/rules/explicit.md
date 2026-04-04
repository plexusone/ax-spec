# Explicit Rules

Rules ensuring all constraints are declared in the specification.

## ax-explicit-request-body-required

**Severity:** warn

Request bodies must set `required` explicitly.

### Why

Ambiguity about whether a request body is required leads to agent errors. Explicit declaration enables proper validation.

### Bad

```yaml
post:
  operationId: createUser
  requestBody:
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/User'
```

### Good

```yaml
post:
  operationId: createUser
  requestBody:
    required: true
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/User'
```

## ax-explicit-parameter-required

**Severity:** warn

Parameters must set `required` explicitly.

### Why

Implicit optionality causes agent confusion. Be explicit about whether parameters are required.

### Bad

```yaml
parameters:
  - name: limit
    in: query
    schema:
      type: integer
```

### Good

```yaml
parameters:
  - name: limit
    in: query
    required: false
    schema:
      type: integer
      default: 20
```

## ax-explicit-enum-values

**Severity:** info

Constrained strings should use `enum`.

### Why

Free-form strings that only accept certain values should declare those values as an enum. This enables:

- Agent validation before API calls
- Code generation with type-safe constants
- Documentation of valid values

### Bad

```yaml
properties:
  status:
    type: string
    description: One of "pending", "completed", or "failed"
```

### Good

```yaml
properties:
  status:
    type: string
    enum:
      - pending
      - completed
      - failed
```
