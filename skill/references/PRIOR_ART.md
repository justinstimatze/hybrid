# Prior art for hybrid loops

Cite these when defending the architecture. Most relevant first.

## Practitioner prior art — Manuel Odendahl ("wesen")

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

Wesen's public contributions concentrate on engineering-side infrastructure for typed LLM workflows — typed-step frameworks, session-streaming substrates, evidence databases, prompt-context libraries, document managers. The portfolio in this repo concentrates on domain-applied tools that sit on top of similar infrastructure (humor, behavioral mechanisms, dramatic arcs, AI-conversation analysis, accountability data, multi-agent deliberation). The two bodies of work are complementary rather than competitive; both are concrete instances of the same architectural pattern at different layers of the stack.

One area of the pattern that neither body of work has yet shipped as a standalone primitive (as of April 2026) is a **calibration / prediction-logging layer** — a tool that closes the loop between an evaluator's intended judgment and the eventual outcome it can be checked against. Minitrace, bucheron, and the diary essay each gesture at parts of this; nothing assembles them into a per-evaluator hit-rate primitive that other projects can drop in. The N3 north star in this repo is to ship that primitive. If wesen is building toward the same thing, this would be a natural place for the two threads to converge.

### Acknowledgments

The framing of "generalization shaping," the deliberate choice of "diary" over "log," the use of "substrate" for typed event-streaming layers, and the Blackboard-Systems-not-agents architectural reading are all wesen's. This repository's pattern writeup is meaningfully shaped by his prior work; any public presentation of hybrid loops should credit his contributions explicitly.

## AlphaGo / AlphaZero

Silver, Huang, Maddison, et al. *Mastering the game of Go with deep neural networks and tree search*. Nature, 2016.
Silver, Schrittwieser, Simonyan, et al. *Mastering the game of Go without human knowledge*. Nature, 2017.

Architectural template for hybrid loops. Policy network (fuzzy/learned) proposes moves; Monte Carlo Tree Search (hard/symbolic) explores and validates; MCTS outputs become training data for the policy. Mutual bootstrapping — neither does well alone, together is superhuman.

Difference from hybrid loops as the user uses the term: AlphaGo's structural prior (rules of Go, board) is fixed. The user's pattern operates over a structural prior that an earlier LLM call generated. That's the load-bearing novelty.

## DreamCoder

Ellis, Wong, Nye, Sablé-Meyer, Morales, Hewitt, Cary, Solar-Lezama, Tenenbaum. *DreamCoder: Bootstrapping inductive program synthesis with wake-sleep library learning*. Nature Communications, 2021. arXiv:2006.08381.

Closest direct lineage. Wake phase (compose library functions to solve tasks) + abstraction sleep (extract recurring patterns into new library functions) + dream sleep (sample from library to generate synthetic training data for a recognition model). Iterates to bootstrap a domain-specific language from a small primitive set.

Maps directly onto:
- The user's "metabolism" → DreamCoder's wake/sleep
- Lamina/poc/dense's compress+verify loop → DreamCoder's wake + abstraction
- Schema discovery → library learning by MDL

DreamCoder limitations to acknowledge: pre-LLM (recognition is small neural net), works in toy domains, library compression can collapse to golf-y abstractions.

## LILO

Grand, Wong, Bowers, Olausson, Liu, Tenenbaum, Andreas. *LILO: Learning Interpretable Libraries by Compressing and Documenting Code*. NeurIPS 2024. arXiv:2310.19791.

LLM-era DreamCoder descendant. Closest cognate to lamina/poc/dense's notation discovery in the published literature.

## Voyager

Wang, Xie, Jiang, Mandlekar, Xiao, Zhu, Fan, Anandkumar. *Voyager: An Open-Ended Embodied Agent with Large Language Models*. arXiv:2305.16291. 2023.

Skill library learning for Minecraft agents. LLM proposes new skills; successful skills enter library; library available for future tasks. Direct DreamCoder descendant in agent context. Demonstrates hybrid loops outside program synthesis.

## OpenCog / Hyperon (Goertzel et al.)

Cite to *distinguish*, not to align. Goertzel's patternist architecture (AtomSpace + PLN + MOSES + ECAN) had the right architectural intuition — typed substrate that metabolizes — and the wrong bet. Tried to do symbolic *reasoning* (PLN) when statistical learning was about to dominate. Failed for the bitter-lesson reason.

Hybrid loops invert OpenCog's bet: keep the typed substrate, let LLMs do the reasoning. Same architecture, different targets, finally tractable. Worth claiming the lineage; worth distinguishing the bet.

## Burroughs / Gysin: The Third Mind

Burroughs and Gysin. *The Third Mind*. 1978.

Cite when discussing the social/team version of hybrid loops. The third mind was the emergent entity from two minds collaborating; a team-shared substrate with metabolism (winze-style) becomes that emergent entity in the AI era. The agency criterion is the load-bearing distinguisher between "passive store" (not a third mind) and "third mind proper."

## Engelbart: Augmenting Human Intellect

Engelbart, Douglas. *Augmenting Human Intellect: A Conceptual Framework*. 1962.

Cite when discussing collective IQ / shared external substrate. Engelbart's vision of structured shared artifacts as collective-intelligence amplifier never fully shipped because the substrate was too expensive to build and maintain. LLMs as the substrate-authoring layer change that cost structure. Gemot's MAGI framing is closer to Engelbart's vision than to Burroughs's.

## Active inference / predictive coding (Friston et al.)

Friston, Karl. *The free-energy principle: a unified brain theory?*. Nature Reviews Neuroscience, 2010.

Loosely relevant. Hybrid loops have a flavor of bidirectional inference (top-down predictions constrain bottom-up perception, and vice versa). Don't lean on this citation hard — the formal connection is thin — but it's a useful pointer for readers from cognitive science.

## What is genuinely new vs cited prior art

Honest accounting. The architecture itself is recovered from blackboards, frames, and Soar (1970s-90s). The bootstrap loop pattern is in AlphaGo (2016) and DreamCoder (2021). What is plausibly novel about hybrid loops as the user uses the term:

1. **The structural prior is generated by an earlier LLM call.** AlphaGo's structure is Go's rules. DreamCoder's structure is the programming language syntax. In hybrid loops, the substrate's schema can be discovered by an earlier loop (lamina/poc/dense). This extends bootstrapping to domains where there is no pre-existing structural prior.

2. **The cognitive-bias self-audit move (winze-specific).** Running known cognitive-bias signatures against the structural metrics of the substrate itself: provenance HHI as availability-bias proxy, irrelevant-to-challenged ratio as survivorship bias, predicate entropy as base-rate-neglect. This specific move I cannot find precedent for. If it generalizes — if any typed substrate can be audited this way — it's the genuine substrate-level primitive.

3. **Solo/small-team scale.** Classical bootstrapping architectures were enterprise/research scale. Hybrid loops at solo-developer scale is downstream of LLMs being cheap, not architecturally novel — but it changes who can build them, which changes what gets built.

These three together are the defensible claim. Everything else is recovered prior art with a coat of LLMs on top.
