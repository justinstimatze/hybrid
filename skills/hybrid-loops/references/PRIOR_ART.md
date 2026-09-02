# Prior art for hybrid loops

Cite these when defending the architecture. Four tiers below:

1. *Directly informed the design* — central for any defense
2. *Contemporaneous practitioner work and 2024-2026 ecosystem* — adjacent active work; nobody arriving at this space from inside one of these tools should be surprised to discover the framework
3. *Cite to distinguish* — same architecture, wrong bet (currently just OpenCog)
4. *Further reading and lineage* — older traditions, cybernetics, conceptual orientation

If you only want the essential references, the first tier is enough. The second tier is what you need to position the framework against the ecosystem a 2026 reader is already inside.

---

## Tier 1 — directly informed the design

Citations that directly shaped how this repo describes the pattern, the primitives it ships, or the architectural choices behind both.

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

In his [`go-go-workshop`](https://github.com/go-go-golems/go-go-workshop) materials, wesen notes that he does not use agents and zero-shot prompting for most of his use cases, and points readers toward the Blackboard System (Hayes-Roth 1985) as a more useful conceptual frame than "agents." Two practitioners arriving independently at the same architectural lineage from different starting points is meaningful evidence about the lineage itself; cite this when the question is whether hybrid loops are well-grounded in classical AI architectures. (Reading wesen's framing as "independent corroboration" is this repo's interpretation — wesen has not made this claim about other practitioners.)

#### Complementarity with this work

Wesen's public contributions concentrate on engineering-side infrastructure for typed LLM workflows — typed-step frameworks, session-streaming substrates, evidence databases, prompt-context libraries, document managers. This repo documents the pattern itself and ships a Claude Code skill that helps reach for it. Applied artifacts — engineering and non-engineering tools that instantiate the pattern — live in their own repositories. The two bodies of work are complementary rather than competitive — both are concrete instances of the same architectural pattern at different layers of the stack.

One area neither body of work has yet shipped as a standalone primitive (as of April 2026) is a *calibration / prediction-logging layer* — a tool that closes the loop between an evaluator's intended judgment and the eventual outcome it can be checked against. Minitrace, bucheron, and the diary essay each gesture at parts of this; nothing assembles them into a per-evaluator hit-rate primitive other projects can drop in.

For wesen's own manifesto on the design philosophy of his ecosystem, see ["I want my software to be visionary — the go-go-golems ecosystem"](https://the.scapegoat.dev/i-want-my-software-to-be-visionary-the-go-go-golems-ecosystem/). Notable principles: rich data representation (applications preserve the structural knowledge embedded in their data rather than reducing everything to printf-style output), discoverability, relentless refinement (willingness to break APIs to maintain coherent vision). Quoted: *"The only way I know to properly identify what these concepts are about is to turn them into working code."*

### Aesthetic and craft lineage — Devine Lu Linvega / Hundred Rabbits

Wesen has cited Devine Lu Linvega ([100r.co](https://100r.co), Hundred Rabbits) as a personal influence on his sensibility. Devine builds small, opinionated, typed software tools — Orca (live-coded sequencer), Left (text editor), Dotgrid (vector tool), Ronin (image processing), uxn (a small virtual machine in the permacomputing tradition) — that prioritize craft, ownership, locality, and minimalism. None of this work is LLM-augmented; none of it has to be. Devine's aesthetic is what hybrid loops aspire to *for the deterministic-shell half* of the pattern.

Cite Devine when defending design choices around: small tool size, single-purpose primitives, typed I/O between tools, permacomputing / locality (compute stays on the user's machine and isn't a cloud service), and the deliberate rejection of platform-scale frameworks in favor of assemblies of focused tools. The Hundred Rabbits collective (Devine + Rek Bell), the uxn ecosystem, and the Merveilles network more broadly are the canonical references for *a personal collection of typed tools the user actually owns*.

### Pattern languages — Christopher Alexander

Alexander, Christopher. *A Pattern Language: Towns, Buildings, Construction*. 1977. Companion volume: *The Timeless Way of Building*. 1979.

The right structural reference for *what hybrid loops is, as a unit of design*. A pattern in Alexander's sense has a recurring problem, a context where it applies, a solution structure, and named consequences for downstream patterns. Hybrid loops is a pattern in this strict sense; the five roles plus meta-layers form a small pattern language with internal nesting (a substrate pattern, a gate pattern, a calibration pattern).

When writing for an audience that includes designers (not just engineers), Alexander's framing lands more cleanly than AI-engineering vocabulary. Cite *A Pattern Language* for the structural argument; *The Timeless Way* for the philosophical one (the "wholeness" thesis distinguishing living pattern languages from catalogs of tricks). The standard software adaptation — Gamma, Helm, Johnson, Vlissides's *Design Patterns* (1994) — preserves Alexander's *structure* but not his *sensibility*; reading Alexander directly is the thing.

### AlphaGo / AlphaZero

Silver, Huang, Maddison, et al. *Mastering the game of Go with deep neural networks and tree search*. Nature, 2016.
Silver, Schrittwieser, Simonyan, et al. *Mastering the game of Go without human knowledge*. Nature, 2017.

Architectural template for hybrid loops. Policy network (fuzzy/learned) proposes moves; Monte Carlo Tree Search (hard/symbolic) explores and validates; MCTS outputs become training data for the policy. Mutual bootstrapping — neither does well alone, together is superhuman.

Difference from hybrid loops as the term is used here: AlphaGo's structural prior (rules of Go, board) is fixed. The user's pattern operates over a structural prior an earlier LLM call generated. That's the defining novelty.

### DreamCoder

Ellis, Wong, Nye, Sablé-Meyer, Morales, Hewitt, Cary, Solar-Lezama, Tenenbaum. *DreamCoder: Bootstrapping inductive program synthesis with wake-sleep library learning*. Nature Communications, 2021. arXiv:2006.08381.

Closest direct lineage. Wake phase (compose library functions to solve tasks) + abstraction sleep (extract recurring patterns into new library functions) + dream sleep (sample from library to generate synthetic training data for a recognition model). Iterates to bootstrap a domain-specific language from a small primitive set.

Maps onto: this repo's *metabolism* → DreamCoder's wake/sleep; the compress+verify loop → wake + abstraction; schema discovery → library learning by MDL. DreamCoder limitations to acknowledge: pre-LLM, works in toy domains, library compression can collapse to golf-y abstractions.

### LILO

Grand, Wong, Bowers, Olausson, Liu, Tenenbaum, Andreas. *LILO: Learning Interpretable Libraries by Compressing and Documenting Code*. NeurIPS 2024. arXiv:2310.19791.

LLM-era DreamCoder descendant. Closest published cognate to the framework's compress-and-verify approach to notation discovery.

### Voyager

Wang, Xie, Jiang, Mandlekar, Xiao, Zhu, Fan, Anandkumar. *Voyager: An Open-Ended Embodied Agent with Large Language Models*. arXiv:2305.16291. 2023.

Skill library learning for Minecraft agents. LLM proposes new skills; successful ones enter the library, available for future tasks. Direct DreamCoder descendant in agent context. Demonstrates hybrid loops outside program synthesis.

### Anthropic — Building Effective Agents (Dec 2024)

[Anthropic's product team's blog post](https://www.anthropic.com/engineering/building-effective-agents) on agentic patterns. Names a four-tier hierarchy: augmented LLM → workflow → agent → multi-agent. The "augmented LLM" base case is the lens-block-with-tool-use shape; the "workflow" tier is the canonical hybrid-loop runtime cycle; the agent tier is sub-loop-with-its-own-graph. The closest official-Anthropic alignment with the framework's vocabulary. What the framework adds is naming the disciplines (calibration, context-as-code, dev-time loop) Anthropic's post leaves implicit, and applying the pattern beyond engineering use cases.

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

Maps onto: REASONS Canvas → context-as-code as core infrastructure (highest-impact flavor); prompt-first vs code-first refactor → operational rule for the dev-time loop; "Reject chat-and-drift" → the calibration / discipline argument.

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

- **Compound engineering (Every.to / Kieran Klaassen)** — practitioner methodology for AI-assisted dev; current articulation is a seven-step loop **Ideate → Brainstorm → Plan → Work → Review → Polish → Compound**. The final step embeds learnings into searchable artifacts so subsequent work is easier; structurally the dev-time hybrid loop wrapping the runtime. See `AGENT_FRAMEWORKS.md` for honest overlap-and-gap treatment; eight-beliefs / five-stages framing reads as consultancy packaging the framework declines to adopt.
- **Harness engineering (OpenAI / Ryan Lopopolo, 2026)** — platform vendor's engineering-practice writeup for building production systems with agents-as-engineers. Thesis: *"design environments, specify intent, and build feedback loops."* AGENTS.md as table of contents + recurring "garbage collection" by cleanup agents are the two distinctive disciplines. See `AGENT_FRAMEWORKS.md` for full mapping; convergent vocabulary from a non-Anthropic source, partially offsets the Anthropic-Karpathy platform-vendor consolidation noted in `ORCHESTRATION_SHAPES.md` source-tiering.
- **Rule of Five (Jeffrey Emanuel, named and operationalized by Yegge in *Welcome to Gas Town*, Jan 1 2026, https://steve-yegge.medium.com/welcome-to-gas-town-4f25ee16dd04)** — folk self-review heuristic: *"if you make an LLM review something five times, with different focus areas each time though, it generates superior outcomes and artifacts."* Yegge implements it as a formula that wraps any workflow molecule so each step is reviewed 4 times (the implementation counts as the first review). Adjacent to the adversarial-panel-process shape but distinct: panel reviews are by multiple independent reviewers, Rule of Five is sequential self-review by the same agent under rotating focus areas. Folk-register, not a framework.
- **Loop engineering (Addy Osmani, *Loop Engineering*, June 7 2026, [addyosmani.com/blog/loop-engineering](https://addyosmani.com/blog/loop-engineering))** — popular-press articulation of the level-shift past prompt engineering: *"Loop engineering is replacing yourself as the person who prompts the agent. You design the system that does it instead."* Catalogs the five primitives a loop needs (Automations, Worktrees, Skills, Plugins/Connectors, Sub-agents) plus a memory file, all now shipping inside 2026 platform-vendor tooling (Codex, Claude Code); names the *maker/checker* split as the central structural move ("the model that wrote the code is way too nice grading its own homework"). The Skills section says intent written outside the conversation "kind of compounds" across sessions — this framework's *context-as-code* discipline (`THE_CASE.md`). The framework's distinction is making the *calibration* and *dev-time-loop* disciplines explicit, treating the guard as a deterministic typed check (compile-time mesh-checkability in stull) rather than an LLM-grading-LLM checker, and adding the *discoverable vs ambient* axis — every primitive Osmani catalogs is human-invoked (`/goal`, cron, triage inbox), not substrate-hook-fired the way an ambient hybrid loop is. Cite as a 2026 cousin on the discoverable side of the framework; pairs naturally with Compound Engineering and Harness Engineering.

### The Phoenix Architecture — Fowler

Chad Fowler, *The Phoenix Architecture* ([aicoding.leaflet.pub](https://aicoding.leaflet.pub/)). An ongoing essay series (Dec 2025–Jun 2026) arguing that when generation is cheap, durability comes from *regeneration, not preservation*: treat implementations as disposable and keep intent / spec / architecture as the permanent artifact. The thesis sits *beside* hybrid loops rather than inside it — Fowler's central claim is code disposability, which the framework takes no position on — but several essays reach the framework's core disciplines from the code-architecture vantage. The mappings below are suggestive cousins: the value is that a 2026 reader steeped in this series should recognize the rhymes.

Strongest cross-references:

- **["The Gradient of Trust"](https://aicoding.leaflet.pub/3mb2qb6odxc2d)** (subtitle *"better shapes beat better prompts"*) — *constraints as trust* (a strong type system, purity, explicit effects "dramatically shrink the space of possible mistakes") rhymes with the framework's deterministic-layer argument: structure, not prompt-craft, is what makes a fuzzy actor's output trustworthy. The trust-gradient itself — some code trusted on sight, some only after review, some never — is the developer intuition that per-block **calibration** formalizes into a persistent hit-rate. Closest of the series to the framework's actual thesis.
- **["Production Is a Compiler Input"](https://aicoding.leaflet.pub/3mjx4erlboc2l)** — production telemetry becomes an input to what the system generates next, not merely a human-debugging aid; names *evidence decay / technical drift* (a component "can satisfy the spec today and fail it three months from now even if nobody touches the code" — the world changed). Maps onto the framework's two meta-layers — **calibration** (does the lens still hold?) and **metabolism** (is the accumulated record still valid?) — and onto the cycle-not-pipeline feedback edge that separates the pattern from a one-shot codegen pipeline.
- **["The Generative Stack"](https://aicoding.leaflet.pub/3miwhqqvwxc2x)** — spec inputs → canonicalized clauses → requirements/invariants → evaluations → implementation units → code → feedback; *"independent tools and representations coexist at every phase, composable"* is a codegen-with-verification pipeline (`BLOCK_GRAPHS.md`) carrying the same no-single-framework-bet argument as `STACKING.md` / `THE_CASE.md`.
- **["Provenance Is the New Version Control"](https://aicoding.leaflet.pub/3mcbiyal7jc2y)** + **"The Conversation Is the Commit"** — the unit of change is reasons: capture the decisions / tradeoffs / rejected approaches because AI generation "severs" the diff-as-decision-record link and breaks the "docs are an optional tax" tradeoff. Both are the **context-as-code** discipline from the version-control angle: intent, not implementation, is the durable artifact.
- **"The Phoenix Primitives"** + **"Compile to Architecture"** — the file/framework stops being the primitive and the spec/architecture becomes one; *"the file is a cognitive container … the module boundary reflects what a developer can reason about in a single sitting,"* so file layout falls out of the design rather than shaping it upfront. This is `THE_CASE.md`'s cognitive-load argument from the other side: once the LLM holds the graph, the human-sized containers (7±2 files, modules) stop being structural.

Where it diverges / tensions worth keeping: the disposability thesis is orthogonal to hybrid loops, and **"The Implementation Remembers"** is Fowler's own counterweight — mature implementations encode undocumented scar tissue (oddly specific timeouts, defensive checks) that pure regeneration can erase; the honest limit on "code is disposable," and a cousin of the framework's metabolism-audit caution. Lower-relevance entries: **"Regenerative Software"** (foundational framing — *"the limiting factor is no longer writing software, but understanding, evaluating, and governing it"*); **"Pace Layers and AI Integration"** (Stewart Brand's pace layering as a fast-experiments / slow-stabilizes altitude split — cousin to the layering intuition in `STACKING.md`); **"The Regenerative Grain"** (small = "safe to delete"; verification as the scarce resource); **"The Deletion Test"** (`rm -rf src/` as a regenerability probe — diagnostic heuristic); **"UI Is a Conservation Layer"** (bounded replacement behind stable interfaces).

*Naming caution:* unrelated to **Arize Phoenix** (the eval/observability tool listed under LLM observability platforms below) and to the **Phoenix** web framework — "Phoenix Architecture" here is Fowler's regeneration thesis.

### LLM observability and calibration platforms

Production-scale implementations of the calibration discipline named in `THE_CASE.md`. Teams running hybrid loops in production would reach for one of these rather than rolling their own append-only JSONL hit-rate logger.

- **Braintrust** ([braintrust.dev](https://www.braintrust.dev/)) — eval + tracing + regression suites for LLM apps.
- **Langfuse** ([langfuse.com](https://langfuse.com/)) — open-source LLM observability + evals.
- **Langsmith** (LangChain) — evals + traces + datasets, tightly integrated with LangChain ecosystem.
- **Weights & Biases (Weave / Traces)** — extension of W&B's experiment tracking into LLM observability.
- **Arize Phoenix** ([phoenix.arize.com](https://phoenix.arize.com/)) — open-source LLM evaluation + monitoring.
- **Helicone** ([helicone.ai](https://www.helicone.ai/)) — LLM gateway + observability proxy.
- **PromptLayer** ([promptlayer.com](https://www.promptlayer.com/)) — prompt versioning + observability + eval.

These tools are calibration-first; they also cover dataset management, regression detection, prompt versioning, multi-metric eval, and per-cohort A/B comparison. They don't have opinions on graph design, substrate-as-vocabulary, or decline-when. Complementary to the framework, not competitive — and any minimal calibration-logger sketch you might write yourself is only a starter for what these tools provide.

### Safety / evals research

Empirical work whose claims bear on the framework's positioning. Listed here for discourse-awareness, not as central evidence the framework leans on. This sub-section is intentionally narrow at present — entries earn inclusion by being widely-cited AND by the framework being able to honestly characterize their limitations.

- **METR — Task-Completion Time Horizons of Frontier AI Models** ([metr.org](https://metr.org/time-horizons/); original paper [March 19 2025](https://metr.org/blog/2025-03-19-measuring-ai-ability-to-complete-long-tasks/); v1.1 update [Jan 29 2026](https://metr.org/blog/2026-1-29-time-horizon-1-1/); self-issued limitations note [Jan 22 2026](https://metr.org/notes/2026-01-22-time-horizon-limitations/)). Independent autonomy-evaluation nonprofit measuring the 50%-task-completion time horizon — the time a human typically takes on tasks the model completes with 50% success. Original headline: doubling every ~7 months (212 days); v1.1 since-2023 figure is 131 days, since-2024 is 89 days. **Cite with same-breath caveats, not standalone.** METR's own limitations note concedes: ~2x error bars in each direction, *"improving performance on 20% horizon tasks can lower 80% horizon"* (artifact of the two-parameter logistic fit), human baselines measured for only 5 of 31 long tasks (remainder estimated), and explicitly: *"Speculating about the effects of a months- or years-long time horizon is fraught."* Wegner's independent follow-up ([*Are AI time-horizons (still) doubling every 7 months?*, Mar 11 2026](https://medium.com/@AIchats/are-ai-time-horizons-still-doubling-every-7-months-6262ed2bcc6a)) argues the apparent acceleration is methodological — logarithmic binning sensitivity, sophisticated scaffolding given disproportionately to newer models, false precision. The framework declines to lean on the headline number; METR is included because (a) it's the most-cited public benchmark on long-horizon agent capability, (b) METR publishes its own limitations note, and (c) the user of this framework will encounter "doubling every N months" talking-points in the discourse and should know what's underneath them.

### Multi-agent orchestration projects

The exemplars below are the practitioner projects the four shapes in `ORCHESTRATION_SHAPES.md` were named from; arriving at the framework from inside any of them, a reader should recognize the same architectural disciplines under different vocabulary.

- **Orca** ([onorca.dev](https://onorca.dev); Stably AI, YC-backed; [github.com/stablyai/orca](https://github.com/stablyai/orca) — TypeScript, MIT, 26.5k stars July 2026) — *shape 2* exemplar demonstrating the discipline is substrate-independent. Positions itself as an *"ADE"* (its contrast with IDEs; *"An ADE is built for you and your agents"*) for *"working with a fleet of parallel agents"*: Claude Code, Codex, Grok, Cursor, Amp, Cline, Goose et al. run side by side in isolated git worktrees on desktop (macOS / Windows / Linux) with Ghostty-inspired terminals, plus iOS / Android companion and VPS-remote execution. Distinctive versus Conductor below: same human-as-attester discipline surviving multi-CLI, cross-OS, and remote-substrate (mobile + VPS) broadening — evidence the discipline isn't tied to any single tool stack. Its *fleet* framing is also where shape 2 empirically strains toward shape 3: at fan-out that exceeds one human's review capacity, the Simon / Charity "trust-account" migration trigger fires. *Not to be confused with Devine Lu Linvega's Orca live-coded sequencer above in the `wesen` section — different tool.*

- **Conductor** (Charlie, co-founder; YC S24; macOS desktop app for orchestrating coding agents; practitioner walkthrough at https://www.youtube.com/watch?v=fQmlML9Lay4) — narrower-substrate exemplar of *shape 2 — human-orchestrated parallel multiplexer* (single-CLI Claude Code, macOS-only, no mobile / no VPS): the shape's earlier articulation before Orca above showed the discipline is substrate-independent. The human keeps planning and approval; agents run in parallel git worktrees on sliced-up tasks. Charlie's central folk rule *"don't let the AI be your architect"* rhymes with the substantive/procedural split (the mapping is loosely suggestive). "Slot free zones" / "do not touch if you are an AI" comments are a cousin to bounded action (`PROCEDURE.md` §4) — both constrain by exclusion, different protocols. *Distinct from Netflix Conductor*, the workflow orchestrator listed under "Adjacent ecosystems" below.

- **Gas Town** ([github.com/gastownhall/gastown](https://github.com/gastownhall/gastown), Steve Yegge, Go, v1.0 April 2026) — multi-agent coordination workspace with persistent state, git-backed worktrees ("Hooks"), three-tier watchdog system (Witness/Deacon/Dogs), targets coordinating 20-30 agents. Solves "agents lose context on restart" with durable state — same problem Temporal / Netflix Conductor solve at workflow scale, with agent-specific abstractions and a "town" metaphor (Mayor / Rigs / Crews / Polecats / Convoys / Beads). Earliest deployed exemplar of *shape 3 — engineered-resilience autonomous, with per-pack adversarial verification*.

- **Gas City** (Yegge; SDK v1.0 April 24, 2026; primary implementor Julian Knutsen) — Gas Town's successor framework for building custom agent packs. Introduces MEOW (Molecular Expression of Work) with Formulas (reusable templates) and Molecules (instances), version-controlled in Dolt and forkable across an org. Codifies the *two-or-three agents per pack* deployment rule (adversarial verification as deployment policy, not guideline). Yegge: *"Reliability, friends, is a dial. You choose where to set it. More rounds of review, more backstops, more guardrails, more judges..."* See `ORCHESTRATION_SHAPES.md` shape 3 for the per-protocol mapping.

- **Wasteland** (Yegge et al.; *Welcome to the Wasteland: A Thousand Gas Towns*, March 2026; built by Julian Knutsen + Matt Beane + others with Yegge's vision) — federated network of Gas Towns sharing a Wanted Board and trust-graded validation atop Dolt (Git-versioned SQL). Stamps as multi-dimensional verdicts on an append-only ledger; trust levels gate what new rigs can do; the *"yearbook rule"* forbids self-stamping. Earliest documented exemplar of *shape 4 — socially-validated federation*. The protocol-overlap with `PROCEDURE.md` (standing, recusal, substrate-as-record) is suggestive — three of five protocols rhyme — but the mapping rests on Yegge's essay and has not been validated by independent inspection of the running system.

- **Claude Code dynamic workflows** (Anthropic — Thariq Shihipar & Sid Bidasaria, *"A harness for every task: dynamic workflows in Claude Code,"* blog, June 2 2026; https://claude.com/blog/a-harness-for-every-task-dynamic-workflows-in-claude-code — direct WebFetch returned 403 on 2026-07-01; full text confirmed same day via user-supplied copy) — the orchestrating model writes its own orchestration script (a plain JS file using `agent()` / `parallel()` / `pipeline()` / `phase()` primitives, plus ordinary JS built-ins) rather than following a pre-built graph; a deterministic runtime then executes that authored script, spawning and coordinating the resulting subagent fleet (per-run concurrency and total-agent ceilings are set in the shipped tool's own primitive spec, not stated in the blog itself). The blog names six composable patterns — *classify-and-act*, *fan-out-and-synthesize*, *adversarial verification*, *generate-and-filter*, *tournament*, *loop-until-done* — that are a shipped, named vocabulary for compositions this framework's `BUILDING_BLOCKS.md` already catalogs abstractly: classify-and-act ≡ the `LLM-classify + code-dispatch` pair listed there; fan-out-and-synthesize ≡ the branching/joining operators named in that file's closing section, with the synthesis step explicitly named as "a barrier — it waits for all the fan-out agents." The blog's stated reason for isolating subagents into separate context windows is three named failure modes — *agentic laziness* (declaring partial progress done on a multi-part task), *self-preferential bias* (an LLM over-trusting its own prior output when asked to verify it), *goal drift* (constraints lost to lossy compaction over a long single-context run) — which is independent, practitioner-facing evidence for this framework's calibration job (b): a fluent, self-consistent-but-wrong output a rolling hit-rate never catches, only grounding against independent re-derivation does. Reported case study: Jarred Sumner's Bun Zig→Rust port — one subagent per fix/module in its own worktree, "another agent" adversarially reviewing before merge, then a fix loop driving build + test suite clean (fuller quantitative claims — exact reviewer count, pass-rate — live in Sumner's own X thread, cited by the blog but not independently checked here). The blog's own "dynamic vs. static workflows" distinction — a static graph "needs to work for all edge cases" and is "usually more generic," a dynamic workflow authored per-task — restates, in the platform vendor's own words, the distinction `SKILL.md` already draws between a fixed pipeline of MCP tools and runtime dynamic dispatch. See `ORCHESTRATION_SHAPES.md` shape 1 (the "third path" update) and `STACKING.md` §"recursive harness authoring" for where this sits relative to the four shapes and the tier-stacking regime; see `AGENT_FRAMEWORKS.md` for the discipline-coverage comparison against DSPy / LangGraph / AutoGen / Ralph Wiggum and the rest of the implementation-toolkit ecosystem.

- **Software Survival 3.0** (Yegge, Jan 2026) — the theoretical framing under the Yegge stack. The Survival Ratio `Survival(T) ∝ (Savings × Usage × H) / (Awareness_cost + Friction_cost)` and its six levers (Insight Compression, Substrate Efficiency, Broad Utility, Publicity, Friction Minimization, Human Coefficient) give a different organizing principle from hybrid-loops' discipline-pattern decomposition but reach overlapping conclusions about tool-block fitness. Lever 1 (Insight Compression) ≈ substrate-as-crystallized-knowledge; Lever 2 (Substrate Efficiency) ≈ the gate. Cited from `BUILDING_BLOCKS.md` as non-self-referential evidence the *discipline of tools for agents* has economic stakes.

- **Maggie Appleton — *Gas Town's Agent Patterns, Design Bottlenecks, and Vibecoding at Scale*** ([maggieappleton.com/gastown](https://maggieappleton.com/gastown), Jan 2026) — the most useful adversarial vantage on Yegge's claims. Frames Gas Town as *"speculative design fiction"* rather than a present-day shippable tool. The tempering reading is what keeps shape 3's framework treatment honest: the most-elaborated articulation of what engineered-resilience autonomy *could* look like, not ablation-validated infrastructure that's currently easy to adopt. Also names *"design becomes the limiting factor"* once execution gets cheap — the framework's discipline is what makes design-as-bottleneck tractable.

- **Ringer** ([github.com/NateBJones-Projects/ringer](https://github.com/NateBJones-Projects/ringer), Nate B Jones, Python, created July 2026) — parallel cheap-worker swarm where the expensive model plans + reviews and each worker's artifact is verified by a deterministically executed shell check; per the README, *"exit 0 is the only thing Ringer believes."* Failures retry once with the failure context injected; every attempt logged. Distinct from Gas Town / Gas City above in checker discipline: Yegge's stack uses LLM-adversarial-verification packs; ringer uses an executed oracle — same shape-3-adjacent fan-out shape, different checker. The model-performance log slices pass-rate by `(model, task_type)` — a *worker-selection scoreboard* (evidence-based routing, promotion ladder), sibling to the per-evaluator calibration primitive Tier 1's wesen *Complementarity with this work* subsection flagged as not yet shipped; different axis — per-worker track record for routing, not per-evaluator prediction vs outcome. Shallow single-gate / one-retry loop, not a rich multi-block composition. Third-party arrival at the deterministic-checker discipline from outside the SF/blog cluster; cite for the executed-gate and per-worker-per-task routing scoreboard, not as a full shape-3 exemplar.

### Enterprise-SaaS productizations of the loop primitive

Non-cluster industry-adoption signal: SaaS products shipping *loops* as a first-class user-authored primitive, not just an internal orchestration pattern. Distinct from workflow-orchestration engines (Temporal, Netflix Conductor) because the loop is authored in natural language by a domain user, not by an engineer in a graph editor; distinct from multi-agent orchestration frameworks above because the target user is an issue-tracker admin, distinct from the infrastructure builders the frameworks above serve.

- **Linear Loops** ([linear.app/docs/loops](https://linear.app/docs/loops); Linear; Business + Enterprise plans, AI-credit billing effective July 20 2026) — enterprise-SaaS productization of the discoverable-loop primitive inside an issue tracker. Loops trigger on a schedule OR when issues match a set of conditions; instructions in natural language; optional tools for external services (Slack, Code Intelligence, coding sessions); permissions scoped per team or workspace. Example loop verbatim from docs: *"When an issue enters this team's triage queue, investigate its likely root cause with Code Intelligence. If you think you can fix it, start a coding session."* Cite as the first-party non-cluster industry-adoption signal for the shape Osmani cataloged as *"Automations"* — verifies the discoverable-vs-ambient meta-loop distinction at a SaaS-primitive layer well outside the SF/blog-circle consolidation.

### Enterprise workflow, evaluation, and automation adjacencies (2026)

Enterprise tools engaging with the hybrid-loop pattern from adjacent angles: workflow-orchestration engines with AI-agent features, platform-vendor posts on maker-checker calibration, LLM-ops evaluation platforms, and low-code / SaaS-integration automation. Positive-corroboration bullets come first, arranged as deterministic-gating side → platform-vendor argument for calibrated-LLM-judge tractability → concrete productizations of the calibrated-LLM-judge pattern (whose fit with the framework's discipline remains the deferred question flagged in the Academic subsection below); the negative citation (n8n Validator) is last.

- **Camunda 8 Agentic Orchestration — *Guardrails and Best Practices for Agentic Orchestration*** ([camunda.com/blog/2026/01/guardrails-and-best-practices-for-agentic-orchestration](https://camunda.com/blog/2026/01/guardrails-and-best-practices-for-agentic-orchestration/), Jan 2026) — enterprise workflow-orchestration engine articulating the same deterministic-gating discipline the framework calls for, at BPMN-activity granularity. The post's central rule: *"Separate deterministic from dynamic logic"* — *"Use BPMN and DMN to capture deterministic logic and decisions visually,"* let agents handle the dynamic segments, and *"hand the process back to deterministic flow as fast as possible."* Positive corroboration from the workflow-orchestration ecosystem the framework already cites (Temporal / Netflix Conductor / Airflow) — the same architectural move enacted at process granularity rather than block-in-loop granularity. The framework's additional typed-record substrate discipline (what the LLM writes into) has no direct analog here; the BPMN-activity mapping covers the gating side only.

- **Anthropic — *Harness design for long-running application development*** ([anthropic.com/engineering/harness-design-long-running-apps](https://www.anthropic.com/engineering/harness-design-long-running-apps), Prithvi Rajasekaran / Anthropic Labs, Mar 24 2026) — platform-vendor engineering post detailing a planner + generator + evaluator three-agent architecture for autonomous coding sessions. Defends the LLM-judge substrate on tractability grounds: *"tuning a standalone evaluator to be skeptical turns out to be far more tractable than making a generator critical of its own work, and once that external feedback exists, the generator has something concrete to iterate against."* Names the calibration mechanism explicitly: *"I calibrated the evaluator using few-shot examples with detailed score breakdowns. This ensured the evaluator's judgment aligned with my preferences, and reduced score drift across iterations."* Platform-vendor self-report tier; central to the deferred task-class question flagged in the Academic subsection. The source's own argument is *tractability* — evaluator-side skepticism is easier to build than self-critical generation. The framework reads this as its strongest platform-vendor argument for admitting calibrated-LLM-judges into the fuzzy-quality tier; that reading belongs to this framework — Rajasekaran's own post doesn't make it.

- **Cursor — *Governing agent autonomy with Auto-review*** ([cursor.com/blog/agent-autonomy-auto-review](https://cursor.com/blog/agent-autonomy-auto-review), Cursor Research, Jun 11 2026) — dev-tools platform vendor productizing the maker-checker split with an LLM-classifier checker calibrated on a labeled corpus. Target metric is *"flapping"* — the source names the failure mode verbatim: *"If the same case allowed six times and blocked four times, that usually meant the policy or prompt was underspecified."* The per-block calibration primitive the framework's Tier 1 wesen *Complementarity with this work* subsection flagged as not yet shipped as a drop-in per-evaluator hit-rate primitive, now shipped: 6,122 labeled rows deduplicated from internal developer sessions, plus synthetic high-risk cases. Written policy governs stakes-based leniency: *"more lenient when the security stakes are lower, and more cautious when they're higher."* Cite as the clearest *in-loop* productization of the calibration primitive the framework has been arguing for since Tier 1 — the ambient counterpart to LangSmith Align Evals's discoverable batch-UI shape below.

- **LangSmith Align Evals — *How to Calibrate LLM-as-Judge with Human Corrections*** ([langchain.com/resources/llm-as-a-judge](https://www.langchain.com/resources/llm-as-a-judge), Mar 10 2026) — LangChain's productization of per-evaluator calibration, and the cleanest live example of the *discoverable* meta-loop shape at the evaluation layer. Align Evals *"replaces prompt-guessing with a measurable loop: collect human corrections, build few-shot examples, and track agreement over time"*; treats prompt iteration as insufficient (*"Prompt iteration alone won't close the gap between a technically correct evaluator and a reliable one"*) and requires *"systematic alignment to human corrections."* Same underlying substrate move as the framework's per-block calibration discipline, but shipped as a *discoverable product surface* (a batch UI you visit) rather than the ambient in-loop discipline that is the framework's separate contribution. The value here is one clean exemplar of one side of the discoverable-vs-ambient cut in mainstream tooling.

- **n8n Validator Agent pattern** (community-documented at [chronexa.io/blog/n8n-ai-agent-node-build-multi-agent-systems-in-2026](https://chronexa.io/blog/n8n-ai-agent-node-build-multi-agent-systems-in-2026), 2026) — *negative citation.* In the low-code / SaaS-integration ecosystem's canonical multi-agent architecture (*"Research Agent → Action Agent → Validator Pattern"*), the *"Validator Agent (The 'Quality Gate')"* is another LLM checking the first — the LLM-grading-LLM pattern the framework's deterministic-gating discipline currently rejects. Documented positioning: agents run *reasoning loops* while *"The Chain node is deterministic—it runs once and stops"* — the same deterministic-vs-agentic split as Camunda above, but the guard is LLM-authored with no calibration mechanism documented in the source (contrast Cursor and Anthropic above, where the LLM-judge substrate is explicitly calibrated). Distinct from n8n's presence in the *Adjacent ecosystems* list below (which names the tool generically); this bullet names the specific Validator-Agent architectural pattern the community teaches. See intro note: whether this rejection holds for all task classes is the open question flagged in the Academic subsection.

### Academic prior art (2026 contemporaneous)

Recent (2026) academic work that arrived at the same or adjacent architectural moves independently, and forms the framework's academic neighborhood. Papers below did not inform the framework's design (all postdate the framework's core writing) but strengthen or reframe specific claims. AgentLTL and Lean-Agent Protocol sit on the deterministic-guard side of the framework's typed-substrate discipline; SemaClaw and Yue et al. contribute at the taxonomy / harness-articulation layer. Whether the framework's typed-substrate discipline admits a task-class carve-out is an open question these citations frame but do not resolve; that amendment, if it happens, belongs in `THE_CASE.md`, not this section's intro.

- **Dynamic Runtime Graphs survey — Yue et al., *From Static Templates to Dynamic Runtime Graphs: A Survey of Workflow Optimization for LLM Agents*** (arXiv:2603.22386, Mar 23 2026) — best-single 2026 academic survey to sit alongside `BLOCK_GRAPHS.md`. Treats LLM-agent workflows as *"agentic computation graphs (ACGs)"* and organizes the literature along three axes: *when* structure is determined (static / dynamic), *what part* is optimized, and *which evaluation signals guide optimization* (task metrics, verifier signals, preferences, or trace-derived feedback). Separates *reusable workflow templates* from *run-specific realized graphs* from *execution traces*. Cite when someone wants an academic vocabulary for positioning hybrid loops within the workflow-optimization literature.

- **AgentLTL — Elkoussy & Perez (EPITA/LRE), *A Trace-Verification Framework for Measuring, Enforcing, and Training Procedural Compliance in Tool-Using LLM Agents*** (arXiv:2607.02599, Jul 1 2026) — academic articulation of the framework's alphabet claim for the strict-verifiability task class. A First-Order Linear Temporal Logic (FO-LTL) language over agent traces yields *"a deterministic, judge-free compliance score,"* used dually as an *online harness gate* (blocking tool calls whose trace prefix violates the spec) and as a *reward signal for finetuning*. Reports block-and-warn harnessing improving compliance on 5 of 7 evaluated models; finetuning with the same spec yields +38 / +17.5 percentage-point gains in accuracy / compliance. The single-spec-drives-both-runtime-and-training pattern parallels the framework's dev-time-loop-wrapping-runtime discipline; the deterministic-judge-free gate is the framework's typed-substrate discipline in academic form.

- **Lean-Agent Protocol — Rashie & Rashi, *Type-Checked Compliance: Deterministic Guardrails for Agentic Financial Systems Using Lean 4 Theorem Proving*** (arXiv:2604.01483, Apr 1 2026) — formal-methods extreme of the deterministic-guard end of the substrate spectrum. Each proposed agentic action is *"treated as a mathematical conjecture: execution is permitted if and only if the Lean 4 kernel proves that the action satisfies pre-compiled regulatory axioms."* Institutional policies auto-formalized into Lean 4 via Harmonic AI's Aristotle neural-symbolic model; targeted at SEC 15c3-5, FINRA 3110, OCC 2011-12, and CFPB compliance. Cite as the reference point for "how far can the alphabet go" — the framework's typed-guard discipline taken to its provable-correctness limit for a high-stakes-verifiable domain.

- **SemaClaw — Zhu et al. (Midea AIRC), *A Step Towards General-Purpose Personal AI Agents through Harness Engineering*** (arXiv:2604.11548, Apr 13 2026; source at [github.com/midea-ai/SemaClaw](https://github.com/midea-ai/SemaClaw)) — a Midea AIRC paper citing OpenClaw (see Personal AI section below) as the population-scale inflection point when *"millions of users began deploying personal AI agents into their daily lives"* — an academic touchpoint for the harness-engineering discipline alongside the framework's Ryan Lopopolo *Harness Engineering* citation (Tier 1). Contributions include a DAG-based two-phase hybrid agent team orchestration method, PermissionBridge behavioral safety, three-tier context management, and an agentic wiki skill for automated personal knowledge base construction.

### Personal AI / local-first

- **OpenClaw** ([github.com/openclaw/openclaw](https://github.com/openclaw/openclaw)) — local-first self-hosted personal AI assistant. Gateway control plane routing across messaging surfaces (WhatsApp, Telegram, Slack, Discord). Emphasizes "always-on / local / fast" personal automation over multi-agent orchestration. Sits in the *deployment shape* corner (substrate-on-user's-device, single-user) rather than the architecture corner. See also SemaClaw (Academic prior art above) — a Midea AIRC paper citing OpenClaw as its motivating case.

### Books, guides, and pattern catalogs

- **Chip Huyen, *AI Engineering* (O'Reilly, 2024)** — broad textbook for LLM-application engineering.
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
- *Workflow orchestration*: Temporal, Netflix Conductor, AWS Step Functions, Airflow
- *Visual LLM-app builders*: Dify, LangFlow, Flowise
- *Low-code / SaaS-integration automation*: n8n, Zapier, Make
- *Structured-output / typed I/O tools*: pydantic, instructor, Anthropic tool use, OpenAI structured outputs

---

## Tier 3 — cite to distinguish

Same architecture, different bet. Useful for showing the lineage and where this work explicitly disagrees.

### OpenCog / Hyperon (Goertzel et al.)

Cite to *distinguish*, not to align. Goertzel's patternist architecture (AtomSpace + PLN + MOSES + ECAN) had the right architectural intuition — typed substrate that metabolizes — and the wrong bet. Tried to do symbolic *reasoning* (PLN) when statistical learning was about to dominate. Failed for the bitter-lesson reason.

Hybrid loops inverts OpenCog's bet: keep the typed substrate, let LLMs do the reasoning. Same architecture, different targets, finally tractable. Worth claiming the lineage. The bet, though, is worth distinguishing.

---

## A note on naming

This repo uses "hybrid loops" as the working name for the pattern. The broader field has no settled name; adjacent terms with partial coverage include "compound AI systems" (Zaharia et al., BAIR 2024), "generalization shaping" (wesen), "schemaed cognition" (this repo, earlier draft, retired), "structured introspection" (informal). Citing the pattern by *any* of these names is fine.

---

## Tier 4 — further reading and lineage

Loosely related citations for orienting readers from adjacent fields. Not essential to any defense of the architecture; cite when the audience comes from these traditions and benefits from the pointer.

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

Cite when the architecture-recovery claim needs grounding. The framework's typed substrate descends from the 1970s frame-and-script tradition (Schank's scripts and the KL-ONE / structured-semantic-network family). Its cycle structure traces to Soar's production-system + working-memory architecture (filtered through 50 years of cost-structure changes). Already implicit in the Hayes-Roth / Buchanan & Feigenbaum citations in Tier 1; named here for completeness when the audience knows Soar specifically.

### Burroughs / Gysin: The Third Mind

Burroughs and Gysin. *The Third Mind*. 1978.

Cite when discussing the social/team version of hybrid loops. The third mind was the emergent entity from two minds collaborating; a team-shared substrate with periodic metabolism phases becomes that emergent entity in the AI era. The agency criterion is the defining distinguisher between "passive store" (not a third mind) and "third mind proper."

### Engelbart: Augmenting Human Intellect

Engelbart, Douglas. *Augmenting Human Intellect: A Conceptual Framework*. 1962.

Cite when discussing collective IQ and shared external substrate. Engelbart's vision of structured shared artifacts as collective-intelligence amplifier never fully shipped because the substrate was too expensive to build and maintain. LLMs as the substrate-authoring layer change that cost structure. Team-shared collective-IQ deployments are closer to Engelbart's vision than to Burroughs's.

### Active inference / predictive coding (Friston et al.)

Friston, Karl. *The free-energy principle: a unified brain theory?*. Nature Reviews Neuroscience, 2010.

Loosely relevant. Hybrid loops have a flavor of bidirectional inference (top-down predictions constrain bottom-up perception, and vice versa). Don't lean on this hard — the formal connection is thin — but it's a useful pointer for readers from cognitive science.

---

