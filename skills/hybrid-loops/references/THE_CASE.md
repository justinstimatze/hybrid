# The case for hybrid loops

What's new, what isn't, and why the framework needs to exist.

## The algebra is 1945. The alphabet is new.

Eighty years of composing software has meant arranging deterministic blocks connected by typed I/O — pipelines, substrates, codegen, macros, compilers, linters. The graph algebra of hybrid loops is, plainly, von Neumann. Saying "block A's output type matches block B's input type" doesn't differentiate this work from Unix pipes or from any modular-software-architecture text written since 1975.

The set of block types in the alphabet just expanded by one, and the new one has properties no prior block type had:

1. **Soft-input → structured-output natively.** Drop a transcript in, get typed records out. Pre-LLM, this required brittle regex, hand-trained classifiers per task, or human review — none of which composed cleanly. The LLM is the first general-purpose soft-input parser, and it slots into a deterministic graph as a single block type.
2. **Behavior reconfigurable by prose at runtime.** Same block, different system prompt, different behavior. No recompile, no plugin system. Plugin architectures gave a sliver of this; the LLM-as-block generalizes it. *Prose-as-program* at the block level.
3. **Multi-modal generator from one box.** One block can produce text, JSON, regex, SQL, code, schemas, prompts for other LLMs. Pre-LLM, generators were tied to specific input/output shapes — codegen tools read specs, classifiers read features.

Combine those three and you have a block that fits a wider range of positions than any prior single technology — anywhere a human would otherwise author specialized code, *as long as the position tolerates non-determinism*. That qualifier matters. Avionics flight-control, medical-device firmware, cryptographic primitives, kernel-level concurrency, real-time control loops: those domains' correctness bars require deterministic computation as a precondition, and no LLM block belongs in their runtime path. Outside that corner — most of what engineers work on day-to-day — the same block-graph shapes keep surfacing, different projects exploiting the same three affordances in different proportions.

## LLMs are fuzzy pattern mappers

