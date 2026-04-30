# schemaforge — schema-discovery MCP server

Compress + verify loop. Designs a dense notation for a corpus, translates each item into the notation, expands the notation back to the target form, scores the roundtrip, and evolves the notation across rounds based on per-round metrics.

This is the **shippable primitive for Conjecture 3** of the hybrid-loops repo: the claim that the dense-notation discovery loop generalizes to non-program corpora. Until this is run on at least one program corpus and at least two non-program corpora with comparable correctness curves, that's a conjecture, not a confirmed contribution.

The descent is from DreamCoder (Ellis et al. 2021) and LILO (Grand et al. 2024) — wake/sleep loops that grow a library of reusable abstractions by compressing programs and verifying behavior. schemaforge keeps the compress+verify shape and drops the program-specific assumptions: input is a JSONL corpus of `{id, spec_text}`; the target language and roundtrip rubric are free-text parameters.

Written in Go using [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) and [`anthropic-sdk-go`](https://github.com/anthropics/anthropic-sdk-go).

## What it does

| phase | LLM op | what flows |
|---|---|---|
| design | `design_notation` | seed items + target → notation spec |
| compress | `compress` | spec_text + notation → per-item notation text |
| expand | `expand` | per-item notation → expanded text |
| evaluate | `evaluate_roundtrip` | original_spec + expanded + rubric → score 0-1 |
| aggregate | `score_round` (pure) | per-item JSON files → mean ER, mean correctness, summary |
| evolve | `design_notation` (round N+1) | previous notation + previous metrics → improved notation |

The two metrics that matter:

- **Expansion ratio (ER)** = `expanded_tokens / notation_tokens`. How much each notation token implies. Higher is denser.
- **Roundtrip correctness** = mean of per-item LLM-rated 0-1 scores against the original spec. Has to stay above ~0.7 for the notation to be useful.

A successful run shows ER climbing across rounds while correctness holds. A run where correctness collapses as ER climbs has an over-aggressive notation; a run where ER plateaus immediately has a notation that's already saturated.

## Substrate format

Corpus is JSONL — one JSON record per line:

```json
{"id": "blog_app", "spec_text": "freeform spec — could be a CRUD description, a behavioral mechanism, a humor template, a knowledge claim, anything"}
```

`id` must be unique (used as the per-item filename in the round directory). `spec_text` is opaque to schemaforge; it's whatever the caller wants to find a dense schema for.

## Tools (MCP)

| Tool | What it does | Cost |
|---|---|---|
| `design_notation(corpus_path, target, model, previous_notation?, previous_metrics_summary?)` | Designs (round 1) or evolves (round N+1) a notation. Reads up to 2 longest items as design seeds. | 1 LLM call |
| `compress(spec_text, notation_spec, target, model)` | Translates one item into the notation. | 1 LLM call |
| `expand(notation_spec, item_notation, target, model)` | Expands one notation back to the target form. | 1 LLM call |
| `evaluate_roundtrip(original_spec, expanded, rubric, model)` | Scores 0-1; returns score + one-sentence reasoning. | 1 LLM call |
| `score_round(round_dir)` | Pure compute. Aggregates per-item JSON files into RoundMetrics. | local IO |
| `run_round(corpus_path, target, rubric, model, output_dir, round_number, ...)` | Convenience: full round end-to-end. Writes `round{N}/notation.txt`, `round{N}/items/{id}.json`, `round{N}/metrics.json`. Resumable. | 1 + 3N LLM calls for N items |

Default model: `claude-sonnet-4-6`. The system prompt of every per-item call includes the notation spec — prompt caching means the bulk of N+N round-trips is cached.

## Why fine-grained tools, not one monolithic discover

The whole point of hybrid-loops is letting an LLM orchestrate deterministic primitives. A single multi-hour `discover()` tool would invert that. Per-step tools mean the caller (a Claude Code session, a script, anything) can read the notation after design, decide whether to continue, redesign the rubric mid-run, kill a bad round early. Each tool call is bounded — minutes, not hours.

`run_round` exists for the "just do a full round" case but emits per-item files as it goes, so you can ctrl-C and re-invoke to resume.

## Storage layout

```
<output_dir>/
  round1/
    notation.txt           # the notation specification (text)
    items/
      <id>.json            # ItemResult per item
    metrics.json           # RoundMetrics aggregate
  round2/
    notation.txt
    items/...
    metrics.json
  ...
```

No global state, no goroutines, no lifecycle. Resume = caller checks what's there, reads `metrics.json` for the summary, passes it back into `run_round` as `previous_notation` + `previous_metrics_path` for the next round.

## Requirements

- `ANTHROPIC_API_KEY` in env (every tool except `score_round` makes API calls).
- Go 1.25+ for standalone build.

## Running standalone

```bash
go test ./...                       # all 18 unit tests
go run .                            # stdio MCP server
go build -o schemaforge .           # compiled binary
```

## Running as part of the hybrid-loops plugin

When the `hybrid-loops` Claude Code plugin is installed, this server is auto-registered as `schemaforge`. Tools become callable as `mcp__schemaforge__{design_notation,compress,expand,evaluate_roundtrip,score_round,run_round}`.

## Conjecture-3 testing — falsifying experiment

The Claim: the dense-notation compress+verify loop generalizes to non-program corpora. The falsifying experiment:

1. Pick three corpora of different shapes — one programmatic (e.g. CRUD app specs), two non-programmatic (e.g. behavioral mechanism descriptions; structured humor templates; knowledge claims with provenance).
2. Convert each to the JSONL format.
3. For each corpus, run 3+ rounds of `run_round`. Pick a target and rubric appropriate to that domain.
4. Plot ER and correctness across rounds.
5. **Predicted shape if Conjecture 3 holds:** ER climbs across rounds while correctness stays above ~0.7. The trajectories should be qualitatively similar across all three corpora.
6. **Predicted shape if Conjecture 3 is refuted:** correctness collapses on the non-programmatic corpora as ER climbs (the notation becomes lossy on natural language in a way it doesn't on programs); or ER plateaus immediately because a useful dense notation requires programmatic structure that natural language doesn't have.

Either outcome is informative. A Conjecture-3-positive result means the primitive is domain-general. A Conjecture-3-negative result narrows what compress+verify is good for, which is also useful.

### Pilot results (April 2026, naturalistic-prose meme corpus, 10 items × 3 rounds)

A 10-item smoke against a non-program corpus (memes serialized to naturalistic prose, 11 latent fields) showed:

- **Round 1 alone produces useful dense notation.** Two independent runs hit ER 45.6x at 0.89 correctness and ER 15.9x at 0.71 correctness; both clearly compressing, just at different densities depending on whether the model picked a sigil-based or pipe-delimited grammar by chance.
- **The evolve prompt is load-bearing.** The first version of the evolve prompt explicitly invited improvement (`Create compositional primitives for repeated patterns`) and produced a 25% spec-size growth and 26% ER drop in round 2 with no correctness gain. The current version defaults to "output unchanged" unless concrete failures (items < 0.70) justify a targeted fix; round 2 in the same configuration produced a 1-character grammar edit, round 3 was byte-identical, while correctness drifted up modestly across rounds via compress/expand variance against a stable target.
- **Limitation surfaced.** The conservative evolve prompt cannot escape a poor round-1 design — if round 1 picks a verbose notation philosophy, subsequent rounds hold the line rather than redesign. A future enhancement would trigger redesign-from-scratch when correctness < some-threshold for 2 consecutive rounds.

## Cost notes

Each round on N items costs 1 + 3N LLM calls (1 design + N compress + N expand + N evaluate). Prompt caching on the notation spec brings the marginal cost of compress+expand calls down within a round. For corpora >20 items a single `run_round` call may still take tens of minutes — caller's responsibility to budget accordingly.

## Future work

- **Streaming progress.** Currently `run_round` returns only on completion. A future variant could stream per-item events through the MCP `progress` channel so callers can react without polling the filesystem.
- **Baseline measurement.** Optional "generate target without notation, count tokens" pass to compute compression ratio against an unaided baseline. Lamina's reference implementation has this; cut from the primitive to keep scope tight.
- **Corpus chunking.** Very large corpora (100s of items) should run rounds in batches to keep individual `run_round` calls bounded.
- **Fine-tuning extraction.** `ExtractTrainingPairs` produces (notation → expanded) pairs above a score threshold; not yet exposed as an MCP tool but trivially could be.
- **Redesign-from-scratch path.** Currently the evolve prompt only modifies an existing notation. A `redesign_notation` op that re-runs the round-1 design call (ignoring `previous_notation`) would let the loop escape a poor initial design when correctness stays below threshold across rounds.

## Module path

`github.com/justinstimatze/hybrid/mcp_servers/schemaforge`
