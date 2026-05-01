# The case for hybrid loops

What's new, what isn't, and why the framework needs to exist.

## The algebra is 1945. The alphabet is new.

For 80 years, composing software has meant arranging deterministic nodes connected by typed I/O. Pipelines, substrates, codegen, macros, compilers, linters — all of it. **Hybrid loops as a graph algebra is genuinely von Neumann.** Saying "block A's output type matches block B's input type" doesn't differentiate from Unix pipes or any modular software-architecture text from the last five decades.

If we stop there, we've earned the "yeah, no shit" reaction.

What's actually new isn't the graph algebra. It's that **the set of node types in the alphabet just expanded by one** — and the new one has properties no prior node type had:

1. **Soft-input → structured-output natively.** Drop a transcript in, get typed records out. Pre-LLM: brittle regex, hand-trained classifiers per task, or human review. None composed cleanly. The LLM is the first general-purpose soft-input parser, and it slots into a deterministic graph as a single node type.
2. **Behavior reconfigurable by prose at runtime.** Same node, different system prompt → different behavior. No recompile, no plugin system. Plugin architectures gave a sliver of this; the LLM-as-node generalizes it. Prose-as-program at the node level.
3. **Multi-modal generator from one box.** One node can produce text, JSON, regex, SQL, code, schemas, prompts for other LLMs. Pre-LLM, generators were tied to specific input/output shapes (codegen reads spec; classifier reads features).

Combine those three and you have a node that fits anywhere a human would otherwise have to author specialized code. **That's the alphabet expansion**, and it's why the same kinds of graph shapes keep showing up in different projects — different surfaces are exploiting the same three affordances in different proportions.

## LLMs are fuzzy pattern mappers

