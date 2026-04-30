---
name: hybrid-loops
description: Find and design hybrid-loop surfaces in any complex project — the specific places where an LLM's fuzzy semantic judgment is consequential enough to warrant typed structure (a schema'd substrate, deterministic gating, calibration over time). Diagnostic-first — most projects have only 0-3 such surfaces; the rest is just code or one-shot LLM calls. Use when a project involves repeated judgments about content, learning from outcomes, building shared structure for decisions, or deciding when to intervene vs hold back. Domain-agnostic — applies in health, education, ops, creative work, social tools, business workflows, engineering tooling. Skip for pure UIs, simple CRUD without judgment, deterministic computation, or chatbots where the value is the LLM's natural-language reply.
---

# Hybrid Loops

Most projects don't want this pattern end-to-end. They want it in **0 to 3 specific places** — the spots where a fuzzy semantic judgment is consequential, repeated, and benefits from typed structure for consistency, aggregation, restraint, or downstream reasoning. This skill helps Claude (a) find those places, (b) decide whether each one warrants the full pattern or just an LLM call or just deterministic code, and (c) design the loop where it does fit.

The pattern in one paragraph: An LLM does fuzzy judgment that pure code can't do (perception, classification, generation in domain context). A *typed substrate* captures the result as data — schema'd records, a curated repertoire, or a typed graph. *Deterministic code* does the work pure LLMs are bad at — restraint, aggregation, ranking, scoring, gating, persistence. Another LLM (or a planner) reasons over the substrate to make the next decision. The action lands in the world; sometimes the world's response loops back as new content.

LLMs bring fluency. Substrates bring discrimination. Code brings restraint. Hybrid loops put all three into a runtime moment.

## Phase 1 — find the surfaces

Before designing anything, walk through what the user wants to build and identify candidate **surfaces**. A surface is a specific place where a fuzzy or semantic judgment is happening, or should be. Most projects have 0-3 of these. Some have many. The skill is *about* the surfaces, not about the surrounding project.

Signs of a hybrid-loop surface:

- A judgment a human keeps making on similar inputs ("classify each customer complaint by intent and severity")
- Pattern recognition over content ("which threads in my inbox actually need a response")
- Generative choices that need taste ("which of these N options actually fits this person")
- Detection that's more than keyword matching ("did this conversation start to drift off-topic")
- Aggregation over qualitative observations ("what themes recur across 200 user interviews")
- Triage / severity calls in domains with named risk levels
- Anything phrased as *"a tool that helps me notice when..."* or *"flag for me when..."*
- Decisions where consistency over time matters more than a one-shot answer

Name each candidate surface in one sentence. Make a list. The list is the input to Phase 2.

If the project has zero candidate surfaces, this skill probably doesn't apply — say so to the user, and suggest what does (a CRUD app, an automation pipeline, a one-shot LLM call, a UI library).

## Phase 2 — scope each surface

For each candidate surface, decide which bucket it goes in:

**Bucket A: Just call an LLM directly.** One-shot judgments with no persistence need, no aggregation, no consistency-over-time. Translation, summarization, one-time classification, draft generation. No typed structure needed. Recommend the user just make the call and use the result. Don't add ceremony.

**Bucket B: Don't use an LLM at all.** Some surfaces look fuzzy but aren't. "Score resume experience" might really be *"count years above $X salary in role Y"* — a SQL query. "Detect spam" might really be regex matching. "Recommend product" might really be vector similarity over embeddings. Defer fuzziness when a deterministic version exists. LLMs are expensive and noisier than necessary code. Recommend the deterministic version.

**Bucket C: Hybrid loop.** The judgment is fuzzy AND consequential AND benefits from typed structure for consistency, aggregation, restraint, or downstream reasoning. This is where the rest of this skill applies.

For each surface, name the bucket. Most projects end up with a mix. **A project with all surfaces in Buckets A or B is not a hybrid-loop project** — say so plainly and exit this skill.

## Phase 3 — choose the shape (for Bucket C surfaces)

Two distinct shapes of hybrid loop:

