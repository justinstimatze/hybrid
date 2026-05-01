# hybrid

**A cycle, not a pipeline.** The *hybrid-loops* design pattern places LLM judgment and deterministic code in **alternating layers that mutually generate each other's working surface** — not just constraining each other, but *producing* the very inputs the other half operates over. The LLM generates typed records (and often the schema, notation, or code those records live in). The deterministic layer takes those records and generates filtered, scored, ranked context that becomes the next LLM call's input. Each half makes the other possible.

The cycle: an LLM extracts → typed structure accumulates → deterministic gates filter, score, rank → another LLM reasons over the gated substrate → an action lands (often as new content the lens reads next turn) → calibration closes the loop on whether the evaluators actually worked → the substrate metabolizes as a whole.

Most projects have **0-3 specific places** that warrant the full pattern. The skill in this repo helps find them and decline where it doesn't fit; the three MCP servers ship the primitives that recur once you do.

## What's in it for you

If you're building anything where an LLM reads non-deterministic content (transcripts, dialogues, plans, documents, screenshots, behavior logs) and produces decisions or recommendations another part of the system acts on, you probably have one of those places. **Pure LLM** does fluent extraction but drifts and can't restrain itself. **Pure code** is unyielding but can't read the soft input. **Pure pipeline** (LLM-then-code-then-done) misses the load-bearing trick: the LLM and the code keep generating new working surfaces for each other. This repo helps you design what goes at each surface:

- **Typed schemas** that capture LLM judgments as structured records
- **Deterministic gates** that handle restraint, scoring, ranking — what LLMs are bad at
- **Calibration logs** that track whether your evaluators actually work (the `cal_log` MCP server, included)
- **Cognitive-bias self-audit** on substrate structure (the `metacog` MCP server, included)
- **Schema discovery** for finding dense notations over a corpus via compress+verify (the `schemaforge` MCP server, included)
- **Anti-pattern detection** so you don't over-apply the architecture where it doesn't fit

The skill is **diagnostic-first** — most of any project isn't this pattern, and identifying which part is (and which isn't) is half the value.

## Who it's for

- Solo developers and small teams building tools that involve LLM judgment
- AI engineers shipping production LLM features who need typed observability
- Domain experts (teachers, advocates, writers, coaches, anyone who makes repeated judgments) who want personal typed-judgment tools rather than chatbots

## Install — Claude Code

