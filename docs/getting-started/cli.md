# CLI Reference

The `ax-spec` CLI provides commands for linting, enriching, and generating code from OpenAPI specifications.

## Installation

```bash
go install github.com/plexusone/ax-spec/cmd/ax-spec@latest
```

## Commands

### lint

Check OpenAPI specs against AX rules.

```bash
ax-spec lint <openapi-spec> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--level` | Compliance level (l1, l2, l3) |
| `--format` | Output format (text, json) |

**Example:**

```bash
ax-spec lint api.yaml --level l2
```

### enrich

Add `x-ax-*` extensions to OpenAPI specs.

```bash
ax-spec enrich <openapi-spec> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file (default: stdout) |
| `--discover` | Make API calls to discover error codes |
| `--api-base` | API base URL for discovery |
| `--api-token` | API token for discovery |
| `--dry-run` | Preview changes without writing |

**Example:**

```bash
# Infer extensions from spec structure
ax-spec enrich api.yaml -o api-ax.yaml

# Also discover via API calls
ax-spec enrich api.yaml -o api-ax.yaml \
  --discover \
  --api-base https://api.example.com \
  --api-token $API_TOKEN
```

### gen

Generate Go code from `x-ax-*` extensions.

```bash
ax-spec gen <openapi-spec> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `-o, --output` | Output directory (default: .) |
| `-p, --package` | Go package name (default: ax) |
| `--only` | Generate specific files (errors,retry,capabilities) |
| `--dry-run` | Preview what would be generated |

**Example:**

```bash
ax-spec gen api-ax.yaml -o pkg/ax --package ax
```

**Generated Files:**

| File | Description |
|------|-------------|
| `errors.go` | Error code constants and metadata |
| `retry.go` | Retry policy mappings |
| `capabilities.go` | Operation capability mappings |
| `validation.go` | Required field validators |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `API_TOKEN` | Default API token for `--discover` |
