# Building blocks of hybrid loops

A hybrid-loop project is a *graph of blocks* that snap together by type. This file enumerates the eight primitive blocks, gives concrete examples, then shows how they compose into pairs and triples — and how any longer chain is just more of the same.

> *Vocabulary note: this repo uses **block** consistently — for the typed primitives, for their instances in a specific graph, and for what other frameworks variously call "nodes" (LangGraph), "modules" (DSPy), "agents" (AutoGen), "steps" (Temporal). When citing those frameworks' own vocabularies the original term is kept; everywhere else, "block."*

## Setup: two actors, one medium

Two actors — *deterministic code* and *LLM* — connected through one universal medium, *data*. Code is also data (the executor knows how to run it), so the only meaningful distinction at any block is *who's acting* and *what shape the data has at the input and output*.

This isn't the only useful breakdown — a project might further partition the actor space (small classifiers, vector search, optimizers, simulation engines, human reviewers). The 2×2 below is the lowest-friction starting partition; substitute richer ones as the project demands.

The "code" / "data" axis hides a real category between them: **context-as-code**. A markdown rule sheet is *data* to deterministic code (parse it, version it, lint it) and *code* to an LLM (the LLM interprets the rules and changes its behavior accordingly). This is why "LLM writes a system prompt for another LLM" is a real primitive, not a curiosity — it's code generating code, where both sides happen to be natural language and the executor is the LLM. `THE_CASE.md` unpacks the discipline of treating context-as-code as production infrastructure; the table below uses the simpler {code, data} distinction.

| | reads data | reads code |
|---|---|---|
| **code** produces data | filters, queries, aggregates, ETL, math | compilers, linters, type checkers, test runners |
| **code** produces code | codegen from spec, schema-to-types, parser generators | source-to-source transforms, macros, refactoring tools, optimizers |
| **LLM** produces data | classify, summarize, extract, score, translate, answer | code review, explain, document, audit |
| **LLM** produces code | NL → SQL/regex/function, design notation/schema | refactor, fix bug, translate language, optimize, port |

## Eight primitive blocks

### 1. Code reads data → produces data

The dominant case in all of computing.

- SQL: `SELECT * FROM events WHERE ts > now() - interval '1 day'`
- Filter: `[r for r in records if r.score > 0.7]`
- Aggregate transactions by category, monthly P&L from a ledger
- HTTP request handler: incoming JSON → outgoing JSON
- Hash a file, sort by key, join two tables, deduplicate by ID

### 2. Code reads data → produces code

Less common but well-established.

- Protobuf compiler: `.proto` schema → generated Go/Python/TypeScript
- OpenAPI generator: API spec → client SDK
- JSON Schema → TypeScript types (`quicktype`, `json-schema-to-typescript`)
- ORM scaffolding from database introspection
- Table-driven test generators reading fixture data
- Build-system targets generated from manifest files (Bazel, Make patterns)

### 3. Code reads code → produces data

The static-analysis family.

- Compilers (parse phase): source → AST → diagnostic data
- Linters: code → list of style/correctness findings
- Type checkers: code → list of type errors
- Test runners: test code → pass/fail records, coverage data
- Static analyzers: SonarQube, CodeQL
- Profilers and instrumented runtimes: running code → execution traces

### 4. Code reads code → produces code

The transformation family. Source-to-source.

- Optimizing compilers: source → optimized binary
- Lisp / Rust / Scala macros: code → expanded code at compile time
- Babel / TypeScript / SWC: modern JS → older JS / WASM
- Refactoring tools: rename variable, extract function, inline (`gofmt -r`)
- Codemods for large-scale API migrations (`jscodeshift`)
- Prettier / gofmt / black: code → reformatted code

### 5. LLM reads data → produces data

The popular case since ChatGPT. The lens role in the five-role default.

- Classify intent: "is this a complaint, question, or compliment?"
- Summarize: 10K-token doc → 200-token TL;DR
- Extract typed entities: news article → `[{person, role, organization}]`
- Score sentiment: review text → `{negative: 0.8, positive: 0.1, neutral: 0.1}`
- Translate language
- Answer questions about provided context (the LLM half of RAG)
- Pick from options

### 6. LLM reads data → produces code

The schema-and-notation-design half of generative authoring.

- NL → SQL, regex, API call
- Function from docstring (copilot-style autocomplete)
- Design a schema from sample records: 50 JSON examples → proposed JSON Schema
- Design a dense notation for a corpus (compress-and-verify shape)
- Author a gate function: "write a Python predicate that filters this set the way we just discussed"
- Synthesize a prompt — meta-prompting where one LLM writes prompts for another

