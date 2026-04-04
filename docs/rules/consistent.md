# Consistent Rules

Rules ensuring uniform patterns across endpoints.

## ax-consistent-response-content-type

**Severity:** info

Responses should use `application/json`.

### Why

Consistent content types simplify agent implementation. JSON is the standard for API responses and has universal parsing support.

### Bad

```yaml
responses:
  '200':
    content:
      text/plain:
        schema:
          type: string
```

### Good

```yaml
responses:
  '200':
    content:
      application/json:
        schema:
          type: object
          properties:
            message:
              type: string
```

### Exceptions

Some endpoints legitimately return non-JSON content:

- File downloads (`application/octet-stream`)
- Images (`image/png`, `image/jpeg`)
- HTML pages (`text/html`)

For these, the rule is informational.

## ax-consistent-pagination

**Severity:** info

List operations should have pagination parameters.

### Why

Unpaginated list endpoints can return unbounded data, causing:

- Memory issues in agents
- Timeout errors
- Rate limiting

### Bad

```yaml
get:
  operationId: listUsers
  parameters: []
```

### Good

```yaml
get:
  operationId: listUsers
  parameters:
    - name: limit
      in: query
      required: false
      schema:
        type: integer
        default: 20
        maximum: 100
    - name: offset
      in: query
      required: false
      schema:
        type: integer
        default: 0
```

### Alternative: Cursor-based

```yaml
parameters:
  - name: cursor
    in: query
    required: false
    schema:
      type: string
  - name: limit
    in: query
    required: false
    schema:
      type: integer
      default: 20
```
