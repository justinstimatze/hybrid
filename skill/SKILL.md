---
name: hybrid-loops
description: Plan and scaffold hybrid-loop architectures — projects where an LLM extracts typed structure from non-deterministic content (transcripts, dialogues, plans, docs, screenshots), deterministic code operates on that typed surface (gates, scores, aggregates, persists), and another LLM call reasons over it to drive decisions or actions. Triggers on prompts like "build a tool that analyzes/evaluates/extracts X over time", "I want to make sense of Y across many Zs", "track patterns in W", "score/rank/compare conversations or documents", or any project where the value is data ABOUT content, not natural-language output. Do NOT trigger for one-shot transforms, pure UIs, code refactors, or chatbots.
---

# Hybrid Loops

A design pattern for projects where an LLM's fuzziness and code's hardness mutually bootstrap something neither does alone. The user has built ~15 projects in this shape (slimemold, drivermap, score, plancheck, lucida, winze, gemot, etc.) and wants Claude to reach for this architecture by default when a new project fits.

## When this fits — telltale signals

Strong signals (any two → invoke this skill):
- Project consumes non-deterministic content (text, dialogue, image, plan)
- Output value is *data ABOUT content*, not a natural-language reply
- Aggregation, comparison, ranking, or learning across many inputs is wanted
- Downstream reasoning needs typed handles into the artifact
- Persistence across sessions matters
- The user mentions: "evaluate", "score", "extract", "track patterns", "audit", "consensus"

Weak signals (suggestive but not load-bearing):
- User mentions schemas, taxonomies, types, structured outputs
- User wants a tool that "gets better over time"
- User wants to analyze AI-generated content (transcripts, traces, agent logs)

## When this does NOT fit — anti-patterns

Do NOT apply this pattern when:
- The project is a one-shot transform (translate, summarize, format)
- The value is the LLM's natural-language output (chatbot, advisor)
- The project is pure code refactoring or pure UI work
- There is no non-deterministic input — only structured data
- The output is consumed once and discarded
- Calibration / aggregation across runs is not wanted

If a project fits two anti-patterns, stop. Tell the user the pattern doesn't fit and explain why. Do not retrofit a hybrid loop where one isn't useful — that's exactly the over-application failure mode this skill is meant to prevent.

## Activation surfaces

A hybrid loop needs a trigger — something that fires it. The trigger shapes architectural choices as much as the schema does. Common patterns across the user's existing stack:

**Claude Code hooks (UserPromptSubmit, Stop, PostToolUse, etc.)** — the loop fires automatically based on Claude Code's lifecycle. Slimemold uses Stop (every 3-5 turns to inject claim-graph findings). Crowdwork uses UserPromptSubmit (before Claude responds, to inject humor signal). This is the *invisible* activation pattern — the user never invokes the loop, it's in the air around the conversation.

**MCP tools** — the loop is callable on demand by Claude or other agents. Winze's `mcp__winze__{claims,disputes,provenance,search,stats,theories}`, drivermap's `predict_mechanisms`, gemot's deliberation tools. Use this when the substrate should be queryable from outside, especially from other agents.

**CLI / cron** — the loop runs as a scheduled job. Winze's metabolism. Use for genuinely background processes that don't need to be reactive.

**Stream watchers** — the loop polls a transcript file or event stream. Lucida watches the Claude Code transcript and mints visualization cells in real time. Use only when polling is unavoidable; they're noisy and brittle compared to hooks.

**Hook-triggered MCP** — a hook fires, calls an MCP server, the result lands in Claude's context. Hook = when, MCP = what. The combination is powerful for ambient-but-rich injections (slimemold's findings).

Activation choice depends on:
- *Who initiates* — user, Claude, external agent, clock?
- *Latency* — synchronous (hooks, MCP) or asynchronous batch (cron)?
- *Visibility* — should the loop fire visibly (action layer) or invisibly (system-message injection)?
- *Cost* — hooks fire on every relevant lifecycle event; budget accordingly. Cache where possible.

Defaults: MCP for query-shaped loops; hooks for ambient-injection loops; CLI/cron only for genuine background processes; avoid stream watchers unless polling is unavoidable.

If the project will live as a Claude Code plugin (bundling hooks + MCP + skills), the form factor is plugin not skill. Start with skill; promote to plugin when the architecture needs to ship hooks alongside the cognitive primitives.

## The pattern in plain language

A hybrid loop has five roles. Most projects use 4-5; some use only 3.

