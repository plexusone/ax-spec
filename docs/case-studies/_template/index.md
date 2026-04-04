# Agent Experience (AX) Integration Case Study

## Executive Summary

<!-- Brief overview: what SDK, what was integrated, key results -->

This case study documents the integration of Agent Experience (AX) principles into [SDK-NAME], demonstrating how machine-readable metadata enables AI agents to [key benefit].

**Key Results:**

- X error codes [discovered/defined]
- Y operations mapped with retry policies
- Z operations have required field definitions
- [Other key metrics]

## Background

### The Challenge

<!-- What problem does this SDK solve? Why do agents need it? -->

### The Project

<!-- SDK overview: endpoints, spec size, domain, use case -->

[SDK-NAME] is a Go SDK for [SERVICE]:

- **X endpoints** covering [features]
- **Y line OpenAPI specification**
- **[Key capability]** for [use case]
- Used by AI agents for [purpose]

## The Problem

<!-- What errors/issues do agents face without AX? Show example error. -->

When an agent encounters an error:

```json
{
  "status": 404,
  "message": "Not found"
}
```

Questions the agent cannot answer:

1. **What wasn't found?** — [Resource types]
2. **Should it retry?** — [Uncertainty]
3. **How to recover?** — [Options unknown]

### Before AX Integration

```go
// Show problematic code pattern
```

## Solution: AX Integration

### Error Code Design

<!-- List error codes by category -->

| Category | Error Codes | HTTP Status |
|----------|-------------|-------------|
| **not_found** | ... | 404 |
| **auth** | ... | 401, 403 |
| **validation** | ... | 400 |

### Retry Policy Mapping

<!-- Summarize retry policy distribution -->

| Category | Count | Retryable |
|----------|-------|-----------|
| GET | X | Yes |
| POST | Y | No |
| DELETE | Z | No |

### Required Fields Extraction

<!-- Show example required fields -->

```go
var RequiredFields = map[string][]string{
    "createResource": {"name", "type"},
    // ...
}
```

### Capability Mapping

<!-- List domain-specific capabilities -->

## Results

### Error Handling Improvement

**Before:**

```go
// Old pattern
```

**After:**

```go
// New pattern with AX
```

### [Other improvements: retry, validation, etc.]

## Metrics

### Code Changes

| Component | Files | Lines |
|-----------|-------|-------|
| ax package | X | Y |
| errors.go | 1 | Z |
| **Total** | **N** | **M** |

### Coverage

| Metadata Type | Count | Coverage |
|---------------|-------|----------|
| Error codes | X | ... |
| Retry policies | Y | 100% |
| Required fields | Z | ...% |

## Key Learnings

1. **Learning 1** — Explanation
2. **Learning 2** — Explanation
3. **Learning 3** — Explanation

## Future Work

- [ ] Future improvement 1
- [ ] Future improvement 2

## Conclusion

<!-- Summary table: before vs after -->

| Aspect | Before | After |
|--------|--------|-------|
| Error handling | ... | ... |
| Retry decisions | ... | ... |

## References

- [AX Spec](https://github.com/plexusone/ax-spec)
- [DIRECT Principles](https://github.com/grokify/direct-principles)
- [SDK Repository](https://github.com/...)
