# AX Spec

Agent Experience (AX) specification and linting rules for building agent-friendly APIs.

AX Spec provides [Spectral](https://github.com/stoplightio/spectral) rules that enforce the [DIRECT principles](https://github.com/grokify/direct-principles) for APIs designed to be consumed by autonomous AI agents.

## Overview

Traditional APIs are designed for human developers who can interpret documentation, handle ambiguity, and debug errors manually. Agent Experience (AX) requires APIs that are:

- **Deterministic** - Predictable behavior with strict schemas
- **Introspectable** - Machine-readable capabilities and metadata
- **Recoverable** - Structured errors with actionable suggestions
- **Explicit** - All constraints declared in the specification
- **Consistent** - Uniform patterns across endpoints
- **Testable** - Safe sandbox environments for agent experimentation

## Quick Start

### Using with Vacuum

```bash
# Install vacuum
brew install daveshanley/vacuum/vacuum

# Lint your OpenAPI spec with AX rules
vacuum lint --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/ax-openapi.json your-api.yaml

# Use a specific compliance level
vacuum lint --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/profiles/ax-l3-agent-ready.yaml your-api.yaml
```

### Using with Spectral CLI

```bash
# Install spectral
npm install -g @stoplight/spectral-cli

# Lint your OpenAPI spec
spectral lint your-api.yaml --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/ax-openapi.json
```

## Compliance Levels

AX Spec defines progressive compliance levels for incremental adoption:

| Level | Name | Description |
|-------|------|-------------|
| **AX-L1** | Structured | Valid OpenAPI with explicit types and operation IDs |
| **AX-L2** | Deterministic | Strict schemas, explicit required fields, no ambiguous types |
| **AX-L3** | Agent-Ready | Full AX extensions, structured errors, capability discovery |

### L1: Structured (Basic)

Minimum requirements for any API agents might consume:

- All schemas have explicit `type` fields
- All operations have `operationId` and `summary`
- Parameters declare `required` explicitly

### L2: Deterministic

Strict typing for predictable agent behavior:

- `additionalProperties: false` on object schemas
- `x-ax-required-fields` on operations
- Request bodies explicitly set `required`

### L3: Agent-Ready

Full agent compatibility with error recovery:

- `x-ax-capabilities` for operation discovery
- `x-ax-error-code` on error responses
- `x-ax-retryable` and `x-ax-idempotent` flags
- Error responses have defined schemas

## x-ax-* Extensions

AX Spec defines OpenAPI extensions for agent metadata:

| Extension | Type | Description |
|-----------|------|-------------|
| `x-ax-required-fields` | `string[]` | Explicit list of fields required for operation |
| `x-ax-capabilities` | `string[]` | Semantic capabilities (e.g., `create_payment`) |
| `x-ax-error-code` | `string` | Machine-readable error code |
| `x-ax-error-suggestion` | `string` | Actionable fix suggestion |
| `x-ax-retryable` | `boolean` | Whether operation is safe to retry |
| `x-ax-idempotent` | `boolean` | Whether repeated calls are safe |
| `x-ax-sandboxable` | `boolean` | Whether safe to call in test environment |
| `x-ax-category` | `string` | Functional grouping for discovery |
| `x-ax-cost-estimate` | `object` | Resource consumption estimate |

### Example

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

## Rules Reference

### Deterministic Rules

| Rule | Severity | Description |
|------|----------|-------------|
| `ax-deterministic-required-fields` | warn | Operations must declare `x-ax-required-fields` |
| `ax-deterministic-schema-type` | error | Schema properties must have explicit `type` |
| `ax-deterministic-no-additional-properties` | warn | Object schemas should set `additionalProperties: false` |

### Introspectable Rules

| Rule | Severity | Description |
|------|----------|-------------|
| `ax-introspectable-capabilities` | info | Operations should declare `x-ax-capabilities` |
| `ax-introspectable-operation-id` | error | Operations must have `operationId` |
| `ax-introspectable-summary` | error | Operations must have `summary` |

### Recoverable Rules

| Rule | Severity | Description |
|------|----------|-------------|
| `ax-recoverable-error-structure` | warn | Error responses should include `x-ax-error-code` |
| `ax-recoverable-error-schema` | error | Error responses must have a defined schema |
| `ax-recoverable-retryable` | info | Mutating operations should indicate retry semantics |
| `ax-recoverable-idempotent` | info | PUT/DELETE should indicate idempotency |

### Explicit Rules

| Rule | Severity | Description |
|------|----------|-------------|
| `ax-explicit-request-body-required` | warn | Request bodies must set `required` explicitly |
| `ax-explicit-parameter-required` | warn | Parameters must set `required` explicitly |
| `ax-explicit-enum-values` | info | Constrained strings should use `enum` |

### Consistent Rules

| Rule | Severity | Description |
|------|----------|-------------|
| `ax-consistent-response-content-type` | info | Responses should use `application/json` |
| `ax-consistent-pagination` | info | List operations should have pagination params |

### Testable Rules

| Rule | Severity | Description |
|------|----------|-------------|
| `ax-testable-sandboxable` | info | Operations should indicate sandbox safety |
| `ax-testable-example-values` | info | Properties should have example values |

## Real-World Analysis

AX Spec includes analysis of production OpenAPI specifications:

| API | Endpoints | Errors | Warnings | Quality Score |
|-----|-----------|--------|----------|---------------|
| **Opik** (observability) | 201 | 85 | 2,734 | 10/100 |
| **ElevenLabs** (TTS) | 204 | 0 | 603 | 25/100 |

**Key findings:**

- Both APIs lack `x-ax-*` extensions (expected - predates AX standard)
- Opik has structural issues (missing types, missing required flags)
- ElevenLabs is well-typed but lacks agent metadata

See [examples/openapi/ANALYSIS.md](examples/openapi/ANALYSIS.md) for detailed breakdown.

## Project Structure

```
ax-spec/
├── rules/
│   ├── ax-openapi.json       # Main ruleset (JSON canonical)
│   └── profiles/
│       ├── ax-l1-structured.yaml
│       ├── ax-l2-deterministic.yaml
│       └── ax-l3-agent-ready.yaml
├── schemas/
│   └── x-ax-extensions.schema.json
├── examples/
│   └── openapi/
│       ├── payments-api-ax.yaml    # AX-compliant example
│       ├── opik-openapi.yaml       # Real-world: Opik API
│       ├── elevenlabs-openapi.json # Real-world: ElevenLabs API
│       ├── ANALYSIS.md             # Detailed analysis
│       └── reports/
└── docs/
```

## Related Tools

- [DIRECT Principles](https://github.com/grokify/direct-principles) - The conceptual foundation for AX
- [Spectral](https://github.com/stoplightio/spectral) - OpenAPI linter
- [Vacuum](https://github.com/daveshanley/vacuum) - Fast Spectral-compatible linter
- [schemalint](https://github.com/grokify/schemalint) - JSON Schema linter for Go compatibility

## Integration with CI/CD

```yaml
# GitHub Actions example
- name: Lint OpenAPI for Agent Experience
  run: |
    vacuum lint \
      --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/profiles/ax-l3-agent-ready.yaml \
      --fail-severity warn \
      api/openapi.yaml
```

## Background

AX Spec emerged from practical experience building SDKs for dozens of APIs at [PlexusOne](https://github.com/plexusone), including:

- Code generation challenges with ogen and OpenAPI Generator
- Nullable/optional field ambiguity causing runtime errors
- Missing error schemas preventing agent error recovery
- Inconsistent patterns requiring per-API workarounds

See the [Agent Experience article](https://github.com/grokify/grokify-articles/tree/master/agent-experience-ax) for the full context.

## License

MIT
