# Design interview — full version

Use when SKILL.md's Phase 4 quick-design (3 questions + auto-defaults) isn't sufficient. Typically: surfaces large enough that getting the schema wrong has material cost; projects where calibration matters from day one; or when the user explicitly wants more rigor before scaffolding.

Default to quick-design for v0 projects. This file is loaded on demand.

## For analytical surfaces (substrate-as-record)

1. **What is the non-deterministic input?** Transcripts, dialogue turns, screenshots, plans, behavior logs, observations, sensor readings, journal entries.
2. **What does the lens extract?** Sketch the schema. Fields, types, enums.
   - Required: `notes` (free-text rationale for graceful failure), `model_id`, `schema_version`.
   - Useful: provenance pointers to source documents, per-field confidence scores, extraction timestamp.
3. **Where does the substrate live?** JSONL file, sqlite, an existing database, a typed module in the host language, an MCP-exposed query API.
4. **What aggregations or queries does the substrate need to support?** Group-by, time-series, similarity search, full-text? This shapes the substrate choice.
5. **What does the reasoner do with the substrate?** Concrete output: ranked list, recommendation, alert, summary, generated content informed by typed structure.
6. **What's the action?** UI render, file write, notification, new record back into substrate, external API call.
7. **What does the calibration log capture?** Predict + verdict per call. What's the verdict signal — follow-up user action, a metric, a manual confirm, a second LLM check?
8. **Does this surface need metabolism?** Almost always: no, in v0. Add only when the substrate accumulates over weeks and needs periodic audit (consolidation, bias-check, speculative connection).

## For interventional surfaces (substrate-as-vocabulary)

1. **What's the moment that needs a move?** What triggers consideration?
2. **What's the repertoire?** Sketch the typed library. Entries, fields per entry, discrimination criteria.
   - Common per-entry fields: `deploy_when` (triggering condition), `contraindications` (when not to fire), `mechanism` (why this works), `tone_or_register`, `last_used_at`.
3. **Where does the repertoire live?** Embedded in code, JSON file, database, MCP-exposed library.
4. **What restraint policy gates the firing?** **Usually the most important question.** Cooldown decay, ripeness window, distress gate, confidence threshold, mandatory abstention conditions. A bad gate makes the tool annoying; a good gate makes it feel respectful.
5. **What does the reasoner consider when choosing a move?** The moment + the repertoire entry's criteria + recent state (what's been used, what didn't land).
6. **What's the action?** Surface to user, take in-context, generate content, log only.
7. **What does success look like?** Often weak signal for interventional surfaces. Sometimes the only signal is the user not pushing back. Name what counts as a landed move and what counts as a missed one.

## For surfaces that are both

If a surface accumulates AND offers a typed repertoire (the substrate-as-record-and-vocabulary case), run both interviews and reconcile. Substrate often has separate "record" tables and "vocabulary" tables; the gate operates on both.

## Schema design tips

- Start with 3-5 fields. Add more only when a missing field causes a downstream confusion.
- Always include `notes` (free-text rationale).
- Always include `model_id` and `schema_version` per record (not per project — bumps will happen).
- Enums: start with a closed list. Open enums only when the closed version repeatedly produces "other" values that turn out to be a real category.
- Confidence/scores: 0-1 floats are fine. Don't over-engineer with calibrated probability distributions until you have calibration data.

## Anti-patterns in the design interview itself

- **Asking all 8 questions before producing a draft.** Produce a draft after Q1-3; let the user revise.
- **Designing for hypothetical future calibration before having any v0 calibration data.** Log first; design later.
- **Over-specifying the gate before observing actual over-firing patterns.** Start permissive; tighten with data.
- **Designing the metabolism phase before the substrate has accumulated 4+ weeks of records.** Premature.
- **Making the schema large because it might-need-this-later.** Smaller schemas evolve better than larger ones.

## Companion to SKILL.md

This file is loaded on demand. Most projects don't need it; the quick-design in SKILL.md Phase 4 is sufficient for v0.
