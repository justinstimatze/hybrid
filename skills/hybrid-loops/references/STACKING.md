# Stacking hybrid loops

The recursive-composition trajectory: hybrid loops compose recursively until the system has so many layers of typed introspection between raw generation and final action that the action is overwhelmingly determined by structure rather than free generation. *"Back and forth forever until what's getting generated and run is on top of so many layers of meta guardrails it basically always does the right thing."*

The reason stacking works at all is **mutual generation**: each layer doesn't just constrain the next, it *generates the working surface the next layer operates over*. The LLM writes typed records — and often writes the schema, notation, gate logic, or code those records live in. The deterministic layer aggregates and shapes those records into the next LLM call's input. Add a layer above and you get a critic-LLM reading transcripts of layer-N's behavior and writing patches that re-shape layers below. Each addition is another generative half of another loop.

This file unpacks that claim — what it buys, what it costs, when it saturates, and the discipline required to make it work.

## Two regimes: runtime and development-time

Stacking shows up at two scales, and the discipline differs:

**Runtime stacking** — multiple cycles fire per single user-facing decision. Latency-bound; token cost compounds multiplicatively in synchronous chains. Worth doing when the marginal layer's reliability gain exceeds the cost. Saturates fast.

**Development-time stacking** — cycles wrap around the runtime, with humans (or LLMs) reading transcripts of runtime behavior and patching the deterministic layers below. *Not* latency-bound; cost is per-iteration not per-decision. The classic shape: a runtime engine + LLM player → transcript log (deterministic) → LLM-critic panel reads transcripts → structured findings (typed records) → patch plan (deterministic prioritization) → LLM writes code/schema/prompt changes → runtime picks up the change next turn. The development loop itself is a full hybrid loop wrapped around the runtime hybrid loop. **This is where most stacks live in practice** — the runtime stays one or two cycles deep while the development loop iterates across many runs.

When someone says "we have a hybrid-loops architecture," they usually mean the runtime cycle. When the system actually works, it's often because there's a development loop above it that has been iterating for weeks or months.

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
- **Schema brittleness.** Layer N+1 only works if layer N's schema is stable. Schema versioning becomes load-bearing.
- **Diagnosis becomes harder.** When the final action is wrong, which layer dropped the ball? Without per-layer calibration, the answer is "all of them."
- **Cross-layer prompt-injection surfaces.** Each typed record passing between layers is a potential injection vector if the LLM that produced it was prompted by untrusted input. Defending requires sanitization at every boundary.

## Saturation — the empirical question

Plausible hypothesis: reliability increases monotonically with N up to some point, then saturates, then declines as cost exceeds guardrail value. The N depends on:

- How well layers compose (mismatched schemas hurt)
- How accurate each layer's calibration is
- How much the task tolerates latency
- Whether errors are systematic (cascade) or independent (cancel)

The user's working hypothesis is testable: stack increasing N on a fixed task, measure final-output error rate vs. N. The shape of that curve is the answer. *Nobody has run that experiment at small scale*; it's the empirical research project hidden inside this work.

## Composition discipline

When stacking is the goal, the architecture needs explicit composition discipline:

### 1. Canonical schemas at interfaces

A "finding" record, a "claim" record, an "entity" record — these should have standard shapes that multiple loops produce and consume. Without canonical schemas, composition is one-off integration that doesn't scale.

If a project family has *almost-canonical* shapes (claim records in one project, mechanism records in another, finding records in a third) but they're not aligned, a v1 of the hybrid pattern in that family would standardize the cross-cutting record shapes.

### 2. MCP-first for substrate providers

If a loop's output is intended for another loop, expose it as MCP tools. This is the protocol that makes composition work at solo-developer scale. Hook-shaped loops (those that fire automatically on lifecycle events) need an MCP layer added if their output should compose with other loops.

### 3. Calibration at every layer

Per-layer hit-rate tracking is the only way to diagnose where the chain breaks. Each LLM call in the chain logs its prediction; verdicts get resolved against downstream outcomes; hit-rate is measured per layer.

### 4. Provenance on every record

Layer N+2 needs to be able to trace a finding back through layers N+1 and N to the ground truth (often a human-authored document). Provenance fields on every typed record: `model_id`, `prompt_hash`, `source_documents`, `parent_record_ids`, `layer_index`.

### 5. Schema versioning

Every record carries `schema_version`. Layer-N schema bumps don't break Layer-N+1 reads; they trigger re-extraction. Layer-N+1 declares which schema versions it accepts.

### 6. Sanitization at boundaries

Each typed record passing between layers should be treated as potentially adversarial input to the next layer's prompts. Wrap content in untrusted-source delimiters; strip prompt-injection patterns; never let upstream-generated text directly construct downstream prompts without sanitation.

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

The bias should be toward fewer layers, not more. Add a layer only when there's a concrete failure mode it catches that the existing layers don't.
