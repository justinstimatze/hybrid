# hybrid

Working repository for the hybrid-loops design pattern: projects where an LLM extracts typed structure from non-deterministic content, deterministic code operates on that typed surface, and another LLM call reasons over it to drive decisions or actions. Mutual bootstrapping between LLM fuzziness and deterministic constraint, where neither does well alone.

## Contents

- `skill/` — the Claude Code skill that lets Claude reach for this architecture by default. Symlinked to `~/.claude/skills/hybrid-loops/` for live use.
  - `SKILL.md` — main skill file (loaded when Claude invokes the skill)
  - `references/PRIOR_ART.md` — citations and how each compares
  - `references/EXAMPLES.md` — the user's existing repos mapped to the five-role schema
  - `references/PRIMITIVES.md` — extractable primitives (mostly not yet packaged)
  - `references/STACKING.md` — recursive composition / multi-layer stacking discussion

## Companion projects

The skill references real projects across the portfolio (slimemold, drivermap, score, plancheck, lucida, winze, gemot, publicrecord, lamina/poc/dense, crowdwork, groupchat, effigy, ismyaialive, seeing). The skill's `references/EXAMPLES.md` maps each to the five-role schema (lens / substrate / gate / reasoner / action) plus optional meta-layers (calibration / metabolism).

## Provenance and honest novelty claim

This pattern is not novel as architecture. It is recovered from frames (Minsky 1974), blackboards (Hayes-Roth 1985), and Soar (Newell 1990); the bootstrapping loop is in AlphaGo (2016) and DreamCoder (2021); LLMs as fuzzy schema-fillers are what makes the architecture finally tractable in 2026.

