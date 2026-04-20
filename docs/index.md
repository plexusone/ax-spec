# AX Spec

Agent Experience (AX) specification and linting rules for building agent-friendly APIs.

AX Spec provides [Spectral](https://github.com/stoplightio/spectral) rules that enforce the [DIRECT principles](https://github.com/grokify/direct-principles) for APIs designed to be consumed by autonomous AI agents.

## Why AX?

Traditional APIs are designed for human developers who can interpret documentation, handle ambiguity, and debug errors manually. Agent Experience (AX) requires APIs that are:

| Principle | Description |
|-----------|-------------|
| **Deterministic** | Predictable behavior with strict schemas |
| **Introspectable** | Machine-readable capabilities and metadata |
| **Recoverable** | Structured errors with actionable suggestions |
| **Explicit** | All constraints declared in the specification |
| **Consistent** | Uniform patterns across endpoints |
| **Testable** | Safe sandbox environments for agent experimentation |

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

### Using the CLI

```bash
# Build the CLI
go install github.com/plexusone/ax-spec/cmd/ax-spec@latest

# Lint your OpenAPI spec
ax-spec lint openapi.yaml

# Enrich spec with x-ax-* extensions
ax-spec enrich openapi.yaml -o openapi-ax.yaml

# Generate Go code from extensions
ax-spec gen openapi-ax.yaml -o pkg/ax
```

## Compliance Levels

AX Spec defines progressive compliance levels for incremental adoption:

| Level | Name | Description |
|-------|------|-------------|
| **AX-L1** | Structured | Valid OpenAPI with explicit types and operation IDs |
| **AX-L2** | Deterministic | Strict schemas, explicit required fields, no ambiguous types |
| **AX-L3** | Agent-Ready | Full AX extensions, structured errors, capability discovery |

## Case Studies

See real-world implementations of AX principles:

| SDK | Domain | Endpoints | Error Codes | Retry Policies |
|-----|--------|-----------|-------------|----------------|
| [elevenlabs-go](case-studies/elevenlabs-go/index.md) | Voice generation | 204 | 9 | 236 |
| [opik-go](case-studies/opik-go/index.md) | LLM observability | 201 | 19 | 201 |

## Releases

- [v0.1.0](releases/v0.1.0.md) — Initial release with CLI, Spectral rules, and code generation

## Resources

- [DIRECT Principles](https://github.com/grokify/direct-principles) — The conceptual foundation for AX
- [Spectral](https://github.com/stoplightio/spectral) — OpenAPI linter
- [Vacuum](https://github.com/daveshanley/vacuum) — Fast Spectral-compatible linter
