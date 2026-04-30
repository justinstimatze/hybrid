# Worked examples — existing repos mapped to the five-role schema

Use these as anchors when explaining the architecture or when scaffolding a new project that resembles one of them.

## slimemold (epistemics on Claude Code transcripts)

- **Lens** — `internal/extract/extract.go`: Claude tool_use call extracting per-claim records (basis, edges, Moore flags) from transcript chunks
- **Substrate** — typed claim graph persisted across sessions; age decay built into storage
- **Gate** — cooldown decay, cold-start floor, age-based priority selection, per-claim cooldown
- **Reasoner** — structural-fragility analysis (eight vulnerability types) over the claim graph
- **Action** — system-message injection during Stop hook, biases future generation
- **Calibration** — partial; load-bearing predictions get logged but verdict signals are weak
- **Metabolism** — none yet

## drivermap (behavioral mechanism KB)

- **Lens** — extract.py: blind-then-guided two-phase Claude extraction of mechanism records from Wikipedia/Kagi text; schema empirically derived from 20 blind extractions
- **Substrate** — ~137 mechanism JSON records + scoring engine
- **Gate** — schema enforcement, prompt-repetition technique for causal bidirectionality, post-extraction verifier
- **Reasoner** — predict_mechanisms (deterministic scoring on person×situation), then verbalize_motivation (LLM)
- **Action** — MCP tool returns predictions; demo.py composes them into dialogue
- **Calibration** — none currently; verbalization predictions could be logged
- **Metabolism** — none

## winze (typed epistemic substrate with sleep cycles)

- **Lens** — Resolve phase (Claude Sonnet classifies sensor signal as corroborated/challenged/irrelevant)
- **Substrate** — typed Go AST as KB; entities, claims, predicates, theories
- **Gate** — bias-audit gates control which phases fire (availability heuristic, survivorship bias)
- **Reasoner** — MCP server exposing claims/disputes/provenance/search/stats/theories to other LLMs
- **Action** — promote claims, schedule sensor queries, write to corpus
- **Calibration** — `.metabolism-calibration.jsonl` per-cycle (this is the canonical implementation; copy from here)
- **Metabolism** — full: dream (consolidation), trip (speculative connections), evolve (sensor), bias-audit (KB self-check)

## plancheck (file-prediction for AI coding agents)

- **Lens** — RunAgentSpike (tool-using Claude agent that explores codebase) producing AgentFile records with confidence scores
- **Substrate** — AgentResult + structural probe outputs (compiler blast radius, comod history, reference graph)
- **Gate** — novelty-weighted confidence ranking (struct weight 0.5→0.1, semantic 0.1→0.4 as novelty rises)
- **Reasoner** — rankCandidateFiles blends spike + structural signals
- **Action** — file-list output to user, optional record_outcome for calibration
- **Calibration** — record_outcome and record_reflection MCP tools (partial)
- **Metabolism** — none

## lucida (live notebook from Claude Code transcript)

- **Lens** — classifier (Haiku, prompt-cached): discourse_move, cell_type, confidence; specialist (Sonnet): produces concrete spec
- **Substrate** — cells.json with provenance per cell
- **Gate** — confidence threshold (>0.8 mint, 0.6-0.8 draft, <0.6 suppress)
- **Reasoner** — reflect.py every 30 cells consumes the cell stream and synthesizes
- **Action** — frontend renders cells; reflection cells appear inline
- **Calibration** — cost/cache stats per cell; cell quality not yet tracked
- **Metabolism** — partial: reflection synthesis has dream-flavor

## gemot (multi-agent deliberation)

- **Lens** — analyze action: LLM extracts taxonomy, claims, cruxes from agent positions (parallel claim extraction)
- **Substrate** — positions, votes, cruxes, clusters; vote matrix as separate substrate (PCA/SVD/k-means)
- **Gate** — round-based protocol; vote-matrix analysis; commitment fulfillment tracking
- **Reasoner** — propose_compromise, reframe (LLM consumes substrate to produce synthesis)
- **Action** — return analysis to participating agents; track commitments
- **Calibration** — convergence metrics across rounds (built in)
- **Metabolism** — reputation graph across deliberations (the long-term moat)

## publicrecord (decision-time RAG for accountability data)

- **Lens** — LLM extracts findings (entity, relationship, severity, verbatim quote) from primary-source documents (court records, enforcement actions, regulatory orders)
- **Substrate** — SQLite DB; 4,656 entities, 6,377 findings, relationships up to 3 levels deep; rebranding-aware entity resolution
- **Gate** — `riskLevel` enum (low/high/medium); temporal awareness; provenance enforcement
- **Reasoner** — *not internal*; downstream agents consume via MCP/REST/SDK and make decisions ("recommend a vendor", "due diligence")
- **Action** — structured response with citations; agent decision flow follows `riskLevel` (stop / offer alternatives / note)
- **Calibration** — partial; "data maturity notice" acknowledges summaries are AI-synthesized and under review
- **Metabolism** — none currently; would benefit from re-extraction as new public records are filed

Canonical example of **recursive composition**: publicrecord is a hybrid loop whose output is intentionally designed to be consumed by other hybrid loops. The substrate (typed entities, riskLevel enum, citations) is shaped for downstream agent consumption, not for end-user reading. This makes it a *substrate provider* rather than a standalone tool.

If a new project is intended as a substrate provider for other agents (rather than an end-user tool), use publicrecord as the template.

## When scaffolding a new project, pick the closest example as the template

- Analyzing transcripts/dialogues for patterns → slimemold-shaped
- Building a typed corpus from documents → drivermap-shaped
- KB that should grow over weeks/months → winze-shaped
- Predicting/scoring something for an AI coding workflow → plancheck-shaped
- Real-time visualization or augmentation of a stream → lucida-shaped
- Multi-agent or multi-party deliberation → gemot-shaped
- Substrate provider for other agents (typed records intended for downstream loops) → publicrecord-shaped

Most new projects will resemble one of these closely enough that the corresponding template is a good starting point. Resist the urge to invent a new architecture; reuse and refine.
