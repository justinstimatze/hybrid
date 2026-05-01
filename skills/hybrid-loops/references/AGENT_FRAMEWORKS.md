# Hybrid loops vs. agent frameworks

A natural critique: *"This sounds like LangGraph / AutoGen / DSPy / pydantic with an extra coat of theory. What's actually different?"* The answer needs care because there's real overlap with each tool, and the positioning depends on getting the *levels of abstraction* right.

## The actual gap these tools leave open

DSPy has *modules* and *signatures* and *optimizers*. LangGraph has *nodes* and *edges* and *checkpoints*. AutoGen has *agents* and *conversations* and *managers*. CrewAI has *crews* and *roles* and *tasks*. pydantic has *models* and *validators*. instructor has *response_model*. Anthropic and OpenAI have *tool use* and *structured outputs*.

Each tool covers one cell of a broader pattern with its own primitives, and none of them point at the others. A senior engineer who's used three of them carries three half-overlapping mental models and reconciles them by intuition. A mid-level engineer who picks any one of them locks into that cell without seeing the rest of the alphabet — the tools don't surface what they don't cover.

The framework is filling that gap. Not "build another competitor to LangGraph" — provide the unifying vocabulary the ecosystem lacks, so a senior with three frameworks under their belt has one coherent mental model and a mid-level engineer has a starting model that points at the broader space.

## TL;DR — different levels

Hybrid loops is a *design pattern* one level of abstraction up from agent frameworks. The frameworks below are *implementation toolkits* for pieces of the pattern.

- The pattern says: design your system as a graph of typed blocks where LLM-actors and code-actors alternate, with three disciplines (calibration, context-as-code, dev-time loop) applied to the LLM blocks.
- The frameworks say: here's a runtime / programming model / structured-output validator for the LLM portion of your graph.

Critique that flattens to "this is just LangGraph" misunderstands the levels. Critique that flattens to "this is just architecture line-noise; the frameworks are real" misses that the disciplines are the framework's content. The comparisons below try to forestall both.

## DSPy (Khattab et al., Stanford)

**What it is**: a programming model for LM programs. Typed *signatures* (input/output fields) define modules; *optimizers* (BootstrapFewShot, MIPROv2, etc.) search over prompts and few-shot demonstrations to maximize a metric. The "compile" step turns a program + metric + examples into an optimized program.

The closest cousin in this work. DSPy's typed-signature module is essentially a typed-block-with-an-LLM-actor; its compile-for-metric loop is structurally the same as the compress-and-verify shape this framework names, just optimizing prompts/demos rather than notation. If you've internalized DSPy you've internalized a chunk of the hybrid-loops algebra already.

**Where they differ**:
- DSPy's alphabet is *mostly LLM modules*. Deterministic non-LLM blocks are second-class — Python interop, not first-class architectural elements. Hybrid loops puts LLM and code on equal footing.
- DSPy optimizes for a user-provided metric. Hybrid loops names *disciplines* — calibration treats per-block hit-rate as ship-blocking, not as one optimization signal among many.
- DSPy doesn't distinguish runtime cycles from development-time cycles. The compile step *is* a dev-time loop, but the framework treats it as tooling, not architecture.
- DSPy doesn't have an explicit substrate concept — programs are stateless w.r.t. each other. Hybrid loops treats typed-records-accumulating-over-time as a first-class role.

**Verdict**: complementary. Use DSPy to implement the LLM blocks of a hybrid-loops graph when prompt/demo optimization is worth the compile-step cost.

## LangGraph (LangChain)

**What it is**: a graph executor for LLM-centric agent workflows. Nodes are typically LLM calls or tool calls; edges express the control flow; supports cycles (hence "graph" not "chain").

**Where they differ**:
- LangGraph nodes are typically LLM calls or tools — the deterministic-substrate / deterministic-gate / typed-substrate-store roles are implemented ad hoc. Hybrid loops names those roles structurally.
- LangGraph's graph IS code (Python). A hybrid-loops graph would be data, in principle — typed specs that code could validate, lint, simulate. That spec and executor don't ship in this repo; the framework points the direction. LangGraph chose the practical alternative of letting Python be the graph.
- LangGraph doesn't have an opinion on calibration, context-as-code as infra, or dev-time-loop discipline. It's a runtime, not a methodology.

**Verdict**: complementary at runtime. A hybrid-loops graph can be executed by LangGraph or by direct Python or by an MCP-tool sequence — the framework doesn't prescribe a runtime. LangGraph fits projects whose graph is mostly LLM blocks; for graphs with heavy deterministic infrastructure, a thinner runtime works better.