### 7. LLM reads code → produces data

Code as soft input. The audit / review family.

- Code review: PR diff → review comments as findings
- Explain what code does: function → English summary
- Estimate complexity / risk
- Security audit: code → suspected vulnerabilities with severity
- Spec compliance check: code + spec → pass/fail per item
- Read traces / logs: instrumented runtime output → root-cause hypotheses
- Map dependencies: code → graph data describing module relationships

### 8. LLM reads code → produces code

The refactor / fix / port family.

- Refactor: clean up naming, restructure for readability
- Fix bug given symptoms: failing test + function → patch
- Translate language: Python → Go (preserve semantics)
- Port across frameworks: Express → Fastify
- Optimize a hot path given a profile
- Add a feature: existing code + description → patched code
- Inline a refactor that conventional tools won't touch (cross-cutting concerns)

## Why neither half collapses to the other

When LLMs got good, the temptation was to collapse the algebra: replace `code acts on data` with `LLM acts on data` (just describe the rules in a prompt) or with `LLM writes code that acts on data` (let the LLM author the deterministic side). Both fail in well-named ways.

### Fallacy 1 — "Just tell the LLM the rules in a prompt"

Replace `code-filter-by-predicate` with `LLM-filter-by-predicate-described-in-prompt`. Tempting because writing rules in English is easier than writing them in code. The failure: LLMs are famously fickle about following directions consistently. Same prompt, same input, different sampling — different result. Variants:

- The rule applies on most inputs but silently bends on the edge cases that matter most
- Rule applies fine per item but drifts across a long context window — by item 80 of a list, item 1's rules have been forgotten
- The rule depends on a definition you didn't think to spell out; the LLM substitutes a related-but-different one
- The LLM is following the rule, but what *you* meant by it turns out to be ambiguous, and the LLM disambiguated against you

Code follows directions perfectly because *directions are the code* — no translation step between rule and execution. With an LLM, there's an approximate translation from your stated rule into whatever the model's distribution actually does. For load-bearing rules, that's a no-go.

### Fallacy 2 — "Have the LLM write the code, then run it deterministically"

Replace `code-filter` with `LLM-author-filter-code → code-run-it`. Tempting because you get determinism from the run step while keeping LLM flexibility in authoring. The failure: LLMs are quirky in their authoring, and a wrong-but-runnable program is worse than no program. Variants:

- The code looks right but has a subtle off-by-one or wrong column name
- The LLM imports a library the runtime doesn't have, or imports the wrong one with the same name
- Idiomatic but semantically wrong code for an edge case the description didn't mention
- Same description produces slightly different code each time — version skew across runs

These are *mitigated*, not eliminated, by a deterministic verification step (typecheck, run tests). What you want is the triple `LLM-author + code-verify + LLM-revise-on-failure` — the agentic-codegen-with-verification pattern. The verification step is what makes the LLM-as-author primitive trustworthy enough to land in production. Without it, you're hoping.

### What hybrid loops do instead

Hybrid loops don't collapse the algebra. They make the deterministic half explicitly responsible for *the consistency the LLM lacks*, and the LLM half responsible for *the fluency the deterministic side lacks*. Per-block calibration is what tells you each block is earning its keep — see `THE_CASE.md` for the full disciplines argument.

## Pairs: where primitives snap together

A pair is two primitives chained where the first's output type matches the second's input. The connection is the load-bearing claim — *"data → data" connects to anything; "data → code" requires the next block to be a code-runner; "code → data" follows naturally from any code-on-code analyzer.*

### Pairs starting from "code → data"

- `code-query` + `LLM-summarize` — RAG's cheap variant: SQL query → LLM-written report
- `code-filter` + `LLM-decide` — narrow candidates by hard rules first, then have the LLM pick from survivors
- `code-aggregate` + `LLM-explain` — compute monthly metrics, LLM writes the dashboard narrative
- `code-instrument` + `LLM-audit` — capture runtime traces, LLM finds patterns in them (development-time)

### Pairs starting from "code → code"

- `code-codegen` + `code-typecheck` — generate from spec, immediately verify it compiles
- `code-refactor` + `LLM-explain` — transform code, LLM writes the commit message

### Pairs starting from "LLM → data"

- `LLM-extract` + `code-persist` — read transcripts, write typed records to substrate (the canonical lens path)
- `LLM-classify` + `code-dispatch` — read intent, route to specialist handler (the canonical agentic shape)
- `LLM-score` + `code-rank` — assign scores, sort by them
- `LLM-summarize` + `LLM-summarize` — multi-stage summarization: per-doc → overall