The closest practitioner prior art is **Manuel Odendahl** (wesen), whose work at [github.com/go-go-golems](https://github.com/go-go-golems) and writing at [the.scapegoat.dev](https://the.scapegoat.dev) has shaped this writeup directly. He has named the design principle of **"generalization shaping"** (deterministic machinery shaping what the LLM has to do); introduced or made canonical the use of **"diary,"** **"evidence database,"** **"substrate,"** and **"step"** in this design space; and identified the Blackboard System lineage as the right architectural frame. His shipped infrastructure — including [geppetto](https://github.com/go-go-golems/geppetto), [sessionstream](https://github.com/go-go-golems/sessionstream), [glazed](https://github.com/go-go-golems/glazed), [pinocchio](https://github.com/go-go-golems/pinocchio), [go-go-agent](https://github.com/wesen/2026-04-29--go-go-agent), and [docmgr](https://github.com/go-go-golems/docmgr) — concentrates on engineering-side typed LLM workflows. The portfolio in this repository concentrates on domain-applied tools that sit on top of similar infrastructure. The two bodies of work are complementary; see `skill/references/PRIOR_ART.md` for the full citation and credit.

## Conjectures (not asserted contributions)

Four conjectures about what this work might contribute beyond the cited prior art. **Each is testable; none has been tested.** They are conjectures, not asserted contributions, until the experiments below have data. The retrofit in `skill/references/PORTFOLIO_RETROFIT.md` provides initial empirical context.

### Conjecture 1 — per-evaluator calibration is undershipped

*Claim.* A standalone primitive that logs predictions and verdicts per typed LLM evaluator, with rolling-window hit-rate aggregation, would generalize across hybrid-loop projects and meaningfully change development decisions.

*Falsifying experiment.* Build the primitive (`cal_log` per `skill/references/PRIMITIVES.md`); deploy on 3+ existing projects (candidates: slimemold, plancheck, drivermap); measure over 60 days whether the hit-rate signal changes any concrete development decision (prompt change, schema bump, gate adjustment). If hit-rate is collected but no decisions are made on it, the primitive is theater.

*Status.* Conjecture. The portfolio retrofit shows every existing project gestures at calibration but none ships per-evaluator hit-rate; wesen's body of work shows the same gap. The empirical absence is consistent with the conjecture but doesn't confirm the primitive would be useful — only the deployment experiment does.

### Conjecture 2 — cognitive-bias self-audit on substrate structure generalizes beyond winze

*Claim.* The nine bias-detection metrics from winze (provenance HHI as availability-heuristic proxy, irrelevant:challenged ratio as survivorship-bias proxy, predicate entropy as base-rate-neglect proxy, etc.) work on any typed substrate, not only winze's typed Go AST.

*Falsifying experiment.* Lift the `metacog` primitive out of winze; run it on slimemold's claim graph and drivermap's mechanism corpus. Have project owners mark which triggered findings correspond to actual structural problems vs. false positives. If false-positive rate is high or findings don't track owner intuition, the metrics are winze-specific rather than substrate-general.

*Status.* Conjecture. No published precedent found; no cross-substrate testing done.

### Conjecture 3 — schema discovery extends to non-program domains

*Claim.* The dense compress+verify loop (DreamCoder/LILO descendant) can discover useful schemas for non-program domains — humor structures, dramatic arcs, behavioral mechanisms, AI-conversation patterns — not just code or notation for code.

*Falsifying experiment.* Run dense on a corpus of 100+ examples in one non-program domain (drivermap's mechanism corpus is the obvious candidate); compare the discovered schema against the hand-authored one on downstream extraction quality, with a domain expert as the evaluator. If hand-authored beats discovered, the loop doesn't generalize past program-shaped domains.

*Status.* Partially supported by dense's existing CRUD/SEC/patent results, but those domains are still program-adjacent. Non-program domain testing missing.

### Conjecture 4 — there is unmet demand for domain-applied substrate-as-vocabulary tooling outside engineering

*Claim.* Users in non-engineering domains (coaching, teaching, parenting, advocacy, creative work) would benefit from typed-repertoire-with-restraint tools and don't currently have them.

*Falsifying experiment.* Build one such tool (the teacher's intervention tracker or the coach's typed library from `EXAMPLES.md`); ship to 5+ domain users; measure 30+ day retention. If retention is below baseline rates for similar consumer tools, the demand isn't there or the tool is wrong-shaped.

*Status.* Conjecture. The author's existing portfolio is engineering-flavored; no shipped non-engineering substrate-as-vocabulary tools yet.

---

These four are the open ground after acknowledging wesen, DreamCoder, AlphaGo, OpenCog, and the classical-AI lineage. **The next material work in this repo is running the experiments above, not making more architectural claims.**

## What this repo is and isn't

**This is research output, not a product.** It documents a design pattern, a Claude Code skill that helps reach for the pattern, retrofit notes on existing projects, and four conjectures with named experiments. None of it is sold; nothing here is trying to be a SaaS or a paid tool. The user's existing portfolio (slimemold, drivermap, gemot, etc.) and wesen's portfolio are the *applied* artifacts; this repo is the *meta-artifact* that names what they have in common.

If a hybrid-loop product is to be built, it would be one of the *applied* tools (a teacher's intervention tracker, a recruiter triage assistant, a small advocacy legislator-tracker — see `skill/references/EXAMPLES.md`), not the skill or the pattern itself. The pattern doesn't sell; specific tools built with it might.

## Naming

"Hybrid loops" is the **working name in this repository**, not a claim of universal nomenclature. The broader field has no settled name; adjacent terms with partial coverage include "compound AI systems" (Zaharia, BAIR 2024), "generalization shaping" (wesen), "schemaed cognition" (this repo, earlier draft, retired), "structured introspection" (informal practitioner usage). The pattern can be cited by any of these names with rough fidelity.

A separate term sometimes used in this conversation is **"third mind"** — a *deployment context* for hybrid loops where the substrate is shared between collaborators (and possibly an AI), distinct from a personal external substrate ("second mind") or the user's own thinking ("first mind"). Third mind has prior occupancy in Burroughs/Gysin's 1978 sense (the emergent entity in collaborative writing); the user's usage extends this to substrates with their own metabolism (winze-style). **Third mind is a deployment shape; hybrid loops is the architectural pattern.** They are related but not synonyms.

## Status

Pre-v1. The skill is the primary deliverable. Primitives (`cal_log`, `metacog`, `schemaforge`, `metabolism`) are mentioned in `skill/references/PRIMITIVES.md` but not yet packaged.

## Acknowledgments

This writeup is meaningfully shaped by **Manuel Odendahl** (wesen)'s prior work in this design space. The "generalization shaping" framing, the deliberate choice of "diary" over "log," the use of "substrate" for typed event-streaming layers, and the Blackboard-Systems architectural reading are all his. Any public presentation of hybrid loops should credit his contributions; a fuller account is in `skill/references/PRIOR_ART.md`.

Thanks also to the maintainers of [DreamCoder](https://github.com/ellisk42/ec), [LILO](https://github.com/gabegrand/lilo), [Voyager](https://github.com/MineDojo/Voyager), and the [Polis](https://pol.is/) and [Talk to the City](https://github.com/AIObjectives/talktothe.city) projects, whose published work is referenced throughout the skill.

## Form factor: skill now, plugin later

V0 ships as a Claude Code skill. When the architecture needs to bundle hooks (auto-firing UserPromptSubmit/Stop scaffolds) alongside the cognitive primitives, promote to a Claude Code plugin. The skill structure is portable to plugin form when that time comes.
