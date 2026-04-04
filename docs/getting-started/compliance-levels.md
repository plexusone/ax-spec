# Compliance Levels

AX Spec defines three progressive compliance levels for incremental adoption.

## Overview

| Level | Name | Description |
|-------|------|-------------|
| **AX-L1** | Structured | Valid OpenAPI with explicit types and operation IDs |
| **AX-L2** | Deterministic | Strict schemas, explicit required fields, no ambiguous types |
| **AX-L3** | Agent-Ready | Full AX extensions, structured errors, capability discovery |

## AX-L1: Structured

Minimum requirements for any API agents might consume.

### Requirements

- All schemas have explicit `type` fields
- All operations have `operationId` and `summary`
- Parameters declare `required` explicitly

### Ruleset

```bash
vacuum lint --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/profiles/ax-l1-structured.yaml api.yaml
```

## AX-L2: Deterministic

Strict typing for predictable agent behavior.

### Requirements

Everything in AX-L1, plus:

- `additionalProperties: false` on object schemas
- `x-ax-required-fields` on operations
- Request bodies explicitly set `required`

### Ruleset

```bash
vacuum lint --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/profiles/ax-l2-deterministic.yaml api.yaml
```

## AX-L3: Agent-Ready

Full agent compatibility with error recovery.

### Requirements

Everything in AX-L2, plus:

- `x-ax-capabilities` for operation discovery
- `x-ax-error-code` on error responses
- `x-ax-retryable` and `x-ax-idempotent` flags
- Error responses have defined schemas

### Ruleset

```bash
vacuum lint --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/profiles/ax-l3-agent-ready.yaml api.yaml
```

## Choosing a Level

| If you... | Start with |
|-----------|------------|
| Have an existing API | AX-L1 |
| Are building a new API | AX-L2 |
| Are building for AI agents | AX-L3 |

Most teams should aim for **AX-L2** as a baseline, with **AX-L3** for agent-facing endpoints.
