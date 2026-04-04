# AX Case Studies

Real-world implementations of Agent Experience (AX) principles in production SDKs.

## Overview

These case studies document the integration of AX metadata into Go SDKs, demonstrating how machine-readable error codes, retry policies, and validation improve agent reliability.

## Case Studies

| SDK | Domain | Endpoints | Error Codes | Retry Policies | Required Fields |
|-----|--------|-----------|-------------|----------------|-----------------|
| [elevenlabs-go](elevenlabs-go/) | Voice generation | 204 | 9 (discovered) | 236 | 72 |
| [opik-go](opik-go/) | LLM observability | 201 | 19 (defined) | 201 | 67 |

## Key Findings

### Common Patterns

Both integrations share common patterns:

1. **Error code constants** — Typed constants for machine-readable errors
2. **Error metadata** — Category, retryability, HTTP status mapping
3. **Retry policy mapping** — All operations classified as safe/unsafe to retry
4. **Required field extraction** — Pre-flight validation before API calls
5. **Capability mapping** — Operation introspection for agents

### Domain-Specific Differences

| Aspect | elevenlabs-go | opik-go |
|--------|---------------|---------|
| Error discovery | API probing | Domain definition |
| Special capabilities | — | Stream, Evaluate, Analytics |
| Self-healing focus | Media errors | Tracing pipeline |
| Retryable errors | Rate limits | Rate limits + server errors |

## Metrics Summary

| Metric | elevenlabs-go | opik-go |
|--------|---------------|---------|
| New code | ~1,075 lines | ~1,030 lines |
| New files | 9 | 7 |
| Test coverage | Full | Full |
| Integration effort | ~2 hours | ~2 hours |

## How to Use These Case Studies

### For SDK Authors

1. Read the case study for an SDK similar to yours
2. Use the [template](_template/) to structure your integration
3. Follow the patterns: error codes → retry policies → validation → capabilities

### For API Providers

1. Review how x-ax-* extensions improve SDK ergonomics
2. Consider adding these extensions to your OpenAPI specs
3. Use ax-spec CLI to enrich existing specs

### For Agent Developers

1. Look for SDKs with AX integration
2. Use the ax package for error handling and retry decisions
3. Implement self-healing patterns based on error categories

## Presentations

Each case study includes a Marp presentation (`presentation.md`). To view as slides:

```bash
# Install marp CLI
npm install -g @marp-team/marp-cli

# Generate HTML slides
marp docs/case-studies/elevenlabs-go/presentation.md -o docs/case-studies/elevenlabs-go/presentation.html
marp docs/case-studies/opik-go/presentation.md -o docs/case-studies/opik-go/presentation.html

# Or serve with live reload
marp -s docs/case-studies/
```

The HTML files are generated during the docs build process and included in the MkDocs site.

## Contributing

To add a new case study:

1. Copy the `_template/` directory
2. Rename to your SDK name
3. Fill in index.md (article) and presentation.md (slides)
4. Submit a pull request

## Resources

- [AX Spec](https://github.com/plexusone/ax-spec) — CLI and rules
- [DIRECT Principles](https://github.com/grokify/direct-principles) — Conceptual foundation
- [elevenlabs-go](https://github.com/plexusone/elevenlabs-go) — Voice SDK with AX
- [opik-go](https://github.com/plexusone/opik-go) — Observability SDK with AX
