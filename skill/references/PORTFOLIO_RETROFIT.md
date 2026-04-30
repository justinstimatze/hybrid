# Portfolio retrofit — applying the skill to existing projects

Run the skill's diagnostic phases retroactively against the user's existing projects. Two purposes: calibrate the skill against real artifacts (would it have helped?), and surface concrete next moves per project (where does the framework suggest expansion?).

Format per project:
- **Surfaces** — what the skill's Phase 1 would find
- **Buckets** — Phase 2 scope (A: direct LLM, B: pure code, C: hybrid loop)
- **Shape** — analytical / interventional / both
- **What the skill would have surfaced early** — decisions made implicitly that the skill would have made explicit
- **Skill-prescribed gaps** — calibration log, ablation test, deployment-ethics, schema versioning, restraint policy
- **Recommended expansions** — concrete next moves per the framework

---

## slimemold

- **Surfaces:** (1) per-claim extraction from transcript chunks; (2) load-bearing-claim detection over the graph; (3) finding-priority selection for injection.
- **Buckets:** all C.
- **Shape:** analytical (substrate-as-record), with a small interventional bias in the priority selector.
- **What the skill would have surfaced early:** that there are *three* surfaces, not one. Currently they're entangled in `internal/extract/` and `internal/analysis/`. The skill would have suggested separating extraction (lens) from fragility analysis (reasoner) into different modules with their own schemas.
- **Skill-prescribed gaps:** calibration log is partial — load-bearing predictions are recorded but verdict signals are weak. No ablation test. Deployment ethics: single-user / Claude-Code-internal, so power-neutral. Schema versioning is in place.
- **Recommended expansions:** **(a)** instrument calibration: log every "load-bearing" prediction with a verdict pulled from "did this claim still appear load-bearing 5 turns later" or "did the user push back when the finding was injected." **(b)** ablation: compare downstream Claude turn quality with vs. without findings injected. **(c)** add a metabolism phase that audits the claim graph for the same biases winze audits its KB for — slimemold's substrate is the natural second test of the metacog primitive.

## drivermap

- **Surfaces:** (1) extraction of mechanism records from Wikipedia/Kagi text; (2) prediction of mechanisms given person×situation; (3) verbalization of motivation for downstream dialogue.
- **Buckets:** all C.
- **Shape:** both — substrate is record (mechanism corpus) AND vocabulary (taxonomy used by predict_mechanisms).
- **What the skill would have surfaced early:** the substrate-as-vocabulary half is half-built. The mechanism taxonomy is hand-curated; the verbalization library is nascent. The skill would have asked "what's the discrimination criterion per mechanism for the predict step?" and gotten a clearer schema for `deploy_when`-style triggers.
- **Skill-prescribed gaps:** no calibration log on predict_mechanisms (verbalization predictions could be logged with verdict = "did the dialogue actually pattern-match the mechanism"). No ablation test. Deployment ethics: when applied to people who haven't consented, has surveillance shape — needs a `consent_recorded` field. Schema versioning is implicit.
- **Recommended expansions:** **(a)** add `consent_recorded` to all records describing third parties. **(b)** add a verbalize-with-verdict calibration log; the verdict is whether the generated dialogue is later judged to fit the predicted mechanism. **(c)** run dense-style schema discovery on a corpus of dialogues to see whether the discovered mechanism schema agrees with or refines the hand-authored 137 — this is the canonical Conjecture 3 experiment.

## winze

- **Surfaces:** (1) Resolve (sensor signal classification); (2) Trip (speculative connection generation); (3) MCP query primitives; (4) Bias-audit; (5) Calibrate.
- **Buckets:** all C.
- **Shape:** both — KB is record AND vocabulary.
- **What the skill would have surfaced early:** the bias-audit and calibrate phases as *first-class metabolism components* rather than late additions. The skill's "deployment ethics" phase would have asked who owns the KB; for a personal-research-substrate, the answer is single-user, power-neutral.
- **Skill-prescribed gaps:** calibration is the most-shipped of any project here; that's the canonical implementation. Ablation test: missing — does removing typed claim structure and giving raw text to the reasoner change downstream answer quality? No deployment-ethics issues for personal use.
- **Recommended expansions:** **(a)** implement the ablation test against a held-out set of theory-of-consciousness questions — measure answer quality with vs. without typed substrate. **(b)** lift the metacog primitive (the nine bias auditors) into a standalone package. **(c)** publish the calibration JSONL format and resolver as a reference implementation; this is the canonical `cal_log` source.

