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

What is plausibly new in this work specifically (that wesen and the cited literature don't already own):

1. **Calibration discipline as a per-evaluator primitive** — every typed LLM judgment logs prediction + verdict over time, hit-rate aggregated. Wesen's minitrace and "diary" essay gesture; nothing ships this as a standalone. Open ground.
2. **Cognitive-bias self-audit on substrate structure** (winze): running known cognitive-bias signatures against the structural metrics of the KB itself (provenance HHI as availability-bias proxy, irrelevant:challenged ratio as survivorship bias). No precedent found.
3. **Schema discovery for cognitive schemas, not just program libraries** (lamina/poc/dense): DreamCoder/LILO discover program libraries; this work extends compress+verify to discovering descriptive notations for unstructured-text domains.
4. **Domain-applied substrate-as-vocabulary tooling** across non-engineering domains (humor, dramatic arcs, behavioral mechanisms, AI-conversation analysis): the engineering-side infrastructure exists; the applied-side does not.

## Status

Pre-v1. The skill is the primary deliverable. Primitives (`cal_log`, `metacog`, `schemaforge`, `metabolism`) are mentioned in `skill/references/PRIMITIVES.md` but not yet packaged.

## Acknowledgments

This writeup is meaningfully shaped by **Manuel Odendahl** (wesen)'s prior work in this design space. The "generalization shaping" framing, the deliberate choice of "diary" over "log," the use of "substrate" for typed event-streaming layers, and the Blackboard-Systems architectural reading are all his. Any public presentation of hybrid loops should credit his contributions; a fuller account is in `skill/references/PRIOR_ART.md`.

Thanks also to the maintainers of [DreamCoder](https://github.com/ellisk42/ec), [LILO](https://github.com/gabegrand/lilo), [Voyager](https://github.com/MineDojo/Voyager), and the [Polis](https://pol.is/) and [Talk to the City](https://github.com/AIObjectives/talktothe.city) projects, whose published work is referenced throughout the skill.

## Form factor: skill now, plugin later

V0 ships as a Claude Code skill. When the architecture needs to bundle hooks (auto-firing UserPromptSubmit/Stop scaffolds) alongside the cognitive primitives, promote to a Claude Code plugin. The skill structure is portable to plugin form when that time comes.
