---
marp: true
theme: default
paginate: true
backgroundColor: #fff
---

# Agent Experience (AX) Case Study

## Integrating AX Spec into [SDK-NAME]

[Tagline: Building Agent-Friendly X SDKs]

---

# The Challenge

<!-- 3-4 bullet points on why agents need this SDK -->

- **Point 1**
- **Point 2**
- **Point 3**

---

# The Problem

When an agent encounters an error:

```json
{
  "status": 404,
  "message": "Not found"
}
```

**What can the agent do?**

- Parse strings? Fragile.
- Retry? Maybe unsafe.
- Recover? How?

---

# The Vision: Agent Experience (AX)

Design interfaces that enable agents to:

> **understand → call → recover → iterate**

---

# DIRECT Principles

| Principle | Application |
|-----------|-------------|
| **D**eterministic | ... |
| **I**ntrospectable | ... |
| **R**ecoverable | ... |
| **E**xplicit | ... |
| **C**onsistent | ... |
| **T**estable | ... |

---

# The Project: [SDK-NAME]

- **X endpoints**
- **Y line OpenAPI spec**
- **[Key feature]**
- Used by AI agents for [purpose]

---

# Implementation Approach

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  OpenAPI Spec   │────▶│   ax-spec CLI   │────▶│  Generated Go   │
│                 │     │                 │     │     Code        │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

---

# Step 1: Define Error Codes

| Category | Codes |
|----------|-------|
| **not_found** | ... |
| **auth** | ... |
| **validation** | ... |

---

# Step 2: Map Retry Policies

| Category | Count | Retryable |
|----------|-------|-----------|
| GET | X | Yes |
| POST | Y | No |
| DELETE | Z | No |

---

# Step 3: Extract Required Fields

```go
var RequiredFields = map[string][]string{
    "createResource": {"name", "type"},
}
```

---

# Before: Generic Error Handling

```go
if err != nil {
    // What kind of error?
    // Should we retry?
    // How to recover?
}
```

---

# After: AX Error Handling

```go
if code, ok := sdk.GetAXErrorCode(err); ok {
    switch code {
    case ax.ErrResourceNotFound:
        // Specific handling
    case ax.ErrUnauthorized:
        // Re-authenticate
    }
}
```

---

# Results Summary

| Metric | Before | After |
|--------|--------|-------|
| Error handling | ... | ... |
| Retry decisions | ... | ... |
| Validation | ... | ... |

---

# Key Learnings

1. **Learning 1**
2. **Learning 2**
3. **Learning 3**

---

# Resources

- **AX Spec:** github.com/plexusone/ax-spec
- **DIRECT Principles:** github.com/grokify/direct-principles
- **[SDK]:** github.com/...

---

# Summary

| Traditional SDK | AX-Enhanced SDK |
|-----------------|-----------------|
| HTTP status codes | Typed error codes |
| Implicit retry rules | Explicit policies |
| Runtime validation | Pre-flight validation |

---

# Thank You

Questions?
