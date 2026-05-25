# Procedural infrastructure for hybrid loops

The main skill names the *substantive* layers — what each block does, what shape the data takes, where the cycle closes. This file names the *procedural* infrastructure that surrounds them — what must happen before a layer is invoked, what constraints bind its output, what review wraps around its decisions.

The distinction is borrowed from legal systems, which have accumulated procedural scaffolding around soft-judgment/hard-rule reasoning over centuries of practical use under real stakes. *Substance is what gets decided. Procedure is the discipline around how it gets decided.* That dichotomy is itself contested in legal philosophy — the American *Erie* doctrine is the canonical battleground over how impossible it is to maintain a clean line, and Lon Fuller's *Morality of Law* argues the two interpenetrate at the root. This doc uses the distinction as a working partition, not a settled one. Even so, the rough partition is generative: legal systems have figured out that procedural patterns are at least as load-bearing as substantive ones — a substantively correct outcome from a procedurally broken process is the legal equivalent of a kangaroo court. The same applies to hybrid loops: a well-designed substrate and a competent reasoner can still produce untrustworthy systems if the procedural scaffolding around them is missing.

The lineage drawn on here is specifically Anglo-American common law procedure. Civil-law systems (most of continental Europe, Latin America, Japan, and beyond) handle several of these problems with different infrastructure — *stare decisis* is barely a doctrine in civil-law tradition; standing works differently under inquisitorial procedure; recusal mechanisms vary substantially; the codified case-or-controversy requirement is parochially American. The engineering content of the protocols below transfers across traditions, but the vocabulary and citations are tradition-specific. Treat them as one well-developed lineage of procedural-infrastructure design, not as the only one.

This doc adds five protocols. Two extend layers that already exist in the main skill; three are genuinely new. None of them is a substantive pattern — they are constraints on how the substantive patterns compose. Most loops need none of these — reach for this doc only when a loop is high-stakes, runs at scale, or has started failing in ways the substantive layers don't explain. Below that bar, the procedural overhead doesn't pay back.

> *Vocabulary note: legal procedure has accumulated specialized terms (standing, recusal, stare decisis, limited remedies) over centuries. This doc uses them where they fit cleanly and aliases them where they don't. The main skill's vocabulary (lens, substrate, gate, reasoner, action, calibration, metabolism) stays primary; legal terms are pointers to a lineage, not replacements.*

## 1. Stratified standards of proof

