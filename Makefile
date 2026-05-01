SERVERS := cal_log metacog schemaforge

.PHONY: check fmt vet test test-race ci-local clean help

# Default target.
help:
	@echo "Targets:"
	@echo "  make check       — run the same gates CI runs (gofmt -l, vet, test -race) on all servers"
	@echo "  make fmt         — auto-fix formatting (gofmt -w)"
	@echo "  make vet         — go vet across all servers"
	@echo "  make test        — go test (no -race) across all servers, faster for iteration"
	@echo "  make test-race   — go test -race across all servers (matches CI)"
	@echo "  make ci-local    — alias for 'check', mirrors the CI workflow"
	@echo
	@echo "CI also runs against Go 1.21, 1.22, 1.23 — to replicate that matrix locally"
	@echo "you need a multi-version setup (gvm, asdf, or containers). Single-version 'check'"
	@echo "catches the vast majority of issues."

# 'check' is the gate to run before pushing. Matches CI exactly except for the
# Go version matrix (which CI handles).
check: gofmt-check vet test-race

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
