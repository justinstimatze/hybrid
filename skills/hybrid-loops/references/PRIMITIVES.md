# Primitives for hybrid loops

Extractable primitives that recur across hybrid-loop projects. Three are now shipped as working MCP servers in this repo (`mcp_servers/`); the other two remain sketches. When scaffolding a new project: reach for the shipped servers where they apply, write a local version of the still-sketched primitives where they don't.

## cal_log — calibration logger (shipped)

Append-only JSONL of predictions and verdicts. The simplest and highest-leverage primitive — every typed evaluator should have one of these from day one.

Implementation: see `mcp_servers/cal_log/`. Stdio MCP server with tools `predict`, `resolve`, `hit_rate`, `list_pending`, `list_recent`, `stats`. Default storage at `$CAL_LOG_PATH` (default `~/.cal_log/calibration.jsonl`).

When in a Claude Code session with the plugin installed, call as `mcp__cal_log__predict(...)` and `mcp__cal_log__resolve(...)` from the lens or reasoner code. When building outside Claude Code, append-events directly to the JSONL using the documented schema is sufficient.

Per Conjecture 1 in the README, this primitive's value claim — that it changes development decisions when run at scale — is conjectured but untested.

## metacog — cognitive-bias auditor (shipped)

Run cognitive-bias signature checks against any typed substrate. The metrics check the *structure* of the substrate, not its content.

Auditors:
- **confirmation_bias** — corroboration rate among resolved cycles
- **anchoring** — Spearman(record_age, density) correlation
- **clustering_illusion** — group vs cluster Jaccard
- **availability_heuristic** — provenance HHI (concentration index over source domains)
- **survivorship_bias** — irrelevant:challenged signal ratio
- **framing_effect** — evaluative-adjective regex frequency in summary fields
- **dunning_kruger** — low-complexity zero-edge rate
- **base_rate_neglect** — predicate distribution entropy
- **premature_closure** — thought-terminating cliché detection

Implementation: see `mcp_servers/metacog/`. Substrate format is JSONL with documented schema; auditors gracefully skip when their required fields are absent. Per Conjecture 2 in the README, the metrics' substrate-generality is untested across 3+ independent substrates.

## schemaforge — DreamCoder-style schema discovery (shipped)

Given a corpus, discover a dense notation that compresses items while preserving roundtrip fidelity. Compress + verify loop over multiple rounds with conservative evolution.

Tools: `design_notation`, `compress`, `expand`, `evaluate_roundtrip`, `score_round`, `run_round` (convenience). Fine-grained per-LLM-call tools mean a Claude session can drive the loop directly and react to per-round metrics; `run_round` is the one-call-per-round convenience. Cite DreamCoder (Ellis et al., 2021) and LILO (Grand et al., 2024) as direct lineage.

Implementation: see `mcp_servers/schemaforge/`. Per Conjecture 3 in the README, generalization to non-program corpora is partially supported by a 10-item pilot; replication on 2+ more corpora is needed.

## metabolism — phase scheduler (sketch only)

Periodic execution of dream / trip / audit / evolve phases with phase-gate dependencies on a typed substrate.

API sketch:
```python
metabolism = Metabolism(
    substrate=...,
    phases=[
        Phase("dream", run_dream, gate=lambda ctx: ctx.last_dream_age > 1d),
        Phase("trip",  run_trip,  gate=lambda ctx: ctx.bias_audit.passed),
        Phase("audit", run_audit, gate=lambda _: True),
    ]
)
metabolism.run_cycle()
```

Generalize so phases are pluggable. Phase outputs themselves should log via `cal_log` so phase predictions can be calibrated alongside lens/reasoner predictions.

Not yet shipped as a standalone server. Most projects past v1 want this; if you build the third one by hand, package it.

## mcp_substrate — MCP server template (sketch only)

Wrap a typed substrate as an MCP server with read tools (search, get, stats, query). One template generator that takes a substrate schema + sample queries and emits a working MCP server. Most hybrid-loop projects that grow past v0 want this; once a developer has built two of them by hand, packaging is worth it.

Not yet shipped as a standalone server.

## When to package vs. inline

- v0 prototype: inline a simplified version of `cal_log` directly. Skip the others.
- v1 with > 1 month of expected life: install the shipped servers and use them; they are bounded and don't pull in heavy framework code.
- Once the same primitive is written twice across projects, lift it into a real package.

The bottleneck for most builders isn't "missing packages." It's "Claude doesn't reach for the pattern by default." This skill addresses that. Packaging is downstream.
