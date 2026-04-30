# hybrid

A design pattern + Claude Code plugin for the **specific places in any project** where an LLM doing fuzzy semantic judgment needs to feed structured decisions downstream. Ships a diagnostic skill, a working calibration MCP server, and stub manifests for other coding agents.

## What's in it for you

If you're building anything where an LLM reads non-deterministic content (transcripts, dialogues, plans, documents, screenshots, behavior logs) and produces decisions or recommendations another part of the system acts on, you probably have **0-3 specific places** where the architecture matters. This repo helps you find those places and design what goes there:

- **Typed schemas** that capture LLM judgments as structured records
- **Deterministic gates** that handle restraint, scoring, ranking — what LLMs are bad at
- **Calibration logs** that track whether your evaluators actually work (the `cal_log` MCP server, included)
- **Anti-pattern detection** so you don't over-apply the architecture where it doesn't fit

The skill is **diagnostic-first** — most of any project isn't this pattern, and identifying which part is (and which isn't) is half the value.

## Who it's for

- Solo developers and small teams building tools that involve LLM judgment
- AI engineers shipping production LLM features who need typed observability
- Domain experts (teachers, advocates, writers, coaches, anyone who makes repeated judgments) who want personal typed-judgment tools rather than chatbots

## Install — Claude Code

```bash
# Add this repo as a marketplace, then install the plugin
/plugin marketplace add justinstimatze/hybrid
/plugin install hybrid-loops@hybrid-loops
```

