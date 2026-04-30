# Worked examples

The hybrid-loops pattern applies across many domains, not just engineering tooling. This file leads with non-engineering examples to show the range, then lists the user's existing portfolio as a second section. When scaffolding a new project, pick whichever example is closest in *shape* — domain doesn't matter, structure does.

---

## Non-engineering examples

### A teacher tracking which interventions work for which kinds of students (substrate-as-record, education)

- **Surface scope** — only the post-interaction reflection. The rest of the teacher's day (lesson planning, grading, classroom management) is not a hybrid loop.
- **Lens** — LLM extracts typed records from the teacher's brief notes after each student interaction: `{student_id, situation_type, intervention_used, response_pattern, notes, schema_version, model_id}`
- **Substrate** — sqlite; per-student-per-week records; schema versioned because `intervention_used` taxonomy will evolve as the teacher learns
- **Gate** — flags unusual patterns (a student whose responses don't match their typical cluster); aggregates per-intervention success rates per student-type
- **Reasoner** — given a current situation, suggests interventions ranked by historical fit, citing specific past records
- **Action** — note in the teacher's planning doc; optional weekly digest
- **Calibration** — each suggestion logged; verdict comes from the next interaction's `response_pattern`

This is the canonical "single-domain-expert wants to learn from their own past" shape. The whole project is the surface; no separate UI complexity needed. A notebook + sqlite is enough.

### A parent reflecting on interactions with their child (substrate-as-record, personal)

- **Surface scope** — the reflection-after-the-fact pass. Real-time parenting is not a hybrid loop and shouldn't be — the loop is too slow for in-the-moment use.
- **Lens** — voice-note transcription + LLM extraction of `{moment_type, child_state, parent_response, emotion_in_self, what_worked, what_didnt, schema_version}`
- **Substrate** — sqlite locally; tagged by date, context, emotion
- **Gate** — clusters similar moments (recurring scenarios); flags when child_state escalates over a week
- **Reasoner** — when asked, surfaces patterns ("you tend to respond with X when child is Y; here's a moment where Z worked better")
- **Action** — written reflection in a journal note; never auto-intervention — the user invokes
- **Calibration** — verdict signal is weak (no objective "was that response right?"); use self-reported "looking back" verdict, manually entered

Privacy note: data stays local; never sent beyond the LLM call. Schema includes a `consent_recorded` field per interaction if other family members are described. This is a case where the deployment shape (local-only, no cloud) is load-bearing.

### A small advocacy group tracking legislators on an issue (substrate-as-record + metabolism, civic)

- **Surface scope** — the legislator-position-tracking pass. Fundraising, event planning, member communications are separate (probably non-hybrid-loop) parts of the project.
- **Lens** — LLM extracts position-on-issue from voting records, press releases, social posts: `{legislator_id, issue, position, evidence_type, evidence_quote, date, source_url}`
- **Substrate** — sqlite; per-legislator-per-issue records over time; provenance enforced
- **Gate** — trajectory detection (drift over months); confidence threshold for inclusion; source-diversity audit
- **Reasoner** — when planning advocacy strategy, suggests targets by movability + influence
- **Action** — strategy memo; alerts when a legislator's trajectory changes
- **Calibration** — verdict via subsequent voting record matching predicted trajectory
- **Metabolism** — re-extract weekly; bias audit against source distribution (don't over-weight one outlet)

Substrate is record AND vocabulary (the issue taxonomy is the vocabulary the system uses to discriminate). This is the publicrecord shape applied to a different domain.

### A coach with a typed intervention library (substrate-as-vocabulary, professional/coaching)

- **Surface scope** — the intervention-selection moment during a session. Session notes, scheduling, billing are not hybrid loops.
- **Lens** — at the session moment, an LLM classifies the conversation's current state: `{client_emotional_register, conversation_topic, stuckness_signal, recent_breakthrough}`
- **Substrate** — a *curated* repertoire of typed intervention questions: `[{question_text, deploy_when, contraindications, mechanism, depth_level}]` — maybe 30-100 entries, hand-authored or distilled from training
- **Gate** — restraint policy: don't suggest the same intervention twice in one session; honor `contraindications`; only fire when stuckness_signal is high enough
- **Reasoner** — picks the best-fit intervention given conversation state and recent history
- **Action** — surfaces the suggestion to the coach (not to the client) as a card during the session
- **Calibration** — verdict signal is whether the coach used the suggestion; success rate is whether the session unstuck after

Substrate-as-vocabulary projects almost always have human authorship of the repertoire as a load-bearing step. The coach designs the library; the system picks from it. This is a partnership, not an automation.

### A writer maintaining voice consistency across drafts (substrate-as-record, creative)

- **Surface scope** — the voice-checking pass on a finished draft. Generation of new content is not a hybrid loop here (it's the writer writing).
- **Lens** — LLM extracts voice features from finished pieces: `{piece_id, sentence_complexity, lexical_register, rhetorical_moves, cadence_features, stylistic_quirks}`
- **Substrate** — JSONL; one record per piece; the writer's "voice corpus"
- **Gate** — for a new draft, computes deviation from rolling-window average; flags passages that drift significantly
- **Reasoner** — when invoked on a draft, identifies passages reading differently from the writer's baseline and explains the deviation
- **Action** — annotations in markdown comments; never auto-edits
- **Calibration** — writer accepts/rejects each annotation; rejection rate is inverse hit-rate

The lens schema is the writer's *implicit theory of their own voice*, made explicit. Discovering that schema (which fields, which enums) is itself a sub-project worth the dense-style schema-discovery treatment if the writer has a sufficient corpus.

### A recruiter screening resumes (multiple surfaces in one project)

This project has *three* surfaces, illustrating Phase 2 scope decisions:

- **Surface 1 — resume parsing** (Bucket B, *not* hybrid loop): extracting `{name, education, experience_entries[]}` is mostly deterministic if resume format is consistent. Use a parser library, not an LLM, unless the formats vary wildly.
- **Surface 2 — fit-scoring** (Bucket C, hybrid loop, analytical): lens extracts candidate-criteria-fit `{years_in_role, domain_match, level, signal_strength}`; substrate over all candidates; gate filters below threshold; reasoner ranks for human review; calibration via "did we interview, did they pass."
- **Surface 3 — outreach composition** (Bucket C, hybrid loop, interventional): typed library of message templates `[{template, deploy_when, tone, length}]`; reasoner picks given candidate context; gate restrains template reuse within same-day; action drafts message.

The "project" is one tool; the surfaces are three with different shapes. This is the realistic case for most non-toy projects.

---

## Picking a template by shape

| Shape | Use when... | Example anchor |
|---|---|---|
| Substrate-as-record (analytical) | Value comes from making sense of accumulated data over time | Teacher's intervention tracker; writer's voice corpus |
| Substrate-as-vocabulary (interventional) | System needs to discriminate the right move from a typed repertoire | Coach's intervention library; D&D game-master's plot-device picker |
| Both (record AND vocabulary) | Substrate is queried by present moment AND grows over time | Advocacy legislator-tracker (record of positions, vocabulary of issue types) |
| Substrate provider for downstream agents | Output is intentionally typed for other agents to consume | Decision-time RAG (publicrecord-style) |
| Multiple surfaces in one project | Project has 2-3 distinct fuzzy-judgment places | Recruiter tool above |

---

## Engineering examples — the user's existing portfolio

These are the projects that informed the original pattern. They lean engineering/AI-tooling because that's what the user has built. Use them as anchors when scaffolding a similar engineering tool, but don't assume new projects need to look like them.

### slimemold (epistemics on Claude Code transcripts)

- **Lens** — `internal/extract/extract.go`: Claude tool_use call extracting per-claim records (basis, edges, Moore flags) from transcript chunks
- **Substrate** — typed claim graph persisted across sessions; age decay built into storage
- **Gate** — cooldown decay, cold-start floor, age-based priority selection, per-claim cooldown
- **Reasoner** — structural-fragility analysis (eight vulnerability types) over the claim graph
- **Action** — system-message injection during Stop hook, biases future generation
- **Calibration** — partial; load-bearing predictions get logged but verdict signals are weak
- **Metabolism** — none yet

### drivermap (behavioral mechanism KB)

- **Lens** — extract.py: blind-then-guided two-phase Claude extraction of mechanism records from Wikipedia/Kagi text; schema empirically derived from 20 blind extractions
- **Substrate** — ~137 mechanism JSON records + scoring engine
- **Gate** — schema enforcement, prompt-repetition technique for causal bidirectionality, post-extraction verifier
- **Reasoner** — predict_mechanisms (deterministic scoring on person×situation), then verbalize_motivation (LLM)
- **Action** — MCP tool returns predictions; demo.py composes them into dialogue
- **Calibration** — none currently; verbalization predictions could be logged
- **Metabolism** — none

### winze (typed epistemic substrate with sleep cycles)

- **Lens** — Resolve phase (Claude Sonnet classifies sensor signal as corroborated/challenged/irrelevant)
- **Substrate** — typed Go AST as KB; entities, claims, predicates, theories
- **Gate** — bias-audit gates control which phases fire (availability heuristic, survivorship bias)
- **Reasoner** — MCP server exposing claims/disputes/provenance/search/stats/theories to other LLMs
- **Action** — promote claims, schedule sensor queries, write to corpus
- **Calibration** — `.metabolism-calibration.jsonl` per-cycle (canonical implementation; copy from here)
- **Metabolism** — full: dream (consolidation), trip (speculative connections), evolve (sensor), bias-audit (KB self-check)

### plancheck (file-prediction for AI coding agents)

- **Lens** — RunAgentSpike (tool-using Claude agent that explores codebase) producing AgentFile records with confidence scores
- **Substrate** — AgentResult + structural probe outputs (compiler blast radius, comod history, reference graph)
- **Gate** — novelty-weighted confidence ranking (struct weight 0.5→0.1, semantic 0.1→0.4 as novelty rises)
- **Reasoner** — rankCandidateFiles blends spike + structural signals
- **Action** — file-list output to user, optional record_outcome for calibration
- **Calibration** — record_outcome and record_reflection MCP tools (partial)
- **Metabolism** — none

### lucida (live notebook from Claude Code transcript)

- **Lens** — classifier (Haiku, prompt-cached): discourse_move, cell_type, confidence; specialist (Sonnet): produces concrete spec
- **Substrate** — cells.json with provenance per cell
- **Gate** — confidence threshold (>0.8 mint, 0.6-0.8 draft, <0.6 suppress)
- **Reasoner** — reflect.py every 30 cells consumes the cell stream and synthesizes
- **Action** — frontend renders cells; reflection cells appear inline
- **Calibration** — cost/cache stats per cell; cell quality not yet tracked
- **Metabolism** — partial: reflection synthesis has dream-flavor

### gemot (multi-agent deliberation)

- **Lens** — analyze action: LLM extracts taxonomy, claims, cruxes from agent positions (parallel claim extraction)
- **Substrate** — positions, votes, cruxes, clusters; vote matrix as separate substrate (PCA/SVD/k-means)
- **Gate** — round-based protocol; vote-matrix analysis; commitment fulfillment tracking
- **Reasoner** — propose_compromise, reframe (LLM consumes substrate to produce synthesis)
- **Action** — return analysis to participating agents; track commitments
- **Calibration** — convergence metrics across rounds (built in)
- **Metabolism** — reputation graph across deliberations (the long-term moat)

### publicrecord (decision-time RAG for accountability data)

- **Lens** — LLM extracts findings (entity, relationship, severity, verbatim quote) from primary-source documents
- **Substrate** — SQLite DB; 4,656 entities, 6,377 findings, relationships up to 3 levels deep; rebranding-aware entity resolution
- **Gate** — `riskLevel` enum (low/high/medium); temporal awareness; provenance enforcement
- **Reasoner** — *not internal*; downstream agents consume via MCP/REST/SDK
- **Action** — structured response with citations
- **Calibration** — partial; "data maturity notice" acknowledges summaries are AI-synthesized
- **Metabolism** — none currently

Canonical example of **recursive composition**: publicrecord is a hybrid loop whose output is intentionally designed to be consumed by other hybrid loops.

### groupchat (meme deployment), crowdwork (humor signal), effigy (NPC voice), score (dramatic arcs), ismyaialive (AI-conversation pattern analysis)

These five lean substrate-as-vocabulary or hybrid record/vocabulary. Read the respective READMEs for the full mapping when scaffolding a similar interventional project. The common thread: a typed repertoire (memes / frame domains / character traits / dramatic plays / sycophancy patterns) plus a runtime restraint policy.

---

## When scaffolding, pick the closest match

- Single-domain expert tracking their own observations over time → teacher / parent / writer template
- Civic / policy / accountability tracking with provenance → advocacy / publicrecord template
- Coach / clinician / facilitator with a typed move library → coach / D&D-GM template
- Recruiter / triage / customer-service routing → multi-surface (recruiter) template
- Engineering tool that augments AI agents → slimemold / lucida / plancheck template
- Multi-agent or multi-party deliberation → gemot template
- Substrate provider for downstream agents → publicrecord template
- Building a typed corpus from documents → drivermap template
- KB that should grow over weeks/months with metabolism → winze template

Resist inventing a new architecture. Pick the closest existing example, name what's different, design the difference.
