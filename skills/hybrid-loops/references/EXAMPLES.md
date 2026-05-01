# Worked examples

The hybrid-loops pattern applies across many domains. These are illustrative examples — fictional or generic — for someone scaffolding a project. When designing your own, pick whichever is closest in *shape*; domain doesn't matter, structure does.

---

## A teacher tracking which interventions work for which kinds of students (substrate-as-record, education)

- **Surface scope** — only the post-interaction reflection. The rest of the teacher's day (lesson planning, grading, classroom management) is not a hybrid loop.
- **Lens** — LLM extracts typed records from the teacher's brief notes after each student interaction: `{student_id, situation_type, intervention_used, response_pattern, notes, schema_version, model_id}`
- **Substrate** — sqlite; per-student-per-week records; schema versioned because `intervention_used` taxonomy will evolve as the teacher learns
- **Gate** — flags unusual patterns (a student whose responses don't match their typical cluster); aggregates per-intervention success rates per student-type
- **Reasoner** — given a current situation, suggests interventions ranked by historical fit, citing specific past records
- **Action** — note in the teacher's planning doc; optional weekly digest
- **Calibration** — each suggestion logged; verdict comes from the next interaction's `response_pattern`

The canonical "single-domain-expert wants to learn from their own past" shape. The whole project is the surface; no separate UI complexity needed. A notebook + sqlite is enough.

## A parent reflecting on interactions with their child (substrate-as-record, personal)

- **Surface scope** — the reflection-after-the-fact pass. Real-time parenting is not a hybrid loop and shouldn't be.
- **Lens** — voice-note transcription + LLM extraction of `{moment_type, child_state, parent_response, emotion_in_self, what_worked, what_didnt, schema_version}`
- **Substrate** — sqlite locally; tagged by date, context, emotion
- **Gate** — clusters similar moments (recurring scenarios); flags when child_state escalates over a week
- **Reasoner** — when asked, surfaces patterns ("you tend to respond with X when child is Y; here's a moment where Z worked better")
- **Action** — written reflection in a journal note; never auto-intervention — the user invokes
- **Calibration** — verdict signal is weak; use self-reported "looking back" verdict, manually entered

Privacy note: data stays local; never sent beyond the LLM call. Schema includes a `consent_recorded` field per interaction if other family members are described. The deployment shape (local-only, no cloud) is load-bearing.

## A small advocacy group tracking legislators on an issue (substrate-as-record + metabolism, civic)

- **Surface scope** — the legislator-position-tracking pass. Fundraising, event planning, member communications are separate (probably non-hybrid-loop) parts of the project.
- **Lens** — LLM extracts position-on-issue from voting records, press releases, social posts: `{legislator_id, issue, position, evidence_type, evidence_quote, date, source_url}`
- **Substrate** — sqlite; per-legislator-per-issue records over time; provenance enforced
- **Gate** — trajectory detection (drift over months); confidence threshold for inclusion; source-diversity audit
- **Reasoner** — when planning advocacy strategy, suggests targets by movability + influence
- **Action** — strategy memo; alerts when a legislator's trajectory changes
- **Calibration** — verdict via subsequent voting record matching predicted trajectory
- **Metabolism** — re-extract weekly; bias audit against source distribution (don't over-weight one outlet)

Substrate is record AND vocabulary (the issue taxonomy is the vocabulary the system uses to discriminate).

## A coach with a typed intervention library (substrate-as-vocabulary, professional/coaching)

- **Surface scope** — the intervention-selection moment during a session. Session notes, scheduling, billing are not hybrid loops.
- **Lens** — at the session moment, an LLM classifies the conversation's current state: `{client_emotional_register, conversation_topic, stuckness_signal, recent_breakthrough}`
- **Substrate** — a *curated* repertoire of typed intervention questions: `[{question_text, deploy_when, contraindications, mechanism, depth_level}]` — maybe 30-100 entries, hand-authored or distilled from training
- **Gate** — restraint policy: don't suggest the same intervention twice in one session; honor `contraindications`; only fire when stuckness_signal is high enough
- **Reasoner** — picks the best-fit intervention given conversation state and recent history
- **Action** — surfaces the suggestion to the coach (not to the client) as a card during the session
- **Calibration** — verdict signal is whether the coach used the suggestion; success rate is whether the session unstuck after

Substrate-as-vocabulary projects almost always have human authorship of the repertoire as a load-bearing step. The coach designs the library; the system picks from it. Partnership, not automation.

## A writer maintaining voice consistency across drafts (substrate-as-record, creative)

- **Surface scope** — the voice-checking pass on a finished draft. Generation of new content is not a hybrid loop here (it's the writer writing).
- **Lens** — LLM extracts voice features from finished pieces: `{piece_id, sentence_complexity, lexical_register, rhetorical_moves, cadence_features, stylistic_quirks}`
- **Substrate** — JSONL; one record per piece; the writer's "voice corpus"
- **Gate** — for a new draft, computes deviation from rolling-window average; flags passages that drift significantly
- **Reasoner** — when invoked on a draft, identifies passages reading differently from the writer's baseline and explains the deviation
- **Action** — annotations in markdown comments; never auto-edits
- **Calibration** — writer accepts/rejects each annotation; rejection rate is inverse hit-rate

The lens schema is the writer's *implicit theory of their own voice*, made explicit. Discovering that schema (which fields, which enums) is itself a sub-project worth schema-discovery treatment if the writer has a sufficient corpus.

## A recruiter screening resumes (multiple surfaces in one project)

This project has *three* surfaces, illustrating Phase 2 scope decisions:

- **Surface 1 — resume parsing** (Bucket B, *not* hybrid loop): extracting `{name, education, experience_entries[]}` is mostly deterministic if resume format is consistent. Use a parser library, not an LLM, unless the formats vary wildly.
- **Surface 2 — fit-scoring** (Bucket C, hybrid loop, analytical): lens extracts candidate-criteria-fit `{years_in_role, domain_match, level, signal_strength}`; substrate over all candidates; gate filters below threshold; reasoner ranks for human review; calibration via "did we interview, did they pass."
- **Surface 3 — outreach composition** (Bucket C, hybrid loop, interventional): typed library of message templates `[{template, deploy_when, tone, length}]`; reasoner picks given candidate context; gate restrains template reuse within same-day; action drafts message.

The "project" is one tool; the surfaces are three with different shapes. Realistic case for most non-toy projects.

---

## Picking a template by shape

| Shape | Use when... | Example anchor |
|---|---|---|
| Substrate-as-record (analytical) | Value comes from making sense of accumulated data over time | Teacher's intervention tracker; writer's voice corpus |
| Substrate-as-vocabulary (interventional) | System needs to discriminate the right move from a typed repertoire | Coach's intervention library |
| Both (record AND vocabulary) | Substrate is queried by present moment AND grows over time | Advocacy legislator-tracker (record of positions, vocabulary of issue types) |
| Substrate provider for downstream agents | Output is intentionally typed for other agents to consume | A typed accountability-data corpus exposed via MCP |
| Multiple surfaces in one project | Project has 2-3 distinct fuzzy-judgment places | Recruiter tool above |

---

## When scaffolding, pick the closest match

- Single-domain expert tracking their own observations over time → teacher / parent / writer template
- Civic / policy / accountability tracking with provenance → advocacy template
- Coach / clinician / facilitator with a typed move library → coach template
- Recruiter / triage / customer-service routing → multi-surface (recruiter) template
- Multi-agent or multi-party deliberation → typed-deliberation template (positions, votes, cruxes, cluster discovery)
- Substrate provider for downstream agents → typed-corpus-as-MCP template

Resist inventing a new architecture. Pick the closest existing example, name what's different, design the difference.