## plancheck

- **Surfaces:** (1) agent spike (file-prediction lens); (2) structural probes (compiler, comod, reference graph — these are *gates*, not lenses); (3) ranking blend.
- **Buckets:** spike is C. Structural probes are B (deterministic). Ranking is C-adjacent (gate logic).
- **Shape:** analytical.
- **What the skill would have surfaced early:** the structural probes are correctly identified as deterministic (Bucket B), not lenses. The spike is the only LLM lens. The skill's auto-default for the gate is "confidence threshold + chronological"; plancheck's actual gate is sophisticated (novelty-weighted blending) — this is exactly the kind of evolved gate that requires observing over-firing patterns first.
- **Skill-prescribed gaps:** calibration is partial (record_outcome MCP tool exists). No formal ablation. Deployment ethics: single-user dev tool, neutral. Schema versioning in place.
- **Recommended expansions:** **(a)** systematically resolve recorded outcomes against actual diffs to compute hit-rate per spike model and per ranking weight. **(b)** publish a small benchmark dataset (50+ tasks with ground-truth file lists) for ablation against substrate-less LLM file prediction. **(c)** the cascade-risk and forecast tools could become a substrate-as-vocabulary library (typed risk patterns, deploy_when criteria) rather than current ad-hoc form.

## lucida

- **Surfaces:** (1) classifier (discourse_move + cell_type); (2) specialist (spec generation per cell_type); (3) reflect (synthesis every 30 cells).
- **Buckets:** all C.
- **Shape:** mostly analytical; reflect has metabolism flavor.
- **What the skill would have surfaced early:** the two-LLM staging (classifier → specialist) is the right shape and skill would have prescribed it. The skill's "calibration from day one" prescription would have caught that cell quality is currently not tracked.
- **Skill-prescribed gaps:** no quality-of-cell calibration. No ablation. Deployment ethics: single-user dashboard, neutral. Schema versioning in place via cells.json provenance.
- **Recommended expansions:** **(a)** track cell-quality verdicts: thumbs-up/down on rendered cells, OR length-of-time-the-user-keeps-the-cell-rendered as a passive signal. **(b)** ablation: how much does the classifier's discourse_move tag actually help the specialist vs. specialist-given-prose-directly? **(c)** turn reflect.py into a real metabolism phase with bias-audit on the cell stream (over-representing certain visualization types? missing whole categories of discourse?).

## gemot

- **Surfaces:** (1) analyze.run (claim/crux extraction from positions); (2) propose_compromise (reasoner consuming substrate); (3) reframe (mediator); (4) commitment tracking; (5) reputation graph.
- **Buckets:** all C except commitment tracking (B — pure state machine).
- **Shape:** both, with strong vocabulary side (cruxes, opinion clusters as discriminating taxonomy).
- **What the skill would have surfaced early:** the two-engine pipeline (LLM text + vote-matrix PCA) is a substrate-shaping move that the skill's gate definition naturally accommodates. The reputation graph is correctly identified as the long-term metabolism layer.
- **Skill-prescribed gaps:** convergence metrics are calibration-adjacent but not predict-and-verdict shaped. No formal ablation against in-band agent disagreement (no-substrate baseline). Deployment ethics: multi-party by design, with cluster-discovery shaping the substrate — this is the *positive* deployment shape (substrate emerges from participation, not imposed).
- **Recommended expansions:** **(a)** define a verdict signal: "did the proposed compromise actually get adopted?" — log every propose_compromise call with this verdict. **(b)** ablation: compare gemot-mediated deliberation against in-band CrewAI/AutoGen baseline on a fixed dispute corpus. **(c)** the reputation graph is the natural recursive-composition substrate — surface it as a primary product feature rather than a long-term moat.