1. **LENS** — an LLM call that extracts/generates typed records from non-deterministic content. Schema is fixed; LLM fills it with tool_use or JSON mode. The schema includes a `notes` or `rationale` field for graceful failure.
2. **SUBSTRATE** — the typed surface produced by the lens. Stored as JSON, JSONL, sqlite, or a typed Go AST (winze-style). Has schema versioning. Records carry provenance (which model, which prompt hash).
3. **GATE** — deterministic code that operates on the substrate: filtering, ranking, scoring, cooldowns, ripeness windows, audit triggers. This is where opinionated policy lives. Not generic plumbing.
4. **REASONER** — an LLM call that consumes substrate records (often via MCP tools) and produces decisions, recommendations, or generated content informed by the typed structure.
5. **ACTION** — deterministic effect: write to the substrate, send a notification, trigger a hook, render UI, modify code. Sometimes loops back to step 1 by producing new content.

Two meta-layers turn a one-shot loop into a hybrid system:

6. **CALIBRATION** — every prediction the lens or reasoner makes is logged with a verdict to be resolved later. JSONL append-only. Hit-rate over time = whether the architecture earns its keep.
7. **METABOLISM** — periodic phases that operate on the substrate as a whole: dream (consolidation), trip (speculative connections), bias-audit (run cognitive-bias detectors against KB structure), evolve (sensor queries to grow KB). Each phase logs its own predictions for calibration.

## The design interview

When invoking this skill on a new project, work through these questions with the user. Present them as a checklist; let the user answer in any order. Do not proceed to scaffolding until at least the first five are clear.

