---
name: hybrid-loops
description: Find and design hybrid-loop surfaces in any project — places where an LLM's fuzzy semantic judgment is consequential enough to warrant typed structure (schema'd substrate, deterministic gating, calibration). Triggers on prompts like "build a tool that tracks/analyzes/evaluates/extracts X over time", "make sense of Y across many Zs", "detect patterns in W", "notice when X is happening", "score/rank/compare documents", "suggests when I'm repeating an anti-pattern", "flag for me when...", "build a system that learns from outcomes", or any project where the value is DATA ABOUT content rather than the LLM's reply. Diagnostic-first — most projects have only 0-3 such surfaces. Domain-agnostic. Do NOT trigger for one-shot transforms (translate, summarize, format), pure UIs, CRUD without judgment, code refactors, or chatbots.
---

# Hybrid Loops

## TL;DR (one screen)

> A **cycle** of alternating LLM-and-code layers that mutually constrain each other. **LLMs bring fluency. Substrates bring discrimination. Code brings restraint.** The point isn't LLM-as-pipeline-stage. The point is *LLM-as-half-of-a-loop* — what one half can't do, the other carries.
>
> 5-phase diagnostic:
> 1. Find candidate **surfaces** in the project (places where fuzzy judgment is happening or should be)
> 2. **Scope** each: A (just call an LLM), B (don't use an LLM), or C (hybrid loop)
> 3. Choose **shape**: substrate-as-record (analytical) or substrate-as-vocabulary (interventional)
> 4. **Quick design** in 3 questions (input, schema, action) with sane defaults for the rest
> 5. **Scaffold** to the surface, not the project. Always include a calibration log and an ablation test.
>
> Five roles in a cycle: **lens** (LLM extracts) → **substrate** (typed records accumulate) → **gate** (deterministic policy filters/scores/ranks) → **reasoner** (LLM consumes substrate) → **action** (deterministic effect; often loops back as new content). Plus two meta-layers that close the loop: **calibration** (predict + verdict log — does the lens actually work?) and **metabolism** (substrate-wide audit — is the accumulated record drifting?).
>
> Decline when: one-shot transform, chatbot, pure UI, no fuzziness in input, output discarded once, or a deployment shape that imposes a substrate on workers who can't edit it.

## Full skill

A design pattern for the **specific places** in a project where a fuzzy semantic judgment benefits from typed structure. Most projects have 0-3 such places. The skill helps Claude identify them, decide whether each warrants the full pattern, and design what's there. Domain-agnostic — applies in health, education, ops, creative, social, business, engineering.

In one sentence: an LLM does fuzzy judgment, a typed substrate captures the result as data, deterministic code does aggregation/restraint/scoring, another LLM reasons over the substrate, an action lands. **LLMs bring fluency. Substrates bring discrimination. Code brings restraint.**

*Note on naming and synthesis.* "Hybrid loops" is the working name in this repository; the broader field has no settled name (see PRIOR_ART.md for adjacent terms — "compound AI systems," "generalization shaping," etc.). The why-it-works-now account in this skill (soft-input/hard-output dispatch fabric, pre-loaded world knowledge, cheap schema iteration, free-text rationales, MCP composition) is this repository's synthesis of tractability factors, not a single citation.

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
- Some surfaces are both (a substrate that grows AND offers a vocabulary).

## Phase 4 — quick design (3 questions)

Three minimum questions. Everything else gets sane defaults; refine after a draft scaffold is on the table.

1. **What's the non-deterministic input?**
2. **What does the lens extract — sketch the schema** (3-7 fields including `notes`, `model_id`, `schema_version`).
3. **What's the action — where does value land?**

Defaults if not specified:
- *Substrate*: JSONL file in project; sqlite when growing past ~1000 records.
- *Gate*: confidence threshold + chronological ordering. Add restraint policies (cooldown, ripeness window) only when over-firing is observed.
- *Reasoner*: read recent records via simple query; produce structured output for the action.
- *Calibration log*: append-only JSONL, predict + verdict. Add from day one even when verdict signals don't yet exist — the log surfaces which verdicts are reachable.
- *Metabolism*: skip in v0.

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

The roles **alternate between fluency (LLM) and discrimination (code)** so each constrains the other. Read the arrows below as a cycle, not a pipeline:

1. **LENS** *(LLM)* — produces typed records from soft input. Fixed schema, `notes` field for graceful failure.
2. **SUBSTRATE** *(typed)* — accumulating record, carries `model_id` + `schema_version`. The substrate is what makes the loop *learn*: each turn's records constrain the next.
3. **GATE** *(code)* — deterministic policy: filtering, scoring, cooldowns, ranking. Where opinionated policy lives. **Code restrains the LLM here so the LLM doesn't have to restrain itself.**
4. **REASONER** *(LLM)* — consumes substrate, produces decisions or generated content. Reads what the lens accumulated and what the gate prioritized; reasons across them.
5. **ACTION** *(code)* — deterministic effect. Often loops back as new content the lens reads next time. **This is what closes the cycle.**

Plus two meta-layers that close *different* loops:

- **CALIBRATION** — predict + verdict log per evaluator. Closes the loop on *whether the lens/reasoner is actually working* (rolling hit-rate). Without it, the architecture is theater.
- **METABOLISM** — periodic substrate-wide phases (audit, prune, refactor). Closes the loop on *substrate quality over time*. Skip until v1+.

The lens may be staged or parallel — treat lens as a *role*, not a single LLM call. Same for the reasoner. The cycle is the structural invariant; how many calls fill each role is project-specific.

## Activation surface

How does the loop fire? Pick one:

- **In-process call** — the surrounding project invokes the lens/reasoner directly. Default for embedded surfaces.
- **MCP tool** — callable on demand by Claude or other agents. Use when the substrate should be queryable from outside.
- **Hooks** (Claude Code lifecycle, request handlers, browser events) — fire automatically. Use for ambient injection.
- **CLI / cron** — scheduled batch. Use only for genuine background processes.
- **Stream watcher** — polls a transcript or stream. Avoid unless polling is unavoidable.

For non-Claude-Code contexts (API apps, web services, mobile, no-code), ask: *where in the host project's normal control flow does this surface get called?* That's the activation point.

## Deployment ethics

Hybrid loops separate *judgment substrate* (the typed library / record) from *judgment execution* (the LLM that picks or produces from it). That separation has political consequences depending on who owns each.

Ask, before scaffolding:

- **Who owns the substrate?** Users, the team, a platform, a service?
- **Who executes the move?** Same party as substrate owner, or different?
- **Can the executor edit the substrate?** Add entries, refine criteria, override?

Power-neutral deployments: single user owns and executes; team shares ownership and execution; platform owns substrate but executor can edit.

Concerning: platform owns substrate, gig-worker or low-status executor must apply the imposed taxonomy without editing. This is a deskilling architecture and a misapplication of the pattern. Recommend redesign — give the executor edit authority, or refuse the project shape.

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

Also decline when the deployment ethics check (above) flags a deskilling shape that the user won't fix.

## When in doubt — one question

Ask the user: *Is the value of this surface the data it produces about content, or the natural-language output the LLM gives you?*

Data → hybrid loop fits. Natural-language → it doesn't.

## Naming

Use these terms internally:

- **lens** / **substrate** / **gate** / **reasoner** / **action**
- **calibration log** / **metabolism**
- **substrate-as-record** vs **substrate-as-vocabulary**
- **hybrid loop** (umbrella; not "compound AI system" or "schemaed cognition")

When talking to the user, prefer their domain vocabulary. Keep the pattern terms internal.

## Citations

- **wesen / Manuel Odendahl** ([the.scapegoat.dev](https://the.scapegoat.dev), [github.com/go-go-golems](https://github.com/go-go-golems)) — closest practitioner prior art. Coined "generalization shaping"; uses "diary," "evidence database," "substrate," "step." Cites Blackboard Systems lineage. See `references/PRIOR_ART.md` for the full citation and credit.
- **Devine Lu Linvega / Hundred Rabbits** ([100r.co](https://100r.co)) — the small-typed-tools aesthetic in pure form.
- **AlphaGo / AlphaZero** (Silver et al., 2016, 2017) — fuzzy+hard mutual-constraint template.
- **DreamCoder** (Ellis et al., Nature 2021) — wake/sleep library learning.
- **LILO** (Grand et al., NeurIPS 2024) — LLM-era DreamCoder descendant.
- **Voyager** (Wang et al., 2023) — skill libraries for agents.
- **Christopher Alexander, *A Pattern Language*** (1977) — structural reference for hybrid loops as a unit of design.
- **Blackboard Systems** (Hayes-Roth 1985) — recovered architectural lineage.
- **OpenCog / Hyperon** (Goertzel) — *cite to distinguish*: same intuition, different bet (symbolic reasoning where LLMs now win).

What's plausibly new is **conjectured, not asserted** — see `../README.md` for the four falsifiable conjectures (per-evaluator calibration discipline, cognitive-bias self-audit on substrate structure, schema discovery for non-program domains, domain-applied substrate-as-vocabulary tooling) and the experiments that would falsify each.

## See also and reference loading

The main file is enough for project planning and most design conversations. Load references only when needed:

- `references/EXAMPLES.md` — when scaffolding by analogy to a known shape
- `references/DESIGN_INTERVIEW.md` — when Phase 4's quick design isn't sufficient
- `references/PRIOR_ART.md` — when defending novelty or citing lineage
- `references/STACKING.md` — only when project is past v0 and explicitly stacks
- `references/PRIMITIVES.md` — when scaffolding and looking for what's already packaged