## publicrecord

- **Surfaces:** (1) finding extraction from primary sources; (2) entity resolution (rebranding-aware); (3) decision-time agent recommendation via riskLevel.
- **Buckets:** all C.
- **Shape:** both — record (findings over time) AND vocabulary (riskLevel enum + finding types as discrimination).
- **What the skill would have surfaced early:** publicrecord is the *canonical recursive-composition example* — explicitly designed for downstream agent consumption. The skill would have flagged this as a "substrate provider" pattern requiring extra discipline on canonical schemas at interfaces.
- **Skill-prescribed gaps:** "data maturity notice" gestures at calibration but no shipped per-finding hit-rate. No ablation against agents-without-publicrecord baseline. Deployment ethics: substrate-for-agents shape is power-neutral; if humans-applying-imposed-riskLevels became a deployment, that would shift.
- **Recommended expansions:** **(a)** ship a calibration log: every agent recommendation that consulted publicrecord, with verdict = "did the user follow it / was the warning acted on." **(b)** ablation: deploy two versions of an agent (with vs. without publicrecord access) on a shared task corpus, measure recommendation quality difference. **(c)** publish the finding schema as the canonical "decision-time RAG record shape" for the broader ecosystem to compose against.

## crowdwork

- **Surfaces:** (1) element extraction (regex + LLM); (2) tension detection (deterministic over elements); (3) frame suggestion (LLM with closed-domain menu); (4) signal injection.
- **Buckets:** elements/tensions/frames are C. Tension detection itself is B (deterministic rules). Signal injection is the action layer.
- **Shape:** mostly interventional (substrate-as-vocabulary: the closed-domain menu of frames). Element extraction is record-flavored.
- **What the skill would have surfaced early:** the two-phase staging (element extraction → frame suggestion) is the right move. The closed-domain menu (10 non-obvious domains) is exactly the substrate-as-vocabulary shape the skill prescribes.
- **Skill-prescribed gaps:** no calibration log. No ablation. The ripeness window (peak at 5 turns, decay by 14) is a sophisticated gate that the skill would have prescribed observing over-firing first; this gate clearly evolved from real over-firing observations.
- **Recommended expansions:** **(a)** track whether Claude actually used the injected signal in its response (text-overlap heuristic, or explicit flag in subsequent turns). **(b)** ablation: turn off signal injection for half of test turns; compare humor density and quality. **(c)** discover the frame-domain library via dense-style schema discovery on a corpus of successful humor exchanges.

## groupchat

- **Surfaces:** (1) meme selection given conversational context (LLM matches roster to moment); (2) cooldown gate; (3) rendering pipeline.
- **Buckets:** selection is C. Cooldown and rendering are B.
- **Shape:** purely interventional. Substrate is the curated meme roster (66 entries with rich metadata).
- **What the skill would have surfaced early:** the substrate-as-vocabulary shape with hand-curated repertoire is exactly the pattern; skill would have prescribed the "deploy_when / too_much_if / mechanism / key" schema groupchat already has. The adaptive cooldown θ(t) = e^(−λt) is a model gate that the skill would have noted as evolved from observing over-firing.
- **Skill-prescribed gaps:** no calibration log on whether dropped memes actually land. No ablation against "drop randomly" or "drop never" baselines. Deployment ethics: single-user / personal terminal use, power-neutral.
- **Recommended expansions:** **(a)** every drop_meme call logs "did the user laugh / acknowledge / push back" — verdict signal is weak but not zero (text reactions in subsequent turns). **(b)** ablation: A/B between current selection logic and random selection from the roster; measure user-keep rate or laughing emoji rate. **(c)** the meme-info / list-memes endpoint is the natural place to add metacog-style audit (over-using mech=I, under-using mech=R; clusters of memes that never fire).

## lamina/poc/dense

