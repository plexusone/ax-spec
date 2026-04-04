# AX Spec Analysis: Real-World OpenAPI Specifications

This document analyzes two production OpenAPI specifications against the AX Spec ruleset to demonstrate Agent Experience compliance gaps in real APIs.

## Summary

| API | Endpoints | Errors | Warnings | Info | Quality Score |
|-----|-----------|--------|----------|------|---------------|
| **Opik** | 201 | 85 | 2,734 | 2,732 | 10/100 |
| **ElevenLabs** | 204 | 0 | 603 | 3,136 | 25/100 |

## Opik API Analysis

**Source:** [opik-go](https://github.com/plexusone/opik-go) (Comet ML observability platform)

### Top AX Violations

| Rule | Violations | Category |
|------|------------|----------|
| `ax-testable-example-values` | 1,744 | Testable |
| `ax-explicit-parameter-required` | 167 | Explicit |
| `ax-introspectable-capabilities-post` | 81 | Introspectable |
| `ax-recoverable-retryable-post` | 81 | Recoverable |
| `ax-deterministic-required-fields-post` | 81 | Deterministic |
| `ax-testable-sandboxable-post` | 81 | Testable |

### Key Issues

1. **Missing Example Values (1,744 violations)**
   - Schema properties lack `example` fields
   - Agents cannot generate test data without examples
   - Fix: Add `example` to all properties

2. **Parameters Without Required Flag (167 violations)**
   - Query/path parameters don't explicitly set `required: true/false`
   - Agents must guess which parameters are optional
   - Fix: Add `required` to all parameters

3. **Missing Agent Metadata (81 violations each)**
   - No `x-ax-capabilities` for agent discovery
   - No `x-ax-retryable` for retry semantics
   - No `x-ax-required-fields` for deterministic inputs
   - No `x-ax-sandboxable` for test safety

### Schema Errors (38)

The `ax-deterministic-schema-type` rule caught properties without explicit types:

```yaml
# Bad - agents can't infer type
properties:
  metadata: {}

# Good - explicit type
properties:
  metadata:
    type: object
```

## ElevenLabs API Analysis

**Source:** [elevenlabs-go](https://github.com/plexusone/elevenlabs-go) (Text-to-speech platform)

### Top AX Violations

| Rule | Violations | Category |
|------|------------|----------|
| `ax-introspectable-capabilities-post` | 109 | Introspectable |
| `ax-deterministic-required-fields-post` | 109 | Deterministic |
| `ax-testable-sandboxable-post` | 109 | Testable |
| `ax-recoverable-retryable-post` | 109 | Recoverable |
| `ax-introspectable-capabilities-get` | 86 | Introspectable |

### Key Issues

1. **No Agent Metadata Extensions**
   - All 109 POST endpoints lack `x-ax-*` extensions
   - 86 GET endpoints lack capability tags
   - This is expected - the spec predates AX conventions

2. **No Schema Type Errors**
   - ElevenLabs spec is well-typed (0 errors from `ax-deterministic-schema-type`)
   - All properties have explicit types

3. **Tags Not Defined (240 warnings)**
   - Operations reference tags not declared in `tags` array
   - Not an AX issue per se, but indicates spec hygiene

### Why ElevenLabs Scored Higher

- **Complete type definitions** - All schemas have explicit types
- **Well-structured spec** - 54K lines, professionally maintained
- **Consistent patterns** - Response shapes are uniform

## DIRECT Principle Mapping

### Violations by Principle

| Principle | Opik | ElevenLabs | Common Issue |
|-----------|------|------------|--------------|
| **D**eterministic | 81 | 109 | Missing `x-ax-required-fields` |
| **I**ntrospectable | 81 | 195 | Missing `x-ax-capabilities` |
| **R**ecoverable | 81 | 133 | Missing `x-ax-retryable` |
| **E**xplicit | 167 | 0 | Missing `required` on params |
| **C**onsistent | - | - | (naming rules) |
| **T**estable | 1,825 | 109 | Missing examples, sandboxable |

### Analysis

1. **Both APIs lack agent metadata** - Neither spec uses `x-ax-*` extensions because they predate the AX standard. This is the primary gap.

2. **Opik has structural issues** - Missing explicit types, missing required flags on parameters.

3. **ElevenLabs is structurally sound** - Well-typed schemas, but no agent-specific metadata.

## Recommendations

### For Opik

1. Add `type` to all schema properties (fixes 38 errors)
2. Add `required: true/false` to all parameters (fixes 167 warnings)
3. Add `example` values throughout (fixes 1,744 info)
4. Consider adding `x-ax-*` extensions for agent workflows

### For ElevenLabs

1. Add `x-ax-capabilities` to help agents discover endpoints
2. Add `x-ax-retryable: false` to TTS endpoints (idempotent concerns)
3. Add `x-ax-sandboxable: true` for safe testing

### For API Providers Generally

1. **Start with AX-L1** - Ensure all types explicit, all operationIds present
2. **Progress to AX-L2** - Add `additionalProperties: false`, fix required arrays
3. **Target AX-L3** - Add `x-ax-*` extensions for full agent compatibility

## Running This Analysis

```bash
# Install vacuum
npm install -g @quobix/vacuum

# Run AX Spec rules
vacuum lint \
  --ruleset ax-spec/rules/ax-openapi.json \
  your-openapi.yaml
```

## Files in This Directory

- `opik-openapi.yaml` - Opik API spec (201 endpoints)
- `elevenlabs-openapi.json` - ElevenLabs API spec (204 endpoints)
- `payments-api-ax.yaml` - Example AX-compliant spec
- `reports/opik-report.txt` - Full Opik lint output
- `reports/elevenlabs-report.txt` - Full ElevenLabs lint output