The cleanest framing for what an LLM *is*, at the system-design level: a *fuzzy pattern mapper* — a function from input to output where the mapping is non-deterministic and approximate. Manuel Odendahl ("wesen") uses "mapping" and "interface-mapping" routinely — see ["Tool use and notation as shaping LLM generalization"](https://the.scapegoat.dev/tool-use-and-notation-as-generalization-shaping/) — and the framing here builds on that vocabulary. Inputs and outputs don't have to be structured on either side; the fuzziness is in *how* the mapping happens, not in the shape of the data.

Any working engineer has been using *deterministic* pattern mappers their whole career — compilers, linters, type checkers, SQL, regex, JSON-to-Protobuf, ETL stages, shell pipes. A few have written the tools; most just compose them. Either way the abstraction is familiar: a function from structured input to structured output.[^transformer-note] LLMs are the same family with different defaults.

### Complementary cost gradients

Each actor has an effort-cost gradient running the opposite direction along the determinism axis.

```
 effort
  cost
   │\                /
   │ \              /
   │  \            /          \  = code (cheap when deterministic;
   │   \          /                expensive when random — you need
   │    \        /                 crypto RNG, hardware entropy, etc.)
   │     \      /
   │      \    /              /  = LLM  (cheap when fuzzy / free-form;
   │       \  /                    expensive when deterministic — you
   │        \/   ← crossover       need temp 0, structured outputs,
   │        /\                     constrained decoding, retry-until-
   │       /  \                    stable, ensembling, calibration)
   │      /    \
   │     /      \
   │    /        \
   │   /          \
   │  /            \
   │ /              \
   │/                \
   └─────────────────────→ determinism
    random            deterministic
    ↳ LLM cheaper   ↳ code cheaper
```

The curves cross. Below the crossover (toward random / fuzzy), the LLM is the cheaper actor; above it (toward deterministic), code is. The design move at each block is to figure out which side of the crossover the work sits on, and place the corresponding actor there. Forcing an actor across is a known move with real uses — deterministic LLM blocks for safety-critical outputs, randomized code for cryptography — but it costs, and the framework's job is naming when the cost is worth it.

The chart isn't a full taxonomy. Most pre-LLM technologies (Bayesian nets at inference, fuzzy controllers, HMMs in Viterbi mode, random forests) cluster with code: they're *deterministic algorithms computing over uncertainty-valued types*. What distinguishes them from each other lives on other axes — output representation, generality, authoring cost, state accumulation — not the determinism axis. (Cluster walk-through in the footnote.[^cluster-walkthrough]) LLMs are unusual on this axis specifically because they're the first *general-purpose* technology whose native mode is non-deterministic. Monte Carlo methods are the closest historical comparator, but narrow. The gap LLMs fill is general-purpose-ness combined with non-deterministic-native operation, cheaply, for the first time.

### What lifts isn't just authoring; it's modification

The classical critique of expert systems, Cyc, Bayesian nets, and fuzzy controllers was the *knowledge-acquisition bottleneck* — initial authoring of typed structures was prohibitively expensive. LLMs lift that. They lift a second cost less often discussed: *modification and maintenance*. A rule sheet you edit by changing a markdown file and immediately see different LLM behavior is materially different from one you recompile, redeploy, and re-version. Both halves of that collapse are what makes the previously-uneconomic methods now economic. Authoring cost wasn't the only thing that killed those traditions — rules brittle on edge cases, the world messier than first-order logic, edge-case proliferation outpacing patches — and LLMs cover those defects too: the fuzziness that makes them unreliable as the only actor is the same fuzziness that makes them tolerant of the messiness pure rules can't survive.

Neither family generalizes the other in the formal sense. Compilers have correctness guarantees, bounded resource use, and verification properties LLMs lack; LLMs handle soft input and prose-as-config that compilers can't touch. The chatbot UI and agent-as-autonomous-worker framing buries this — it presents LLMs as a category apart, when "fuzzy sibling of the compiler" is closer to how an LLM actually slots into a system graph.

[^transformer-note]: The framing is at the systems-design level, not the architecture level. What's inside the actor — attention, state-space models, mixture-of-experts, whatever comes next — is unchanged by the systems-level claim. *Pattern mapper* rather than *pattern transformer* is chosen to avoid the Vaswani-2017 attention-architecture collision; the compiler-tradition meaning of "transformer" (a function from one structured representation to another) gets lost the moment ML readers see the word, so "mapper" keeps the bridge useful.

[^cluster-walkthrough]: Where pre-LLM technologies sit on the cost-vs-determinism chart.

    **Cluster with code** (cheap when deterministic): classical fuzzy logic / fuzzy controllers (Zadeh) — *deterministic* computation over *fuzzy values*; Bayesian networks at inference; HMMs in Viterbi mode; random forests and deterministic ensemble ML at inference; compilers, regex, SQL, parsers, linters, type checkers, codegen.

    **Cluster with LLMs** (cheap when stochastic): Monte Carlo methods, simulated annealing, MCMC, particle filters; genetic algorithms / GP at search time; HMMs in sampling mode.

    **Mixed** (one technology, two curves depending on mode): reinforcement learning (stochastic training, deterministic inference); HMMs (sampling vs Viterbi).

## Three new disciplines

The new block type has a tax: *non-determinism*. Same input can produce different output across sampling, drift across long contexts, prose instructions interpreted approximately. That tax means three disciplines the old algebra didn't need to compel. The disciplines themselves are recovered — calibration is unit testing for non-deterministic functions; context-as-code is configuration management; the dev-time loop is CI/CD with an LLM in the critic seat. What's new is that you can't skip any of them once an LLM block is in the graph.

### Per-block calibration

A non-deterministic block embedded in a deterministic graph means rolling hit-rate per LLM block is no longer optional. Without it, when the system misbehaves you can't tell which block dropped the ball. Every LLM block in the algebra is a candidate for replacement by code if its calibration says it should be — wherever you put one, you also need the deterministic half nearby and a way to tell whether the LLM half is earning its keep. The attribution isn't always clean. End-to-end agent failures often have joint-failure modes where blame doesn't decompose neatly per block; the discipline holds the bar where attribution is recoverable (most blocks most of the time) and acknowledges the residual where it isn't.

### Context-as-code as load-bearing infrastructure

A block configurable by prose means the prose is *production code*. A markdown rule sheet that conditions LLM behavior is no longer documentation — it's a binary you ship. Version it, lint it, audit it, calibrate against it. Schemas, DSLs, and structured-output specs are the highest-leverage flavor: same artifact serves as LLM-output, LLM-input-constraint, and code-side validator.

### Dev-time hybrid loop wrapping the runtime

A block that can author other blocks — LLM writes code, schemas, notation, or prompts — means the *development cycle* is itself a hybrid loop. The runtime stays small (one or two cycles per user-facing decision); the dev-time loop iterates many times across runs (LLM-critic reads transcripts → finds patterns → patches the runtime layers). Most real systems live in the dev-time regime. Compress-and-verify schema discovery is one shape of this loop.

These three disciplines are the framework's actual content. The cycle and block-graph machinery is the carrier.

## The cognitive-load argument

Beyond a certain complexity threshold, a human can't hold a full hybrid-loop graph in their head. Working memory bounds the number of typed blocks, edges, and feedback paths a person can simultaneously reason about — beyond roughly 7±2 components, the architect externalizes into diagrams, naming conventions, or just accepts that nobody understands the whole system.

LLMs don't have that bound at the same scale. Context windows hold thousands of typed records and hundreds of blocks; the LLM can simultaneously hold the runtime cycle, the dev-time loop, the calibration history, the substrate's findings, and the patch plan. The LLM is not just the new block type in the graph; for graphs above a certain size, it's the natural place to *hold the model of the graph itself*.

Four consequences:

1. *More-complex hybrid-loop systems become economically viable.* The architect offloads the whole-graph view to the LLM and reviews specific decisions. The complexity ceiling rose.
2. *The experience floor for non-trivial architectures dropped.* Designing a 12-block hybrid loop with three feedback paths used to require years of senior-architect experience to fit in one head — which made these shapes the privilege of a few. With the LLM holding the graph, a less-experienced engineer (or a domain expert who isn't an engineer at all — teacher, coach, advocate, parent) can ship systems whose architectural complexity previously required hiring a senior. The reason these tools haven't been built outside engineering isn't that the affordances were missing; it's that the people who feel the need couldn't hold the architecture.
3. *Calibration grows more load-bearing in proportion to complexity.* If the LLM holds the only complete view that no human reviews end-to-end, you'd better trust each block individually. Calibration is the discipline that makes the offloaded view safe.
4. *The typed substrate stays structurally necessary, not just a UX affordance for humans.* Holding the graph in context is read-only. The moment the LLM acts on the graph — updates a record, refactors the schema, regenerates a section — that act is fuzzy. Fields drop, joins hallucinate, refactors drift. Deterministic code is what mutates the substrate losslessly and repeatably; without it, the graph the LLM "holds" is a fog of its own making, and the next turn's read isn't the same artifact as the last turn's write. The disciplines don't erode as models scale — the fuzziness is in the architecture, not the capability ceiling.

The complexity ceiling argument applies at the upper end. Most production hybrid loops are 4–7 blocks with one or two feedback paths; for those, the LLM holding the graph isn't the architect's bottleneck. The benefit kicks in at the larger projects (12+ blocks, three or more feedback paths, multiple stacked dev-time loops) — which used to be the privilege of senior architects and are now reachable by less-experienced builders. The systems we can build are now larger than the systems we can entirely comprehend; the systems available to non-experts are now structurally richer than what they could previously author alone.

## Why the framework has to exist

The pattern resolves a tension. The new block type's three affordances — soft-input parsing, prose-as-program, multi-modal generation — pull toward hybrid-loop shapes wherever non-determinism is tolerable. Four layers of conditioning pull the other way:

1. *Training data.* Pre-2023, ~all software engineering looked like conventional engineering — pipelines of typed transforms that return a result. Hybrid-loops-style architectures barely existed in the corpus because the affordances didn't exist. The LLM's priors pull toward the shape it has the most evidence for.
2. *Harness.* Claude Code, Cursor, Copilot all chose familiar IDE/terminal interfaces — slash commands, file edits, REPL sessions. The choice was correct for adoption, but the harness is itself context-as-code that conditions the LLM toward IDE-shaped work.
3. *User expectation.* People interact with LLMs the way they interacted with their tools — write code, run tests, fix bugs. The harness doesn't surface graph-of-blocks as a primitive interaction.
4. *Ecosystem fragmentation.* The agent-framework ecosystem has been swirling for two years without a consistent vocabulary. DSPy has *modules* and *signatures*; LangGraph has *nodes* and *edges*; AutoGen has *agents* and *conversations*; pydantic has *models*. Each tool covers one cell with its own primitives. A mid-level engineer who picks any one of them locks into that cell without seeing the rest of the alphabet. (See `references/AGENT_FRAMEWORKS.md` for the per-tool comparison.)

The conventional shape is gravitational. The skill is a piece of context-as-code that pushes back. It's not primarily teaching the user (though it does); it's *counter-conditioning the LLM* — pulling it out of "what would my training data do here" and into "what does this surface actually want."

Patterns don't propagate primarily by being explained. They propagate by being proposed in moments where someone has to choose. A senior engineer who's been building pipelines for twenty years doesn't read a hybrid-loops doc and reorganize their thinking. They encounter the LLM proposing a hybrid-loops shape for their next project, push back or accept, and either way they now have a pattern in their working vocabulary they didn't have before.

That's the theory of change. The skill counter-conditions the default; the catalog of named block-graphs in `BLOCK_GRAPHS.md` makes proposals concrete enough that humans say yes; code written in this shape enters public corpora; future generations train on it; the gravitational pull rebalances. The success metric isn't "everyone reads SKILL.md." It's *"LLMs propose hybrid-loops shapes when relevant, and users encounter the proposal as a real option."*

## What's not new, and what's still open

The framework isn't claiming novelty in:

- The graph algebra (von Neumann, Unix pipes, every modular-software text since)
- The five-role decomposition (recovered from blackboards, scripts/frames, and Soar)
- The bootstrap-loop pattern (AlphaGo 2016, DreamCoder 2021)
- "Generalization shaping" as a design principle (wesen)
- The small-typed-tools aesthetic (Devine Lu Linvega, Hundred Rabbits)
- Pattern-language structure (Christopher Alexander, 1977)

Everything else is recovered prior art with a coat of LLMs on top. The recovery isn't the contribution; the disciplines for safely placing the new block type in the recovered architecture are. Whether those disciplines materially change development decisions across the kinds of projects where the pattern fits, and whether they generalize beyond the engineering surfaces this writeup was sketched against, are open empirical questions the writeup does not pretend to have answered.

## Lit review still wanted

The lit review currently in `PRIOR_ART.md` covers blackboards, scripts/frames, AlphaGo, DreamCoder, OpenCog, the knowledge-acquisition bottleneck, and the practitioner work directly informing this writeup. There's a broader sweep worth doing: architectural ideas from the last 50 years that required cheap typed-structure authoring to work, where authoring cost was the limiting factor. Candidates include argumentation frameworks (Dung 1995), behavior trees (game AI), Petri nets, multi-agent systems (Wooldridge / Jennings), inductive logic programming (Muggleton 1991), genetic programming / program synthesis (Koza onward), constraint logic programming, semantic-web / topic-maps, structured wikis, BPMN / workflow engines, actor-model concurrent systems, cybernetics (Wiener / Ashby). Most have the algebra right and were limited by the same bottleneck the knowledge-acquisition critique named. LLMs lift the bottleneck.