### Pairs starting from "LLM → code"

- `LLM-author-SQL` + `code-run` — text-to-SQL with deterministic execution
- `LLM-author-regex` + `code-apply` — LLM writes the pattern, code runs it on the corpus
- `LLM-author-function` + `code-typecheck` — LLM proposes a function, type-checker verifies
- `LLM-design-schema` + `code-validate` — LLM proposes JSON Schema, code validates samples
- `LLM-design-notation` + `code-compress` — design a notation, deterministically apply it (the compress-and-verify shape)

### Pairs from "LLM reads code"

- `LLM-review-code` + `code-persist` — review findings as typed records (calibration substrate for review quality)
- `LLM-explain-code` + `code-render-docs` — LLM produces explanations, code stitches them into a doc site
- `LLM-audit-code` + `code-prioritize` — LLM emits findings, code prioritizes by severity / cooldown / age
- `LLM-refactor` + `code-test` — LLM rewrites, test suite verifies
- `LLM-port` + `code-typecheck` — LLM translates Python → Go, type-checker enforces
- `LLM-fix-bug` + `code-bisect` — LLM proposes patch, bisect confirms it resolves the regression

These pairs aren't exhaustive. They cover the well-trodden compositions; anywhere two primitives have type-compatible I/O, they snap together.

## Triples: recognizable patterns

Three blocks chained close enough to a loop that you recognize them as named patterns. Diagrams and per-shape commentary live in `BLOCK_GRAPHS.md`; the sketches below are just the type-compositions.

### RAG — Retrieval-Augmented Generation

```
code-query  →  data  →  LLM-reason-with-context  →  data
```

The simplest RAG is a pair (query + reason); it becomes a triple when the answer is post-processed by code (citation linking, fact-checking, persisting).

### ReAct — Reason + Act, with verification

```
LLM-decide-next-action  →  code  →  code-run-tool  →  data  →  LLM-reason-on-result  →  ...
```

Agent decides, deterministic tool runs, LLM reads the result and decides again. The deterministic tool-run is what makes the agent grounded.

### Codegen with verification

```
LLM-author-code  →  code  →  code-typecheck-or-test  →  data  →  LLM-revise-on-failures
```

LLM proposes, type checker (or test suite) verifies, LLM iterates on failures. The deterministic verifier is the gate that prevents drift.

### The five-role canonical hybrid loop

```
soft-data  →  LLM-extract  →  data  →  code-filter+score+rank  →  data  →  LLM-reason  →  data  →  code-action
                                                                                                        │
                                                                                              new soft-data ←┘
```

What `SKILL.md` describes as the default arrangement. The cycle closes when the action produces new content the lens reads next turn.

### Development-time critique loop

```
code-instrument  →  data  →  LLM-critic-panel  →  data  →  code-prioritize  →  data  →  LLM-write-patch  →  code  →  code-apply
                                                                                                                        │
                                                                                                            modified runtime ←┘
```

Wraps around the runtime hybrid loop. Most real systems live here — the runtime is small, the dev loop iterates many times across runs, and the dev loop is what earns the system's reliability over time.

## How this generalizes

The above are named patterns because someone already named them. Every other useful chain is more of the same algebra: pick the next primitive whose input type matches the previous primitive's output, and you have a valid hybrid-loop fragment.

Branching, joining, dynamic dispatch, recursion (a block calling a sub-loop), parallel fan-out / fan-in — these are higher-order composition operators on the same primitive set. They don't change the primitive vocabulary; they change the topology of how primitives connect.

The disciplines that keep composition from collapsing:

1. *Schema versioning at every typed edge.* When the upstream block's output schema changes, every downstream block sees it.
2. *Calibration at every load-bearing LLM block.* Without it, you can't tell which block is dropping the ball when a chain misbehaves.
3. *Restraint where the LLM block isn't naturally restrained.* The deterministic block immediately downstream of an LLM block is where most projects' opinionated policy lives — filtering, thresholding, normalizing, ranking.
4. *Loop-closure visibility.* Which arrows feed back, when, on what cadence? In a real system there might be three different feedback paths. Diagram them.

The combinatorial space looks daunting until you realize you're never visiting all of it. A typical hybrid-loop app is 5–15 blocks arranged in 1–3 cycles. The skill in `../SKILL.md` exists because *most programmers were trained on three or four of the eight primitives, not the whole eight, and certainly not the compositional algebra that connects them.*
