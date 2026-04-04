# Rules Reference

AX Spec provides Spectral/Vacuum rules organized by the DIRECT principles.

## Overview

| Principle | Rules | Focus |
|-----------|-------|-------|
| [Deterministic](deterministic.md) | 3 | Strict schemas, explicit types |
| [Introspectable](introspectable.md) | 3 | Operation metadata, capabilities |
| [Recoverable](recoverable.md) | 4 | Error handling, retry semantics |
| [Explicit](explicit.md) | 3 | Required fields, constraints |
| [Consistent](consistent.md) | 2 | Patterns, conventions |
| [Testable](testable.md) | 2 | Sandbox safety, examples |

## Using the Rules

### With Vacuum

```bash
vacuum lint --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/ax-openapi.json api.yaml
```

### With Spectral

```bash
spectral lint api.yaml --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/ax-openapi.json
```

## Severity Levels

| Severity | Description |
|----------|-------------|
| `error` | Must fix — blocks agent usage |
| `warn` | Should fix — degrades agent experience |
| `info` | Consider — improves agent experience |

## Quick Reference

| Rule | Severity | Description |
|------|----------|-------------|
| `ax-deterministic-schema-type` | error | Schemas must have explicit `type` |
| `ax-deterministic-required-fields` | warn | Operations should declare `x-ax-required-fields` |
| `ax-introspectable-operation-id` | error | Operations must have `operationId` |
| `ax-introspectable-summary` | error | Operations must have `summary` |
| `ax-recoverable-error-schema` | error | Error responses must have schemas |
| `ax-recoverable-error-structure` | warn | Errors should have `x-ax-error-code` |
| `ax-explicit-request-body-required` | warn | Request bodies must set `required` |