*Different action consequences demand different confidence thresholds.* Legal systems run multiple standards — preponderance of evidence for civil judgments, clear and convincing evidence for fraud or termination of parental rights, beyond reasonable doubt for criminal conviction. Courts have explicitly declined to assign these standards specific probability values; surveys of federal judges (canonically C.M.A. McCauliff's 1982 study, since replicated) find wide variance in how individual judges quantify them, and a strand of legal scholarship argues that beyond reasonable doubt is qualitatively different from a high preponderance rather than merely further along the same axis. The point worth porting isn't the specific numbers — it's that the standards are *ordered* by stakes, and the system explicitly recognizes that a single "act vs. don't act" threshold is wrong.

**Slot:** extension to *action*. The main skill mentions confidence thresholds in the gate's defaults but doesn't tier them to action reversibility.

**What it adds:** an explicit tiering of action consequence to required confidence. Suggested taxonomy:

- *Reversible-suggestion tier:* "I think this would help." Low cost if wrong. Examples: highlighting a candidate, surfacing a tag, drafting a message for review.
- *Persistent-edit tier:* file edits, refactors, records committed to the substrate. Reversible but expensive to undo. Examples: a calibration verdict that updates the schema, an autonomous code change.
- *Irreversible-action tier:* `rm -rf`, `force-push`, schema migrations, external API side effects that can't be retracted, public-facing posts. Examples: deleting from the substrate, sending an email, applying a database migration.

The tiers are *ordered* — each tier requires the loop's confidence in its own judgment to be higher than the tier below — but the calibration of what counts as "enough" for each tier is project-specific and gets refined against the calibration log. The discipline is *making the tiers explicit and refusing to act on irreversible-tier actions with reversible-tier confidence*, not assigning specific numbers up front.

**Diagnostic:** if the loop uses one confidence threshold for all actions, it's mis-calibrated. The fix is at the action layer, not the reasoner.

**Anti-diagnostic:** pure-suggestion loops where nothing is committed don't need this — overhead doesn't pay back. Single-tier projects are fine if all actions are in the same tier.

## 2. Standing (runtime extension to gate)

*Not every input deserves invocation.* Legal standing asks three questions before the reasoner is even invoked: was there injury-in-fact, is it traceable to the defendant, can the court actually redress it. The system refuses to decide cases where these aren't met — not because the reasoner can't, but because the reasoner shouldn't.

The main skill's Phase 2 (scope each surface as A/B/C) is design-time standing — does this surface warrant the hybrid pattern at all. What's missing is the *runtime* version: does this specific invocation warrant the loop firing.

**Slot:** extension to *gate*, specifically a pre-reasoner check distinct from the existing post-reasoner output filtering.

**What it adds:** an explicit gate that runs before the lens, asking three questions of the input:

- *Is the input real?* (injury-in-fact) — is there actual content to reason about, or is the loop being invoked on whitespace, on a test fixture, on cached output that already exists?
- *Is the reasoner the right address?* (traceability) — could substrate alone answer this, or does the user actually need a deterministic check or a human?
- *Can the action layer redress it?* (redressability) — if the reasoner produces output, does the action layer have authority to act on it, or will it just produce text that goes nowhere?

If any answer is no, the loop should not fire. Return early, route to a fallback, or surface the standing failure as the output.

**Diagnostic:** most coding-agent failures attributed to "the reasoner was wrong" are actually standing failures — the reasoner was invoked on an input it shouldn't have run on. Look for these before tuning the reasoner.

**Anti-diagnostic:** in agentic loops where the agent is expected to handle whatever input it gets, hard standing checks suppress the value. The standing check belongs at user-facing surfaces, not at internal loop edges.

## 3. Precedent retrieval (extension to substrate)

*The reasoner is constrained by past decisions on similar inputs.* Stare decisis: courts must follow binding precedent. This is structurally distinct from generic substrate read — it's not "retrieve semantically similar context" but "retrieve past decisions this decision must be consistent with."

The main skill's substrate-as-record covers the storage shape and the reasoner-reads-recent default covers basic retrieval. What's missing is the *consistency obligation* — the reasoner being told not just "here is relevant context" but "you are bound to decide consistently with these prior records unless you can articulate why this case differs."

**Slot:** extension to *substrate*, specifically a retrieval mode distinct from semantic similarity or chronological recency.

**What it adds:** an explicit precedent retrieval step where the reasoner must surface and consider past decisions on structurally similar inputs *before* producing output. The retrieval function isn't "what's semantically close" — it's "what past decisions does this decision need to be consistent with."

The output of the precedent step has to feed into the reasoner's prompt as a binding constraint, not just as additional context. The reasoner is asked: *(a)* is this case materially distinguishable from the retrieved precedents, and if not, *(b)* does your proposed decision follow them. If it doesn't, the reasoner has to explicitly articulate the deviation (the legal version: *distinguishing* the case).

Two constraints keep this from backfiring. Precedent must bind to *verified* outcomes, not the reasoner's own prior outputs — otherwise the loop anchors on its own unverified history and a first-cycle error ossifies into binding "precedent," a self-justifying drift. And consistency is necessary but not sufficient for a *correct* substrate: a reasoner perfectly consistent with a biased record produces stable bias, not quality. Precedent retrieval enforces consistency; whether the precedent set itself is sound is a substrate-quality question for the metabolism layer, not something the consistency obligation can settle. In consequential domains (anything that scores or sorts people) this distinction is the whole game: consistency guards against arbitrary inconsistency (in legal terms, disparate *treatment*) but does nothing about a uniformly-applied biased criterion (disparate *impact*) — procedural fairness over a biased baseline is exactly how substantive bias persists.

**Diagnostic:** when consistency across the substrate matters (architectural style, naming conventions, classification taxonomies, scoring rubrics), substrate access has to include precedent retrieval — not just relevant context retrieval. Loops that produce inconsistent decisions on structurally similar inputs are missing this step.

**Anti-diagnostic:** for genuinely novel problems where past decisions are uninformative, precedent retrieval is noise. Also skip in loops where the substrate is too small (~<20 records) for precedent to mean anything.

## 4. Bounded action (limited remedies)

*The action layer is bounded by what was requested.* Courts can only grant the relief asked for, within their jurisdiction. They cannot impose remedies the parties didn't request. The principle: action is constrained by input, not by what the reasoner thinks would also be good.

The single most common coding-agent failure mode is scope drift — the agent fixes the bug *and* refactors three unrelated files *and* updates the README *and* renames variables it didn't like. Each of those might be correct in isolation; together they violate the limited-remedy principle.

**Slot:** extension to *action*, as a constraint relating action output to lens input.

**What it adds:** an explicit gate at the action layer that compares what was requested in the input (the lens's read of the prompt) with what the action is about to do. Anything in the action that isn't traceable to the input is either rejected (strict mode) or surfaced as a suggestion rather than committed (advisory mode).

Adjacent improvements still get visibility — they're recorded in the substrate, surfaced in the next reasoner cycle, presented to the user — but they don't take effect in the action layer of the current cycle. The discipline: *suggesting and acting are different things*. Legal analog: an advisory opinion isn't a judgment.

**Diagnostic:** if the action's output set is larger than the input's request set, scope is drifting. Tighten the action gate, don't make the reasoner smarter.

**Anti-diagnostic:** exploratory loops where unbounded scope is the point (research agents, brainstorming surfaces) lose value with a strict bounded-action rule. Use advisory mode there at most.

## 5. Recusal protocol (cross-layer constraint)

*Some reasoner-input combinations are structurally compromised regardless of reasoner quality.* Judges recuse from cases where they have known bias — financial interest, prior involvement, family relationship. The system acknowledges that the right answer here isn't "try harder" but "this reasoner shouldn't be the one deciding."

**Slot:** cross-layer protocol that binds before lens and constrains reasoner choice. Closest existing concept is the deployment power-balance check, but recusal is about reasoner-input compatibility, not deployment shape.

**What it adds:** an explicit enumeration of inputs the reasoner should not run on, with routing to a fallback (a different model, a human, a different procedure):

- *Self-evaluation* — the reasoner doesn't grade its own output. The critic in a critique loop is a different instance, ideally a different model.
- *Known-failure inputs* — patterns the project has documented as reasoner blind spots route to fallback rather than getting another attempt.
- *Conflict of training* — when the input is plausibly in the reasoner's training data in a way that would bias output (asking a model to evaluate its own provider's product, asking about its own training corpus), the result is suspect regardless of reasoner quality.
- *Affective interference* — the reasoner is asked something that activates known sycophancy patterns (high-stakes emotional framing, user-stated belief that pre-commits to a position). Route through a deliberately adversarial framing.

The protocol is a checklist that runs at the gate, not a piece of reasoning. Either the input matches a recusal condition (deterministic) or it doesn't; the reasoner is never asked "should you recuse yourself" because that question has the same sycophancy failure mode it's meant to prevent.

**Diagnostic:** any time a loop's failure mode is "the reasoner was too agreeable / too confident / too aligned with the user's frame," recusal is missing or under-specified. The fix isn't a smarter reasoner; it's a route-around.

**Anti-diagnostic:** low-stakes loops where recusal-overhead exceeds failure-cost can skip. Recusal protocols are infrastructure investments that pay back at scale.

## Integration with existing meta-layers

The five protocols above interact with the two meta-layers already in the main skill:

**Calibration log.** The calibration log already records predict + verdict per loop cycle. The procedural protocols add new things worth recording:

- Standing failures (input rejected before lens) as a distinct verdict class. A loop with a high standing-rejection rate is over-invoked, which is a different problem from a high reasoner-error rate.
- Precedent deviations (cases the reasoner distinguished from retrieved precedent) as flagged records for review. A loop that distinguishes precedent frequently is either handling genuinely novel cases or is drifting from its own past decisions; the calibration log surfaces which.
- Recusal triggers as their own log stream. A loop whose recusal protocol fires often is either over-specified (too many false-positive recusals) or operating in a domain that's largely outside its competence.

**Metabolism.** The substrate-wide audit already looks for drift. The procedural protocols give it more things to check:

- Are the standards-of-proof tiers calibrated against actual outcomes? A loop whose irreversible-tier actions have the same error rate as its reversible-tier actions has mis-calibrated tiers.
- Is the precedent set growing without coherence — multiple precedents that contradict each other on structurally identical inputs? That's a substrate-shape problem, not a per-decision problem, and metabolism is where it gets caught.
- Are the recusal conditions getting hit on inputs that look very similar to inputs the loop *should* be handling? That's a sign the recusal definitions are over-broad and should be tightened.

## What this doc does not claim

The legal-procedural analogy is generative but not exact. Judges have real-world consequences for their decisions — career, impeachment, reputation — and that backpressure shapes the whole system's incentives. Reasoner instances do not face that backpressure, and neither half of the courtroom's accountability transfers: the subject of an automated adverse action, unlike a litigant, has no motion to file. The five protocols above port what's portable; they do not claim to reproduce the accountability infrastructure that makes legal procedure work. A perfectly-procedured but unaccountable loop — procedure outrunning accountability — is precisely the object this caveat warns the analogy cannot produce; the more consequential the action, the more the irreversible tier belongs with a human.

The protocols also don't replace substantive design. A well-procedured loop with a bad substrate is still a bad loop. Procedure is necessary scaffolding around competent substance, not a substitute for it. *Substrate brings discrimination, code brings restraint, procedure brings discipline.* All three are required at scale.

And like the rest of the skill, these protocols are reasoned conjecture, not ablation-validated: by the skill's own standard, a procedural overlay earns its keep only when a loop that has it measurably beats one that doesn't. Treat them as a hypothesis to test, not a result.

## Citations

- **Federal Rules of Civil Procedure** and **Federal Rules of Evidence** — the canonical American codifications of procedural infrastructure. Worth reading at least the table of contents as a checklist of "what procedural decisions has the legal system already made."
- **Article III, U.S. Constitution** — origin of the standing doctrine; *Lujan v. Defenders of Wildlife* (1992) is the modern three-prong articulation (injury-in-fact, traceability, redressability).
- **28 U.S.C. § 455** — federal judicial recusal statute; the canonical list of structurally compromised reasoner-input combinations.
- **Stare decisis** — common law principle; *Burnet v. Coronado Oil & Gas Co.* (1932), Brandeis dissent, is the canonical articulation of when precedent should and shouldn't bind.
- **Lon Fuller, *The Morality of Law*** (1964) — articulates eight "principles of the inner morality of law": *generality* (rules apply to classes, not single cases), *promulgation* (publicly knowable), *prospectivity* (no retroactive rules), *clarity*, *non-contradiction*, *possibility of compliance*, *constancy through time*, and *congruence between the declared rule and official action*. Fuller's stronger claim is that these procedural-formal requirements are *themselves* substantively moral — a "law" that violates enough of them isn't a defective law, it's not law at all. The mapping to hybrid loops is direct and worth taking seriously: schemas should apply consistently across similar inputs (generality), be inspectable (promulgation), not change in ways that retroactively invalidate past records (prospectivity, constancy — already in scope of the metabolism layer), be readable (clarity), not contradict themselves (non-contradiction — addressed by precedent retrieval), be possible to actually comply with (calibration log surfaces violations of this), and match what the system actually does (congruence — addressed by the calibration log's predict-versus-verdict structure). Fuller's claim that procedural integrity is itself a property of the system rather than overhead on top of it is the deepest available defense of why a procedurally-broken hybrid loop isn't merely a suboptimal system but a fundamentally different and worse kind of object.

## See also

- Main `SKILL.md` — substantive layers (lens, substrate, gate, reasoner, action) and the calibration / metabolism meta-layers
- `references/THE_CASE.md` — the algebra-vs-alphabet-vs-disciplines argument; procedural infrastructure is closest to the *disciplines* leg
- `references/BUILDING_BLOCKS.md` — primitive blocks the procedural protocols constrain
- `references/EXAMPLES.md` — worked examples; the protocols here would be layered on top of any of them in production. The recruiter fit-scoring overlay there shows all five applied to one surface in compact form, with the depth kept in this file