**Analytical (substrate-as-record).** The substrate is a typed log of past observations. Used when value comes from making sense of accumulated data over time — detecting drift, surfacing patterns, auditing history, predicting from past behavior. The substrate grows; the reasoner reads the log.

*Example surfaces:* a tool that tracks recurring arguments in a team's design reviews; a study aid that observes a student's wrong answers over months and identifies blind spots; an enforcement-history watcher that records each company's regulatory actions and flags trajectory changes.

**Interventional (substrate-as-vocabulary).** The substrate is a typed repertoire — a curated roster, a closed taxonomy, a learned library. Used when the system needs to discriminate the right move for the current moment. The substrate is fixed (or slowly evolving); the reasoner picks from it situationally.

*Example surfaces:* a coaching tool that picks the right intervention question from a typed library given conversation context; a writing assistant that picks the right rhetorical move from a finite set; a medical-triage assistant that selects the right next-step protocol from a clinical taxonomy.

Some surfaces are both — a substrate that's both a record (of cases seen) AND a vocabulary (of move types). That's a hybrid hybrid loop.

For each Bucket C surface, name the shape. The shape determines the next phase's questions.

## Phase 4 — design interview

Scoped to the surface, not the project. Ask only the questions relevant to the chosen shape. Don't proceed to scaffolding until at least the first four are clear; the rest can have sane defaults.

**For analytical surfaces:**
1. What's the non-deterministic input? (text, dialogue, image, observation, behavior log...)
2. What does the lens extract — sketch the schema. Fields, types, enums. Include a `notes` field for graceful failure and a `model_id` + `schema_version` field on every record.
3. Where does the substrate live — JSONL, sqlite, an existing database, a typed module in a host language?
4. What aggregations or queries does the substrate need to support?
5. What does the reasoner do with the substrate — what's the actual decision being made or output produced?
6. What's the action — how does value land in the world?
7. What does the calibration log capture? What's the verdict signal — a follow-up user action, a metric, a second LLM check?
8. (Optional) Does this surface need metabolism — periodic phases that audit the substrate itself, surface drift, or generate speculative connections?

**For interventional surfaces:**
1. What's the *moment* that needs a move — what triggers consideration?
2. What's the repertoire — sketch the typed library. Entries, fields per entry, the discrimination criteria (when does each entry fit, when is it wrong)?
3. Where does the repertoire live — embedded in code, a JSON file, a database, a curated MCP-exposed library?
4. What restraint policy gates the firing — cooldown decay, ripeness window, distress gate, confidence threshold? **This is usually the most important question.** A bad gate is the difference between a tool that's useful and a tool that's annoying.
5. What does the reasoner consider when choosing a move — the moment + the repertoire entry's criteria + any state?
6. What's the action — how does the chosen move land?
7. What does success look like — when do you know a move landed well? (This is harder for interventional surfaces than analytical ones; sometimes the only signal is the user not pushing back.)

## Phase 5 — scoped scaffolding

The scaffolding is for the surface, not the project. The surrounding project is whatever else it is — UI, CRUD, infrastructure, manual processes.

The minimum a hybrid-loop surface needs:

