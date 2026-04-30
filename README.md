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

The closest practitioner prior art is **Manuel Odendahl's (wesen) work** at [github.com/go-go-golems](https://github.com/go-go-golems) and his blog [the.scapegoat.dev](https://the.scapegoat.dev). He coined **"generalization shaping"** for the design principle inside hybrid loops; he uses **"diary"**, **"evidence database"**, and **"substrate"** as terms; he explicitly cites Blackboard Systems as the right lineage. His shipped infrastructure (geppetto, sessionstream, glazed, go-go-agent) is engineering-focused; this repo's portfolio is more domain-applied. See `skill/references/PRIOR_ART.md` for the full mapping.

What is plausibly new in this work specifically (that wesen and the cited literature don't already own):

1. **Calibration discipline as a per-evaluator primitive** — every typed LLM judgment logs prediction + verdict over time, hit-rate aggregated. Wesen's minitrace and "diary" essay gesture; nothing ships this as a standalone. Open ground.
2. **Cognitive-bias self-audit on substrate structure** (winze): running known cognitive-bias signatures against the structural metrics of the KB itself (provenance HHI as availability-bias proxy, irrelevant:challenged ratio as survivorship bias). No precedent found.
3. **Schema discovery for cognitive schemas, not just program libraries** (lamina/poc/dense): DreamCoder/LILO discover program libraries; this work extends compress+verify to discovering descriptive notations for unstructured-text domains.
4. **Domain-applied substrate-as-vocabulary tooling** across non-engineering domains (humor, dramatic arcs, behavioral mechanisms, AI-conversation analysis): the engineering-side infrastructure exists; the applied-side does not.

## Status

Pre-v1. The skill is the primary deliverable. Primitives (`cal_log`, `metacog`, `schemaforge`, `metabolism`) are mentioned in `skill/references/PRIMITIVES.md` but not yet packaged.

## Form factor: skill now, plugin later

V0 ships as a Claude Code skill. When the architecture needs to bundle hooks (auto-firing UserPromptSubmit/Stop scaffolds) alongside the cognitive primitives, promote to a Claude Code plugin. The skill structure is portable to plugin form when that time comes.
