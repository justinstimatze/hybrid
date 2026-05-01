# cal_log — calibration logger MCP server

Per-evaluator prediction logging with rolling-window hit-rate aggregation. Append-only event log; minimal dependencies; runs as a stdio MCP server.

This is the **shippable primitive for Conjecture 1** of the hybrid-loops repo: a standalone calibration logger that any hybrid-loop project can drop in to track whether its typed LLM evaluators are actually working over time.

Written in Go using [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go), the community Go MCP SDK widely used in the practitioner ecosystem.

## Why

Every hybrid-loop project has at least one typed LLM evaluator (a "lens" producing structured records, or a "reasoner" consuming them). Without a calibration log, there is no way to tell whether the evaluator's predictions track reality over time. The architecture is theater unless you can show hit-rate per evaluator.

This server makes the discipline cheap: log a prediction in one tool call, mark a verdict later in another, ask for hit-rate when you want it. JSONL on disk; no database; no setup beyond having Go installed.

## Tools

| Tool | Purpose |
|---|---|
| `predict(loop, input_hash, prediction, model_id, ...)` | Log a typed evaluator's prediction. Returns `prediction_id`. |
| `resolve(prediction_id, verdict, verdict_source)` | Mark the verdict (`confirmed` / `refuted` / `partial` / `irrelevant` or whatever you define). |
| `hit_rate(loop, window_days=30)` | Aggregated hit-rate over the past N days for one loop. |
| `list_pending(loop?, limit=50)` | Predictions awaiting verdicts, oldest-due first. |
| `list_recent(loop?, limit=50)` | Recent predictions, most recent first. |
| `stats()` | Top-level summary across all loops. |

## Storage

Append-only JSONL events at `$CAL_LOG_PATH` (default: `~/.cal_log/calibration.jsonl`). Two event types: `predict` and `resolve`. The current state of a prediction is the fold of its events.

Means: no destructive writes, easy backup, trivially auditable, concurrent-safe at OS append granularity.

## Running standalone

```bash
# from this directory
go run .

# or build a binary once and run from anywhere
go build -o cal_log .
./cal_log
```

The server speaks stdio JSON-RPC; pair with any MCP client.

## Running as part of the hybrid-loops plugin

When the `hybrid-loops` Claude Code plugin is installed, this MCP server is auto-registered as `cal_log`. The plugin invokes `go run ${CLAUDE_PLUGIN_ROOT}/mcp_servers/cal_log` on stdio. Tools become callable as `mcp__cal_log__{predict,resolve,hit_rate,list_pending,list_recent,stats}`.

For lower-latency invocation: `go build -o cal_log .` once, then change the plugin manifest's `mcpServers.cal_log` entry to point at the built binary.

## Conventions for `loop` and `verdict`

- **`loop`** — a stable identifier for the project + evaluator combination. Examples: `"intervention-tracker-fit-scorer"`, `"voice-corpus-deviation-detector"`, `"recruiter-fit-scorer"`. Pick a name and stick with it; renames cost you the time-series.
- **`verdict`** — the resolution category. The default vocabulary `confirmed` / `refuted` is the simplest. Some projects benefit from richer verdicts (`partial` for half-right, `irrelevant` for "the prediction was about a thing we didn't end up evaluating"). Document your project's verdict vocabulary somewhere stable.

## Testing

```bash
go test -v ./...
```

Seven tests covering the storage layer (predict, resolve, fold, orphan handling, malformed-line handling, last-write-wins). The MCP tool wiring is shallow and exercised by mark3labs/mcp-go's own test suite.

## Module path

`github.com/justinstimatze/hybrid/mcp_servers/cal_log`

If you want to use the storage primitives in your own Go code without the MCP wrapper:

```go
import calstore "github.com/justinstimatze/hybrid/mcp_servers/cal_log"
// ... but it's a `package main` right now; lift to a `pkg/` if reuse becomes a thing.
```

(Current shape is single-file `main` package; if reusable Go-library use materializes, lift `Store`, `Record`, etc. into a sub-package.)

## Status

v0.1. Single-process; stdio only; no auth (it's a personal-use logger). The schema is stable; major version bump if it changes incompatibly.
