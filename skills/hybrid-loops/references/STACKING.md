# Stacking hybrid loops

The recursive-composition trajectory: hybrid loops compose recursively until the system has so many layers of typed introspection between raw generation and final action that the action is overwhelmingly determined by structure rather than free generation. *"Back and forth forever until what's getting generated and run is on top of so many layers of meta guardrails it basically always does the right thing."*

The reason stacking works at all is **mutual generation**: each layer doesn't just constrain the next, it *generates the working surface the next layer operates over*. The LLM writes typed records — and often writes the schema, notation, gate logic, or code those records live in. The deterministic layer aggregates and shapes those records into the next LLM call's input. Add a layer above and you get a critic-LLM reading transcripts of layer-N's behavior and writing patches that re-shape layers below. Each addition is another generative half of another loop.

This file unpacks that claim — what it buys, what it costs, when it saturates, and the discipline required to make it work.

## Two regimes: runtime and development-time

Stacking shows up at two scales, and the discipline differs:

**Runtime stacking** — multiple cycles fire per single user-facing decision. Latency-bound; token cost compounds multiplicatively in synchronous chains. Worth doing when the marginal layer's reliability gain exceeds the cost. Saturates fast.

**Development-time stacking** — cycles wrap around the runtime, with humans (or LLMs) reading transcripts of runtime behavior and patching the deterministic layers below. *Not* latency-bound; cost is per-iteration not per-decision. The classic shape: a runtime engine + LLM player → transcript log (deterministic) → LLM-critic panel reads transcripts → structured findings (typed records) → patch plan (deterministic prioritization) → LLM writes code/schema/prompt changes → runtime picks up the change next turn. The development loop itself is a full hybrid loop wrapped around the runtime hybrid loop. **This is where most stacks live in practice** — the runtime stays one or two cycles deep while the development loop iterates across many runs.

When someone says "we have a hybrid-loops architecture," they usually mean the runtime cycle. When the system actually works, it's often because there's a development loop above it that has been iterating for weeks or months.

### A third regime: recursive harness authoring (tier-stacking)

Both regimes above keep the *authoring* outside the loop — a human or critic-LLM patches a system someone else built. A stronger recursion is **vertical**: a more capable tier authors the *entire* harness a cheaper tier runs inside — its prompts, schema, gate, verifier — and re-authors it as signals propagate up, rather than just patching a prompt. The repeated unit is the whole loop.

Two demands it adds: the guardrail set is itself tier-authored and must be re-authored as failure modes mutate (a static verifier silently loses coverage); and the up-channel a human silently provided in ordinary dev-time loops now has to be carried by the structure. *If you can't articulate what carries the up-signal once the human is gone, you don't have a tier-stack — you have a dev-time loop with the human elided and the failures about to go silent.* Single-hop versions of this are already shipping; what's open is stacking it across depth.

**Shipped single-hop instance (2026).** Anthropic's Claude Code dynamic workflows (Shihipar & Bidasaria, *"A harness for every task,"* blog, Jun 2 2026) put this in production at the smallest possible grain: the orchestrating tier (Claude, on Opus 4.8 per the announcement) authors a short JavaScript harness — subagent prompts, structured outputs, gate logic via `parallel()`/`pipeline()` barriers, named verification stages (adversarial verification, tournament, generate-and-filter, loop-until-done) — that a cheaper, more numerous tier of subagents then executes in their own isolated context windows. It's a real instance of "a more capable tier authors the entire harness a cheaper tier runs inside," not a patched prompt. The blog's own stated motivation doubles as evidence for two of this file's disciplines: the three failure modes it names — *agentic laziness*, *self-preferential bias*, *goal drift* — are exactly what accumulates in one unbroken context, so isolating subagents into fresh windows is this file's #7 (lossless up-signal) applied at the authoring layer rather than the inter-layer-record layer. Its *quarantine* pattern for triage workflows — agents that read untrusted public content are barred from high-privilege actions, which a separate agent performs based only on the first agent's typed output — is a shipped instance of #6 (sanitization at boundaries) and rhymes with `PROCEDURE.md` §4's bounded action. Read literally, the tool also partially answers "what's open is stacking it across depth": it refuses recursion past one hop (a workflow invoked from inside a child workflow throws), a live data point that the depth question is currently closed by design choice rather than opened by demonstrated capability. See `ORCHESTRATION_SHAPES.md` (shape 1's "third path" update) and `PRIOR_ART.md` for the fuller citation and hedging.

## What stacking buys

Each additional typed constraint layer can:

- Catch errors the layer below missed
- Constrain free generation to typed dispatch
- Add audit-trail provenance
- Enable calibration on intermediate outputs (not just final ones)
- Make the system's reasoning legible to downstream agents

In the limit, stacking approaches *correctness by construction*: the final generation is so heavily scaffolded that the LLM's freedom is mostly between near-equivalent typed outputs.

## What stacking costs

- **Token cost compounds.** N layers ≈ N LLM calls per action. Synchronous chains compound multiplicatively in latency; async chains compound additively in dollars.
- **Schema brittleness.** Layer N+1 only works if layer N's schema is stable. Schema versioning becomes critical.
- **Diagnosis becomes harder.** When the final action is wrong, which layer dropped the ball? Without per-layer calibration, the answer is "all of them."
- **Cross-layer prompt-injection surfaces.** Each typed record passing between layers is a potential injection vector if the LLM that produced it was prompted by untrusted input. Defending requires sanitization at every boundary.

