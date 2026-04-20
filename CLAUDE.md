# CLAUDE.md

Project-specific guidelines for Claude Code working on ax-spec.

## Project Overview

AX Spec is a specification and CLI toolset for building Agent Experience (AX) - APIs designed for autonomous AI agents. It provides:

- Spectral-compatible linting rules enforcing DIRECT principles
- CLI tools for linting, enriching, and generating code from OpenAPI specs
- `x-ax-*` OpenAPI extensions for agent metadata

## Architecture

```
ax-spec/
├── cmd/ax-spec/        # CLI commands (lint, enrich, gen)
├── pkg/
│   ├── gen/            # Code generation from x-ax-* extensions
│   └── enrich/         # OpenAPI enrichment and API discovery
├── rules/              # Spectral rulesets and compliance profiles
├── schemas/            # JSON Schema for x-ax-* extensions
├── examples/           # Real-world OpenAPI specs and generated code
└── docs/               # MkDocs documentation site
```

## Build & Test Commands

```bash
# Build
go build ./...

# Test
go test -v ./...

# Lint
golangci-lint run

# Build CLI
go build -o ax-spec ./cmd/ax-spec

# Run CLI
./ax-spec lint examples/openapi/payments-api-ax.yaml
./ax-spec enrich examples/openapi/opik-openapi.yaml -o /tmp/opik-ax.yaml
./ax-spec gen examples/openapi/opik-openapi-ax.yaml -o /tmp/ax --package ax
```

## Dependencies

The `lint` command requires `vacuum` CLI to be installed:

```bash
brew install daveshanley/vacuum/vacuum
```

## Key Concepts

### DIRECT Principles

The foundation of AX - APIs should be:

- **D**eterministic - Predictable behavior with strict schemas
- **I**ntrospectable - Machine-readable capabilities and metadata
- **R**ecoverable - Structured errors with actionable suggestions
- **E**xplicit - All constraints declared in the specification
- **C**onsistent - Uniform patterns across endpoints
- **T**estable - Safe sandbox environments

### Compliance Levels

- **L1 Structured** - Valid OpenAPI with explicit types and operationIds
- **L2 Deterministic** - Strict schemas, `additionalProperties: false`
- **L3 Agent-Ready** - Full `x-ax-*` extensions, error codes, capabilities

### x-ax-* Extensions

| Extension | Type | Purpose |
|-----------|------|---------|
| `x-ax-required-fields` | `string[]` | Explicit required input fields |
| `x-ax-capabilities` | `string[]` | Semantic capabilities |
| `x-ax-error-code` | `string` | Machine-readable error code |
| `x-ax-error-suggestion` | `string` | Actionable fix suggestion |
| `x-ax-retryable` | `boolean` | Safe to retry? |
| `x-ax-idempotent` | `boolean` | Repeated calls safe? |
| `x-ax-sandboxable` | `boolean` | Safe in test environment? |

## Code Generation

The `gen` command produces four Go files:

- `errors.go` - Error code constants and metadata
- `retry.go` - Retry policy mappings by operation
- `capabilities.go` - Capability constants and operation mappings
- `validation.go` - Required field validators

## Documentation

The MkDocs site is in `docs/`. To serve locally:

```bash
mkdocs serve
```

To build:

```bash
mkdocs build
```

## Changelog

Use structured-changelog format:

```bash
# Validate
schangelog validate CHANGELOG.json

# Generate markdown
schangelog generate CHANGELOG.json -o CHANGELOG.md
```

## Related Repositories

- [DIRECT Principles](https://github.com/grokify/direct-principles) - Conceptual foundation
- [elevenlabs-go](https://github.com/plexusone/elevenlabs-go) - SDK with AX integration
- [opik-go](https://github.com/plexusone/opik-go) - SDK with AX integration
- [omnivoice](https://github.com/plexusone/omnivoice) - Voice abstraction using AX
