# Contributing

This is research-stage software. The pattern itself is the artifact; the primitives are reference implementations. Contributions are welcome but the bar for what changes is higher than the bar for what gets discussed.

## Easy first contributions

- **Run a falsifying experiment.** Each of the four conjectures in the README has a named experiment. If you run one of them — `cal_log` deployed in a real project for 30+ days, `metacog` audit against a substrate that isn't the one it was prototyped on, `schemaforge` round-loop on a non-program corpus you care about — open an issue or PR with the result. Negative results count.
- **Implement a sketched primitive.** `metabolism` (phase scheduler) and `mcp_substrate` (MCP server template generator) are sketched in `skills/hybrid-loops/references/PRIMITIVES.md` but not yet shipped. Working implementations are welcome.
- **Cross-agent ports.** The skill is portable; the cross-agent stub manifests under `.codex-plugin/`, `.cursor-plugin/`, and `gemini-extension.json` need adjustment per each agent's actual current spec. PRs from users on those platforms are the right path.

## Working on the code

A `Makefile` at the repo root mirrors what CI does. Run before pushing:

```bash
make check       # gofmt + vet + golangci-lint + go test -race across all servers
make fmt         # auto-fix formatting (gofmt -w)
make lint        # golangci-lint only
make test        # tests without -race, faster for iteration
```

`make check` matches the CI workflow exactly. Both pin to the Go version declared in each module's `go.mod` (currently 1.25.5).

`make lint` requires `golangci-lint` installed locally:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

Lint configuration lives at `.golangci.yml` in the repo root. CI also runs **CodeQL** (security analysis, weekly cron) and **Dependabot** (weekly updates for Go modules and GitHub Actions).

The schemaforge integration test (`pilot_test.go`) hits the Anthropic API and is build-tagged behind `-tags pilot`. It's skipped by default. To run it:

```bash
cd mcp_servers/schemaforge
SCHEMAFORGE_CORPUS=/path/to/corpus.jsonl \
SCHEMAFORGE_OUT=/tmp/forge-run \
SCHEMAFORGE_TARGET="what the notation should expand to" \
go test -tags pilot -timeout 60m -v -run TestPilot .
```

## What changes the bar is high for

- **Expanding the conjecture list.** Four is intentional. New conjectures need a named falsifier and a runnable primitive to check them with. Speculation without a falsifier is not in scope.
- **Renaming the pattern or its roles.** "Hybrid loops" / "lens / substrate / gate / reasoner / action" are the working internal vocabulary. They're not perfect; changing them is disruptive enough to need real argument.
- **Major skill structural changes.** The 5-phase diagnostic (find → scope → shape → quick design → scaffold) has been validated by trigger tests. Restructuring needs an equally rigorous validation.

## Commit style

Semantic commits, one logical change per commit. Rough conventions used so far:
- `Ship X: Y` for new servers/skills
- `Fix X` for bug fixes  
- `Update X` / `Tighten X` / `Clean up X` for content edits

Co-authored-by trailers welcome when the work was meaningfully assisted by an AI tool.

## Code of conduct

See `CODE_OF_CONDUCT.md`. Short version: be kind, assume good faith, don't drag personal disputes into Issues.

## License

By contributing, you agree your contributions are licensed under the MIT License (see `LICENSE`).
