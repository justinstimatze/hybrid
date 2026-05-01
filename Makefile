SERVERS := cal_log metacog schemaforge

.PHONY: check fmt vet lint test test-race ci-local clean help

# Default target.
help:
	@echo "Targets:"
	@echo "  make check       — run the same gates CI runs (gofmt, vet, lint, test -race) on all servers"
	@echo "  make fmt         — auto-fix formatting (gofmt -w)"
	@echo "  make vet         — go vet across all servers"
	@echo "  make lint        — golangci-lint across all servers (uses .golangci.yml)"
	@echo "  make test        — go test (no -race) across all servers, faster for iteration"
	@echo "  make test-race   — go test -race across all servers (matches CI)"
	@echo "  make ci-local    — alias for 'check', mirrors the CI workflow"
	@echo
	@echo "Lint requires golangci-lint installed locally. If absent, install with:"
	@echo "  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"

# 'check' is the gate to run before pushing. Matches CI exactly.
check: gofmt-check vet lint test-race

ci-local: check

gofmt-check:
	@echo "==> gofmt -l mcp_servers/"
	@unformatted=$$(gofmt -l mcp_servers/); \
	if [ -n "$$unformatted" ]; then \
		echo "Files need gofmt:"; \
		echo "$$unformatted"; \
		echo; \
		echo "Run 'make fmt' to auto-fix."; \
		exit 1; \
	fi

fmt:
	@echo "==> gofmt -w mcp_servers/"
	@gofmt -w mcp_servers/

vet:
	@for s in $(SERVERS); do \
		echo "==> go vet ./mcp_servers/$$s/..."; \
		(cd mcp_servers/$$s && go vet ./...) || exit 1; \
	done

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not installed. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"; \
		exit 1; \
	}
	@for s in $(SERVERS); do \
		echo "==> golangci-lint run ./mcp_servers/$$s/..."; \
		(cd mcp_servers/$$s && golangci-lint run --config $(CURDIR)/.golangci.yml ./...) || exit 1; \
	done

test:
	@for s in $(SERVERS); do \
		echo "==> go test ./mcp_servers/$$s/..."; \
		(cd mcp_servers/$$s && go test ./...) || exit 1; \
	done

test-race:
	@for s in $(SERVERS); do \
		echo "==> go test -race ./mcp_servers/$$s/..."; \
		(cd mcp_servers/$$s && go test -race ./...) || exit 1; \
	done

clean:
	@for s in $(SERVERS); do \
		(cd mcp_servers/$$s && go clean -testcache); \
	done