**Prerequisites**:
- Go 1.21+ on `$PATH` (the MCP servers compile and run via `go run`).
- `ANTHROPIC_API_KEY` in your environment if you want to use `schemaforge` tools (the other two servers don't need it).

```bash
# Add this repo as a marketplace, then install the plugin
/plugin marketplace add justinstimatze/hybrid
/plugin install hybrid-loops@hybrid-loops
```

Plugin installs the skill (`hybrid-loops`) plus three MCP servers (auto-registered): `cal_log`, `metacog`, and `schemaforge`. Skill auto-triggers on relevant prompts; tools are callable as `mcp__cal_log__*`, `mcp__metacog__*`, `mcp__schemaforge__*`.

The marketplace command requires this GitHub repo to be reachable. While it's still under early review the repo may be private — in that case clone the repo and use the local-path forms: `/plugin marketplace add /path/to/hybrid` and `/plugin install hybrid-loops@hybrid-loops`.

If you prefer just the skill without the plugin scaffolding: symlink `skills/hybrid-loops/` into `~/.claude/skills/hybrid-loops/`.

## Install — other agents

The skill content and MCP server are model-agnostic. Stub manifests are included for OpenAI Codex, Cursor, and Gemini — see `CROSS_AGENT.md`. The maintainer's primary platform is Claude Code; PRs from users on other agents are welcome.

## What's in the repo

```
hybrid/
├── .claude-plugin/         Claude Code plugin + marketplace manifests
├── skills/
│   └── hybrid-loops/
│       ├── SKILL.md           the skill (one-screen TL;DR + 5-phase diagnostic)
│       └── references/        loaded on demand
├── mcp_servers/
│   ├── cal_log/            calibration logger (Go, 8 tests passing)
│   ├── metacog/            cognitive-bias auditor (Go, 20 tests passing)
│   └── schemaforge/        schema-discovery (compress+verify) loop (Go, 18 tests passing, pilot run)
├── .codex-plugin/          Codex stub
├── .cursor-plugin/         Cursor stub
├── gemini-extension.json   Gemini stub
├── CROSS_AGENT.md          portability notes
└── README.md               (this file)
```

## What's in the primitives

Three stdio MCP servers, one per shippable conjecture-falsifier:

- **`cal_log`** — calibration logger. Append-only JSONL at `$CAL_LOG_PATH` (default `~/.cal_log/calibration.jsonl`). Tools: `predict`, `resolve`, `hit_rate`, `list_pending`, `list_recent`, `stats`. 8/8 unit tests. See `mcp_servers/cal_log/README.md`.
- **`metacog`** — cognitive-bias auditor. Nine bias-signature checks against any typed substrate (provenance HHI as availability-heuristic proxy, predicate entropy as base-rate-neglect proxy, etc.). Tools: `audit`, `audit_one`, `list_auditors`. 20/20 unit tests. See `mcp_servers/metacog/README.md`.
- **`schemaforge`** — schema-discovery (compress+verify) loop. Designs dense notation for a corpus, translates items into it, expands back, scores roundtrip, evolves across rounds. Tools: `design_notation`, `compress`, `expand`, `evaluate_roundtrip`, `score_round`, `run_round`. 18/18 unit tests, pilot run on a 10-item non-program corpus (3 rounds, results in `mcp_servers/schemaforge/README.md`). Requires `ANTHROPIC_API_KEY`.

## Status

**Research output, not a product.** Documents a design pattern, ships a working skill + three working MCP servers. Nothing is sold and nothing is trying to be a SaaS. The pattern itself doesn't sell; specific tools built with it might.

The four claims about what's plausibly new are framed as **conjectures with named falsifying experiments** — see "Conjectures" below. Each of `cal_log`, `metacog`, `schemaforge` is the runnable primitive for one conjecture; running them at scale on real projects is the experiment. As of April 2026, only Conjecture 3 has been partially exercised — see `mcp_servers/schemaforge/README.md` for the pilot results.

## Naming

"Hybrid loops" is the **working name in this repository**, not a claim of universal nomenclature. The broader field has no settled name. Adjacent terms with partial coverage:

- **"Compound AI systems"** (Zaharia et al., BAIR 2024) — broader umbrella; this pattern is one shape inside it
- **"Generalization shaping"** (Manuel Odendahl / wesen, 2026) — the design principle inside hybrid loops; closest practitioner framing
- **"Structured introspection"** — informal practitioner usage; partial overlap

The pattern can be cited by any of these names.

A separate term used here: **"third mind"** — a *deployment context* where the substrate is shared between collaborators (and possibly an AI), distinct from a personal external substrate or one's own thinking. Burroughs/Gysin's 1978 sense (the emergent entity in collaborative writing) extended to substrates that themselves metabolize. **Third mind is a deployment shape; hybrid loops is the architectural pattern.**

## Conjectures

Four conjectures about what this work might contribute beyond the cited prior art. **Each is testable; only Conjecture 3 has been partially exercised so far.**

### Conjecture 1 — per-evaluator calibration is undershipped

*Claim.* A standalone primitive that logs predictions and verdicts per typed LLM evaluator, with rolling-window hit-rate aggregation, would generalize across hybrid-loop projects and meaningfully change development decisions.

*Falsifying experiment.* Use the `cal_log` MCP server (in this repo) on 3+ existing projects of varied shape; measure over 60 days whether the hit-rate signal changes any concrete development decision (prompt change, schema bump, gate adjustment). If hit-rate is collected but no decisions are made on it, the primitive is theater.

### Conjecture 2 — cognitive-bias self-audit on substrate structure generalizes

*Claim.* Cognitive-bias signature checks (provenance HHI as availability-heuristic proxy, irrelevant:challenged ratio as survivorship-bias proxy, predicate entropy as base-rate-neglect proxy, etc.) work on any typed substrate, not only the substrate they were prototyped on.

*Falsifying experiment.* Run the `metacog` MCP server (in this repo) on 3+ independent typed substrates of varied shape; have project owners mark which triggered findings correspond to actual structural problems vs. false positives. If false-positive rate is high or findings don't track owner intuition, the metrics are substrate-specific rather than substrate-general. Untested as of April 2026.

### Conjecture 3 — schema discovery extends to non-program domains

*Claim.* Compress+verify schema-discovery loops (DreamCoder/LILO descendants) can discover useful schemas for non-program domains — humor structures, dramatic arcs, behavioral mechanisms, AI-conversation patterns — not only for code or notation for code.

*Falsifying experiment.* Run the `schemaforge` MCP server (in this repo) on a non-program corpus over multiple rounds; check that mean expansion ratio and roundtrip correctness behave like they do on program corpora — ER climbs while correctness stays above ~0.7. **Partial result (April 2026):** a 10-item naturalistic-prose corpus (memes serialized to prose with 11 latent fields) produced ER 15.9–45.6x at correctness 0.71–0.89 in round 1, depending on which notation philosophy the model picked. Over 3 rounds with a conservative evolve prompt, correctness drifted up modestly while notation stayed stable. Positive on this single corpus; replication on 2+ more independent non-program corpora is needed before declaring the conjecture supported. See `mcp_servers/schemaforge/README.md` for full pilot writeup, including the round-2 regression caused by the original evolve prompt.

### Conjecture 4 — there is unmet demand for domain-applied substrate-as-vocabulary tooling outside engineering

*Claim.* Users in non-engineering domains (coaching, teaching, parenting, advocacy, creative work) would benefit from typed-repertoire-with-restraint tools and don't currently have them.

*Falsifying experiment.* Build one such tool (e.g. the teacher's intervention tracker or the coach's typed library from `skills/hybrid-loops/references/EXAMPLES.md`); ship to 5+ domain users; measure 30+ day retention. If retention is below baseline rates for similar consumer tools, the demand isn't there or the tool is wrong-shaped.

---

These are the open ground after acknowledging the cited prior art. **The next material work is running these experiments, not making more architectural claims.** All three runnable primitives are shipped; the C-3 pilot is the first datapoint.

## Acknowledgments

This writeup is meaningfully shaped by **Manuel Odendahl** ("wesen"), whose work in this design space at [github.com/go-go-golems](https://github.com/go-go-golems) and writing at [the.scapegoat.dev](https://the.scapegoat.dev) directly informed the pattern as documented here. The "generalization shaping" framing, the deliberate use of "diary" over "log," the term "substrate" for typed event-streaming layers, and the Blackboard-Systems architectural reading are all his. Any public presentation of hybrid loops should credit his contributions; a fuller account is in `skills/hybrid-loops/references/PRIOR_ART.md`.

Thanks also to the published work of [DreamCoder](https://github.com/ellisk42/ec) (Ellis et al., 2021), [LILO](https://github.com/gabegrand/lilo) (Grand et al., 2024), [Voyager](https://github.com/MineDojo/Voyager) (Wang et al., 2023), and the [Polis](https://pol.is/) and [Talk to the City](https://github.com/AIObjectives/talktothe.city) projects, all referenced throughout the skill. Devine Lu Linvega ([100r.co](https://100r.co)) and the Hundred Rabbits collective inform the small-tools aesthetic that the deterministic-shell half of the pattern aspires to. Christopher Alexander's *A Pattern Language* (1977) is the structural reference for what the pattern *is* as a unit of design.

## License

MIT.
