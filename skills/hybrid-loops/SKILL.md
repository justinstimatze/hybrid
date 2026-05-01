---
name: hybrid-loops
description: Find and design hybrid-loop surfaces in any project — places where an LLM's fuzzy semantic judgment is consequential enough to warrant typed structure (schema'd substrate, deterministic gating, calibration). Triggers on prompts like "build a tool that tracks/analyzes/evaluates/extracts X over time", "make sense of Y across many Zs", "detect patterns in W", "notice when X is happening", "score/rank/compare documents", "monitor my agent for X", "watch this stream and alert when...", "build a regression detector for...", "log and analyze evaluator outputs over time", "suggests when I'm repeating an anti-pattern", "flag for me when...", "build a system that learns from outcomes", or any project where the value is DATA ABOUT content rather than the LLM's reply. Diagnostic-first — most projects have only 0-3 such surfaces. Domain-agnostic. Do NOT trigger for one-shot transforms (translate, summarize, format), pure UIs, CRUD without judgment, code refactors, chatbots, agent harness setup, or generic prompt engineering.
---

# Hybrid Loops

## TL;DR (one screen)

> A *cycle* of alternating LLM-and-code layers that *mutually generate each other's working surface* — not just constraining each other, but producing the very inputs the other half operates over. *LLMs bring fluency. Substrates bring discrimination. Code brings restraint.* The LLM writes typed records (often the schema or notation itself); the deterministic layer aggregates and shapes those records into the input the next LLM call sees. They don't just gate each other — they manufacture each other.
>
> The point isn't LLM-as-pipeline-stage. It's *LLM-as-half-of-a-loop* — and at scale, *layered loops that wrap around each other.* Runtime: one cycle resolving one judgment. Development-time: a critique-patch loop wraps around the runtime, with an LLM-panel reading transcripts of runtime behavior and patching the deterministic layer (or the lens prompts, or the substrate schema) below. The system grows by stacking such loops.
>
> 5-phase diagnostic:
> 1. Find candidate **surfaces** in the project (places where fuzzy judgment is happening or should be)
> 2. **Scope** each: A (just call an LLM), B (don't use an LLM), or C (hybrid loop)
> 3. Choose **shape**: substrate-as-record (analytical) or substrate-as-vocabulary (interventional)
> 4. **Quick design** in 3 questions (input, schema, action) with sane defaults for the rest
> 5. **Scaffold** to the surface, not the project. Always include a calibration log and an ablation test.
>
> Five roles in a cycle: **lens** (LLM extracts) → **substrate** (typed records accumulate) → **gate** (deterministic policy filters/scores/ranks) → **reasoner** (LLM consumes substrate) → **action** (deterministic effect; often loops back as new content). Plus two meta-layers: **calibration** (predict + verdict log — does the lens actually work?) and **metabolism** (substrate-wide audit — is the accumulated record drifting?).
>
> Decline when: one-shot transform, chatbot, pure UI, no fuzziness in input, output discarded once, or a deployment shape that imposes a substrate on workers who can't edit it.

## Why this skill exists

Most working programmers carry three mental primitives:

- *Code operates on data* (the classical view; what school teaches)
- *An LLM operates on data* (recent; chatbots, evaluators, copilots)
- *An LLM writes code* (newer still; codegen, autonomous agents)

These are three cells in a much larger combinatorial space. Every mix of `{LLM | code}` as actor × `{data | code}` as input × `{data | code}` as output is a valid block, and most useful systems built today are *graphs that span many cells* rather than pipelines that occupy one. Almost nobody was trained for the combinatorial space. Schooling and working experience produce strong intuitions for the three classical cases and almost none for the multi-block dynamic-graph cases.

The LLM that's actually building the system is in roughly the same position as the programmer asking for it. Without explicit guidance, the LLM also defaults to pipeline thinking — extract once, decide once, return. This skill exists to push back on that default: to put the broader space in front of the LLM as a working option, and to scaffold the multi-layered dynamic graph the project actually wants. *What blocks does this surface need? What should they generate for each other? Where does the cycle close? What wraps around it?*

Systems that come out of this kind of design tend to feel a little organic — they grow rather than getting authored top-down, they adapt as they run and surprise you, they have metabolic phases (audit, prune, evolve) that aren't part of any single decision but keep the substrate fit over time. That's not poetic. It comes from the cycles being mutually generative: each layer keeps remaking the surface the others act on, and the system as a whole behaves more like an ecology than an engineered artifact.

## Full skill

A design pattern for the *specific places* in a project where a fuzzy semantic judgment benefits from typed structure. Most projects have 0–3 such places. The skill helps Claude identify them, decide whether each warrants the full pattern, and design what's there. Domain-agnostic — applies in health, education, ops, creative, social, business, engineering.

In one sentence: an LLM does fuzzy judgment, a typed substrate captures the result as data, deterministic code does aggregation/restraint/scoring, another LLM reasons over the substrate, an action lands. *LLMs bring fluency. Substrates bring discrimination. Code brings restraint.*

*Note on naming.* "Hybrid loops" is the working name in this repo; the broader field has no settled name (see `references/PRIOR_ART.md` for adjacent terms — "compound AI systems," "generalization shaping," etc.).

## Phase 1 — find the surfaces

Walk the project. Name each candidate surface in one sentence. Signs:

- A judgment a human keeps making on similar inputs
- Pattern recognition over content (more than keyword matching)
- Generative choices needing taste
- Aggregation over qualitative observations
- Triage / severity calls
- Anything phrased *"a tool that helps me notice when..."*

Zero candidates → skill probably doesn't apply. Say so plainly.

## Phase 2 — scope each surface

Three buckets:

- **A: Just call an LLM directly.** One-shot, no persistence, no aggregation. Don't add ceremony.
- **B: Don't use an LLM at all.** Looks fuzzy, actually deterministic (regex, vector similarity, SQL).
- **C: Hybrid loop.** Fuzzy AND consequential AND benefits from typed structure.

If all surfaces are A or B, this isn't a hybrid-loops project. Exit.

## Phase 3 — choose shape (for C surfaces)

- **Substrate-as-record (analytical):** typed log of past observations. For making sense of accumulated data over time.
- **Substrate-as-vocabulary (interventional):** typed repertoire (curated roster, closed taxonomy). For discriminating the right move now.

This is a diagnostic heuristic, not a hard binary. Most real surfaces end up as both — a typed library that *also* accumulates per-instance records — and the distinction is mainly useful for picking which side to design first. Lead with whichever the user-facing value is closer to (analytical → record-first; interventional → vocabulary-first), then add the other side once the first is working.

## Phase 4 — quick design (3 questions)

Three minimum questions. Everything else gets sane defaults; refine after a draft scaffold is on the table.

1. *What's the non-deterministic input?*
2. *What does the lens extract — sketch the schema* (3-7 fields including `notes`, `model_id`, `schema_version`).
3. *What's the action — where does value land?*

Defaults if not specified:
- *Substrate*: JSONL file in project; sqlite when growing past ~1000 records.
- *Gate*: confidence threshold + chronological ordering. Add restraint policies (cooldown, ripeness window) only when over-firing is observed.
- *Reasoner*: read recent records via simple query; produce structured output for the action.
- *Calibration log*: append-only JSONL, predict + verdict. Add from day one even when verdict signals don't yet exist — the log surfaces which verdicts are reachable.
- *Metabolism*: skip in v0. Triggers for v1+: substrate above ~1000 records, schema-version churn (≥2 bumps), or the calibration log showing a recurring class of failure that points at a substrate-shape issue rather than a per-block one. Until one of those fires, metabolism is over-engineering.

Produce a draft scaffold from the three answers. Iterate from there.

For the full design interview (when the surface is large enough that getting it wrong has material cost), see `references/DESIGN_INTERVIEW.md`.

## Phase 5 — scaffold to the surface, not the project

The scaffolding is bounded to the surface; the rest of the project is whatever else it is. Minimum:

- Versioned schema definition
- Lens code (typed LLM call)
- Substrate code (storage + read API)
- Gate code (deterministic policy)
- Reasoner code (LLM consuming substrate)
- Action wiring
- Calibration log

Implementation language follows the surrounding project. Deployment shape options: *embedded in existing codebase, standalone module, MCP server, notebook/script, no-code wiring* (Airtable + Zapier + an LLM call). Pick by user fit, not by familiarity.

## Five roles, reference

The roles *alternate between fluency (LLM) and discrimination (code)*, with each role generating the input the next role consumes. The arrows close a cycle:

1. **LENS** *(LLM generates → typed records)* — produces typed records from soft input. Fixed schema, `notes` field for graceful failure. Often the LLM also generates the schema itself (one tier up in development time).
2. **SUBSTRATE** *(typed records accumulate)* — the accumulating record. Carries `model_id` + `schema_version`. The substrate is what makes the loop *learn* — each turn's records become the constraint surface for the next.
3. **GATE** *(code generates → filtered context)* — deterministic policy: filtering, scoring, cooldowns, ranking. Code restrains the LLM here so the LLM doesn't have to restrain itself; equally important, the gate manufactures the next LLM call's input by deciding which records get through.
4. **REASONER** *(LLM generates → decisions / new content)* — consumes the gate's output, produces decisions or generated content.
5. **ACTION** *(code generates → state change, often new content)* — deterministic effect. Often produces new content the lens reads next turn. This is what closes the cycle.

Plus two meta-layers that close *different* loops:

- **CALIBRATION** — predict + verdict log per evaluator. Closes the loop on *whether the lens/reasoner is actually working* (rolling hit-rate). Without it, the architecture is theater.
- **METABOLISM** — periodic substrate-wide phases (audit, prune, refactor). Closes the loop on *substrate quality over time*. Skip until v1+.

The lens may be staged or parallel — treat lens as a *role*, not a single LLM call. Same for the reasoner.

**Stacked loops.** Many real systems have multiple cycles wrapping around each other. A common shape: a runtime cycle (engine + player + lens), wrapped by a development-time cycle (LLM-critic reads runtime transcripts, generates a patch plan, the patch plan modifies the lens prompts / substrate schema / gate policy / engine code, and the next runtime turn picks up the change). The development-time loop is itself a hybrid loop. See `references/STACKING.md`.

## Building blocks

The five roles are an *opinionated default arrangement* of more general primitives. The eight primitive blocks come from `{LLM | code}` × `{reads data | reads code}` × `{produces data | produces code}`; pairs and triples compose into recognizable patterns. The full algebra and worked examples live in `references/BUILDING_BLOCKS.md`. The catalog of named graphs (RAG, ReAct, codegen-with-verification, the canonical 5-role loop, dev-time critique loop, knowledge-base auditor, teacher's intervention tracker, etc.) lives in `references/BLOCK_GRAPHS.md`.

The diagnostic in this skill defaults to the five-role shape because it covers most analytical and interventional cases. When your case wants something else — `LLM-as-architect`, `code-as-perceiver`, `LLM-audits-code`, `LLM-generates-prompts-for-LLM` — name the blocks you need and snap their I/O together. The cycle is the structural invariant; *which blocks fill it is project-specific*.

### What this looks like in practice

The dominant workflow for hybrid-loop projects isn't authoring a formal architecture document upfront and then implementing it. It's iterative:

1. *Brainstorm the graph with an LLM.* "I have this kind of input, this kind of decision needs to come out, what's between?" The LLM proposes blocks (some LLM, some code, some pure data shapes), draws connections, asks about restraint policies, suggests where calibration belongs. The graph emerges turn-by-turn.
2. *Glue the blocks together.* Sometimes statically (a fixed pipeline of MCP tools). Often *dynamically* — routing logic decides at runtime which subgraph fires for which input.
3. *Subagents are nested hybrid loops.* A reasoner block sometimes decomposes into a smaller hybrid loop running inside it (its own lens / substrate / reasoner). Claude Code's subagent spawn is one common implementation.
4. *The graph-as-data stays alive.* Edited as the system runs and surprises you. New block when an audit shows a recurring failure. Block removed when calibration shows it never earned its keep.

When this skill fires, it's because the user is at step 1. The role of the skill is to help that brainstorming converge on a shape that's likely to work — naming the recurring blocks, flagging anti-patterns, suggesting where to put the calibration log so step 2 onward goes well.

## Activation surface

How does the loop fire? Pick one:

- **In-process call** — the surrounding project invokes the lens/reasoner directly. Default for embedded surfaces.
- **MCP tool** — callable on demand by Claude or other agents. Use when the substrate should be queryable from outside.
- **Hooks** (Claude Code lifecycle, request handlers, browser events) — fire automatically. Use for ambient injection.
- **CLI / cron** — scheduled batch. Use only for genuine background processes.
- **Stream watcher** — polls a transcript or stream. Avoid unless polling is unavoidable.

For non-Claude-Code contexts (API apps, web services, mobile, no-code), ask: *where in the host project's normal control flow does this surface get called?* That's the activation point.

## Deployment power-balance check

Hybrid loops separate *judgment substrate* (the typed library / record) from *judgment execution* (the LLM that picks or produces from it). That separation has political consequences depending on who owns each. This section addresses *one* axis — substrate-owner-vs-executor power balance — and is named accordingly. Real deployment ethics in any domain (privacy, accountability, accessibility, recourse, ongoing monitoring) is broader than what's here and lives outside the framework's scope; treat this section as a single decision-point at scaffolding time, not as full ethics coverage.

Ask, before scaffolding:

- *Who owns the substrate?* Users, the team, a platform, a service?
- *Who executes the move?* Same party as substrate owner, or different?
- *Can the executor edit the substrate?* Add entries, refine criteria, override?

Power-neutral deployments: single user owns and executes; team shares ownership and execution; platform owns substrate but executor can edit.

Concerning: platform owns substrate, gig-worker or low-status executor must apply the imposed taxonomy without editing. That's a deskilling architecture and a misapplication of the pattern. Recommend redesign — give the executor edit authority, or refuse the project shape.

Substrate-for-agents deployments (where the typed library is consumed by automated agents rather than imposed on human workers) avoid this. Platform-owned interventional libraries imposed on workers do not.

## Ablation discipline

Every Bucket-C surface should be able to answer: *if I removed the typed substrate and just gave the same LLM raw content, would performance drop?*

```python
def test_ablation_substrate_helps():
    typed_score = run_with_substrate(test_input)
    raw_score   = run_without_substrate(test_input)
    assert typed_score > raw_score, "Substrate is not earning its keep"
```

Define "performance" per project. For a recruiter-triage tool: rate at which top-5 recommended candidates pass screen vs. random-5 baseline. For a writer's voice-checker: rate at which annotations are accepted by the writer. For an ambient finding-injection hook: whether downstream model turns produce meaningfully different reasoning when given vs. not given the findings.

Add the test from day one. Without it, the architecture is theater.

## Anti-patterns and refusal

Decline the pattern when:

- One-shot transform (translate, summarize, format)
- Value is the LLM's natural-language output (chatbot, advisor)
- Pure code refactoring or pure UI
- All inputs already structured — no fuzziness
- Output consumed once and discarded
- Calibration / aggregation across runs is not wanted

Decline template:

> *This sounds like a [transform / chatbot / refactor] task — the value is the LLM's natural-language output, not data about content. A simpler approach is [alternative]. The hybrid-loops pattern would add complexity without buying reliability here.*

Also decline when the deployment power-balance check (above) flags a deskilling shape that the user won't fix.

## When in doubt — one question

Ask the user: *Is the value of this surface the data it produces about content, or the natural-language output the LLM gives you?*

Data → hybrid loop fits. Natural-language → it doesn't.

## Naming

Use these terms internally:

- *lens* / *substrate* / *gate* / *reasoner* / *action*
- *calibration log* / *metabolism*
- *substrate-as-record* vs *substrate-as-vocabulary*
- *hybrid loop* (umbrella; not "compound AI system" or "schemaed cognition")

When talking to the user, prefer their domain vocabulary. Keep the pattern terms internal.

## Citations

- **wesen / Manuel Odendahl** ([the.scapegoat.dev](https://the.scapegoat.dev), [github.com/go-go-golems](https://github.com/go-go-golems)) — closest practitioner prior art. Coined "generalization shaping"; uses "diary," "evidence database," "substrate," "step." Cites Blackboard Systems lineage. Full credit in `references/PRIOR_ART.md`.
- **Devine Lu Linvega / Hundred Rabbits** ([100r.co](https://100r.co)) — the small-typed-tools aesthetic in pure form.
- **AlphaGo / AlphaZero** (Silver et al., 2016, 2017) — fuzzy+hard mutual-constraint template.
- **DreamCoder** (Ellis et al., Nature 2021) — wake/sleep library learning.
- **LILO** (Grand et al., NeurIPS 2024) — LLM-era DreamCoder descendant.
- **Voyager** (Wang et al., 2023) — skill libraries for agents.
- **Christopher Alexander, *A Pattern Language*** (1977) — structural reference for hybrid loops as a unit of design.
- **Blackboard Systems** (Hayes-Roth 1985) — recovered architectural lineage.
- **OpenCog / Hyperon** (Goertzel) — *cite to distinguish*: same intuition, different bet (symbolic reasoning where LLMs now win).

What's plausibly new — empirical generalization of the disciplines beyond the engineering surfaces they were sketched on — remains untested.

## See also and reference loading

The main file is enough for project planning and most design conversations. Load references only when needed:

- `references/BLOCK_GRAPHS.md` — when scaffolding by analogy to a recognizable shape; library of canonical hybrid-loop block-graphs each shown as a Mermaid diagram with brief unpacking
- `references/EXAMPLES.md` — six worked examples with full schema sketches, gate policies, and calibration paths (teacher/parent/advocacy/coach/writer/recruiter); complementary to BLOCK_GRAPHS.md (which gives diagrams; EXAMPLES.md gives implementation-level detail)
- `references/DESIGN_INTERVIEW.md` — when Phase 4's quick design isn't sufficient
- `references/BUILDING_BLOCKS.md` — when the user wants to think about the lower-level algebra (eight primitive block-shapes, pairs, triples, why neither half collapses to the other)
- `references/THE_CASE.md` — when the user is skeptical that hybrid loops are anything more than 1945 von Neumann; the algebra-vs-alphabet-vs-disciplines argument lives here
- `references/AGENT_FRAMEWORKS.md` — when the user thinks this is just DSPy / LangGraph / AutoGen / CrewAI / pydantic with a coat of theory
- `references/PRIOR_ART.md` — when defending novelty or citing lineage
- `references/STACKING.md` — only when project is past v0 and explicitly stacks
