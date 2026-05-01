# Hybrid loops vs. agent frameworks

A natural critique: *"This sounds like LangGraph / AutoGen / DSPy / pydantic with an extra coat of theory. What's actually different?"* The answer needs care because there's real overlap with each tool, and the positioning depends on getting the **levels of abstraction** right.

## The actual gap these tools leave open

DSPy has *modules* and *signatures* and *optimizers*. LangGraph has *nodes* and *edges* and *checkpoints*. AutoGen has *agents* and *conversations* and *managers*. CrewAI has *crews* and *roles* and *tasks*. pydantic has *models* and *validators*. instructor has *response_model*. Anthropic and OpenAI have *tool use* and *structured outputs*.

**The frameworks are all swirling around the same broad space without a consistent vocabulary or taxonomy.** Each tool has its own primitives, its own mental model, its own opinions about what the basic unit is. None of them tell you how their primitives relate to the others'. A senior engineer who's used three of them carries three half-overlapping mental models in their head and reconciles them by intuition.

**A mid-level engineer who hasn't used any of them is in worse shape**: they're still operating in the conventional pipeline mental model (extract → process → return), haven't encountered the broader space yet, and when they do, the first framework they pick will lock them into one specific cell of the alphabet without showing them the rest. They won't see what they're missing because the tools don't surface the missing parts.

This is the actual gap the framework is filling. Not "build another competitor to LangGraph." **Provide the unifying vocabulary and taxonomy the ecosystem lacks**, so:
- A senior who's used three frameworks has one coherent mental model that names what each tool addressed and what each left untouched
- A mid-level engineer who's used none of them has a starting model that points at the broader space and tells them which tool to reach for at each block

## TL;DR — different levels

Hybrid loops is a **design pattern** at one level of abstraction *up* from agent frameworks. The frameworks below are **implementation toolkits** for pieces of the pattern.

- The pattern says: *"design your system as a graph of typed blocks where LLM-actors and code-actors alternate, with three disciplines (calibration, context-as-code, dev-time loop) applied to the LLM blocks."*
- The frameworks say: *"here's a runtime / programming model / structured-output validator for the LLM portion of your graph."*

Critique that flattens to *"this is just LangGraph"* misunderstands the levels. Critique that flattens to *"this is just architecture line-noise; the frameworks are real"* misses that the disciplines are the framework's content. **Both critiques are wrong, both are tempting, and the comparison below tries to forestall both.**

## DSPy (Khattab et al., Stanford)

**What it is**: a programming model for LM programs. Typed *signatures* (input/output fields) define modules; *optimizers* (BootstrapFewShot, MIPROv2, etc.) search over prompts and few-shot demonstrations to maximize a metric. The "compile" step turns a program + metric + examples into an optimized program.

**Closest cousin in this work**. DSPy's typed-signature module is essentially a typed-block-with-an-LLM-actor; its compile-for-metric loop is structurally the same as `schemaforge`'s compress-and-verify loop, just optimizing prompts/demos rather than notation. **If you've internalized DSPy you've internalized a chunk of the hybrid-loops algebra already.**

**Where they differ**:
- DSPy's alphabet is *mostly LLM modules*. Deterministic non-LLM nodes are second-class — you reach for them via Python interop, not as first-class architectural elements. Hybrid-loops puts LLM and code on equal footing in the alphabet.
- DSPy optimizes for a metric the user provides. Hybrid-loops *names disciplines* — calibration is one of them — and treats per-block hit-rate as ship-blocking rather than as one optimization signal among many.
- DSPy doesn't distinguish runtime cycles from development-time cycles; the compile step *is* a dev-time loop, but the framework treats it as a tooling concern, not an architectural one.
- DSPy doesn't have an explicit substrate concept — programs are stateless w.r.t. each other. Hybrid-loops treats the substrate (typed-records-accumulating-over-time) as a first-class role.

**Verdict**: complementary. **Use DSPy to implement the LLM blocks of a hybrid-loop graph when prompt/demo optimization is worth the compile-step cost.** The frameworks are not at war.

## LangGraph (LangChain)

**What it is**: a graph executor for LLM-centric agent workflows. Nodes are typically LLM calls or tool calls; edges express the control flow; supports cycles (hence "graph" not "chain").

**Where they differ**:
- LangGraph nodes are typically *LLM calls or tools*. The deterministic-code-as-substrate / deterministic-gate / typed-substrate-store roles are implemented ad hoc. Hybrid-loops names those roles structurally.
- LangGraph's graph IS code (Python). Hybrid-loops's graph IS data (when realized) — the "graph-as-data" recursion in BUILDING_BLOCKS.md is a deliberate move toward graph specs that code can validate, lint, simulate. LangGraph doesn't go there.
- LangGraph doesn't have an opinion on calibration, context-as-code as infra, or dev-time-loop discipline. It's a runtime, not a methodology.

**Verdict**: complementary at runtime. **A hybrid-loop graph could be executed by LangGraph or by direct Python or by an MCP-tool sequence — the framework doesn't prescribe a runtime.** LangGraph is the most natural fit for projects whose graph is mostly LLM nodes; for graphs with heavy deterministic infrastructure, a thinner runtime works better.

## AutoGen (Microsoft)

**What it is**: a framework for multi-agent conversations. Agents have roles, messaging works as conversation, group chats coordinate via a "manager" agent. Strong on agent-talks-to-agent shapes.