Plugin installs the skill (`hybrid-loops`) plus the `cal_log` MCP server (auto-registered). Skill auto-triggers on relevant prompts; cal_log tools are callable as `mcp__cal_log__{predict,resolve,hit_rate,list_pending,list_recent,stats}`.

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
│   └── cal_log/            calibration logger MCP server (working, tested)
├── .codex-plugin/          Codex stub
├── .cursor-plugin/         Cursor stub
├── gemini-extension.json   Gemini stub
├── CROSS_AGENT.md          portability notes
└── README.md               (this file)
```

## What `cal_log` does

`cal_log` is the calibration primitive every hybrid-loop project should have but most don't ship. It's a stdio MCP server with five tools:

- `predict(loop, input_hash, prediction, model_id, ...)` → log a typed evaluator's prediction
- `resolve(prediction_id, verdict, verdict_source)` → mark the verdict later
- `hit_rate(loop, window_days=30)` → aggregated hit-rate over the window
- `list_pending(loop?, limit)` → unresolved predictions ordered by due date
- `list_recent(loop?, limit)` / `stats()` → introspection

Append-only JSONL events at `$CAL_LOG_PATH` (default: `~/.cal_log/calibration.jsonl`). 8/8 unit tests passing. See `mcp_servers/cal_log/README.md`.

## Status

**Research output, not a product.** Documents a design pattern, ships a working skill + a working MCP server. Nothing is sold and nothing is trying to be a SaaS. The pattern itself doesn't sell; specific tools built with it might.

The four claims about what's plausibly new are framed as **conjectures with named falsifying experiments** — see "Conjectures" below. `cal_log` is the runnable primitive that makes Conjecture 1 testable; running it on real projects for 30+ days is the experiment.

## Naming

"Hybrid loops" is the **working name in this repository**, not a claim of universal nomenclature. The broader field has no settled name. Adjacent terms with partial coverage:

- **"Compound AI systems"** (Zaharia et al., BAIR 2024) — broader umbrella; this pattern is one shape inside it
- **"Generalization shaping"** (Manuel Odendahl / wesen, 2026) — the design principle inside hybrid loops; closest practitioner framing
- **"Structured introspection"** — informal practitioner usage; partial overlap

The pattern can be cited by any of these names.

A separate term used here: **"third mind"** — a *deployment context* where the substrate is shared between collaborators (and possibly an AI), distinct from a personal external substrate or one's own thinking. Burroughs/Gysin's 1978 sense (the emergent entity in collaborative writing) extended to substrates that themselves metabolize. **Third mind is a deployment shape; hybrid loops is the architectural pattern.**

## Conjectures

Four conjectures about what this work might contribute beyond the cited prior art. **Each is testable; none has been tested.**

### Conjecture 1 — per-evaluator calibration is undershipped

*Claim.* A standalone primitive that logs predictions and verdicts per typed LLM evaluator, with rolling-window hit-rate aggregation, would generalize across hybrid-loop projects and meaningfully change development decisions.

*Falsifying experiment.* Use the `cal_log` MCP server (in this repo) on 3+ existing projects of varied shape; measure over 60 days whether the hit-rate signal changes any concrete development decision (prompt change, schema bump, gate adjustment). If hit-rate is collected but no decisions are made on it, the primitive is theater.

### Conjecture 2 — cognitive-bias self-audit on substrate structure generalizes

*Claim.* Cognitive-bias signature checks (provenance HHI as availability-heuristic proxy, irrelevant:challenged ratio as survivorship-bias proxy, predicate entropy as base-rate-neglect proxy, etc.) work on any typed substrate, not only the substrate they were prototyped on.

*Falsifying experiment.* Lift the audit primitive (`metacog`, sketched in `skills/hybrid-loops/references/PRIMITIVES.md`) into a standalone library; run it on 3+ independent typed substrates; have project owners mark which triggered findings correspond to actual structural problems vs. false positives. If false-positive rate is high or findings don't track owner intuition, the metrics are substrate-specific rather than substrate-general.

### Conjecture 3 — schema discovery extends to non-program domains

*Claim.* Compress+verify schema-discovery loops (DreamCoder/LILO descendants) can discover useful schemas for non-program domains — humor structures, dramatic arcs, behavioral mechanisms, AI-conversation patterns — not only for code or notation for code.

*Falsifying experiment.* Run a compress+verify loop on a corpus of 100+ examples in one non-program domain; compare the discovered schema against a hand-authored one on downstream extraction quality, with a domain expert as the evaluator. If hand-authored beats discovered, the loop doesn't generalize past program-shaped domains.

### Conjecture 4 — there is unmet demand for domain-applied substrate-as-vocabulary tooling outside engineering

*Claim.* Users in non-engineering domains (coaching, teaching, parenting, advocacy, creative work) would benefit from typed-repertoire-with-restraint tools and don't currently have them.

*Falsifying experiment.* Build one such tool (e.g. the teacher's intervention tracker or the coach's typed library from `skills/hybrid-loops/references/EXAMPLES.md`); ship to 5+ domain users; measure 30+ day retention. If retention is below baseline rates for similar consumer tools, the demand isn't there or the tool is wrong-shaped.

---

These are the open ground after acknowledging the cited prior art. **The next material work is running these experiments, not making more architectural claims.** `cal_log` is the first runnable piece.

## Acknowledgments

This writeup is meaningfully shaped by **Manuel Odendahl** ("wesen"), whose work in this design space at [github.com/go-go-golems](https://github.com/go-go-golems) and writing at [the.scapegoat.dev](https://the.scapegoat.dev) directly informed the pattern as documented here. The "generalization shaping" framing, the deliberate use of "diary" over "log," the term "substrate" for typed event-streaming layers, and the Blackboard-Systems architectural reading are all his. Any public presentation of hybrid loops should credit his contributions; a fuller account is in `skills/hybrid-loops/references/PRIOR_ART.md`.

Thanks also to the published work of [DreamCoder](https://github.com/ellisk42/ec) (Ellis et al., 2021), [LILO](https://github.com/gabegrand/lilo) (Grand et al., 2024), [Voyager](https://github.com/MineDojo/Voyager) (Wang et al., 2023), and the [Polis](https://pol.is/) and [Talk to the City](https://github.com/AIObjectives/talktothe.city) projects, all referenced throughout the skill. Devine Lu Linvega ([100r.co](https://100r.co)) and the Hundred Rabbits collective inform the small-tools aesthetic that the deterministic-shell half of the pattern aspires to. Christopher Alexander's *A Pattern Language* (1977) is the structural reference for what the pattern *is* as a unit of design.

## License

MIT.