The cleanest framing for what an LLM *is*, at the system-design level: an LLM is a **fuzzy pattern mapper** — it reads a structured input and produces a structured output, where the *fuzzy* qualifier marks that the mapping is non-deterministic and approximate. (Manuel Odendahl ("wesen") uses "mapping" and "interface-mapping" routinely in his writing — e.g. ["Tool use and notation as shaping LLM generalization"](https://the.scapegoat.dev/tool-use-and-notation-as-generalization-shaping/) — to describe LLM behavior at this level. The "fuzzy pattern mapper" framing here is built on that vocabulary.)

The word **"transformer"** is deliberately avoided at this level of abstraction — it carries too much baggage from the Vaswani 2017 attention-architecture paper. **"Pattern mapper"** keeps the input-pattern → output-pattern meaning without the ML-jargon collision, and is closer to how compiler-veterans already think.

A senior engineer who's spent twenty years writing compilers, transpilers, linters, query optimizers, codegen tools, and parser generators has been building pattern mappers the whole time. The species they're used to are **deterministic pattern mappers**:

- Hand-coded
- Narrow (one specific mapping per tool)
- Deterministic (same source → same output)
- Errors are typed (failure mode is a structured diagnostic)

LLMs are **fuzzy pattern mappers** — same family, different properties:

- Learned (not hand-coded)
- General (any input pattern → any output pattern, conditioned by prompt)
- Probabilistic (same input → variable output)
- Errors are emergent (no structured diagnostic; failure mode is "drifted from rules" or "made stuff up")

**Neither generalizes the other** in the formal sense — compilers have correctness guarantees, bounded resource use, and verification properties LLMs lack; LLMs handle soft input and prose-as-config compilers can't touch. They're sibling species in the *family of pattern-mapping functions you've been composing your whole career*. The chatbot UI and agent hype framed LLMs as conversational partners or autonomous workers — categories that bury the pattern-mapper view that's actually closer to how an LLM slots into a system graph. **For someone whose mental model is already "pattern mappers connected by typed I/O," an LLM doesn't require a new mental model. It requires adding the *fuzzy* species to the existing diagram, with three properties (general-domain, learned, probabilistic) and three associated disciplines.**

> *Note for ML readers: this framing is deliberately at the systems-design level, not the architecture level. The framework operates at the level of how systems compose typed actors; what's inside the actor (attention, state-space models, mixture-of-experts, whatever comes next) is unchanged by the systems-level claim. The word "mapper" is chosen specifically to avoid the Vaswani-2017 "transformer" collision.*

## Three new disciplines

The new node type has a tax: **non-determinism.** Same input → different output across sampling, drift across long contexts, prose instructions interpreted approximately. That tax means three disciplines the old algebra didn't need to specify:

### Per-block calibration

A non-deterministic node embedded in a deterministic graph means **rolling hit-rate per LLM block** is no longer optional. Without it, when the system misbehaves you can't tell which node dropped the ball. **Every LLM block in the algebra is a candidate for replacement by code if its calibration says it should be** — the framework doesn't claim LLMs are necessary anywhere; it claims that wherever you put one, you also need the deterministic half nearby. That's `cal_log`, and that's Conjecture 1.

### Context-as-code as load-bearing infrastructure

A node configurable by prose means the prose is *production code*. A markdown rule sheet that conditions LLM behavior is no longer "documentation" — it's a binary you ship. Version it, lint it, audit it, calibrate against it. Schemas, DSLs, and structured-output specs are the highest-leverage flavor: same artifact serves as LLM-output, LLM-input-constraint, and code-side validator. That's the audit / substrate-structure-checking conjecture (`metacog`, Conjecture 2).

### Dev-time hybrid loop wrapping the runtime

A node that can author other nodes — LLM writes code, schemas, notation, or prompts — means your *development cycle* is itself a hybrid loop. The runtime stays small (one or two cycles per user-facing decision); the dev-time loop iterates many times across runs (LLM-critic reads transcripts → finds patterns → patches the runtime layers). **Most real systems live in the dev-time regime.** Compress-and-verify schema discovery is one shape of this loop (`schemaforge`, Conjecture 3).

The graph algebra is 1945. The alphabet is new. **The disciplines for the new alphabet are the framework's actual content** — the cycle/block-graph stuff is just the carrier.

## The cognitive-load argument

Beyond a certain complexity threshold, **a human can't hold a full hybrid-loop graph in their head**. Working memory bounds the number of typed blocks + edges + feedback paths a human can simultaneously reason about. That's been a hard limit on architectural complexity for the entire history of software engineering — beyond roughly 7±2 typed components, the architect has to externalize state into diagrams, comments, naming conventions, or just accept that nobody understands the whole system.

LLMs don't have that bound at the same scale. Context windows hold thousands of typed records and hundreds of blocks; the LLM can simultaneously hold the runtime cycle, the dev-time loop, the calibration history, the substrate's metabolism findings, and the patch plan. **The LLM is not just the new node type in the graph; for graphs above a certain size, it's the natural place to *hold the model of the graph itself*.**

This has three consequences:

1. **More-complex hybrid-loop systems become economically viable.** Pre-LLM, system complexity was bounded by what the architect could mentally maintain. Post-LLM, the architect can offload the whole-graph view to the LLM and review specific decisions. The complexity ceiling rose.
2. **The experience floor for non-trivial architectures dropped.** Designing a 12-block hybrid loop with three feedback paths used to require years of senior-architect experience to fit in your head — which made these shapes the privilege of the few engineers who'd accumulated that experience. With the LLM holding the graph, a less-experienced engineer (or a domain expert who isn't an engineer at all — teacher, coach, advocate, parent) can now ship systems whose architectural complexity previously required hiring a senior. **This is the democratization vector** — and it's load-bearing for Conjecture 4 (substrate-as-vocabulary in non-engineering domains): the reason those tools haven't been built isn't that the affordances were missing, it's that the people who feel the need couldn't hold the architecture.
3. **Per-block calibration becomes more load-bearing in proportion to complexity.** If the LLM holds the only complete view of a system that no human reviews end-to-end, you'd better trust each block individually. Calibration is the discipline that makes the offloaded view safe.

This is also a "what's new" — not just an alphabet expansion, but a working-memory expansion *both* for the architect (ceiling rises) *and* for the non-architect (floor drops). It's the second-order effect of the LLM being good at holding patterns: the systems we can build are now larger than the systems we can entirely comprehend, *and* the systems available to non-experts are now structurally richer than the systems they could previously author alone.

## Why the framework has to exist

Three layers of bias all point engineers toward the conventional shape:

1. **Training data shape.** Pre-2023, ~all software engineering looked like conventional engineering. That's what LLMs saw the most of. Hybrid-loops-style architectures barely existed in the corpus because the affordances didn't exist. So the LLM's priors pull toward the shape it has the most evidence for — *pipeline of typed transforms that returns a result*.
2. **Harness shape.** Claude Code, Cursor, Copilot all chose familiar IDE/terminal interfaces — slash commands, file edits, REPL sessions. The choice was correct for adoption, but it means the harness itself is a piece of context-as-code that further conditions the LLM toward IDE-shaped work.
3. **User expectation shape.** People interact with LLMs the way they used to interact with their tools — write code, run tests, fix bugs. The harness doesn't surface graph-of-blocks as a primitive interaction.

Three layers of bias all pointing the same direction. The conventional shape is gravitational.

A fourth layer that's specific to where the LLM-app ecosystem currently sits: **the agent-framework ecosystem itself has been swirling around this space for two years without a consistent vocabulary or taxonomy.** DSPy has *modules* and *signatures*; LangGraph has *nodes* and *edges*; AutoGen has *agents* and *conversations*; pydantic has *models*; structured-outputs has *tool schemas*. Each tool covers one piece of the broader pattern with its own primitives, none of them point at the others, and a mid-level engineer who picks any one of them locks themselves into that tool's cell without seeing the rest of the alphabet. **Senior engineers who've used three of them carry three half-overlapping mental models reconciled by intuition; mid-level engineers carry one and don't know to ask for more.** The ecosystem genuinely lacks a unifying taxonomy. (See `references/AGENT_FRAMEWORKS.md` for the per-tool comparison.)

The skill is a piece of context-as-code that **fights the gravitational pull of the harness toward the LLM's training-data default**. It's not primarily teaching the user anything (though it does). It's *counter-conditioning the LLM* — pulling it out of "what would my training data do here" and into "what does this surface actually want." Every time a hybrid-loops-shaped project gets designed, it's because something pushed back against the default.

This has a sharper consequence: **patterns don't propagate primarily by being explained; they propagate by being proposed in moments where someone has to choose.** The senior engineer who's spent twenty years building pipelines doesn't read a hybrid-loops doc and reorganize their thinking. They encounter the LLM proposing a hybrid-loops shape for their next project, push back or accept, and either way they now have a pattern in their working vocabulary they didn't have before. **The LLM is the dispersal mechanism. The skill is the seed.**

That's the theory of change for the framework:

- The skill counter-conditions the LLM out of pipeline default
- The runnable primitives (`cal_log`, `metacog`, `schemaforge`) make the LLM's proposals concrete enough that humans say yes
- Code written in this shape enters public corpora
- Future model generations train on it
- The gravitational pull eventually rebalances

This is slow — years, not months. But it's the actual path. The skill's success metric isn't "everyone reads SKILL.md." It's *"LLMs propose hybrid-loops shapes when relevant, and users encounter the proposal as a real option."*

## What's not new

The framework isn't claiming novelty in:

- The graph algebra (von Neumann, Unix pipes, every modular-software text since)
- The five-role decomposition (recovered from blackboards, frames, and Soar)
- The bootstrap-loop pattern (AlphaGo 2016, DreamCoder 2021)
- "Generalization shaping" as a design principle (wesen)
- The small-typed-tools aesthetic (Devine Lu Linvega, Hundred Rabbits)
- Pattern-language structure (Christopher Alexander, 1977)

What's new and conjectured (with named falsifiers, all in `../../README.md`):

- That per-evaluator calibration as a shippable primitive materially changes development decisions on real projects (C1)
- That cognitive-bias self-audit on substrate structure generalizes across substrates (C2)
- That compress+verify schema discovery extends to non-program domains (C3)
- That domain-applied substrate-as-vocabulary tooling has unmet demand outside engineering (C4)

Everything else is recovered prior art with a coat of LLMs on top. The recovery itself isn't the contribution; the disciplines for safely placing the new node type in the recovered architecture are.

## Lit review still wanted

The lit review currently in `PRIOR_ART.md` covers blackboards, frames, AlphaGo, DreamCoder, OpenCog, the knowledge-acquisition bottleneck, and the practitioner work directly informing this writeup. There's a broader sweep worth doing: **architectural ideas from the last 50 years that required cheap typed-structure authoring to work, where authoring cost was the limiting factor**. Candidates: argumentation frameworks (Dung 1995), behavior trees (game AI), Petri nets, multi-agent systems (Wooldridge / Jennings), inductive logic programming (Muggleton 1991), genetic programming / program synthesis (Koza onward), constraint logic programming, semantic-web / topic-maps, structured wikis, BPMN / workflow engines, actor-model concurrent systems, cybernetics (Wiener / Ashby). Most of these have the algebra right and were limited by the same bottleneck the knowledge-acquisition critique named. **LLMs lift the bottleneck.** Mapping which old ideas become economically viable for the first time is the next research move; not done yet.
