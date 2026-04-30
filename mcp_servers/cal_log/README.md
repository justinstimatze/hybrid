# cal_log — calibration logger MCP server

Per-evaluator prediction logging with rolling-window hit-rate aggregation. Append-only event log; minimal dependencies; runs as a stdio MCP server.

This is the **shippable primitive for Conjecture 1** of the hybrid-loops repo: a standalone calibration logger that any hybrid-loop project can drop in to track whether its typed LLM evaluators are actually working.

## Why

Every hybrid-loop project has at least one typed LLM evaluator (a "lens" producing structured records). Without a calibration log, there is no way to tell whether the evaluator's predictions track reality over time. The architecture is theater unless you can show hit-rate per evaluator.

This server makes the discipline cheap: log a prediction in one tool call, mark a verdict later in another, ask for hit-rate when you want it. JSONL on disk; no database; no setup beyond installing the MCP server.

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

This means: no destructive writes, easy backup, trivially auditable, and concurrent-safe at the OS append granularity.

## Running standalone (without the plugin)

```bash
# with uv (recommended)
uv run --script server.py

# or install via pip and run as module
pip install -e .
python -m cal_log.server
```

## Running as part of the hybrid-loops plugin

If you've installed the `hybrid-loops` Claude Code plugin, this MCP server is auto-registered as `cal_log`. Tools are callable as `mcp__cal_log__{predict,resolve,hit_rate,list_pending,list_recent,stats}`.

## Conventions for `loop` and `verdict`

- **`loop`** — a stable identifier for the project + evaluator combination. Examples: `"intervention-tracker-fit-scorer"`, `"voice-corpus-deviation-detector"`, `"recruiter-fit-scorer"`. Pick a name and stick with it; renames cost you the time-series.
- **`verdict`** — the resolution category. The default `confirmed` / `refuted` is the simplest. Some projects benefit from richer verdicts (`partial` for half-right, `irrelevant` for "the prediction was about a thing we didn't end up evaluating"). Document your project's verdict vocabulary somewhere stable.

## Testing

```bash
uv run --script tests/test_server.py
```

Or with pip-installed deps, `pytest tests/`.

## Status

v0.1. Single-process; stdio only; no auth (it's a personal-use logger). The schema is stable; I'll bump the major version if it changes incompatibly.