1. **What is the non-deterministic input?** (transcripts? dialogue turns? screenshots? plans? a corpus of documents? agent traces?)
2. **What does the lens extract?** Sketch the schema. What fields, what types, what enums? Include a `notes` field. Include a `model_id` and `schema_version` field on every record.
3. **Where does the substrate live?** JSON file per record? Single JSONL? SQLite? A typed Go module winze-style? MCP server exposing query tools?
4. **What does the gate do?** Filter what? Rank by what? Cooldown how? Score window how? Be specific — generic gates ("filter low-quality") are usually wrong; opinionated gates ("only fire when ripeness ∈ [2, 14] turns") are usually right.
5. **What does the reasoner do with the substrate?** What's the actual decision being made or output being produced?
6. **What's the action?** Where does the value land? UI render? File write? Notification? New record back into substrate?
7. **What does the calibration log capture?** Every reasoner prediction with: input, prediction, model_id, timestamp, verdict_due_by, verdict (filled in later). What is the verdict signal — a follow-up user action? A timer? A second LLM check?
8. **Does the project need metabolism?** Most v1 projects don't. If the substrate is small or per-session, skip metabolism. Add it only when the substrate accumulates over weeks and would benefit from periodic audit/consolidation.
9. **Which existing primitives can it import?** (See references/PRIMITIVES.md — most are not yet packaged. If a primitive doesn't exist, scaffold the local version with a TODO to extract it later.)
10. **What's the ablation experiment?** Pick one component (lens, gate, reasoner) and describe the experiment that would show the project is worse without it. The hybrid claim is empty without this.

## Scaffolding template

For a Python project (most common shape):

```
project-name/
  README.md                    # describe lens/substrate/gate/reasoner/action in 5 lines each
  schema/
    v1.json                    # the lens's output schema, versioned
  src/
    lens.py                    # the LLM extraction call
    substrate.py               # storage adapter (JSONL or sqlite)
    gate.py                    # deterministic filtering/scoring
    reasoner.py                # the downstream LLM call
    action.py                  # the effect
    cal_log.py                 # calibration logger (~50 lines)
  data/
    substrate.jsonl            # the typed surface
    calibration.jsonl          # predictions and verdicts
  tests/
    test_lens.py               # mock LLM, fixed input → expected schema shape
    test_gate.py               # deterministic, table-driven
    test_ablation.py           # the load-bearing claim of the project
  CLAUDE.md                    # project-specific notes for future Claude sessions
```

For a Go project (winze-shaped), the substrate may be a typed AST module rather than JSONL.

For an MCP-first project (gemot-shaped), wrap the reasoner and substrate as MCP tools and skip the action layer — the consumer is another agent.

## Calibration log format

Append-only JSONL. One record per prediction:

```json
{"ts": "2026-04-30T10:00:00Z", "loop": "slimemold", "lens_or_reasoner": "reasoner",
 "input_hash": "...", "prediction": {...}, "model_id": "claude-sonnet-4-6",
 "schema_version": 3, "verdict_due_by": "2026-05-07T10:00:00Z",
 "verdict": null, "verdict_source": null, "verdict_ts": null}
```

Verdicts get filled in by a separate `cal_log.resolve()` call when the verdict signal arrives. The resolver computes hit-rate over rolling windows.

Without a calibration log, the project's hybrid claim is unfalsifiable. Add it from day one even if verdict signals don't exist yet — the log itself surfaces what verdict signals are reachable.

## Ablation discipline

Every hybrid-loop project should answer: *if I removed the typed substrate and just gave the same LLM raw content, would performance drop?* Encode this as a test:

```python
def test_ablation_substrate_helps():
    # Run reasoner with substrate
    typed_score = run_with_substrate(test_input)
    # Run reasoner with raw transcript only
    raw_score = run_without_substrate(test_input)
    assert typed_score > raw_score, "Substrate is not earning its keep"
```

This is the only honest answer to the "you're just using LLMs with extra steps" objection. Without it, the architecture is theater.

## Recursive composition — stacking hybrid loops

Hybrid loops compose. The output of one loop's substrate can be the non-deterministic input of another loop's lens. The decision a reasoner makes can become a typed claim that another loop's gate inspects. This stacking is the trajectory the user has been pushing toward.

A multi-layer recursion example: **publicrecord**.

- Layer 0: human-authored finding (verbatim quote from primary source) — irreducible ground truth
- Layer 1: LLM (lens) extracts typed records (entity, relationship, severity) from primary sources
- Layer 2: deterministic code (gate) builds the SQLite DB, validates schemas, enforces provenance
- Layer 3: another LLM (reasoner) consumes typed records via MCP and makes a recommendation
- Layer 4: deterministic code (action) wraps the recommendation in a structured response with citations
- Layer 5: a downstream agent's loop (with its own lens/substrate/gate) consumes that response

Each layer is itself a hybrid loop. Each output is typed enough to be the next layer's input. The user's framing: *"agents building non-agents that build agents that build non-agents."*

The trajectory: **back and forth forever until what's getting generated and run is on top of so many layers of meta guardrails it basically always does the right thing.**

This is a real architectural claim about reliability. As typed constraint layers stack between raw LLM output and final action, the action is increasingly determined by typed structure rather than free generation. Real concerns to keep honest:

- **Cost and latency compound.** N hybrid loops ≈ N LLM calls per action.
- **Errors at lower layers propagate.** Wrong record at layer 1 corrupts everything above. Calibration logs at *every* layer is the only mitigation.
- **Saturation is empirical.** Reliability probably increases with N up to a point, then diminishes. The N depends on schema-composition quality, per-layer calibration, task latency tolerance.
- **Composition requires discipline.** Canonical schemas at common interfaces (claim, entity, finding); MCP-first when output is meant for downstream loops; provenance on every record; explicit schema versioning. See references/STACKING.md.

When designing a new project, ask: *will this loop's output ever be another loop's input?* If yes, design the schema and the activation surface for composability from the start. If no, ignore this section — most v0 projects don't need recursive composition.

## Naming conventions

Use the user's vocabulary, not academic terms:

- **lens** (not "extractor" or "perceiver")
- **substrate** (not "knowledge graph" or "store")
- **gate** (not "filter" or "policy layer")
- **reasoner** (not "agent" or "synthesizer")
- **calibration log** (not "telemetry" or "observability")
- **metabolism** (not "background job" or "cron")
- **hybrid loop** (not "compound AI system" or "schemaed cognition")

## Citations to use when defending the architecture

If the user or a reviewer asks "haven't expert systems already done this?" or "is this novel?", cite:

- **AlphaGo / AlphaZero** (Silver et al., 2016, 2017) — the architectural template for fuzzy+hard mutual constraint. Policy net (lens) + MCTS (gate) bootstrapping each other.
- **DreamCoder** (Ellis et al., Nature 2021, arXiv:2006.08381) — wake/sleep library learning, the closest direct lineage for schema discovery and metabolism.
- **LILO** (Grand et al., NeurIPS 2024, arXiv:2310.19791) — LLM-era DreamCoder descendant.
- **Voyager** (Wang et al., 2023, arXiv:2305.16291) — skill library learning for agents, hybrid loop in agent context.
- **OpenCog / Hyperon** (Goertzel, ~2008–) — *cite to distinguish*: same architectural intuition (typed substrate that thinks), different bet (OpenCog tried symbolic reasoning where LLMs now win). Hybrid loops are what OpenCog should have been.

What's actually new in this lineage: hybrid loops bootstrap over a *structural prior the LLM itself generates*. AlphaGo's structure (Go's rules) is fixed; DreamCoder's syntax is fixed; here the substrate's schema is generated and refined by an earlier loop. That's the load-bearing novelty claim — keep it small and defensible.

## When in doubt

If you're not sure whether a project fits this pattern, ask the user one question: **"Is the value of this project the data it produces about content, or the natural-language output the LLM gives you?"**

Data about content → hybrid loop fits.
Natural-language output → it doesn't. Use a chat skill or just write the prompt.

## See also

- references/PRIOR_ART.md — full citation list with paper details and how each compares
- references/PRIMITIVES.md — list of primitives to import (mostly not yet packaged; scaffold local versions for now)
- references/EXAMPLES.md — worked examples mapping the user's existing repos to the five-role schema
