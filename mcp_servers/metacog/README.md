# metacog — cognitive-bias auditor MCP server

Runs nine cognitive-bias signature checks against any typed substrate. Lifted from the auditor pattern in [winze](https://github.com/justinstimatze/winze)'s `cmd/metabolism/dreamaudit.go` (1,778 lines of Go-AST traversal abstracted away behind the `Substrate` interface).

This is the **shippable primitive for Conjecture 2** of the hybrid-loops repo: the claim that cognitive-bias self-audit on substrate structure generalizes beyond a single substrate. Until this is run on at least three independent substrates, that's a conjecture, not a confirmed contribution.

Written in Go using [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) — same SDK as cal_log and the rest of the Go MCP servers in this author's stack.

## What it audits

Each auditor checks one cognitive-bias signature against the *structure* of the substrate (not its content). Metrics are general; thresholds were calibrated against winze's KB and may need re-tuning on substrates of very different shape or scale.

| auditor | metric | threshold | needs |
|---|---|---|---|
| confirmation_bias | corroboration_rate | > 0.75 triggers | Verdict field |
| anchoring | spearman_age_density | > 0.5 triggers | CreatedAt or Group |
| clustering_illusion | group_cluster_jaccard | > 0.7 triggers | Group + Cluster hints |
| availability_heuristic | provenance_hhi | > 0.25 triggers | Provenance.Origin or Type |
| survivorship_bias | irrelevant_to_challenged_ratio | > 5.0 triggers | verdicts incl. 'irrelevant' |
| framing_effect | evaluative_summary_fraction | > 0.15 triggers | SummaryText |
| dunning_kruger | low_complexity_zero_edge_rate | > 0.90 triggers | Complexity hint + Edges |
| base_rate_neglect | predicate_entropy_bits | < 3.0 triggers | Edges.Predicate |
| premature_closure | closure_findings | >= 1 triggers | SummaryText |

Auditors return `Skipped: true` instead of garbage when their required substrate field is absent.

## Substrate format

Currently one format: `jsonl` — one JSON record per line. Schema:

```json
{
  "id": "claim_abc",
  "type": "claim",
  "fields": {"any": "structured data"},
  "provenance": [{"origin": "https://en.wikipedia.org/wiki/X", "type": "wikipedia", "quote": "..."}],
  "edges": [{"predicate": "supports", "to": "claim_xyz"}],
  "created_at": "2026-04-15T12:00:00Z",
  "verdict": "confirmed",
  "verdict_time": "2026-05-01T09:00:00Z",
  "summary_text": "Free-text summary used by framing/closure auditors",
  "complexity": 0.3,
  "cluster": "c1",
  "group": "file_a"
}
```

All fields except some core (a record without any of `verdict`, `provenance`, `edges`, `summary_text` is auditable but every auditor will skip it).

## Verdict vocabulary

The `verdict` field is normalized to one of: `corroborated`, `challenged`, `irrelevant`, `partial`, `unknown`. Acceptable inputs:

- *corroborated*: `confirmed`, `corroborated`, `supported`, `validated`, `verified`
- *challenged*: `refuted`, `challenged`, `contradicted`, `rejected`
- *irrelevant*: `irrelevant`, `noise`, `off_topic`
- *partial*: `partial`, `mixed`

Anything else (or empty) becomes `unknown` and doesn't count toward verdict-needing auditors.

## Tools (MCP)

| Tool | Purpose |
|---|---|
| `audit(substrate_path, format)` | Run all nine auditors. Returns per-auditor results + which triggered + which skipped. |
| `audit_one(substrate_path, auditor, format)` | Run a single named auditor. |
| `list_auditors()` | Enumerate the nine auditors and what each measures. |

Default `format` is `jsonl`.

## Running standalone

```bash
go test ./...                       # all 20 unit tests
go run .                            # stdio MCP server
go build -o metacog .               # compiled binary
```

## Running as part of the hybrid-loops plugin

When the `hybrid-loops` Claude Code plugin is installed, this MCP server is auto-registered as `metacog`. Tools become callable as `mcp__metacog__{audit,audit_one,list_auditors}`.

## Conjecture-2 testing — falsifying experiment

The Claim is that the nine bias-detection metrics work on *any* typed substrate. The Falsifying experiment is:

1. Get 3+ independent typed substrates of varied shape (a claim graph; a behavioral mechanism corpus; a calibration log; etc.)
2. Convert each to the JSONL schema above
3. Run `audit(...)` on each
4. Have project owners mark which triggered findings correspond to actual structural problems vs. false positives
5. If false-positive rate is high or findings don't track owner intuition, the metrics are substrate-specific rather than substrate-general — Conjecture 2 is refuted

## Future work

- **Winze parity reader.** Direct AST + `.metabolism-log.json` reader that produces Records from winze's native format without JSONL conversion. Half a day of work; would let us run a head-to-head parity check against winze's own audit numbers (e.g. `AvailabilityHeuristic` triggered at HHI=0.479, `SurvivorshipBias` triggered at ratio 95.5 in winze's current state).
- **Threshold re-calibration.** Current thresholds are winze-derived. Substrates an order of magnitude smaller or larger may need different cutoffs.
- **Per-substrate baselines.** Conjecture-2 testing should compare audit results against random-shuffle baselines per substrate.

## Module path

`github.com/justinstimatze/hybrid/mcp_servers/metacog`