## AutoGen (Microsoft)

**What it is**: a framework for multi-agent conversations. Agents have roles; messaging works as conversation; group chats coordinate via a "manager" agent.

**Where they differ**:
- AutoGen's primary primitive is *the conversation* — agents exchange messages, context accumulates as transcript. Hybrid loops's primary primitive is *the typed substrate* — records accumulate as structure, not as conversation.
- AutoGen treats multi-agent as the default shape. Hybrid loops treats LLM blocks as one block type among others; whether multiple LLMs converse is a question of substrate-shape, not framework default.
- AutoGen doesn't separate runtime from dev-time. Multi-agent conversations are runtime constructs; the dev-time critique-and-patch loop is something users implement themselves.

**Verdict**: AutoGen's shape is one specific cell — *multiple LLM blocks with conversational substrate*. It's powerful for that shape. Hybrid loops is broader; AutoGen is deeper inside it. Use AutoGen when the right pattern is multi-agent conversation; use hybrid loops as the umbrella when it isn't.

## CrewAI

**What it is**: opinionated role-based multi-agent — "researcher" agents, "writer" agents, "critic" agents working together with role descriptions and goals.

Same architectural cell as AutoGen, more opinionated about *how* the multi-agent shape is constructed. Use when the role-described-agents pattern fits.

## pydantic / instructor / structured outputs

**What they are**: structured-output validators / typed-schema enforcers for LLM calls. pydantic provides Python type-checking; instructor wraps LLM calls in pydantic models; Anthropic's tool use and OpenAI's structured outputs let you define a JSON Schema the LLM is constrained to output.

These tools are how you implement *one* of the framework's disciplines — constrain LLM output to a typed schema. They don't tell you which schema to use. The framework's compress-and-verify shape (see `BLOCK_GRAPHS.md`) is one approach to schema discovery; pydantic / instructor / structured-outputs are how you'd enforce the schema at runtime.

**Verdict**: different layers entirely. Use pydantic + instructor + structured-outputs to implement context-as-code constraints in your LLM blocks; the framework names *that this is a discipline* and explains why every LLM block in a non-trivial graph needs one.

## Other adjacent ecosystems

The agent-framework comparison above is the most-load-bearing because that's where engineers most-often flatten "this is just X." Three more ecosystems sit nearby with different emphases.

### Workflow orchestration: Temporal, Conductor, AWS Step Functions, Airflow

Durable workflow engines for long-running, retry-tolerant, replayable processes. Temporal treats workflows as code with deterministic replay; Conductor (Netflix) is similar; Step Functions and Airflow are AWS/Apache versions tuned for batch.

A hybrid-loop graph that needs to survive process restarts, run for hours, retry on failure, or persist intermediate state across blocks is exactly what these engines exist for. They don't say anything about LLM-specific disciplines (calibration, context-as-code, dev-time loops) — they're runtime infrastructure, not architectural methodology. A hybrid-loop runtime can be implemented on Temporal; the framework's disciplines apply orthogonally. Reach for them when the graph needs durability, replay, retries.

### Visual LLM-app builders: Dify, LangFlow, Flowise

Visual canvases for assembling LLM apps from pre-built nodes. Dify is the most LLM-native; LangFlow and Flowise sit on LangChain.

These make hybrid-loop-shaped graphs *visible* — the canvas IS the graph-as-data. Non-engineers can drag-and-drop a structurally-hybrid-loop graph without writing code. The alphabet is constrained to whatever the canvas exposes; deterministic non-LLM blocks are limited to the platform's integrations; calibration / dev-time loops are out of scope. The graph is data inside the platform but not portable across tools (each is a silo). A real on-ramp for non-engineers, with a ceiling.

### Low-code/no-code automation: n8n, Zapier, Make

Visual workflow automation focused on integrating SaaS apps. Zapier is consumer-grade; n8n is open-source and developer-leaning; Make sits between them.

Massive integration libraries (thousands of services) make these the practical choice when most blocks of a hybrid loop are *interactions with external SaaS systems* rather than LLM reasoning. They've added LLM nodes recently (n8n especially), making them increasingly hybrid-loop-shaped. Not designed around the LLM-as-fuzzy-mapper view; LLM nodes are bolted on alongside deterministic integration steps without the framework's disciplines applied. State, substrate, and calibration are absent or minimal. A real entry point when the project is mostly SaaS-integration with some LLM work.

### Compound engineering (Every.to / Kieran Klaassen)

