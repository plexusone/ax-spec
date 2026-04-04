# x-ax-capabilities

Semantic capabilities for agent operation discovery.

## Purpose

Agents need to discover what operations are available and what they do. Capabilities provide:

- Machine-readable operation semantics
- Grouping by functional area
- Discoverability for tool selection

## Schema

```yaml
x-ax-capabilities:
  type: array
  items:
    type: string
    pattern: "^[a-z][a-z0-9_]*$"
```

## Example

```yaml
paths:
  /users:
    get:
      operationId: listUsers
      x-ax-capabilities:
        - list_users
        - search_users
    post:
      operationId: createUser
      x-ax-capabilities:
        - create_user

  /users/{id}:
    get:
      operationId: getUser
      x-ax-capabilities:
        - read_user
        - get_user_details
    delete:
      operationId: deleteUser
      x-ax-capabilities:
        - delete_user
        - remove_user
```

## Generated Code

```go
type Capability string

const (
    CapListUsers   Capability = "list_users"
    CapSearchUsers Capability = "search_users"
    CapCreateUser  Capability = "create_user"
    CapReadUser    Capability = "read_user"
    CapDeleteUser  Capability = "delete_user"
)

// OperationsForCapability returns operations that provide a capability.
var OperationsForCapability = map[Capability][]string{
    CapListUsers:   {"listUsers"},
    CapSearchUsers: {"listUsers"},
    CapCreateUser:  {"createUser"},
    CapReadUser:    {"getUser"},
    CapDeleteUser:  {"deleteUser"},
}
```

## Best Practices

1. **Use snake_case** — Consistent with the pattern requirement
2. **Be specific** — `create_payment` is better than `create`
3. **Include aliases** — Multiple capabilities can map to one operation
4. **Group by domain** — Prefix with domain area (e.g., `payment_`, `user_`)
