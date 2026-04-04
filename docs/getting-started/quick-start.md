# Quick Start

Get started with AX Spec in minutes.

## Prerequisites

- Go 1.21+ (for CLI)
- [Vacuum](https://github.com/daveshanley/vacuum) or [Spectral](https://github.com/stoplightio/spectral) (for linting)

## Installation

### CLI

```bash
go install github.com/plexusone/ax-spec/cmd/ax-spec@latest
```

### Ruleset Only

No installation needed — reference the ruleset URL directly:

```bash
vacuum lint --ruleset https://raw.githubusercontent.com/plexusone/ax-spec/main/rules/ax-openapi.json your-api.yaml
```

## Basic Workflow

### 1. Lint Your OpenAPI Spec

```bash
ax-spec lint openapi.yaml
```

This checks your spec against AX rules and reports issues by DIRECT principle.

### 2. Enrich with AX Extensions

```bash
ax-spec enrich openapi.yaml -o openapi-ax.yaml
```

This adds `x-ax-*` extensions based on spec analysis:

- `x-ax-capabilities` from operation IDs
- `x-ax-required-fields` from schema requirements
- `x-ax-retryable` based on HTTP method

### 3. Generate Go Code

```bash
ax-spec gen openapi-ax.yaml -o pkg/ax --package ax
```

This generates Go code from the extensions:

- `errors.go` — Error code constants
- `retry.go` — Retry policy mappings
- `capabilities.go` — Operation capabilities
- `validation.go` — Required field validators

## Next Steps

- [Compliance Levels](compliance-levels.md) — Understand AX-L1, AX-L2, AX-L3
- [CLI Reference](cli.md) — Full command documentation
- [Case Studies](../case-studies/index.md) — Real-world implementations
