# Deterministic Rules

Rules ensuring predictable behavior with strict schemas.

## ax-deterministic-schema-type

**Severity:** error

Schema properties must have explicit `type` fields.

### Why

Agents cannot handle dynamic typing. Without explicit types, generated code falls back to `interface{}` or `any`, losing type safety.

### Bad

```yaml
properties:
  data:
    description: The response data
```

### Good

```yaml
properties:
  data:
    type: object
    description: The response data
```

## ax-deterministic-required-fields

**Severity:** warn

Operations should declare `x-ax-required-fields`.

### Why

Agents need to know which fields are required before making API calls. This enables pre-flight validation and better error messages.

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
  x-ax-required-fields:
    - email
    - name
  requestBody:
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/User'
```

## ax-deterministic-no-additional-properties

**Severity:** warn

Object schemas should set `additionalProperties: false`.

### Why

Without this, APIs can return unexpected fields that agents don't know how to handle. Strict schemas ensure predictable responses.

### Bad

```yaml
User:
  type: object
  properties:
    id:
      type: string
    name:
      type: string
```

### Good

```yaml
User:
  type: object
  additionalProperties: false
  properties:
    id:
      type: string
    name:
      type: string
```