- **Surfaces:** (1) compress (lens producing notation from spec); (2) verify (adversarial verifier scoring fidelity); (3) refine (notation evolution given scores).
- **Buckets:** compress is C. Verify is C. Refine is C-adjacent (uses LLM in the loop).
- **Shape:** the *meta-shape* — dense is itself the schema-discovery loop. Substrate is the evolving notation grammar.
- **What the skill would have surfaced early:** dense is the canonical example of *Conjecture 3* (schema discovery for cognitive schemas). The skill would have explicitly named the wake/abstraction loop as DreamCoder lineage and recommended cataloging the discovered notations for cross-domain transfer.
- **Skill-prescribed gaps:** combined_score is the calibration-equivalent here, but no per-discovery-round hit-rate aggregated over time. No baseline against random schema mutation. Deployment ethics: single-user research, neutral.
- **Recommended expansions:** **(a)** apply dense to a non-program-adjacent domain (drivermap mechanism corpus, slimemold claim-graph topology) — this is the canonical Conjecture 3 falsification experiment. **(b)** publish the discovered notations + discovery transcripts as a public artifact ("DreamCoder for descriptive notation"). **(c)** package the compressor+verifier+refiner loop as the `schemaforge` primitive in `PRIMITIVES.md`.

---

## Common gaps across the portfolio

Looking across the retrofits, the recurring missing piece is **calibration**. Every project has a partial gesture (slimemold's load-bearing logging, plancheck's record_outcome, gemot's convergence metrics, publicrecord's data-maturity notice, lucida's cost stats) but none has a per-evaluator predict-and-verdict log with hit-rate aggregation that meaningfully shapes development decisions. **This is direct empirical support for Conjecture 1** — if the user's own portfolio doesn't have it, and wesen's body of work doesn't have it (per `PRIOR_ART.md`), the gap is real.

The second recurring gap is **ablation tests**. Only winze has the discipline implicitly through metabolism predictions; nothing else has a documented "performance-with-substrate vs. performance-without-substrate" test. This is the skill's load-bearing prescription against the "you're just using LLMs with extra steps" critique. **Adding one ablation test to one existing project would be the highest-leverage local action** the user could take to make the framework real.

The third gap is **canonical schemas at composition interfaces**. Several projects have records that *could* compose (slimemold claims, drivermap mechanisms, publicrecord findings) but none of them share a canonical claim/finding/mechanism schema. Recursive composition (per `STACKING.md`) requires this; without it, composition is one-off integration. **A v0.2 of this repo could publish a canonical record-schema spec** that the projects above adopt incrementally.

## Where the skill would have changed past decisions

If the skill had existed when the user started building these projects:

- **slimemold** — the surface-separation step would have led to cleaner module boundaries between extraction, analysis, and injection; calibration would have been wired in at the start.
- **drivermap** — `consent_recorded` would have been in the schema from day one (not retrofitted later).
- **plancheck** — the deterministic-vs-LLM bucket assignment would have been explicit, preventing some early conflation.
- **lucida** — cell-quality calibration would be live now.
- **groupchat** / **crowdwork** — the substrate-as-vocabulary shape would have been named, making the closed-menu / cooldown choices feel deliberate rather than discovered.
- **dense** — the Conjecture-3 framing would have positioned the work in the DreamCoder/LILO lineage from the start, sharpening the publication narrative.

If the skill had existed but the user had ignored it: probably half of these effects, because some of these decisions emerged from the project domain rather than from architectural foresight. The skill's value is more in *catching the missed surface* (e.g. lucida's cell-quality calibration) than in changing well-thought-through choices.

## Where the skill confirms the existing portfolio is well-shaped

Most of the architectural choices in the existing projects are correctly hybrid-loop-shaped. The skill is recovering and naming what the user has been doing rather than prescribing radical changes. This is good news — it means the framework is descriptively accurate of working code — and a calibration concern: a framework that endorses everything its author has already done is not falsifying anything. The conjectures in `../README.md` are the falsification surface; the retrofit above is the descriptive surface.
