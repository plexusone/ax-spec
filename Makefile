.PHONY: build test lint docs docs-serve marp clean

# Go build
build:
	go build -o ax-spec ./cmd/ax-spec

test:
	go test -v ./...

lint:
	golangci-lint run

# Documentation
docs: marp
	mkdocs build

docs-serve: marp
	mkdocs serve

# Marp presentations - convert to HTML
marp:
	@command -v marp >/dev/null 2>&1 || { echo "marp not installed. Run: npm install -g @marp-team/marp-cli"; exit 1; }
	marp docs/case-studies/elevenlabs-go/presentation.md -o docs/case-studies/elevenlabs-go/presentation.html
	marp docs/case-studies/opik-go/presentation.md -o docs/case-studies/opik-go/presentation.html

# Clean generated files
clean:
	rm -f ax-spec
	rm -f docs/case-studies/*/presentation.html
	rm -rf site/

# Install dependencies
deps:
	go mod download
	@echo "Optional: npm install -g @marp-team/marp-cli"
	@echo "Optional: pip install mkdocs mkdocs-material"

# Development helpers
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build the ax-spec CLI"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  docs       - Build MkDocs site (includes marp)"
	@echo "  docs-serve - Serve docs locally with live reload"
	@echo "  marp       - Generate HTML from Marp presentations"
	@echo "  clean      - Remove generated files"
	@echo "  deps       - Install Go dependencies"