## Saturation — the empirical question

Plausible hypothesis: reliability increases monotonically with N up to some point, then saturates, then declines as cost exceeds guardrail value. The N depends on:

- How well layers compose (mismatched schemas hurt)
- How accurate each layer's calibration is
- How much the task tolerates latency
- Whether errors are systematic (cascade) or independent (cancel)

The user's working hypothesis is testable: stack increasing N on a fixed task, measure final-output error rate vs. N. The shape of that curve is the answer. *Nobody has run that experiment at small scale*; it's the empirical research project hidden inside this work.

### Saturation is a quality question. Convergence is a structural one — and it bites first.

The reliability-vs-N curve above asks "does each added layer still buy accuracy?" But there's a prior, harder constraint with nothing to do with model quality: in a stack with feedback flowing *up* (drift signals, escalations, calibration verdicts), can the supervising layer still perceive the error it is responsible for correcting? If each upward hop loses fidelity — a prose summary degrades at every boundary — the top sees an attenuated signal and stalls at a floor it can't see below. Past some depth the loop goes *blind* before quality ever saturates.

Two consequences, both diagnostic-first:
1. **Test the structural condition before the quality curve.** You can model the feedback topology — depth, per-hop signal loss, the thresholds that trigger correction — without live model calls, and find where the loop goes blind. The expensive N-sweep is only worth running inside that ceiling.
2. **Attack the up-signal's fidelity, not the layer count** (see #7, Lossless up-signal). Any max-safe-depth this yields is a *band* contingent on your thresholds — the defensible claim is directional: rising per-hop loss collapses safe depth.

## Composition discipline

When stacking is the goal, the architecture needs explicit composition discipline:

### 1. Canonical schemas at interfaces

A "finding" record, a "claim" record, an "entity" record — these should have standard shapes that multiple loops produce and consume. Without canonical schemas, composition is one-off integration that doesn't scale.

If a project family has *almost-canonical* shapes (claim records in one project, mechanism records in another, finding records in a third) but they're not aligned, a v1 of the hybrid pattern in that family would standardize the cross-cutting record shapes.

### 2. MCP-first for store providers

If a loop's output is intended for another loop, expose it as MCP tools. This is the protocol that makes composition work at solo-developer scale. Hook-shaped loops (those that fire automatically on lifecycle events) need an MCP layer added if their output should compose with other loops.

### 3. Calibration at every layer

Per-layer hit-rate tracking is the only way to diagnose where the chain breaks. Each LLM call in the chain logs its prediction; verdicts get resolved against downstream outcomes; hit-rate is measured per layer.

### 4. Provenance on every record

Layer N+2 needs to be able to trace a finding back through layers N+1 and N to the ground truth (often a human-authored document). Provenance fields on every typed record: `model_id`, `prompt_hash`, `source_documents`, `parent_record_ids`, `layer_index`.

### 5. Schema versioning

Every record carries `schema_version`. Layer-N schema bumps don't break Layer-N+1 reads; they trigger re-extraction. Layer-N+1 declares which schema versions it accepts.

### 6. Sanitization at boundaries

Each typed record passing between layers should be treated as potentially adversarial input to the next layer's prompts. Wrap content in untrusted-source delimiters; strip prompt-injection patterns; never let upstream-generated text directly construct downstream prompts without sanitation.

### 7. Lossless up-signal between layers

The signal a layer sends *up* comes in two channels, and they don't age the same. A prose summary ("the layer below seems to be drifting") loses fidelity at every hop. A deterministic verifier's output — a typed divergence record (`{id, reason, similarity}`) — is lossless at any depth for whatever it covers. So the design move is: route the largest checkable fraction of your inter-layer signal through a deterministic comparator rather than a model's prose summary. The honest limit is that the checkable fraction is capped by what you can write a check for, and the un-checkable residual tends to be exactly the high-value drift; holding coverage steady as failure modes mutate means re-authoring guardrails on an ongoing basis. Closely tied to #3 (the comparator is also where calibration lives) and #6 (a typed record is easier to sanitize than free prose).

## Open questions

- **Where does saturation actually fall?** Empirical measurement, not theory.
- **Do canonical cross-cutting record shapes exist?** Or does each domain need its own?
- **Can layer-N-1 errors be reliably caught at layer N?** Or do some classes of error always slip through?
- **What's the right calibration cadence?** Does every prediction get logged, or is sampling sufficient at scale?
- **How does this interact with reasoning models?** A reasoning model at one layer may absorb several discrete loops; does that change the saturation curve?

These are the research questions the architecture opens up. None of them have published answers in 2026.

## When *not* to stack

- v0 prototypes — don't stack. Single-layer hybrid loop until the design is right.
- Latency-critical paths — every layer adds round-trip cost.
- Tasks where the LLM's free generation is the value — stacking constraints destroys what you wanted.
- When you can't measure reliability at each layer — stacking blindly is worse than not stacking.

The bias should be toward fewer layers. Add a layer only when there's a concrete failure mode it catches that the existing layers don't.