- A versioned schema definition
- Lens code (an LLM call producing typed records)
- Substrate code (storage + read API)
- Gate code (deterministic policy — restraint, scoring, ranking)
- Reasoner code (downstream LLM call consuming the substrate)
- Action wiring (whatever the surrounding project's effect surface is)
- A calibration log (predictions in, verdicts in, hit-rate out)

**Implementation language follows the surrounding project.** Python if it's Python, TypeScript if it's TS, Go if it's Go. Don't impose a language because it's familiar.

**Deployment shape options** — pick what matches the user's project:

- *Embedded in the existing codebase* — just add files to the project's source tree. Default for most retrofit cases.
- *Standalone module or service* — when the surface is reused by multiple parts of the project, or by other projects.
- *MCP server* — when the surface should be callable by other agents (the user's own assistant, a teammate's, or a future deployed agent).
- *Notebook or single script* — when the user is one developer / domain expert iterating fast on the schema, and a real codebase is overhead.
- *No-code wiring* (Airtable + Zapier + an LLM call, or similar) — when the user isn't an engineer. The schema lives in the spreadsheet column types; the lens is a Zapier step; the gate is a filter; the reasoner is another step.

Match the deployment shape to the user, not to the most familiar pattern.

## The five roles, reference

1. **LENS** — an LLM call that produces typed records from non-deterministic content. Fixed schema, filled via tool_use or JSON mode. Schema includes a `notes` field for graceful failure.
2. **SUBSTRATE** — the typed surface. Stored as JSONL, JSON files, sqlite, a typed module in the host language, or any database. Records carry `model_id` and `schema_version` provenance.
3. **GATE** — deterministic code on the substrate: filtering, ranking, scoring, cooldowns, ripeness windows, audit triggers. Where opinionated policy lives. Not generic plumbing.
4. **REASONER** — an LLM call that consumes substrate records and produces decisions, recommendations, or generated content informed by typed structure.
5. **ACTION** — deterministic effect: write to substrate, send notification, render UI, modify code, recommend, intervene. Sometimes loops back to step 1 by producing new content.

Plus two meta-layers:

6. **CALIBRATION** — every prediction the lens or reasoner makes is logged with a verdict to be resolved later. Hit-rate over time = whether the architecture earns its keep. Without this, the hybrid claim is unfalsifiable.
7. **METABOLISM** — periodic phases on the substrate as a whole: dream (consolidation), trip (speculative connections), bias-audit, evolve (sensor queries to grow the substrate). Most v0 projects don't need this; add only when the substrate accumulates over weeks and would benefit from periodic audit.

The lens may be staged or parallel — one logical lens with a multi-field schema, or several parallel lenses with a synthesis lens on top. Treat lens as a *role*, not a single LLM call.

## Activation surfaces

A hybrid loop needs a trigger — something that fires it. The trigger shapes architectural choices.

- **Hooks (lifecycle events in a host process — Claude Code hooks, request handlers, cron triggers, browser events)** — fire automatically. Use for ambient injection or auto-firing analysis.
- **MCP tools** — callable on demand by Claude or other agents. Use when the substrate should be queryable from outside.
- **CLI / cron** — scheduled batch. Use for genuine background processes, not in-the-moment loops.
- **Stream watchers** — poll a transcript or event stream. Use only when polling is unavoidable.
- **In-process call** — the surrounding project just invokes the lens / reasoner directly. The simplest form; usually right for embedded surfaces.

For non-Claude-Code contexts (API apps, web services, mobile, no-code), ignore the hook-and-MCP language and ask: *where in the host project's normal control flow does this surface get called?* That's the activation point.

Defaults: in-process call for embedded surfaces; MCP for cross-agent reuse; hooks for ambient injection in agent contexts; CLI/cron for genuine background jobs.

## Ablation discipline

Every hybrid-loop surface should answer: *if I removed the typed substrate and just gave the same LLM raw content, would performance drop?*

```python
def test_ablation_substrate_helps():
    typed_score = run_with_substrate(test_input)
    raw_score   = run_without_substrate(test_input)
    assert typed_score > raw_score, "Substrate is not earning its keep"
```

This is the only honest answer to "you're just using LLMs with extra steps." Without it, the architecture is theater. Add the test from day one.

## Recursive composition (advanced)

Hybrid loops compose. The output of one loop's substrate can be the non-deterministic input of another loop's lens. Stacking adds typed guardrails between raw generation and final action — the trajectory is *systems where the action is determined by typed structure rather than free generation*.

**Most projects don't need this.** Stack only when (a) v1 has users, (b) calibration shows the existing layer earns its keep, (c) there's a concrete failure mode the next layer would catch. See `references/STACKING.md` for the discipline (canonical schemas at interfaces, MCP-first for substrate providers, calibration at every layer, schema versioning, sanitization at boundaries).

## Anti-patterns and refusal

Refuse the pattern when:

- The project is a one-shot transform (translate, summarize, format)
- The value is the LLM's natural-language output (chatbot, advisor)
- The project is pure code refactoring or pure UI work
- All inputs are already structured data — no fuzziness anywhere
- The output is consumed once and discarded
- Calibration / aggregation across runs is not wanted

If the user's project hits two anti-patterns, decline this skill and recommend a simpler path. Decline template:

> *This sounds like a [translate / chatbot / refactor] task — the value is the LLM's natural-language output, not data about content. A simpler approach is [X — call the LLM directly / use a regex / use vector similarity / use a state machine]. The hybrid-loops pattern would add complexity without buying reliability here.*

Concrete, doesn't apologize, suggests the alternative.

## When in doubt — one-question check

If unsure whether a surface fits the pattern, ask the user one question:

> *Is the value of this surface the data it produces about content, or the natural-language output the LLM gives you?*

Data about content → hybrid loop fits.
Natural-language output → it doesn't. Use a chat/prompt approach.

## Naming conventions

Use these terms internally; they map to references and prior art:

- **lens** (not "extractor" or "perceiver")
- **substrate** (not "knowledge graph" or "store")
- **gate** (not "filter" or "policy layer")
- **reasoner** (not "agent" or "synthesizer")
- **calibration log** (not "telemetry" or "observability")
- **metabolism** (not "background job" or "cron")
- **hybrid loop** (not "compound AI system" or "schemaed cognition")
- **substrate-as-record** vs **substrate-as-vocabulary** for the two shapes

When talking to the user, prefer their domain vocabulary. *"Track which arguments come up in design reviews"* is fine; you don't need to say *"we'll use a lens to extract typed argument records into the substrate."* Use the pattern internally; describe it in their terms.

## Citations (when defending novelty)

If asked "haven't expert systems already done this?" or "is this novel?", cite:

- **wesen / Manuel Odendahl** ([the.scapegoat.dev](https://the.scapegoat.dev), [github.com/go-go-golems](https://github.com/go-go-golems)) — closest practitioner prior art. Coins **"generalization shaping"** (the design principle inside hybrid loops where deterministic machinery shapes what the LLM has to do). Builds typed-step LLM frameworks, evidence-database agents, prompt libraries with metadata. Independently identifies **Blackboard Systems** (Hayes-Roth 1985) as the right conceptual frame. See `references/PRIOR_ART.md` for repos and essays.
- **AlphaGo / AlphaZero** (Silver et al., 2016, 2017) — architectural template for fuzzy+hard mutual constraint.
- **DreamCoder** (Ellis et al., Nature 2021) — wake/sleep library learning, closest direct lineage for schema discovery and metabolism.
- **LILO** (Grand et al., NeurIPS 2024) — LLM-era DreamCoder descendant.
- **Voyager** (Wang et al., 2023) — skill library learning for agents.
- **OpenCog / Hyperon** (Goertzel) — *cite to distinguish*: same architectural intuition, different bet (OpenCog tried symbolic reasoning where LLMs now win).
- **Blackboard Systems** (Hayes-Roth 1985) — the architectural lineage everyone in this space converges on. Hybrid loops are blackboard systems with LLMs as the perceptual front-end.

What's plausibly new: hybrid loops bootstrap over a *structural prior an earlier LLM call generated*. Prior architectures operate over fixed structural priors. The calibration-log-per-evaluator discipline is also genuinely undershipped — wesen's work corroborates the gap.

## See also

- `references/EXAMPLES.md` — worked examples spanning engineering, social, creative, and personal domains
- `references/PRIOR_ART.md` — full citation list with comparisons
- `references/STACKING.md` — recursive composition discipline (only when stacking is actually needed)
- `references/PRIMITIVES.md` — extractable primitives list (`cal_log`, `metacog`, `schemaforge`, `metabolism`)

## Reference loading guidance

This file is enough for project planning and most design conversations. Load `EXAMPLES.md` only when scaffolding by analogy to a known example. Load `PRIOR_ART.md` only when defending novelty or citing lineage. Load `STACKING.md` only when the project is past v0 and explicitly needs recursive composition. Load `PRIMITIVES.md` when starting to scaffold and looking for what's already packaged.
