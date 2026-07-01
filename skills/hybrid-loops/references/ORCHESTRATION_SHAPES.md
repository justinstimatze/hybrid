# Orchestration shapes — where the framework's discipline lands in 2026 practice

Multi-agent orchestration in 2026 — broadly construed to include parallel-session multiplexing where the human is the orchestrator, not just genuinely autonomous multi-agent systems — has converged on four distinct shapes. Each elaborates calibration discipline differently because each is responding to a different failure mode. This doc names the four, shows where the framework's existing protocols already apply (sometimes in practitioners' vocabulary, sometimes not), and tracks where the framework currently lacks vocabulary.

This is a vocabulary + research-notes doc, not an adoption guide — use it to recognize which shape a system is in, not to choose whether to build one. The picker table at the end is coarse.

The observed pattern across all four: **each shape's calibration apparatus scales with how far it pushes the autonomy slider away from the human.** Shape 2 (human-as-attester) needs the least apparatus because the human IS the apparatus; shape 3 (LLM-pack-as-attester) needs watchdogs + escalation + per-pack adversarial verification; shape 4 (peer-network-as-attester) needs trust grades, attestation ledgers, anti-collusion topology. Noticed pattern across the four exemplars below, not proven law.

A complementary framing from Maggie Appleton's January 2026 essay on Gas Town — once execution gets cheap, *"design becomes the limiting factor: imagining what you want to create and then figuring out all the gnarly little details required."* The framework's discipline is what makes design-as-the-bottleneck actually tractable; without it, the cheapness-of-execution is just throughput without direction.

**Design dimension above the four shapes — Karpathy's autonomy slider.** The four shapes sit at different positions along a continuous dimension that Karpathy named in a July 2025 talk:

