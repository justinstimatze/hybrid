# Prior art for hybrid loops

Cite these when defending the architecture. Four tiers below:

1. *Directly informed the design* — load-bearing for any defense
2. *Contemporaneous practitioner work and 2024-2026 ecosystem* — adjacent active work; nobody arriving at this space from inside one of these tools should be surprised to discover the framework
3. *Cite to distinguish* — same architecture, wrong bet (currently just OpenCog)
4. *Further reading and lineage* — older traditions, cybernetics, conceptual orientation

If you only want the load-bearing references, the first tier is enough. The second tier is what you need to position the framework against the ecosystem a 2026 reader is already inside.

---

## Tier 1 — directly informed the design

Citations that meaningfully shaped how this repo describes the pattern, the primitives it ships, or the architectural choices behind both.

### Practitioner prior art — Manuel Odendahl ("wesen")

A software developer (open-source author of the [go-go-golems](https://github.com/go-go-golems) toolchain, blogger at [the.scapegoat.dev](https://the.scapegoat.dev)) who has been working in this design space for several years and is one of the clearest writers on it. His public work is the most important practitioner reference for this pattern; his terminology and tooling deserve direct citation in any writeup of hybrid loops.

#### Theoretical framing he has named

***Generalization shaping*** — the design move of *restructuring a problem with notation, tools, and typed interfaces so the LLM does only the in-distribution mapping work and deterministic machinery carries correctness*. Essay: ["Tool use and notation as shaping LLM generalization"](https://the.scapegoat.dev/tool-use-and-notation-as-generalization-shaping/) (Feb 2026). Quoted: *"Tools don't make cognition deeper — they make the world simple in exactly the places we need it to be."*

Generalization shaping is best understood as a *design principle inside hybrid loops* — corresponding to the gate role plus the lens schema design — not a synonym for the whole pattern. Hybrid loops adds the typed substrate, calibration log, metabolism, and recursive composition on top of generalization shaping at the boundary. When defending why a hybrid loop's gate carries the load it does, cite this principle.

#### Vocabulary he has introduced or made canonical

- **diary** — narrative memory artifact, deliberately chosen over "ledger" / "log." See ["Why I Make My Agents Keep Diaries"](https://the.scapegoat.dev/why-i-make-my-agents-keep-diaries/) for the argument that the word itself activates LLM behaviors he wants.
- **evidence database** — the SQLite typed-record store agent runs leave behind. From [`wesen/2026-04-29--go-go-agent`](https://github.com/wesen/2026-04-29--go-go-agent).
- **substrate** — used in the [`go-go-golems/sessionstream`](https://github.com/go-go-golems/sessionstream) README for the typed event-streaming layer; this repo's use is consistent with his.
- **step** — the unit of typed LLM operation in [`go-go-golems/geppetto`](https://github.com/go-go-golems/geppetto). Each step is a typed function from flags+args to structured records.
- **spray test** — empirical variance probe of a prompt (regenerate N times, measure variance). From ["From prompt and pray to prompt engineering"](https://the.scapegoat.dev/from-prompt-and-pray-to-prompt-engineering/) (Apr 2026). Calibration-adjacent.
- **mapping** / **interface-mapping** — wesen's vocabulary for what an LLM does at the system-design level. Direct usage in the "Tool use and notation" essay. At higher abstraction this surfaces in `THE_CASE.md` as *fuzzy pattern mappers* paired with *deterministic pattern mappers* (compilers, transpilers, linters, codegen) as the sibling species compiler-veterans already know.

When using any of these terms, attribute to wesen explicitly.

#### Most relevant repositories

- [**geppetto**](https://github.com/go-go-golems/geppetto) — Go LLM framework built around the typed-step abstraction. Underpins much of his stack.
- [**pinocchio**](https://github.com/go-go-golems/pinocchio) — CLI/REPL frontend; YAML-based prompt-library-with-metadata.
- [**prompto**](https://github.com/go-go-golems/prompto) and [**promptos**](https://github.com/go-go-golems/promptos) — prompt-context library with metadata; scans configured repos for `prompto/` directories and treats files (and executables) as named, retrievable contexts.
- [**go-go-agent**](https://github.com/wesen/2026-04-29--go-go-agent) — terminal agent with an explicit evidence database for replay/inspection. The closest direct parallel in his work to a hybrid loop with calibration-style provenance.
- [**sessionstream**](https://github.com/go-go-golems/sessionstream) — recently extracted (April 2026) generic typed event-streaming substrate, lifted out of pinocchio's evtstream.
- [**minitrace**](https://github.com/wesen/minitrace) and [**go-minitrace**](https://github.com/go-go-golems/go-minitrace) — common JSON trace format unifying multiple agent session formats; query with DuckDB. Upstream is `fukami/minitrace`; wesen maintains a fork and Go port.
- [**docmgr**](https://github.com/go-go-golems/docmgr) — structured document manager for LLM-assisted workflows; PKM with LLM-aware metadata, frontmatter conventions, vocabulary management, code↔doc relations.
- [**Codex-Reflect-Skill**](https://github.com/wesen/Codex-Reflect-Skill) — runs Codex in parallel over past Codex sessions to surface patterns and propose new skills.
- [**bucheron**](https://github.com/go-go-golems/bucheron) — structured-log upload service for client-side bug reporting.
- [**glazed**](https://github.com/go-go-golems/glazed) — foundational typed-rows-and-columns library underpinning the stack.

#### Architectural framing — Blackboard Systems

In his [`go-go-workshop`](https://github.com/go-go-golems/go-go-workshop) materials, wesen notes that he does not use agents and zero-shot prompting for most of his use cases, and points readers toward the Blackboard System (Hayes-Roth 1985) as a more useful conceptual frame than "agents." Two practitioners arriving independently at the same architectural lineage from different starting points is meaningful evidence about the lineage itself; cite this when the question is whether hybrid loops are well-grounded in classical AI architectures. (Reading wesen's framing as "independent corroboration" is this repo's framing, not a claim wesen has made about other practitioners.)

#### Complementarity with this work

Wesen's public contributions concentrate on engineering-side infrastructure for typed LLM workflows — typed-step frameworks, session-streaming substrates, evidence databases, prompt-context libraries, document managers. This repo documents the pattern itself and ships a Claude Code skill that helps reach for it; applied artifacts (engineering and non-engineering tools that instantiate the pattern) live in their own repositories. The two bodies of work are complementary rather than competitive — both are concrete instances of the same architectural pattern at different layers of the stack.

One area neither body of work has yet shipped as a standalone primitive (as of April 2026) is a *calibration / prediction-logging layer* — a tool that closes the loop between an evaluator's intended judgment and the eventual outcome it can be checked against. Minitrace, bucheron, and the diary essay each gesture at parts of this; nothing assembles them into a per-evaluator hit-rate primitive other projects can drop in.

For wesen's own manifesto on the design philosophy of his ecosystem, see ["I want my software to be visionary — the go-go-golems ecosystem"](https://the.scapegoat.dev/i-want-my-software-to-be-visionary-the-go-go-golems-ecosystem/). Notable principles: rich data representation (applications preserve the structural knowledge embedded in their data rather than reducing everything to printf-style output), discoverability, relentless refinement (willingness to break APIs to maintain coherent vision). Quoted: *"The only way I know to properly identify what these concepts are about is to turn them into working code."*

### Aesthetic and craft lineage — Devine Lu Linvega / Hundred Rabbits

Wesen has cited Devine Lu Linvega ([100r.co](https://100r.co), Hundred Rabbits) as a personal influence on his sensibility. Devine builds small, opinionated, typed software tools — Orca (live-coded sequencer), Left (text editor), Dotgrid (vector tool), Ronin (image processing), uxn (a small virtual machine in the permacomputing tradition) — that prioritize craft, ownership, locality, and minimalism. None of this work is LLM-augmented; none of it has to be. Devine's aesthetic is what hybrid loops aspire to *for the deterministic-shell half* of the pattern.

Cite Devine when defending design choices around: small tool size, single-purpose primitives, typed I/O between tools, permacomputing / locality (substrate stays on the user's machine and isn't a cloud service), and the deliberate rejection of platform-scale frameworks in favor of assemblies of focused tools. The Hundred Rabbits collective (Devine + Rek Bell), the uxn ecosystem, and the Merveilles network more broadly are the canonical references for *a personal collection of typed tools the user actually owns*.

### Pattern languages — Christopher Alexander

Alexander, Christopher. *A Pattern Language: Towns, Buildings, Construction*. 1977. Companion volume: *The Timeless Way of Building*. 1979.

The right structural reference for *what hybrid loops is, as a unit of design*. A pattern in Alexander's sense has a recurring problem, a context where it applies, a solution structure, and named consequences for downstream patterns. Hybrid loops is a pattern in this strict sense; the five roles plus meta-layers form a small pattern language with internal nesting (a substrate pattern, a gate pattern, a calibration pattern).

When writing for an audience that includes designers (not just engineers), Alexander's framing lands more cleanly than AI-engineering vocabulary. Cite *A Pattern Language* for the structural argument; *The Timeless Way* for the philosophical one (the "wholeness" thesis distinguishing living pattern languages from catalogs of tricks). The standard software adaptation — Gamma, Helm, Johnson, Vlissides's *Design Patterns* (1994) — preserves Alexander's *structure* but not his *sensibility*; reading Alexander directly is the thing.

### AlphaGo / AlphaZero

Silver, Huang, Maddison, et al. *Mastering the game of Go with deep neural networks and tree search*. Nature, 2016.
Silver, Schrittwieser, Simonyan, et al. *Mastering the game of Go without human knowledge*. Nature, 2017.

Architectural template for hybrid loops. Policy network (fuzzy/learned) proposes moves; Monte Carlo Tree Search (hard/symbolic) explores and validates; MCTS outputs become training data for the policy. Mutual bootstrapping — neither does well alone, together is superhuman.

Difference from hybrid loops as the term is used here: AlphaGo's structural prior (rules of Go, board) is fixed. The user's pattern operates over a structural prior an earlier LLM call generated. That's the load-bearing novelty.

### DreamCoder

Ellis, Wong, Nye, Sablé-Meyer, Morales, Hewitt, Cary, Solar-Lezama, Tenenbaum. *DreamCoder: Bootstrapping inductive program synthesis with wake-sleep library learning*. Nature Communications, 2021. arXiv:2006.08381.

Closest direct lineage. Wake phase (compose library functions to solve tasks) + abstraction sleep (extract recurring patterns into new library functions) + dream sleep (sample from library to generate synthetic training data for a recognition model). Iterates to bootstrap a domain-specific language from a small primitive set.

Maps onto: this repo's *metabolism* → DreamCoder's wake/sleep; the compress+verify loop → wake + abstraction; schema discovery → library learning by MDL. DreamCoder limitations to acknowledge: pre-LLM, works in toy domains, library compression can collapse to golf-y abstractions.

### LILO

Grand, Wong, Bowers, Olausson, Liu, Tenenbaum, Andreas. *LILO: Learning Interpretable Libraries by Compressing and Documenting Code*. NeurIPS 2024. arXiv:2310.19791.

LLM-era DreamCoder descendant. Closest published cognate to the framework's compress-and-verify approach to notation discovery.

### Voyager

Wang, Xie, Jiang, Mandlekar, Xiao, Zhu, Fan, Anandkumar. *Voyager: An Open-Ended Embodied Agent with Large Language Models*. arXiv:2305.16291. 2023.

Skill library learning for Minecraft agents. LLM proposes new skills; successful skills enter library; library available for future tasks. Direct DreamCoder descendant in agent context. Demonstrates hybrid loops outside program synthesis.

### Anthropic — Building Effective Agents (Dec 2024)

[Anthropic's product team's blog post](https://www.anthropic.com/engineering/building-effective-agents) on agentic patterns. Names a four-tier hierarchy: augmented LLM → workflow → agent → multi-agent. The "augmented LLM" base case is the lens-block-with-tool-use shape; the "workflow" tier is the canonical hybrid-loop runtime cycle; the agent tier is sub-loop-with-its-own-graph. The closest official-Anthropic alignment with the framework's vocabulary; the framework's distinction is naming the disciplines (calibration, context-as-code, dev-time loop) Anthropic's post leaves implicit, and applying the pattern beyond engineering use cases.

### CoALA — Cognitive Architectures for Language Agents (Sumers et al., NeurIPS 2024)

Sumers, Yao, Narasimhan, Griffiths. *Cognitive Architectures for Language Agents*. arXiv:2309.02427. NeurIPS 2024.

The most-direct academic taxonomy of language-agent architectures. Maps memory / actions / decision-making onto a Soar-descended cognitive-architecture frame. Treats LLMs as decision-making policies inside a typed agent shell. Closest academic-literature analog to hybrid loops as a system-design pattern; the framework's contribution beyond CoALA is the explicit dev-time-loop discipline and the substrate-as-vocabulary vs substrate-as-record distinction.

### DSPy — Khattab et al. (arXiv 2310.03714, 2023)

Khattab et al. *DSPy: Compiling Declarative Language Model Calls into Self-Improving Pipelines*. arXiv:2310.03714.

The academic foundation for typed-signature LM programming with optimizers. Closest published cousin of the framework's compress-and-verify shape in spirit (different metric: prompt/demo optimization instead of roundtrip score). See `AGENT_FRAMEWORKS.md` for the per-tool comparison.

### Compound AI Systems — Zaharia et al. (BAIR, 2024)

Zaharia et al. *The Shift from Models to Compound AI Systems*. BAIR Blog, Feb 2024. [link](https://bair.berkeley.edu/blog/2024/02/18/compound-ai-systems/).

Names "compound AI systems" as the umbrella for what hybrid loops sits inside. Identifies that production LLM applications increasingly look like *systems* (multiple components, control logic, retrieval, tools) rather than single-model calls. The framework agrees on the umbrella; "hybrid loops" is one specific shape within it that the BAIR post doesn't fully articulate (cycles, mutual generation, dev-time loops, calibration discipline).

### Structured Prompt-Driven Development — Fowler / openspdd

Patel, Sharif, Fowler. ["Structured Prompt-Driven Development with the REASONS Canvas"](https://martinfowler.com/articles/structured-prompt-driven/). martinfowler.com.

The most-aligned practitioner methodology in the 2026 literature. Treats prompts as "first-class delivery artifacts" version-controlled alongside code; defines the REASONS Canvas (Requirements, Entities, Approach, Structure, Operations, Norms, Safeguards) as a typed prompt-spec; enforces "fix the prompt first, then update the code" discipline; provides `openspdd` CLI to automate the workflow.

Maps onto: REASONS Canvas → context-as-code as load-bearing infrastructure (highest-leverage flavor); prompt-first vs code-first refactor → operational rule for the dev-time loop; "Reject chat-and-drift" → the calibration / discipline argument.

What hybrid loops adds beyond SPDD: explicit calibration discipline (SPDD's "alignment checkpoints" stop short of persistent hit-rate); the broader pattern beyond engineering work; the deterministic-vs-fuzzy actor framing. Cite SPDD prominently as the closest engineering-discipline cousin in current practitioner literature.

### Knowledge-acquisition bottleneck

Buchanan and Feigenbaum. *Rule-based expert systems: the MYCIN experiments of the Stanford Heuristic Programming Project.* Addison-Wesley, 1984.
Hayes-Roth, Waterman, Lenat (eds). *Building Expert Systems.* Addison-Wesley, 1983.
Lenat. *CYC: A Large-Scale Investment in Knowledge Infrastructure.* Communications of the ACM, 1995.

Cite when explaining *why* the 1970s frames-and-rules tradition (Schank's scripts, KL-ONE, MYCIN, XCON, Cyc) didn't scale despite having the architecture mostly right. Buchanan & Feigenbaum named the *knowledge-acquisition bottleneck* — the rate-limiting step was knowledge engineers extracting and encoding domain knowledge into formal representations, which scaled poorly. Cyc was the most ambitious and sustained attempt to overcome it through brute force; Lenat's 1995 paper documents the multi-decade investment and the partial nature of progress. The bottleneck didn't go away; it ended the era.

LLMs change the cost structure on the surfaces that killed expert systems. World knowledge that Cyc tried to author by hand is pre-loaded; schema iteration that took knowledge engineers months can take hours with structured-outputs and an evaluation loop. Other defects mattered too — rules didn't compose at scale, edge cases proliferated faster than they could be patched, the world stayed messier than first-order logic — and LLMs cover those as well: the same fuzziness that makes them unreliable as the only actor is what lets them tolerate messiness pure rules couldn't survive. The architecture was right; the costs *and* the brittleness made it uneconomic, and both have lifted.

---

## Tier 2 — contemporaneous practitioner work and 2024-2026 ecosystem

Adjacent active work covering pieces of the same broader pattern. None of these *are* hybrid loops as a unified design pattern, but each occupies one or more cells of the alphabet and a 2026 reader should recognize them.

### Methodologies

- **Compound engineering (Every.to / Kieran Klaassen)** — practitioner methodology for AI-assisted dev with a Plan→Work→Review→**Compound**→Repeat loop. The compound step is structurally the dev-time hybrid loop. See `AGENT_FRAMEWORKS.md` for honest overlap-and-gap treatment; eight-beliefs / five-stages framing reads as consultancy packaging the framework declines to adopt.

### LLM observability and calibration platforms

Production-scale implementations of the calibration discipline named in `THE_CASE.md`. Teams running hybrid loops in production would reach for one of these rather than rolling their own append-only JSONL hit-rate logger.

- **Braintrust** ([braintrust.dev](https://www.braintrust.dev/)) — eval + tracing + regression suites for LLM apps.
- **Langfuse** ([langfuse.com](https://langfuse.com/)) — open-source LLM observability + evals.
- **Langsmith** (LangChain) — evals + traces + datasets, tightly integrated with LangChain ecosystem.
- **Weights & Biases (Weave / Traces)** — extension of W&B's experiment tracking into LLM observability.
- **Arize Phoenix** ([phoenix.arize.com](https://phoenix.arize.com/)) — open-source LLM evaluation + monitoring.
- **Helicone** ([helicone.ai](https://www.helicone.ai/)) — LLM gateway + observability proxy.
- **PromptLayer** ([promptlayer.com](https://www.promptlayer.com/)) — prompt versioning + observability + eval.

These tools are calibration-first; they also cover dataset management, regression detection, prompt versioning, multi-metric eval, and per-cohort A/B comparison. They don't have opinions on graph design, substrate-as-vocabulary, or decline-when. Complementary to the framework, not competitive — and any minimal calibration-logger sketch you might write yourself is a starter, not a substitute.

### Multi-agent orchestration projects

- **Gas Town** ([github.com/gastownhall/gastown](https://github.com/gastownhall/gastown)) — multi-agent coordination workspace with persistent state, git-backed worktrees ("Hooks"), three-tier watchdog system (Witness/Deacon/Dogs), targets coordinating 20-30 agents. Solves "agents lose context on restart" with durable state — same problem Temporal/Conductor solve at workflow scale, with agent-specific abstractions and a "town" metaphor (Mayor / Rigs / Crews / Polecats / Convoys / Beads).

### Personal AI / local-first

- **OpenClaw** ([github.com/openclaw/openclaw](https://github.com/openclaw/openclaw)) — local-first self-hosted personal AI assistant. Gateway control plane routing across messaging surfaces (WhatsApp, Telegram, Slack, Discord). Emphasizes "always-on / local / fast" personal automation over multi-agent orchestration. Sits in the *deployment shape* corner (substrate-on-user's-device, single-user) rather than the architecture corner.

### Books, guides, and pattern catalogs

- **Chip Huyen, *AI Engineering* (O'Reilly, 2024)** — comprehensive textbook for LLM-application engineering.
- **Eugene Yan, ["Patterns for Building LLM-based Systems" (eugeneyan.com, 2024)](https://eugeneyan.com/writing/llm-patterns/)** — explicit catalogue of LLM application patterns (evals, RAG, fine-tuning, caching, guardrails, defensive UX). Adjacent to `BUILDING_BLOCKS.md` at the per-block-pattern level.
- **John Berryman, *Relevant Search* / RAG-adjacent writing** — search and retrieval patterns useful for the substrate-as-record shape.

### Cultural-register practitioners

Useful for tone/onboarding context; not framework-shaping but in the conversation a 2026 engineer is likely already part of.

- **Andrej Karpathy** — "Software 3.0" framing; coined "vibe coding" (2024-2025). The cultural reference for "LLMs as a new computational substrate."
- **Steve Yegge** — "Cheating is All You Need" / various AI-coding writings.
- **swyx (Shawn Wang) / Latent Space podcast** — extensive ecosystem coverage; the practitioner-conversation venue.
- **Simon Willison** ([simonwillison.net](https://simonwillison.net/)) — patterns: tools, structured outputs, prompt injection, llm-CLI work.
- **Jason Liu / instructor** ([jxnl.co](https://jxnl.co/)) — structured-outputs-with-pydantic discipline; one of the clearest writers on typed-LLM-output engineering.

### Adjacent ecosystems (deeper comparisons in `AGENT_FRAMEWORKS.md`)

- *Agent frameworks*: DSPy (also Tier 1 academically), LangGraph, AutoGen, CrewAI
- *Workflow orchestration*: Temporal, Conductor, AWS Step Functions, Airflow
- *Visual LLM-app builders*: Dify, LangFlow, Flowise
- *Low-code / SaaS-integration automation*: n8n, Zapier, Make
- *Structured-output / typed I/O tools*: pydantic, instructor, Anthropic tool use, OpenAI structured outputs

---

## Tier 3 — cite to distinguish

Same architecture, different bet. Useful for showing the lineage and where this work explicitly disagrees.

### OpenCog / Hyperon (Goertzel et al.)

Cite to *distinguish*, not to align. Goertzel's patternist architecture (AtomSpace + PLN + MOSES + ECAN) had the right architectural intuition — typed substrate that metabolizes — and the wrong bet. Tried to do symbolic *reasoning* (PLN) when statistical learning was about to dominate. Failed for the bitter-lesson reason.

Hybrid loops inverts OpenCog's bet: keep the typed substrate, let LLMs do the reasoning. Same architecture, different targets, finally tractable. Worth claiming the lineage; worth distinguishing the bet.

---

## A note on naming

This repo uses "hybrid loops" as the working name for the pattern. The broader field has no settled name; adjacent terms with partial coverage include "compound AI systems" (Zaharia et al., BAIR 2024), "generalization shaping" (wesen), "schemaed cognition" (this repo, earlier draft, retired), "structured introspection" (informal). Citing the pattern by *any* of these names is fine.

---

## Tier 4 — further reading and lineage

Loosely related citations for orienting readers from adjacent fields. Not load-bearing for any defense of the architecture; cite when the audience comes from these traditions and benefits from the pointer.

### Soft computing — Zadeh's umbrella term

Zadeh, Lotfi. *Fuzzy logic, neural networks, and soft computing*. Communications of the ACM, 1994.

Coined "soft computing" as the umbrella for fuzzy logic + neural networks + evolutionary computation (GAs, genetic programming) + probabilistic reasoning (Bayesian networks, HMMs). The genre shares the shape: *use computation to handle uncertainty, search, or optimization, with hand-authored typed components*. Each method ran into a flavor of the knowledge-acquisition bottleneck (membership functions and rule sets for fuzzy; fitness functions for GAs; graph and priors for Bayesian; states and transitions for HMMs). The word "fuzzy" in *fuzzy pattern mapper* (used throughout `THE_CASE.md`) is borrowed from this tradition.

LLMs lift the same authoring bottleneck for all of these. The 90s soft-computing toolkit becomes a library of deterministic actors that LLMs can now author into hybrid-loop graphs — fuzzy controllers with LLM-designed membership functions; GAs with LLM-written fitness functions; Bayesian nets with LLM-proposed graph structures.

Specific citations worth keeping handy:
- **John Holland**, *Adaptation in Natural and Artificial Systems* (MIT Press, 1975) — genetic algorithms.
- **John Koza**, *Genetic Programming* (MIT Press, 1992) — direct ancestor of compress-and-verify-style search over typed programs.
- **Lawrence Rabiner**, *A Tutorial on Hidden Markov Models* (Proc. IEEE, 1989) — canonical HMM reference.
- **Judea Pearl**, *Probabilistic Reasoning in Intelligent Systems* (Morgan Kaufmann, 1988) — Bayesian network foundations.

### Cybernetics, autopoiesis, and self-producing systems

Wiener, Norbert. *Cybernetics: Or Control and Communication in the Animal and the Machine*. MIT Press, 1948.
Ashby, W. Ross. *Design for a Brain*. Chapman & Hall, 1952.
Maturana, Humberto and Varela, Francisco. *Autopoiesis and Cognition: The Realization of the Living*. Reidel, 1980.

Cite when the audience comes from cybernetics or systems theory. The framework's mutually-generative cycles read as second-order cybernetics restated for the LLM era; *autopoiesis* (a system that produces the components it's constituted of) is the closest conceptual ancestor of the mutual-generation claim. The framework wasn't derived from this lineage — it was reasoned from the LLM's affordances — but a careful reader will recognize the family resemblance, and the cybernetics tradition gets the conceptual credit it deserves even though it didn't directly inform the design.

### Soar, scripts, and the classical architecture lineage

Newell, Allen, Laird, John, Rosenbloom, Paul. *Soar: An Architecture for General Intelligence*. Artificial Intelligence, 1987.
Schank, Roger and Abelson, Robert. *Scripts, Plans, Goals and Understanding*. Lawrence Erlbaum, 1977.
Brachman, Ronald. *What's in a Concept: Structural Foundations for Semantic Networks*. International Journal of Man-Machine Studies, 1977 (and the broader KL-ONE family that followed).

Cite when the architecture-recovery claim needs grounding. The framework's typed substrate is descended from the 1970s frame-and-script tradition (Schank's scripts and the KL-ONE / structured-semantic-network family); the cycle structure is descended from Soar's production-system + working-memory architecture (filtered through 50 years of cost-structure changes). Already implicit in the Hayes-Roth / Buchanan & Feigenbaum citations in Tier 1; named here for completeness when the audience knows Soar specifically.

### Burroughs / Gysin: The Third Mind

Burroughs and Gysin. *The Third Mind*. 1978.

Cite when discussing the social/team version of hybrid loops. The third mind was the emergent entity from two minds collaborating; a team-shared substrate with periodic metabolism phases becomes that emergent entity in the AI era. The agency criterion is the load-bearing distinguisher between "passive store" (not a third mind) and "third mind proper."

### Engelbart: Augmenting Human Intellect

Engelbart, Douglas. *Augmenting Human Intellect: A Conceptual Framework*. 1962.

Cite when discussing collective IQ and shared external substrate. Engelbart's vision of structured shared artifacts as collective-intelligence amplifier never fully shipped because the substrate was too expensive to build and maintain. LLMs as the substrate-authoring layer change that cost structure. Team-shared collective-IQ deployments are closer to Engelbart's vision than to Burroughs's.

### Active inference / predictive coding (Friston et al.)

Friston, Karl. *The free-energy principle: a unified brain theory?*. Nature Reviews Neuroscience, 2010.

Loosely relevant. Hybrid loops have a flavor of bidirectional inference (top-down predictions constrain bottom-up perception, and vice versa). Don't lean on this hard — the formal connection is thin — but it's a useful pointer for readers from cognitive science.

---