**Where they differ**:
- AutoGen's primary primitive is *the conversation* — agents exchange messages, context accumulates as transcript. Hybrid-loops's primary primitive is *the typed substrate* — records accumulate as structure, not as conversation.
- AutoGen treats multi-agent as the default shape. Hybrid-loops treats LLM blocks as one node type among others; whether multiple LLM blocks "converse" is a question of substrate-shape, not framework default.
- AutoGen doesn't separate runtime from dev-time loops. Multi-agent conversations are runtime constructs; the dev-time critique-and-patch loop is something users implement themselves.

**Verdict**: AutoGen's shape is one specific cell of the hybrid-loops alphabet — *multiple LLM blocks with conversational substrate*. It's powerful for that shape. **Hybrid-loops is broader than that one shape; AutoGen is deeper in it.** Use AutoGen when the right pattern is multi-agent conversation; use hybrid-loops as the umbrella when it isn't.

## CrewAI

**What it is**: opinionated role-based multi-agent — "researcher" agents, "writer" agents, "critic" agents working together with role descriptions and goals.

**Where they differ**: same as AutoGen at the structural level — CrewAI is one specific shape (role-described agents collaborating) within the broader hybrid-loops alphabet. CrewAI is more opinionated than AutoGen about *how* the multi-agent shape is constructed; both occupy the same architectural cell.

**Verdict**: niche-within-the-shape. Use when the role-described-agents pattern fits.

## pydantic / instructor / Anthropic structured outputs / OpenAI structured outputs

**What they are**: structured-output validators / typed-schema enforcers for LLM calls. pydantic provides Python type-checking; instructor wraps LLM calls in pydantic models; Anthropic's tool use and OpenAI's structured outputs let you define a JSON Schema the LLM is constrained to output.

**Closest match for the schemas-as-context-as-code subsection of BUILDING_BLOCKS.md**. These tools are how you implement the discipline: define a typed schema once, the LLM is constrained by it at output time, validators catch malformed outputs at runtime.

**Where they differ**:
- pydantic/instructor/structured-outputs are *implementations of one block-level discipline* (constrain LLM output to a typed schema). Hybrid-loops names that discipline and explains *why* it's load-bearing — but doesn't compete with the implementations.
- These tools don't tell you which schema to use. Hybrid-loops's `schemaforge` does (in some sense): it discovers a dense notation through compress-and-verify. The two are complementary — schemaforge designs the schema; pydantic enforces it at runtime.

**Verdict**: not even adjacent — different layers entirely. **Use pydantic + instructor + structured-outputs to implement context-as-code constraints in your LLM blocks; the framework names *that this is a discipline* and explains why every LLM block in a non-trivial graph needs one.**

## What makes hybrid-loops distinct as a framework

Tabulating the disciplines named here against what each tool addresses:

| | calibration | context-as-code as infra | dev-time loop | substrate-as-record | substrate-as-vocabulary | decline-when discipline |
|---|---|---|---|---|---|---|
| **DSPy** | partial (metric-driven) | implicit (signatures) | yes (compile) | no | no | no |
| **LangGraph** | no | partial (prompts) | no | partial (memory) | no | no |
| **AutoGen** | no | partial (system messages) | no | partial (transcript) | no | no |
| **CrewAI** | no | partial (role descriptions) | no | partial (memory) | no | no |
| **pydantic / instructor** | no | yes (schemas) | no | no | no | no |
| **hybrid-loops** | yes (`cal_log`) | yes (substrate-audit via `metacog`; schema-discovery via `schemaforge`) | yes (named regime) | yes (first-class) | yes (first-class) | yes (Phase 2 A/B/C scoping) |

Reading: hybrid-loops names a wider set of disciplines as ship-blocking. The frameworks each cover one or two pieces well. **The framework's distinctive content is the union of disciplines plus the diagnostic for when the pattern doesn't fit.**

## How to position the framework against an agent-framework critic

If a senior engineer says *"this is just LangGraph"*:

> "LangGraph is a runtime. The framework is a design pattern. You can execute a hybrid-loops graph in LangGraph if your graph is LLM-heavy; you'd use a thinner runtime if it isn't. The disciplines named in the framework — per-block calibration, context-as-code as load-bearing infrastructure, dev-time hybrid loop wrapping the runtime — apply regardless of which runtime you pick. The framework isn't claiming to replace LangGraph; it's claiming there are disciplines LangGraph users still have to figure out per project."

If they say *"this is just DSPy"*:

> "DSPy is the closest cousin and you've internalized a chunk of this already. The differences are: hybrid-loops puts deterministic non-LLM nodes on equal footing with LLM modules in the alphabet (DSPy's are second-class), names calibration as ship-blocking rather than as one optimization signal, distinguishes runtime cycles from dev-time cycles as separate disciplines, and treats substrate as a first-class architectural role. Use DSPy to implement and tune the LLM blocks of a hybrid-loops graph."

If they say *"pydantic and structured outputs already solved this"*:

> "Those tools implement *one* of the framework's disciplines (context-as-code constraint at LLM output). They don't tell you the rest of the disciplines (calibration, substrate audit, dev-time loop) or how to design the graph the LLM blocks live in. The frameworks aren't at war; they're at different layers."

The shared move in each response: **acknowledge the overlap, then name what's still on the table.**
