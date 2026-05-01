# Prior art for hybrid loops

Cite these when defending the architecture. Three tiers below: **directly informed the design**, **cite to distinguish**, and **further reading / orientation**. If you only want the load-bearing references, the first tier is enough.

---

## Tier 1 — directly informed the design

These citations meaningfully shaped how this repo describes the pattern, the primitives it ships, or the architectural choices behind both. Reach for them in any defense of the design.

### Practitioner prior art — Manuel Odendahl ("wesen")

Manuel Odendahl is a software developer (open-source author of the [go-go-golems](https://github.com/go-go-golems) toolchain, blogger at [the.scapegoat.dev](https://the.scapegoat.dev)) who has been working in this design space for several years and is one of the clearest writers on it. His public work is the most important practitioner reference for this pattern; his terminology and tooling deserve direct citation in any writeup of hybrid loops.

### Theoretical framing he has named

**"Generalization shaping"** — the design move of *restructuring a problem with notation, tools, and typed interfaces so the LLM does only the in-distribution mapping work and deterministic machinery carries correctness*. Essay: ["Tool use and notation as shaping LLM generalization"](https://the.scapegoat.dev/tool-use-and-notation-as-generalization-shaping/) (Feb 2026). Quoted: *"Tools don't make cognition deeper — they make the world simple in exactly the places we need it to be."*

Generalization shaping is best understood as a **design principle inside hybrid loops** (corresponding to the gate role + the lens schema design — the deterministic-machinery side of the loop), not a synonym for the whole pattern. Hybrid loops as used in this repo add the typed substrate, calibration log, metabolism, and recursive composition on top of generalization-shaping at the boundary. When defending why a hybrid loop's gate carries the load it does, cite this principle and credit the framing to him.

### Vocabulary he has introduced or made canonical

- **diary** — narrative memory artifact, deliberately chosen over "ledger" / "log." See ["Why I Make My Agents Keep Diaries"](https://the.scapegoat.dev/why-i-make-my-agents-keep-diaries/) for his argument that the word "diary" itself activates LLM behaviors he wants.
- **evidence database** — the SQLite typed-record store agent runs leave behind. From [`wesen/2026-04-29--go-go-agent`](https://github.com/wesen/2026-04-29--go-go-agent).
- **substrate** — he uses this term in the [`go-go-golems/sessionstream`](https://github.com/go-go-golems/sessionstream) README for the typed event-streaming layer; this repo's use of "substrate" is consistent with his.
- **step** — the unit of typed LLM operation in [`go-go-golems/geppetto`](https://github.com/go-go-golems/geppetto). Each step is a typed function from flags+args to structured records.
- **spray test** — empirical variance probe of a prompt (regenerate N times, measure variance). From ["From prompt and pray to prompt engineering"](https://the.scapegoat.dev/from-prompt-and-pray-to-prompt-engineering/) (Apr 2026). Calibration-adjacent.
- **mapping** / **interface-mapping** — wesen's vocabulary for what an LLM does at the system-design level (an LLM maps a structured input to a structured output, often through an interface designed to constrain the mapping). Direct usage in ["Tool use and notation as shaping LLM generalization"](https://the.scapegoat.dev/tool-use-and-notation-as-generalization-shaping/) and adjacent essays. At higher abstraction this surfaces in `references/THE_CASE.md` as **"LLMs are fuzzy pattern mappers"** — paired with **deterministic pattern mappers** (compilers, transpilers, linters, codegen) as the sibling species compiler-veterans already know. "Pattern mapper" is preferred over "pattern transformer" specifically to avoid the Vaswani-2017 attention-architecture collision; the framing is at the systems-design level, not the architecture level.

When using any of these terms, attribute to wesen explicitly.

### Most relevant repositories

- [**geppetto**](https://github.com/go-go-golems/geppetto) — Go LLM framework built around the typed-step abstraction. Underpins much of his stack.
- [**pinocchio**](https://github.com/go-go-golems/pinocchio) — CLI/REPL frontend; YAML-based prompt-library-with-metadata.
- [**prompto**](https://github.com/go-go-golems/prompto) and [**promptos**](https://github.com/go-go-golems/promptos) — prompt-context library with metadata; scans configured repos for `prompto/` directories and treats files (and executables) as named, retrievable contexts.
- [**go-go-agent**](https://github.com/wesen/2026-04-29--go-go-agent) — terminal agent with an explicit evidence database for replay/inspection. The closest direct parallel in his work to a hybrid loop with calibration-style provenance.
- [**sessionstream**](https://github.com/go-go-golems/sessionstream) — recently extracted (April 2026) generic typed event-streaming substrate, lifted out of pinocchio's evtstream.
- [**minitrace**](https://github.com/wesen/minitrace) and [**go-minitrace**](https://github.com/go-go-golems/go-minitrace) — common JSON trace format unifying multiple agent session formats; query with DuckDB. (Note: upstream is `fukami/minitrace`; wesen maintains a fork and Go port.)
- [**docmgr**](https://github.com/go-go-golems/docmgr) — structured document manager for LLM-assisted workflows; PKM with LLM-aware metadata, frontmatter conventions, vocabulary management, code↔doc relations, health checks.
- [**Codex-Reflect-Skill**](https://github.com/wesen/Codex-Reflect-Skill) — runs Codex in parallel over past Codex sessions to surface patterns and propose new skills.
- [**bucheron**](https://github.com/go-go-golems/bucheron) — structured-log upload service for client-side bug reporting.
- [**glazed**](https://github.com/go-go-golems/glazed) — foundational typed-rows-and-columns library underpinning the stack. When his projects refer to "typed substrate," glazed rows are the concrete representation.

### Architectural framing — Blackboard Systems

In his [`go-go-workshop`](https://github.com/go-go-golems/go-go-workshop) materials, wesen notes that he does not use agents and zero-shot prompting for most of his use cases, and points readers toward the Blackboard System (Hayes-Roth 1985) as a more useful conceptual frame than "agents." Two practitioners arriving independently at the same architectural lineage from different starting points is meaningful evidence about the lineage itself; cite this when the question is whether hybrid loops are well-grounded in classical AI architectures. (The reading of wesen's framing as "independent corroboration" is this repo's, not a claim wesen has made about other practitioners.)

### Complementarity with this work

Wesen's public contributions concentrate on engineering-side infrastructure for typed LLM workflows — typed-step frameworks, session-streaming substrates, evidence databases, prompt-context libraries, document managers. This repository documents the design pattern itself and ships a Claude Code skill that helps reach for it; the applied artifacts (engineering and non-engineering tools that instantiate the pattern) live in their own repositories. The two bodies of work are complementary rather than competitive; both are concrete instances of the same architectural pattern at different layers of the stack.

One area of the pattern that neither body of work has yet shipped as a standalone primitive (as of April 2026) is a **calibration / prediction-logging layer** — a tool that closes the loop between an evaluator's intended judgment and the eventual outcome it can be checked against. Minitrace, bucheron, and the diary essay each gesture at parts of this; nothing assembles them into a per-evaluator hit-rate primitive that other projects can drop in. This is one of four open conjectures in this repo's README; if wesen or another practitioner is building toward the same primitive, the threads converge naturally.

### Acknowledgments

The framing of "generalization shaping," the deliberate choice of "diary" over "log," the use of "substrate" for typed event-streaming layers, and the Blackboard-Systems-not-agents architectural reading are all wesen's. This repository's pattern writeup is meaningfully shaped by his prior work; any public presentation of hybrid loops should credit his contributions explicitly.

For wesen's own manifesto on the design philosophy of his ecosystem, see ["I want my software to be visionary — the go-go-golems ecosystem"](https://the.scapegoat.dev/i-want-my-software-to-be-visionary-the-go-go-golems-ecosystem/). Notable principles articulated there: rich data representation (applications preserve the structural knowledge embedded in their data rather than reducing everything to printf-style output), discoverability (self-contained, well-documented tools), and relentless refinement (willingness to break APIs to maintain coherent vision). Quoted: *"The only way I know to properly identify what these concepts are about is to turn them into working code."*

### Aesthetic and craft lineage — Devine Lu Linvega / Hundred Rabbits

Wesen has cited Devine Lu Linvega ([100r.co](https://100r.co), Hundred Rabbits) as a personal influence on his sensibility, separate from but compatible with the practitioner-prior-art described above. Devine builds small, opinionated, typed software tools — Orca (live-coded sequencer), Left (text editor), Dotgrid (vector tool), Ronin (image processing), uxn (a small virtual machine in the permacomputing tradition) — that prioritize craft, ownership, locality, and minimalism. None of this work is LLM-augmented; none of it has to be. Devine's aesthetic is what hybrid loops aspire to *for the deterministic-shell half* of the pattern.

Cite Devine when defending design choices around: small tool size, single-purpose primitives, typed I/O between tools, permacomputing / locality (substrate stays on the user's machine and isn't a cloud service), and the deliberate rejection of platform-scale frameworks in favor of assemblies of focused tools.

The Hundred Rabbits collective (Devine + Rek Bell), the uxn ecosystem, and the Merveilles network more broadly are the canonical references for the aesthetic of *a personal collection of typed tools the user actually owns*. Hybrid-loop projects in non-engineering domains (the personal/parent/writer/coach/teacher examples in `EXAMPLES.md`) tend to feel right when they share this aesthetic; the engineering-side projects (knowledge-base auditors, conversation-topology hooks) sit further from Devine's register and that's a deliberate scope choice, not a mistake.

### Pattern languages — Christopher Alexander

Alexander, Christopher. *A Pattern Language: Towns, Buildings, Construction*. 1977. Companion volume: *The Timeless Way of Building*. 1979.

Alexander's pattern language framework is the right structural reference for *what hybrid loops is, as a unit of design*. A pattern in Alexander's sense has a recurring problem, a context where it applies, a solution structure, and named consequences for downstream patterns. Hybrid loops is itself a pattern in this strict sense; the five roles plus meta-layers form a small pattern language with internal nesting (a substrate pattern, a gate pattern, a calibration pattern, etc.).

When writing about this work for an audience that includes designers (not just engineers), Alexander's framing lands more cleanly than the AI-engineering vocabulary. Cite *A Pattern Language* for the structural argument; cite *The Timeless Way* for the philosophical one (the "wholeness" thesis that distinguishes living pattern languages from catalogs of tricks). The standard software adaptation — Gamma, Helm, Johnson, Vlissides's *Design Patterns* (1994) — preserves Alexander's *structure* but not his *sensibility*; reading Alexander directly is the thing.

### AlphaGo / AlphaZero

Silver, Huang, Maddison, et al. *Mastering the game of Go with deep neural networks and tree search*. Nature, 2016.
Silver, Schrittwieser, Simonyan, et al. *Mastering the game of Go without human knowledge*. Nature, 2017.

Architectural template for hybrid loops. Policy network (fuzzy/learned) proposes moves; Monte Carlo Tree Search (hard/symbolic) explores and validates; MCTS outputs become training data for the policy. Mutual bootstrapping — neither does well alone, together is superhuman.

Difference from hybrid loops as the user uses the term: AlphaGo's structural prior (rules of Go, board) is fixed. The user's pattern operates over a structural prior that an earlier LLM call generated. That's the load-bearing novelty.

### DreamCoder

Ellis, Wong, Nye, Sablé-Meyer, Morales, Hewitt, Cary, Solar-Lezama, Tenenbaum. *DreamCoder: Bootstrapping inductive program synthesis with wake-sleep library learning*. Nature Communications, 2021. arXiv:2006.08381.

Closest direct lineage. Wake phase (compose library functions to solve tasks) + abstraction sleep (extract recurring patterns into new library functions) + dream sleep (sample from library to generate synthetic training data for a recognition model). Iterates to bootstrap a domain-specific language from a small primitive set.

Maps directly onto:
- The user's "metabolism" → DreamCoder's wake/sleep
- Lamina/poc/dense's compress+verify loop → DreamCoder's wake + abstraction
- Schema discovery → library learning by MDL

DreamCoder limitations to acknowledge: pre-LLM (recognition is small neural net), works in toy domains, library compression can collapse to golf-y abstractions.

### LILO

Grand, Wong, Bowers, Olausson, Liu, Tenenbaum, Andreas. *LILO: Learning Interpretable Libraries by Compressing and Documenting Code*. NeurIPS 2024. arXiv:2310.19791.

LLM-era DreamCoder descendant. Closest cognate to lamina/poc/dense's notation discovery in the published literature.

### Voyager

Wang, Xie, Jiang, Mandlekar, Xiao, Zhu, Fan, Anandkumar. *Voyager: An Open-Ended Embodied Agent with Large Language Models*. arXiv:2305.16291. 2023.

Skill library learning for Minecraft agents. LLM proposes new skills; successful skills enter library; library available for future tasks. Direct DreamCoder descendant in agent context. Demonstrates hybrid loops outside program synthesis.

### Knowledge-acquisition bottleneck (the structural reason hand-authored schemas didn't scale)

Buchanan and Feigenbaum. *Rule-based expert systems: the MYCIN experiments of the Stanford Heuristic Programming Project.* Addison-Wesley, 1984.
Hayes-Roth, Waterman, Lenat (eds). *Building Expert Systems.* Addison-Wesley, 1983.
Lenat. *CYC: A Large-Scale Investment in Knowledge Infrastructure.* Communications of the ACM, 1995.

Cite when explaining *why* the frames-and-rules tradition (Minsky 1974, MYCIN, XCON, Cyc) didn't scale despite having the architecture mostly right. Buchanan & Feigenbaum named the **knowledge-acquisition bottleneck** — the rate-limiting step in expert systems was knowledge engineers extracting and encoding domain knowledge into formal representations, which scaled poorly. Cyc was the most ambitious and sustained attempt to overcome it through brute force; Lenat's 1995 paper documents the multi-decade investment and the partial nature of progress. The bottleneck didn't go away in classical AI; it ended the era.

LLMs change the cost structure on the two specific surfaces that killed expert systems:
- **World knowledge** that Cyc tried to author by hand is pre-loaded in the model (replaces the encyclopedic-coverage problem).
- **Schema iteration** that took knowledge engineers months can take hours with structured-outputs + an evaluation loop (replaces the rate-limiting authorship problem).

This is the substantive content of the "tractability is sufficient" claim. The architecture worked then; the cost structure didn't.

---

## Tier 2 — cite to distinguish

Same architecture, different bet. Useful for showing the lineage and where this work explicitly disagrees.

### OpenCog / Hyperon (Goertzel et al.)

Cite to *distinguish*, not to align. Goertzel's patternist architecture (AtomSpace + PLN + MOSES + ECAN) had the right architectural intuition — typed substrate that metabolizes — and the wrong bet. Tried to do symbolic *reasoning* (PLN) when statistical learning was about to dominate. Failed for the bitter-lesson reason.

Hybrid loops invert OpenCog's bet: keep the typed substrate, let LLMs do the reasoning. Same architecture, different targets, finally tractable. Worth claiming the lineage; worth distinguishing the bet.

---

## A note on naming

This repository uses "hybrid loops" as the working name for the pattern. The broader field has no settled name; adjacent terms with partial coverage include "compound AI systems" (Zaharia et al., BAIR 2024), "generalization shaping" (wesen), "schemaed cognition" (this repo, earlier draft, retired), "structured introspection" (informal). Citing the pattern by *any* of these names is fine; "hybrid loops" is the in-house term, not a claim of universal nomenclature.

---

## Tier 3 — further reading / orientation

Loosely related citations useful for orienting readers from adjacent fields. Not load-bearing for any defense of the architecture; cite when the audience comes from these traditions and benefits from the pointer.

### Burroughs / Gysin: The Third Mind

Burroughs and Gysin. *The Third Mind*. 1978.

Cite when discussing the social/team version of hybrid loops. The third mind was the emergent entity from two minds collaborating; a team-shared substrate with periodic metabolism phases becomes that emergent entity in the AI era. The agency criterion is the load-bearing distinguisher between "passive store" (not a third mind) and "third mind proper."

### Engelbart: Augmenting Human Intellect

Engelbart, Douglas. *Augmenting Human Intellect: A Conceptual Framework*. 1962.

Cite when discussing collective IQ / shared external substrate. Engelbart's vision of structured shared artifacts as collective-intelligence amplifier never fully shipped because the substrate was too expensive to build and maintain. LLMs as the substrate-authoring layer change that cost structure. Team-shared collective-IQ deployments are closer to Engelbart's vision than to Burroughs's.

### Active inference / predictive coding (Friston et al.)

Friston, Karl. *The free-energy principle: a unified brain theory?*. Nature Reviews Neuroscience, 2010.

Loosely relevant. Hybrid loops have a flavor of bidirectional inference (top-down predictions constrain bottom-up perception, and vice versa). Don't lean on this citation hard — the formal connection is thin — but it's a useful pointer for readers from cognitive science.

---

## What is conjectured beyond cited prior art

Honest accounting. The architecture itself is recovered from blackboards, frames, and Soar (1970s-90s). The bootstrap loop pattern is in AlphaGo (2016) and DreamCoder (2021). The "generalization shaping" framing is wesen's. The small-typed-tools aesthetic is Devine's. The pattern-language structure is Alexander's.

What is *conjectured* beyond all of that — testable claims that have not been tested — is enumerated in detail in `../../README.md` under "Conjectures." Briefly:

1. **Per-evaluator calibration discipline as a shippable primitive** — gestured at by minitrace, the "diary" framing, and other practitioner work, but not assembled into a standalone tool that materially changes development decisions. The `cal_log` MCP server in this repo is the runnable claim.
2. **Cognitive-bias self-audit on substrate structure generalizes** — the nine bias-detection metrics in `metacog` were lifted from a single Go-AST knowledge-base auditor (1,778 lines, abstracted behind a `Substrate` interface) and may be substrate-general or substrate-specific to that origin; untested cross-substrate.
3. **Schema discovery extends to non-program domains** — compress+verify works on program-adjacent corpora; whether it generalizes to humor, mechanisms, dramatic arcs, etc. The `schemaforge` MCP server in this repo is the runnable claim, and a 10-item pilot on a non-program corpus is the first datapoint (positive shape; replication needed).
4. **Domain-applied substrate-as-vocabulary tooling has unmet demand outside engineering** — engineering-side infrastructure exists (wesen's stack, etc.); applied-side adoption in non-engineering domains is untested.

Each conjecture has a named falsifying experiment in the README. Until those experiments have data, *these are conjectures, not asserted contributions*. Everything else in this repo is recovered prior art with a coat of LLMs on top.