> "You are in charge of the autonomy slider, and depending on the complexity of the task at hand you can tune the amount of autonomy that you're willing to give up for that task."
> — Karpathy, *Software in the era of AI*, YC AI Startup School, July 2025 (https://www.youtube.com/watch?v=LCEmiRjPEtQ)

The four shapes below are stable design points along that slider, each with its own failure mode and its own engineered response.

---

## The four shapes at a glance

| # | Name | Exemplar | Who carries the judgment | Calibration locus |
|---|---|---|---|---|
| 1 | Naive autonomous decomposition | (cautionary; no recommended exemplar) | The LLM orchestrator | None engineered; the gap bites |
| 2 | Human-orchestrated parallel multiplexer | Conductor (Charlie / YC S24); Karpathy's daily practice | The human | Gut-feel, short feedback loop, macro-action review |
| 3 | Engineered-resilience autonomous w/ per-pack adversarial verification | Gas Town → Gas City (Yegge) | Distributed across infrastructure | Watchdogs, refinery, multi-agent crews |
| 4 | Socially-validated federation | Wasteland (Yegge et al.) | A trust-graded peer network | Stamps on an append-only ledger |

---

## Underlying invariant: rails vs. off-rails

Karpathy's "jaggedness" framing names the same gap the framework's calibration discipline addresses:

> "You're either on rails and you're part of the super intelligence circuits or you're not on Rails and you're outside of the verifiable domains and suddenly everything kind of just like meanders."
> — Karpathy, *No Priors* interview, Q1 2026 (post-December capability flip, pre-May Anthropic move)

The mapping is the framework's, not Karpathy's: he named the bimodal observation; the framework's discipline (substrate, gate, bounded action, verification primitives) is one answer to *how to keep loops on the verifiable side of that line*. Each of the four shapes below is a different engineering of those rails. Treat the alignment as confirmatory rather than as an endorsement — Karpathy doesn't know about hybrid-loops; we're noting that his observation has the shape the framework's discipline is built for.

---

## Shape 1: Naive autonomous decomposition

The orchestrator-as-LLM picks subtasks, fans out, synthesizes back. Anthropic's own self-assessment names the gap:

> "Large performance gaps persist when it comes to Claude exercising judgement in choosing goals in both engineering and research."
> — Anthropic, *When AI builds itself*, § Evidence from within Anthropic (May 2026)

The judgment is required at the *decomposition* step, before any subagent can be useful. The article also names the complementary positive form:

> "An area of human comparative advantage, for now, is research taste and judgment, including choosing which problems matter, which results to trust, and when an approach is a dead end."
> — same source, same section

**Failure mode in production** (per practitioner reports from earlier-vintage attempts): the system handles simple tasks but cannot break complex tasks into simple ones; costs scale with attempts not results.

**Framework verdict (revised, with appropriate hedging):** shape 1 without one of the scaffolding strategies the other three shapes embody tends to hit the article's gap. The evidence we have is mostly negative — practitioner reports of earlier-vintage attempts not panning out — and a single positive existence proof in Karpathy's daily setup, which is closer to shape 2 than shape 1 (the human is still the decomposer; agents do tasks the human assigns). The honest position is: shape 1 alone is hard, and the reliable responses we know of are either *human-as-decomposer* (shape 2, with Karpathy at the operator-skill ceiling) or *engineered-resilience infrastructure* (shape 3, Gas Town class). Whether a third path opens as model capability advances is an open question; today we don't have a worked example.

**Update, 2026 — a partial candidate for the "third path."** Anthropic's Claude Code dynamic workflows (Shihipar & Bidasaria, *"A harness for every task: dynamic workflows in Claude Code,"* blog, Jun 2 2026) have Claude itself author the decomposition script before any execution happens — the orchestrator being the LLM, which is exactly shape 1's structural definition — but the authored script then runs under shape-3-style engineered-resilience conventions (adversarial-verification quorums, tournament judging, hard isolation via separate subagent context windows) that the harness's own documentation steers the authoring model toward by default, and names as the antidote to three concrete failure modes (*agentic laziness*, *self-preferential bias*, *goal drift*) rather than as generic best practice. Whether this counts as the third path, or is better read as shape 1 with shape 3's disciplines re-authored per-task instead of built once into standing infrastructure, is exactly the open question this section flagged before the feature existed. There is no independent ablation showing it closes the "judgment in choosing goals" gap the Anthropic RSI article names above; the evidence (the substrate provider's own blog, plus one high-profile user's reported migration) sits at the same self-report evidentiary tier as that RSI citation. Filed as *suggestive*, not resolved — see `PRIOR_ART.md` and `STACKING.md` §"recursive harness authoring" for the fuller citation.

---

## Shape 2: Human-orchestrated parallel multiplexer

**Exemplar A:** Conductor (Charlie, co-founder; YC S24; desktop app for orchestrating coding agents on macOS). Source: practitioner walkthrough at https://www.youtube.com/watch?v=fQmlML9Lay4.

*Not to be confused with Netflix Conductor, the workflow orchestrator listed in `PRIOR_ART.md` under workflow orchestration.*

**Exemplar B:** Karpathy's own daily practice as of Q1 2026 — same shape, higher throughput:

> "Code's not even the right verb anymore... I have to express my will to my agents for 16 hours a day. Manifest. I don't think I've typed like a line of code probably since December. The name of the game now is to increase your leverage. I put in just very few tokens just once in a while and a huge amount of stuff happens on my behalf."
> — Karpathy, *No Priors* interview, Q1 2026

The human is the decomposer; the multi-agent infrastructure is fan-out for parallel work the human directs and reviews. The article's "judgement in choosing goals" gap doesn't bite because the human carries it. Karpathy's setup (multiple codex agents on multiple repo worktrees, macro actions over the codebase, review-not-write) is shape 2 at the highest operator skill anyone has documented — not a different shape.

**Framework-aligned moves, in practitioners' vocabulary:**

- **"Slot free zones" / "do not touch if you are an AI"** comments (Charlie) — file-level permission management. *Cousin to bounded action* but not the same protocol: bounded action constrains what an autonomous loop can *enact*; slot-free zones constrain which files the loop can *touch*. Same family (constrain by exclusion), different protocol.
- **"Don't let the AI be your architect"** (Charlie) — a folk heuristic for the substantive/procedural split. Rhymes with the framework's distinction; the mapping is suggestive rather than rigorous.
- **Macro actions** (Karpathy) ≈ the autonomy slider dialed per-task; review at the diff level, not the line level
- **Gut-feel calibration via short feedback loops** (Charlie) — a *deployment shape* for the single-domain-expert calibration shape from `EXAMPLES.md`, with the founder as the calibration loop. The match is structural-but-loose: every short-loop human practice arguably fits this template, so the mapping has limited predictive power.

**Distinctive claim worth wrestling with:**

- **"Code is almost like sawdust now"** (Charlie) — the durable artifact is the prompt (and the spec it encodes), not the generated code: *"when the next generation of models come out, you can just rerun your prompts again and you'll get new code."* Tension with the framework's substrate-as-durable-record concept. Resolution: in shape 2, prompts (or the spec they encode) are substrate; code is regenerable artifact. See `THE_CASE.md` §77–81.

**The trust-account framing as a bridge to shape 3.** Simon Willison, summarizing the 2026 enterprise picture: *"shipping code faster than engineers can read it... you are making withdrawals from a trust account."* (https://simonwillison.net/, May 27 2026). Charity Majors makes the bridge concrete: *"What would it take for you to feel comfortable shipping code without reading it?"* — then engineer the prerequisites. Shape 2 is sustainable when the verification gets dialed back to one fast loop the human can manage; when the throughput exceeds what one human can verify, you need shape 3.

---

## Shape 3: Engineered-resilience autonomous, with per-pack adversarial verification

**Exemplar:** Gas Town → Gas City (Steve Yegge, https://steve-yegge.medium.com/). Gas Town: original autonomous orchestrator (Go, v1.0 April 2026). Gas City: SDK for building custom packs; v1.0 April 24, 2026; reportedly run with hundreds of concurrent workers per city (per Julian Knutsen, primary Gas City implementor).

The orchestration is autonomous, but the infrastructure is *built around the judgment gap* rather than waiting for it to close. The architecture is a *map* of where shape (1) fails:

- **Watchdog tier** (Witness → Deacon → Dogs → Boot) catches stuck agents
- **Escalation** routes hard cases to human via severity-graded beads
- **Refinery** (merge queue with verification gates, Bors-style bisecting) handles "did the work actually compose"
- **Persistent identity via Beads** survives session failure; agents are not sessions, they are durable identities with ephemeral session-executors

Gas City's explicit multi-agent rule:

> "You should almost never deploy a single-agent pack for a real business process... You should always have at least two or three working together on a little crew. ... With Gas City you can build any sort of adversarial group structure you like, for a team of collaborating agents. They can watch over each other."
> — Yegge, *Welcome to Gas City*, § The Shape of Things to Come (April 2026)

Adversarial verification per pack as a deployment rule, not a guideline. The framework's verification-primitives discipline (`THE_CASE.md` §81) has a deployed analog in working code.

*Independent convergent variant:* Kieran Klaassen's Compound Engineering guide (Every, https://every.to/guides/compound-engineering, Jan 2026 / updated May 2026) names the same multi-agent-review pattern with explicit severity tiering — *"Multiple specialized reviewers examine the code in parallel"* with findings marked P1 (must fix) / P2 (should fix) / P3 (nice to fix). Two non-overlapping orgs (Yegge / Gas City and Every / Kieran) arriving at the multi-reviewer pattern from different operational vantages, with Every's P1/P2/P3 vocabulary sharper than Gas City's "two-or-three agents" framing for severity routing specifically.

**Framework-aligned moves, in practitioners' vocabulary:**

- **Refinery's batched-merge-with-verification-gates** ≈ §81's "deterministic verifier that re-derives the claim from raw ground truth," moved to the architecturally-critical boundary (merge to main)
- **Two-or-three-agents-per-pack** ≈ verification primitives embodied as a deployment rule
- **MEOW (Molecular Expression of Work) / Formulas → Molecules** ≈ substrate-as-vocabulary at organizational scale (Formula = template; Molecule = instance; version-controlled in Dolt, forkable). Maps onto the coach-with-typed-intervention-library example in `EXAMPLES.md`, scaled from one professional's repertoire to one organization's process inventory.
- **Watchdog tier** ≈ calibration discipline expressed as infrastructure rather than convention
- **"Reliability, friends, is a dial. You choose where to set it. More rounds of review, more backstops, more guardrails, more judges..."** (*Welcome to Gas City*) = the framework's calibration claim in different vocabulary, with concrete knobs

**Yegge's invariance claim worth noting:**

> "Several new model generations have dropped, and Gas Town hasn't changed shape at all. The architecture has shown remarkable resilience."
> — Yegge, *Welcome to the Wasteland*, § Is Gas Town Ready? (March 2026)

~3 months of practitioner self-report; if it holds, it's strong support for the framework's load-bearing position that *discipline is orthogonal to model capability*. Treat as suggestive, not ablation-validated.

**Distinctive claim worth flagging:**

- **"Hallucinations and false memories and forgetting are baked mathematically into all memory systems; there's no avoiding it."** (*Welcome to Gas City*, § The Shape of Things to Come). Strong epistemological prior. The framework should stay agnostic on the strong form — not every memory failure is equally inevitable — and lean on the weak form: build verifiers for what you can catch.

**Tempering reading — Maggie Appleton, January 2026.** Yegge's claims about Gas Town are largely self-evangelism; the most useful adversarial vantage comes from Maggie Appleton's essay *Gas Town's Agent Patterns, Design Bottlenecks, and Vibecoding at Scale* (https://maggieappleton.com/gastown, Jan 2026):

> "We should take Yegge's creation seriously not because it's a serious, working tool for today's developers (it isn't). But because it's a good piece of speculative design fiction."

Appleton's framing — Gas Town as **speculative design fiction** revealing future architecture, rather than a present-day shippable tool — is a useful counterweight to Yegge's "Gas Town is ready" marketing. For the framework's purposes, this means: treat shape (3) as the most-elaborated articulation of what engineered-resilience autonomy *could* look like, not as ablation-validated infrastructure that's currently easy to adopt. The four shapes are still the right map; Gas Town remains the best exemplar; but the "deployed at hundreds of concurrent workers" claim is practitioner self-report from inside the evangelist position.

Appleton also names the actual failure mode in vibecoding-at-scale that the framework would otherwise reach for §81 to address: *"you don't fully consider what you're building at each step."* This is the per-output verification gap, named from a design-anthropology vantage rather than an engineering one.

---

## Shape 4: Socially-validated federation

**Exemplar:** Wasteland (https://steve-yegge.medium.com/, *Welcome to the Wasteland: A Thousand Gas Towns*, March 2026, built by Julian Knutsen + Matt Beane + others with Yegge's vision). Federated network of Gas Towns sharing a Wanted Board and trust-graded validation, atop Dolt (Git-versioned SQL database).

Stamps as multi-dimensional verdicts (quality, reliability, creativity, with confidence and severity scores), anchored to specific evidence, accumulated in an append-only Git-backed ledger. Trust levels gate what you can do; new rigs only browse/claim/submit; validators are rigs that earned standing through stamped work; the "yearbook rule" forbids self-stamping.

**Concretely: three of `PROCEDURE.md`'s five protocols deployed at federation scale in working code, with no apparent prior contact with `PROCEDURE.md`'s framing.**

- **Trust levels** = standing (`PROCEDURE.md` §2). New rigs have limited standing; standing accumulates via validated work.
- **Yearbook rule** ("you can't stamp your own work") = recusal (`PROCEDURE.md` §5).
- **Stamp ledger** = substrate-as-record with full provenance (`PROCEDURE.md` §3): *"the graph is fully traversable. And because the underlying storage is append-only and versioned, the history can't be rewritten — your ledger is permanent."*

Bounded action (`PROCEDURE.md` §4) doesn't have a clean Wasteland analog — that protocol pertains to what an autonomous loop can enact, while the Wasteland's structure is about who may attest.

The Wasteland's anti-collusion frame is metabolism applied to validation patterns rather than to substrate drift:

> "The stamp graph has a shape, and collusion rings have a distinctive topology — lots of mutual stamping, sharp boundaries, no outside critics. The Wasteland system is designed to make fraud unprofitable, not impossible."
> — Yegge, *Welcome to the Wasteland*, § Why is the Wasteland any different (March 2026)

Worth borrowing into `PROCEDURE.md`'s metabolism section as a second flavor: drift in the *validation graph topology* is a distinct audit target from drift in substrate contents.

---

## What the four shapes share

Internal critique worth foregrounding now that the four shapes have been laid out: shapes 2/3/4 can be read as one pattern at three deployment scales. The shared pattern — **separate the actor from the attester, then log the attestation** — varies in who the attester is:

- Shape 2: human-as-attester
- Shape 3: LLM-pack-as-attester
- Shape 4: peer-network-as-attester

Shape 1 is the null where the attester is missing or the actor itself — the "judgment in choosing goals" gap. The four-shape framing in this doc is for the operational layer; the actor-separated-from-attester pattern is the unifying structural claim. Both are doing work — a human-in-the-loop, a polecat fleet, and a federated trust graph are different artifacts to deploy even when their abstract structure rhymes.

---

## What's genuinely new for the framework

Six patterns the framework currently lacks vocabulary for, tiered by evidentiary weight:

- **Strong** (multiple independent practitioners, convergent moves): §1 Agent UX (4 sources across 3 related sub-disciplines; §1a alone has 2), §4 markdown-as-org-substrate (3 sources), §5 compute-as-binding-constraint (3 sources of different kinds).
- **Medium** (single well-grounded source with clean operational claim): §2 autonomy slider (Karpathy), §3 generation-verification loop velocity (Karpathy), §7 organizational-level feedback loops (Charity).

Two adjacent framings (§6 Yegge's Survival Ratio, §8 Karpathy's Iron Man suit metaphor) are noted below as useful vocabulary the framework can borrow, not as newly surfaced disciplines.

Each item below is documented practice with a stated method, not theory.

### 1. Agent UX / designing-for-the-LLM-end-user

The strongest convergence finding in the lit review — but worth disaggregating, because the five practitioner citations bundled here are addressing three distinct sub-disciplines, not one. They share an orientation (LLMs are the end-user; design accordingly) but the moves are different:

**1a. Substrate-side friction minimization** (the agent's interface to your data/state):
- **Yegge (*Software Survival 3.0*, Jan 2026):** *"What I did was make their hallucinations real, over and over, by implementing whatever I saw the agents trying to do with Beads, until nearly every guess by an agent is now correct. I've driven the friction cost term about as low as it can go."* He calls this *Desire Paths* design.
- **Kieran Klaassen (*The Folder Is the Agent*, Every, April 13 2026; https://every.to/source-code/the-folder-is-the-agent):** the substrate IS what individuates one agent from another. *"A project folder with a CLAUDE.md/AGENT.md...that's an agent."* / *"Just by pointing the model at this folder, which contains some of my personality, knowledge, and taste, the model can be a specialist."* / *"Just by changing the folder and not the model, I have a different agent."* The throughput claim that makes substrate-side friction load-bearing: *"I'm running 44 AI agents across multiple projects. Each one is just a model pointed at a folder."* When the folder is the agent, the design move is configuring the folder, not configuring the model — same orientation as Yegge's Desire Paths, stated at the individual-agent altitude rather than the substrate-build altitude.

**1b. Tool/protocol-side friction minimization** (the agent's interface to the broader world):
- **Karpathy (YC, July 2025):** explicit "build for agents" — `lm.txt` files alongside `robots.txt`, markdown docs over HTML, replacing "click X" with "curl this endpoint", MCP as the canonical protocol.

**1c. Documentation as agent-mediated knowledge routing** (the human's interface to your tool, *through* the agent):
- **Karpathy (No Priors, Q1 2026):** *"It used to be that you have documentation for other people... but like you shouldn't do that anymore... if agents get it then they can just explain all the different parts of it. So it's this redirection through agents."*

And separately, the discipline that makes shape (2) sustainable:
- **Charity Majors (Jun 2026):** *"engineer the prerequisites (better evals, tests, feature flags, observability)"* so the verification is fast enough that velocity becomes sustainable. This is verification-velocity engineering, sister to 1a/b/c rather than identical.

The framework currently treats the LLM block as a given and doesn't address any of 1a/1b/1c structurally. The unified discipline-sibling-to-calibration would be **Agent UX**, with the three sub-disciplines as the operational moves. Treat the convergence as: *three practitioners independently arrived at the orientation; the specific moves differ.* That's still a real convergence, just more honestly framed than "four endorsements of the same thing."

### 2. Autonomy slider as product-design dimension

Karpathy's framing names the design choice each shape implicitly makes:

> "We can build augmentations or we can build agents and we kind of want to do a bit of both, but at this stage I would say working with fallible LLMs, it's less Iron Man robots and more Iron Man suits that you want to build... less like building flashy demos of autonomous agents and more building partial autonomy products. And these products have custom GUIs and UI/UX, and this is done so that the generation-verification loop of the human is very very fast."
> — Karpathy, YC, July 2025

The four shapes sit at different slider positions; the framework should make the slider explicit.

### 3. Generation-verification loop velocity

§81 of `THE_CASE.md` currently emphasizes verification *correctness*. Karpathy adds verification *speed* as equally load-bearing:

> "It is in our interest to make this loop go as fast as possible so we're getting a lot of work done. ... GUIs are extremely important... a GUI utilizes your computer-vision GPU in all of our heads. Reading text is effortful and not fun, but looking at stuff is fun."
> — Karpathy, YC, July 2025

The speed of the loop is part of the design target, not just an optimization. This is the disciplined reframing of "vibe coding" — the loop is fine if it's fast and verifiable; it's broken if it's fast and unverified.

### 4. Markdown-as-organizational-substrate

Three independent endorsements of *the markdown describing the loop is the substrate*:

- **Karpathy (No Priors, Q1 2026):** *"A research organization is a set of markdown files that describe all the roles and how the whole thing connects."* His `program.md` is the spec for his auto research loop; meta-optimization over `program.md` is recursive self-improvement at the org-spec level.
- **Yegge (Gas City, Apr 2026):** MEOW with Formulas (reusable templates) and Molecules (instances), version-controlled in Dolt, forkable across an org. *"Your library of formulas becomes a declarative inventory of every business process you've ever automated."*
- **Anthropic Skills primitive** (referenced indirectly by Karpathy in No Priors): *"a skill is just a way to instruct the agent how to teach the thing."*

This generalizes substrate-as-vocabulary (`EXAMPLES.md` coach template) from one professional's repertoire to a whole organization's process inventory, and answers the substrate question shape 2 left open: when code is regenerable, the durable artifact is *the org-spec markdown*, not the code.

### 5. Compute-as-binding-constraint

Three independent sources converging within ~30 days on compute as the new binding constraint:

- **Karpathy (No Priors, Q1 2026):** *"You're not maximizing your subscription at least... it's not about flops, it's about tokens. So what is your token throughput and what token throughput do you command?"* Plus the explicit reframing: *"Is flop the thing that actually everyone cares about in the future?"*
- **Yegge (Survival 3.0, Jan 2026):** the entire Survival Ratio is denominated in tokens-as-cognitive-cost.
- **Uber (Jun 2026, via Simon Willison's coverage):** $1,500/month per-tool spending caps on agentic coding software after budget overruns.

The Survival Ratio's Awareness + Friction denominator gets a hard ceiling at enterprise scale. Yegge's "reliability is a dial" needs the addendum: the dial doesn't go past where the budget allows.

### 6. Survival Ratio as complementary framing

Yegge's tool-fitness formula (*Software Survival 3.0*):

```
Survival(T) ∝ (Savings × Usage × H) / (Awareness_cost + Friction_cost)
```

Six levers — Insight Compression, Substrate Efficiency, Broad Utility, Publicity, Friction Minimization, Human Coefficient. Different organizing principle from hybrid-loops' discipline-pattern decomposition; reaches overlapping conclusions about what makes a good tool block. Worth a See-Also from `BUILDING_BLOCKS.md` — non-self-referential evidence that the *discipline of tools for agents* has economic stakes.

Two levers map directly:
- **Lever 1 (Insight Compression)** ≈ substrate as crystallized knowledge (*"Git represents decades of accumulated wisdom"* / *"crystallized cognition"*)
- **Lever 2 (Substrate Efficiency)** ≈ the gate concept (*"grep saves cognition by doing it on a cheaper substrate: CPUs"*)

### 7. Organizational-level feedback loops

Charity Majors (https://charity.wtf/, *AI enthusiasts are in a race against time, AI skeptics are in a race against entropy*, June 2 2026): wins get celebrated in public forums (talks, blogs, all-hands); downstream costs surface in private forums (retros, oncall, DMs). The *same people don't see both sides*, so each group reinforces its worldview independently. Her prescription:

> "There is no natural feedback loop connecting enthusiasts with skeptics... Designing feedback loops to help mend the gap in shared reality between the two groups is a fascinating organizational design problem."

This is calibration discipline at the *organizational* level — distinct from the per-loop calibration the framework currently addresses. Add as a separate dimension of metabolism.

### 8. Iron Man suit as the product-shape

Worth folding a single Karpathy quote into `THE_CASE.md` as the closest restatement of the framework's central thesis in product-design vocabulary:

> "It's less Iron Man robots and more Iron Man suits that you want to build... less like building flashy demos of autonomous agents and more building partial autonomy products. And these products have custom GUIs and UI/UX, and this is done so that the generation-verification loop of the human is very very fast."

---

## What the framework would NOT lift

Five postures the framework should keep distinct from:

1. **"Believe the curves"** (Yegge's stated method). The framework's stance is ablation discipline — don't add features for hypothetical futures. Acknowledge Yegge's predictive track record (Death of the Junior Developer 2024, Revenge of the Junior Developer 2025, Gas Town existing as he predicted) while keeping the framework's epistemic posture.

2. **Vibe-coding at the individual-output level** (Charlie: *"I've never seen the code"* / Yegge: *"some work gets lost"*). The framework's `THE_CASE.md` §81 calls for per-output deterministic verification. The reconciliation that actually works in production is shape (3)'s answer: per-output verification happens at the *pack* level (2–3 agents reviewing each other), not at the individual-LLM-call level. The framework can adopt this reconciliation without conceding the strong form. Charity's bridge question — *"what would it take for you to feel comfortable shipping code without reading it?"* — names the engineering work that makes shape 2 sustainable when it is.

3. **"The end of teaching each other things"** (Karpathy, strong form). Cite the framing; don't adopt it. The framework keeps the weaker form — documentation should be readable to both agents AND humans; agents are increasingly the router, but humans still need access. The strong form risks foreclosing the auditability the framework cares about.

4. **The outside-the-lab argument** (Karpathy, No Priors Q1 2026: *"I feel a bit more aligned with humanity outside of a frontier lab"*). He joined Anthropic by May 2026. The argument may still be sound, but the cite is no longer durable; Karpathy is a witness to the *pressures*, not a stable independent voice on the inside-vs-outside question.

5. **Speculative distributed-auto-research / "flop democratization"** (Karpathy, No Priors). Interesting; hold as a one-line note, not as framework input. The technical preconditions (untrusted-pool verification with low-cost checking) are not yet routine.

---

## Picking a shape

| If you're... | Reach for... |
|---|---|
| One human doing parallel exploratory coding work | Shape 2 (Conductor or its open analogs, Karpathy-style throughput as the skill ceiling) |
| Replacing a recurring business process with an autonomous crew | Shape 3 (Gas City packs, two-or-three agents minimum) |
| Coordinating work across organizations or contributors | Shape 4 (Wasteland-style federated validation) |
| Tempted to spawn one LLM to "decide what to do" and fan out | Pause. You're in shape 1; the article's gap will likely bite unless you have Karpathy-level operator skill *or* shape-3 infrastructure. Default to dropping back to shape 2 until you have one of those. |

*Caveat on the table's altitude.* These rows are about what the *problem* looks like — what you're trying to accomplish, who's involved. They are not about what the system would cost, how mature your CI needs to be to ship it, or what stakes / regulatory posture you're in. Those would be the next layer of analysis if this doc grew into a real adoption guide. One stakes-driven exception is worth flagging directly: regulatory or high-stakes contexts (finance, healthcare, anything with mandatory audit trails) may need shape 3's per-pack adversarial verification even when scope alone would say shape 2 — because Charlie's "I've never seen the code" stance becomes legally untenable before it becomes operationally untenable.

## Systems this taxonomy reads onto but doesn't categorize

The four shapes were named because each has a distinct calibration apparatus worth talking about. Several other production multi-agent systems in 2026 sit somewhere the taxonomy reads onto rather than picks out a new shape:

- **Vertical-domain agent platforms** (Sierra is one verified example, sierra.ai — a multi-vertical customer-experience agent platform; the broader category includes other domain-focused agent products in customer service, legal, healthcare, and research). Reads as **constrained-scope shape 3** — adversarial verification happens between the LLM and an escalation policy rather than pack-vs-pack, and tight domain scope substitutes for engineered-resilience-infrastructure depth.
- **Compiled pipelines (DSPy, LangGraph, AutoGen).** Typed handoffs in a human-authored graph; the graph IS the orchestrator's decomposition, executed at runtime. Reads as **shape 2 with the graph as the human's design-time work**. Covered at the implementation-toolkit altitude in `AGENT_FRAMEWORKS.md`.
- **Autonomous coding agents** (Devin from Cognition — cognition.ai — is the named example; the category also includes the various IDE-embedded background-agent products). Loops run without per-step human supervision but inside engineered constraints. Reads as **deployment-constrained shape 1** — the architectural specifics of how each product constrains the agent are mostly not publicly documented, so the mapping is suggestive.
- **Handoff / routing patterns** (e.g., OpenAI Swarm). Stateless handoff between LLMs by role. Reads as **shape 1's task-routing variant** — the article's "judgment in choosing goals" gap bites at the routing decision unless the routing itself is constrained.

These systems are real and load-bearing for many readers' actual decisions; not including them as separate shapes is a *scope choice*, not an assertion that they don't exist. If your problem fits one of them better than it fits the four shapes above, the shape-vocabulary in this doc still applies — it just reads onto a system someone else designed.

---

## Migration paths between shapes

The shapes aren't a strict hierarchy — most production setups end up at one and stay there. But two transitions are common enough to name:

**Shape 2 → Shape 3 (vibe-coding hits a verification ceiling).** Signal: the human reviewer is the bottleneck on throughput, and the team is making "withdrawals from a trust account" (Simon) faster than they can repay them with review cycles. The migration is to per-pack adversarial verification (Gas City's "two-or-three agents" rule), which moves verification off the human's critical path. Charity's bridge question — *"what would it take for you to feel comfortable shipping code without reading it?"* — names the engineering work that has to ship before the migration is real.

**Shape 3 → Shape 4 (single-org verification hits a trust ceiling).** Signal: the substrate of attested outcomes outgrows what any single organization can validate, or the work crosses organizational boundaries (open-source contributions, cross-team coordination, contractor work). The migration is to federation with trust-graded stamps. The Wasteland is one design; other federated-validation systems (Git's PR model is the classical example) are others.

What we don't have a worked example of: **Shape 1 ← anywhere**. Genuinely autonomous decomposition without one of the scaffolding strategies is not a stable destination as of 2026.

## What would falsify the doc's claims

A live question the slimemold-style critique should ask of any synthesis: *what observation would change our minds?* Three concrete falsifiers for the load-bearing claims in this doc:

1. **The "discipline-orthogonal-to-capability" claim (Yegge's invariance assertion) would be falsified** by a clear model-capability advance that makes shape-3 infrastructure unnecessary — i.e., an LLM that decomposes complex tasks reliably enough that watchdogs, refineries, and per-pack verification add no value. Anthropic's RSI article says we're not there as of May 2026; if a future model release demonstrably closes the "judgment in choosing goals" gap, the framework's discipline-is-everything posture needs to move.

2. **The four-shape taxonomy would be falsified** by a working production system that doesn't fit any of the four (e.g., a fully autonomous swarm with no per-pack adversarial verification AND no human orchestrator AND no federated validation, succeeding at non-trivial work for sustained periods). Today's failure cases (undercity-style) are consistent with the taxonomy; a success case outside it would force a fifth shape or a re-carving.

3. **The Wasteland-protocols mapping would be falsified** if attestation in the Wasteland turns out to function differently from how the essay describes it — e.g., if stamps are routinely awarded without verified evidence, or if the trust-level gating doesn't actually constrain what new rigs can do. The mapping currently rests on Yegge's essay; it has not been validated by independent inspection of the running system. A field-report contradicting the essay would weaken the mapping.

A useful epistemic practice for anyone applying this doc: when one of its claims becomes load-bearing for a decision you're about to make, ask the corresponding "what would falsify this" question, then check whether the falsifier holds for your situation.

---

## Citations and source tiering

**Source independence — honest accounting.** The lit review draws on five vantages, but two of them now ladder up to the same organization:

| Source | Role | Independence |
|---|---|---|
| Anthropic, *When AI builds itself* (May 2026) | Substrate provider's own framing | Same org as #2 |
| Karpathy (YC talk Jul 2025; No Priors Q1 2026) | Practitioner / theorist, joined Anthropic May 2026 | Same org as #1 |
| Ryan Lopopolo / OpenAI, *Harness Engineering* (2026; https://openai.com/index/harness-engineering/) | Second substrate provider's framing of agent-native development | Distinct from Anthropic; partially offsets the substrate-provider consolidation noted below |
| Yegge (× 4 essays: Survival 3.0, Gas Town, Gas City, Wasteland; Jan–Apr 2026) | Independent practitioner — also exemplar of shapes 3 and 4 and source of three of eight "what's new" findings | Distinct, but heavily weighted (~45% of evidentiary base by line count) |
| Charlie / Conductor (YT walkthrough) | Independent practitioner — shape 2 exemplar | Distinct |
| Charity Majors (Jun 2026) | Independent observer | Distinct |
| Maggie Appleton (Gas Town essay, Jan 2026) | Independent design-anthropology critic of Yegge's claims | Distinct |
| Simon Willison (synthesis posts, May–Jun 2026) | Synthesis observer; second-hand cites for Uber caps, trust-account framing | Distinct (but synthesist) |
| Anthropic — Shihipar & Bidasaria, *A harness for every task: dynamic workflows in Claude Code* (blog, Jun 2 2026) | Substrate provider's own product announcement, cited above under Shape 1's "third path" update | Same org as #1 |

**Two consolidations the doc carries, not one.**

The Anthropic-Karpathy ladder is named above. A second consolidation is the **YC / blog-circle cluster** that supplies most of the "what's new" findings: Charlie (YC S24), Karpathy (YC AI Startup School speaker, then No Priors podcast), Yegge (regular Discord-and-blog presence in the same ecosystem), Charity Majors (frequent Yegge-adjacent voice in the ops-blog circuit), Simon Willison (active citer of all of them). These sources don't share an employer, but they share a discourse — they cite each other, attend each other's events, podcast on each other's shows. The convergence findings (Agent UX, compute-as-binding-constraint, markdown-as-org-substrate) are real *within this cluster*; whether they generalize to multi-agent practice outside of US developer-tools Twitter is an open question this doc does not answer.

**Yegge weighting is materially heavier than the row count suggests.** He contributes four essays, two of the four shapes' canonical exemplars (Gas Town / Gas City for shape 3; Wasteland for shape 4), and three of the eight "what's new" findings. Where this doc treats Yegge-attested facts (e.g., "hundreds of concurrent workers per Gas City," "Gas Town hasn't changed shape across model generations," the Wasteland's stamp-protocol behavior) as evidence, those are practitioner self-reports from the evangelist position. Maggie Appleton's *speculative-design-fiction* framing tempers shape 3; nothing in the current draft equivalently tempers the shape-4 protocol mapping, which still rests on Yegge's essay without independent field verification.

That's effectively **four independent sources + the substrate provider's framing + one synthesist**, but the four are not equally weighted: Yegge is closer to 3-4 sources, and one peer-cluster connects most of the rest.

**Asymmetric Karpathy note.** Karpathy's pre-Anthropic quotes (YC July 2025, No Priors Q1 2026) were made independently and remain independent evidence; joining Anthropic in May 2026 doesn't retroactively make them substrate-provider statements. The asymmetry matters going forward: *future* Karpathy quotes should be tiered with Anthropic; *past* Karpathy quotes stay where they were said. The doc treats his existing quotes as independent and would not lift any future post-May-2026 Karpathy statement at the same evidentiary weight without flagging the consolidation. The convergence findings (Agent UX, markdown-as-org-substrate, compute-as-binding-constraint) still triangulate across genuinely independent vantages. Power-balance and outside-the-lab arguments derived from Karpathy alone are weaker than they read in isolation.

**Verification status per source:**
- Anthropic RSI article: WebFetched 2026-06-04; quotes and section attributions confirmed verbatim.
- Karpathy YC talk: confirmed via user paste from https://www.youtube.com/watch?v=LCEmiRjPEtQ; transcript not independently fetched.
- Karpathy No Priors interview: confirmed via user paste from https://www.youtube.com/watch?v=kwSVtQ7dziU; transcript not independently fetched. Date is post-December-2025 capability flip, pre-May-2026 Anthropic move.
- Yegge essays: read in full from local Downloads archive copies of HTML originals. Live URLs at https://steve-yegge.medium.com/.
- Conductor walkthrough: transcript pasted by user; not independently verified against video.
- Charity Majors: WebFetched 2026-06-04 from https://charity.wtf/; key quotes confirmed.
- Anthropic dynamic-workflows blog: direct WebFetch returned HTTP 403 on 2026-07-01; full article text confirmed same day via user-supplied copy.
- Simon Willison: WebFetched 2026-06-04 from https://simonwillison.net/; secondary cites (Uber, Anthropic sandboxing) noted but not fetched at primary source. Anthropic sandboxing post URL is unverified — 404 on first guess.
- Maggie Appleton: WebFetched 2026-06-04 from https://maggieappleton.com/gastown via summary report. Two short phrases ("design fiction", "design becomes the limiting factor") were explicitly identified as direct quotes; the longer "We should take Yegge's creation seriously..." blockquote and the "you don't fully consider what you're building" phrase were rendered as quotes by the summary but have not been independently verified against the essay's original wording. Date approximated to Jan 2026 from in-text social-media references in the essay (not from a header date), per the WebFetch summary; the main narrative says "January 2026" with that caveat.
- *Karpathy moved to Anthropic as of May 2026 per user.*

The Yegge invariance claim about Gas Town's shape persisting across model generations is ~3 months of practitioner self-report; treat as suggestive evidence for discipline-orthogonal-to-capability, not as ablation-validated.

---

## What we did not sample

The lit review behind this doc is dominated by US developer-tools voices in the YC / blog-circle cluster. Several perspective-classes are not represented and would likely change parts of the synthesis if they were:

- **Academic agent researchers.** Major-university agent labs and the recent ICML / NeurIPS literature on agent verification, formal evaluation, and capability elicitation. The framework's empirical claims about discipline-vs-capability would benefit from their evidence.
- **Safety / evals researchers.** Independent nonprofits doing capability and autonomy evaluation — METR (metr.org, focused on autonomous-capability and AI-R&D-acceleration evaluation; see `PRIOR_ART.md` *Safety / evals research* for the framework's hedged treatment of their time-horizon headline number) and Apollo Research (apolloresearch.ai, focused on scheming detection in frontier models) are two verified examples; there are others. Strikingly absent for a doc about autonomous orchestration; the "judgment-in-choosing-goals" framing comes from Anthropic's own self-assessment rather than from external safety evaluation.
- **Non-US voices.** Non-US frontier-model labs and their agent work, European AI-Act / regulatory framing, Asian enterprise deployment reports. The synthesis is currently US-flavored; specific labs and reports were not surveyed for this draft.
- **AI-skeptic / labor voices.** Researchers and writers analyzing the labor and accountability dimensions of full-time agent supervision. Their concerns are at a different altitude (sociotechnical critique rather than orchestration pattern), but their absence means the doc reads as more *settled* than the actual discourse is.
- **Non-developer-tools enterprise practice.** Finance, healthcare, legal: production-grade autonomous-ish orchestration with very different calibration apparatus (compliance audits, regulatory attestation, mandatory human review at specific decision points). The doc's "production scale" is coding agents specifically.
- **Ops / SRE voices beyond Charity Majors.** The broader resilience-engineering literature and ops-blog circuit. Charity's "feedback loop" framing is one good vantage; the discipline she points at has decades of unrelated infrastructure.

If you adopt one of the shapes here based on the framing in this doc, the absent perspectives are where the framing is most likely to mislead you. A future revision should at least source-spot-check the safety-evals literature and one non-US enterprise practitioner; everything else is "research notes worth knowing."

---

## See also

- `THE_CASE.md` §81 (calibration's two jobs; deterministic verifier) — verification-primitives discipline the shapes elaborate. Worth folding the Karpathy *"Iron Man suits"* quote here.
- `PROCEDURE.md` — the five procedural protocols that shape 4 deploys in working code without naming them as such. Worth adding metabolism-as-validation-graph-topology (Wasteland's anti-collusion frame) and Charity's organizational-feedback-loop framing.
- `EXAMPLES.md` — single-loop worked examples; this doc covers multi-agent compositions of those loops.
- `AGENT_FRAMEWORKS.md` — comparison to LangGraph / AutoGen / CrewAI / DSPy (different level of abstraction: implementation toolkits, not orchestration shapes).
- `PRIOR_ART.md` — Multi-agent orchestration projects section needs expansion: currently only one shallow Gas Town bullet; Wasteland and Gas City missing; Charlie's Conductor not yet distinguished from Netflix Conductor.
- `BUILDING_BLOCKS.md` — Worth a See-Also pointer to Yegge's Survival Ratio as a complementary framing for tool-block fitness.