A practitioner methodology for AI-assisted software development. Core loop: Plan → Work → Review → Compound → Repeat. The fourth step embeds learnings into searchable artifacts (CLAUDE.md updates, YAML-metadata repos) so subsequent work is "easier, not harder." Recommends 80% time on planning + review, 20% on implementation; multi-agent parallel review; "teach the system, don't do the work yourself."

The "compound" step is structurally the **dev-time hybrid loop wrapping the runtime** — feedback from runtime behavior reshapes the layers below. "Teach the system" maps onto *context-as-code as load-bearing infrastructure*. Multi-agent parallel review maps onto the *adversarial-panel-process* shape.

**Where it differs**:

- *Engineering-only.* Narrowly aimed at software development with AI agents. Hybrid loops is domain-agnostic — the same shape applies in teaching, coaching, advocacy, creative work. Compound engineering reads as the engineering-shaped instance of the broader pattern.
- *No calibration discipline.* The Review step assumes multi-agent critique substitutes for calibration. Review without persistent hit-rate tracking is just review — you can't tell whether reviewers are getting better or worse over time, or which agent has drifted.
- *Limited lineage engagement.* Presented as new without acknowledging that "feedback into the system improves the system" is a 60-year-old idea (cybernetics, autopoiesis, learning systems). The eight-beliefs-to-unlearn framing assumes a specific old-way it's overturning, narrower than the actual prior art.
- *Methodology vs design pattern.* "Five adoption stages, eight beliefs to unlearn" reads as consultancy-product packaging — a specific *prescribed* path. Hybrid loops is a design pattern with explicit decline-when criteria; the diagnostic phase says most projects DON'T need this.

Cite as a contemporaneous practitioner instance of the dev-time-loop discipline; don't adopt the methodology vocabulary.

### Discipline-coverage across the ecosystem

The disciplines named in `THE_CASE.md` (calibration, context-as-code, dev-time loop) plus the framework's additions (substrate-as-record, substrate-as-vocabulary, decline-when), tabulated against what each tool addresses:

| | calibration | context-as-code as infra | dev-time loop | substrate-as-record | substrate-as-vocabulary | decline-when |
|---|---|---|---|---|---|---|
| **DSPy** | yes (optimization-shaped) | implicit (signatures) | yes (compile) | no | no | no |
| **LangGraph** | no | partial (prompts) | no | partial (memory) | no | no |
| **AutoGen** | no | partial (system messages) | no | partial (transcript) | no | no |
| **CrewAI** | no | partial (role descriptions) | no | partial (memory) | no | no |
| **pydantic / instructor** | no | yes (schemas) | no | no | no | no |
| **Temporal / Conductor** | no | no | no | yes (durable state) | no | no |
| **Dify / LangFlow / Flowise** | no | partial | no | partial (canvas state) | no | no |
| **n8n / Zapier / Make** | no | partial | no | partial | no | no |
| **Compound engineering (Every.to)** | no | yes | yes | partial | no | no |
| **SPDD (Fowler)** | no | yes (REASONS canvas) | yes | no | no | no |
| **hybrid-loops** | yes (ship-blocking) | yes | yes | yes | yes | yes |

A note on the calibration column. DSPy's compile step optimizes against a user-provided metric, which is calibration in the engineering sense — the table marks it "yes (optimization-shaped)." Hybrid loops marks calibration as ship-blocking rather than as one optimization signal among many; the framework's distinctive move is treating per-block hit-rate as a separate gate, not a knob inside an optimizer. Different shape, same family. The frameworks each cover one or two pieces well; none name the same set as ship-blocking together.

## Positioning against the critic

Three lines that hold up:

> **"This is just LangGraph."** LangGraph is a runtime; this is a design pattern. You can execute a hybrid-loops graph in LangGraph if your graph is LLM-heavy; you'd use a thinner runtime if it isn't. The disciplines apply regardless of which runtime you pick.

> **"This is just DSPy."** DSPy is the closest cousin and you've internalized a chunk of this already. Hybrid loops puts deterministic non-LLM blocks on equal footing with LLM modules (DSPy's are second-class), names calibration as ship-blocking rather than as one optimization signal, distinguishes runtime cycles from dev-time cycles, and treats substrate as a first-class role. Use DSPy to implement and tune the LLM blocks of a hybrid-loops graph.

> **"pydantic and structured outputs already solved this."** Those tools implement *one* of the framework's disciplines (context-as-code constraint at LLM output). They don't tell you the rest (calibration, substrate audit, dev-time loop) or how to design the graph the LLM blocks live in. Different layers.

The shared move: acknowledge the overlap, then name what's still on the table.
