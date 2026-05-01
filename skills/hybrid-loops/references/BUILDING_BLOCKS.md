# Building blocks of hybrid loops

A hybrid-loop project is a **graph of blocks** that snap together by type. This file enumerates the blocks, gives concrete examples for each, then shows how they compose into pairs and triples — and how any longer chain is just more of the same.

## Setup: two technologies, one medium

The breakdown that follows uses two actors — **the old way (deterministic code)** and **the new way (LLMs)** — connected through one universal medium, **data**. Code is also data (the executor knows how to run it), so the only meaningful distinction at any block is *who's acting* and *what shape the data has at the input and output*.

This isn't the only useful breakdown. A given project might further partition the actor space (e.g. small classifiers, vector search, deterministic optimizers, simulation engines, human reviewers). The 2×2 below is the lowest-friction starting partition; substitute richer ones as the project demands.

### A note on context-as-code

The "code" / "data" axis hides a real category that lives between them: **context-as-code**. Markdown docs with rules (`CLAUDE.md`, skill descriptions, system prompts), prompt templates, persona definitions, behavioral specs in natural language — these are *data* when handled by deterministic code (you can read, parse, version-control, lint, or summarize a `CLAUDE.md` with code), but they're *code* when fed to an LLM (the LLM interprets the rules and changes its behavior accordingly).

The unifying view: **an LLM is an interpreter of natural-language code.** Just as a CPU interprets binary instructions, an LLM interprets context. The text-to-binary distinction in classical computing has an analog here in the text-to-LLM-behavior interpretation. This is what makes "LLM writes a system prompt for another LLM" — meta-prompting — a real primitive, not a curiosity: it's just code generating code, where both sides happen to be natural language and the executor is the LLM.

This means several block-shapes that look like "LLM produces data" or "code produces data" are really **producing context-as-code** that downstream blocks treat as code. The examples below note where this distinction matters.

### Schemas, DSLs, and the highest-leverage context-as-code

A specific case worth naming: **schemas, DSLs, and structured-output specs are the highest-leverage context-as-code artifacts**. They hit all three benefits at once:

- An LLM can author them (the user-facing-feature half of "LLM writes code")
- Other LLMs can be *constrained by them* at inference time (structured-outputs mode, tool-use schemas, "respond in this JSON shape" prompts)
- Deterministic code can validate against them (JSON Schema validators, Protobuf compilers, type checkers, parser generators)

So the same artifact is simultaneously: an LLM output, an LLM input, and a code-side constraint. **One file, three readers, three behaviors.** That's the multiplicatively-powerful position the framework keeps converging on.

This is precisely the operating mode of practitioners like Manuel Odendahl ("wesen") whose toolchain (`geppetto`'s typed-step abstraction, `pinocchio`'s YAML-with-metadata prompt libraries, `glazed`'s typed rows-and-columns, `prompto` / `promptos` retrieval-by-typed-context) consists of LLM-authored structure definitions that downstream code and downstream LLMs both consume. He's been operating in this mode for years; the framework here describes what he's been doing.

| | reads data | reads code |
|---|---|---|
| **code (old way)** produces data | filters, queries, aggregates, ETL, math | compilers, linters, type-checkers, test runners |
| **code (old way)** produces code | codegen from spec, schema-to-types, parser generators | source-to-source transforms, macros, refactoring tools, optimizers |
| **LLM (new way)** produces data | classify, summarize, extract, score, translate, answer | code review, explain, document, audit |
| **LLM (new way)** produces code | NL → SQL/regex/function, design notation/schema | refactor, fix bug, translate language, optimize, port |

## Eight primitive blocks

### 1. Code reads data → produces data

The dominant case in all of computing. Includes most of what's been written.

- **SQL query**: `SELECT * FROM events WHERE ts > now() - interval '1 day'` — request data → result rows
- **Filter records by predicate**: `[r for r in records if r.score > 0.7]`
- **Aggregate transactions by category**: monthly P&L from a ledger
- **HTTP request handler**: incoming JSON → outgoing JSON
- **Hash a file**: bytes in, hex digest out
- **Sort by key**, **join two tables**, **deduplicate by ID** — the standard library of data work

### 2. Code reads data → produces code

Less common but well-established in tooling.

- **Protobuf compiler**: `.proto` schema (data) → generated Go/Python/TypeScript code
- **OpenAPI generator**: API spec (YAML/JSON) → client SDK code
- **JSON Schema → TypeScript types**: `quicktype`, `json-schema-to-typescript`
- **ORM scaffolding from database introspection**: connect to DB, dump schema, emit model classes
- **Table-driven test generators**: read fixture data, emit Go tests
- **Build-system targets generated from manifest files**: Bazel, Make patterns

### 3. Code reads code → produces data

The static-analysis family. Reads source as input, produces records describing it.

- **Compilers (parse phase)**: source code → AST → diagnostic data (errors, warnings)
- **Linters**: code → list of style/correctness findings
- **Type checkers**: code → list of type errors
- **Test runners**: test code → pass/fail records, coverage data
- **Static analyzers**: code → security/complexity findings (SonarQube, CodeQL)
- **Profilers / instrumented runtimes**: running code → execution traces

### 4. Code reads code → produces code

The transformation family. Source-to-source.

- **Optimizing compilers**: source → optimized binary (binary is also code)
- **Lisp macros, Rust macros, Scala macros**: code → expanded code at compile time
- **Babel / TypeScript / SWC**: modern JS → older JS / WASM
- **Refactoring tools**: rename variable, extract function, inline (IDE refactors, `gofmt -r`)
- **Codemods**: large-scale API migrations (e.g., `jscodeshift`)
- **Prettier / gofmt / black**: code → reformatted code

### 5. LLM reads data → produces data

The popular case since ChatGPT. The lens role in the five-role default.

- **Classify intent of a user message**: "is this a complaint, question, or compliment?"
- **Summarize a long document**: 10K-token doc → 200-token TL;DR
- **Extract typed entities**: news article → `[{person, role, organization}]` records
- **Score sentiment**: review text → `{negative: 0.8, positive: 0.1, neutral: 0.1}`
- **Translate language**: English → Spanish
- **Answer questions about provided context** (the LLM half of RAG)
- **Pick from options**: "given these 5 candidates, which fits the criterion best?"

### 6. LLM reads data → produces code

The schema-and-notation-design half of generative authoring. **`schemaforge` in this repo is exactly this primitive applied iteratively.**

- **Natural language → SQL**: "show me April revenue by region" → `SELECT ...`
- **NL → regex**: "match phone numbers in any format" → `\+?\d[\d\-\s\(\)]{8,}`
- **NL → API call**: "find papers on hybrid systems by Andreas" → `arxiv.search(...)`
- **Function from docstring**: copilot-style autocomplete
- **Design a schema from sample records**: ingest 50 JSON examples → propose JSON Schema
- **Design a dense notation for a corpus**: schemaforge's `design_notation` tool
- **Author a gate function**: "write a Python predicate that filters this set the way we just discussed"
- **Synthesize a prompt**: meta-prompting where one LLM writes prompts for another

### 7. LLM reads code → produces data

Code as soft input. The audit / review family.

- **Code review**: read PR diff, output review comments as findings
- **Explain what code does**: function → English summary for a doc
- **Estimate complexity / risk**: code → ranked findings
- **Security audit**: code → list of suspected vulnerabilities (with CVE-style severity)
- **Spec compliance check**: code + spec → pass/fail per spec item
- **Read traces / logs**: instrumented runtime output → root-cause hypothesis records
- **Map dependencies**: code → graph data describing module relationships

### 8. LLM reads code → produces code

The refactor / fix / port family. The new way doing what the old way's refactoring tools did, but with semantic understanding.

- **Refactor**: clean up naming, restructure for readability
- **Fix bug given symptoms**: "this test fails with X, here's the function" → patch
- **Translate language**: Python → Go (preserve semantics)
- **Port across frameworks**: Express → Fastify
- **Optimize a hot path**: rewrite for performance given a profile
- **Add a feature**: existing code + feature description → patched code
- **Inline a refactor that conventional tools won't touch**: cross-cutting concerns

## Why neither half collapses to the other

The temptation when LLMs got good was to collapse the algebra: replace `code acts on data` with `LLM acts on data` (just describe the rules in a prompt) or with `LLM writes code that acts on data` (let the LLM author the deterministic side). Both fail in well-named ways, and both failures are why hybrid loops aren't optional.

### Fallacy 1 — "Just tell the LLM the rules in a prompt"

Replacing `code-filter-by-predicate` with `LLM-filter-by-predicate-described-in-prompt`. Tempting because writing rules in English is easier than writing them in code. The failure: **LLMs are famously fickle about following directions consistently.** Same prompt, same input, different sampling — different result. Variants:

- The rule applies on most inputs but silently bends on the edge cases that matter most
- The rule applies fine on individual items but drifts across a long context window — by item 80 of a list, the rules from item 1 have been forgotten or down-weighted
- The rule depends on a definition you didn't think to spell out; the LLM substitutes a related-but-different one
- The LLM is following the rule, but what *you* meant by the rule turns out to be ambiguous, and the LLM disambiguated against you

Code follows directions perfectly because *directions are the code* — no translation step between the rule and its execution. With an LLM, there's an approximate translation from your stated rule into whatever the model's distribution actually does. **For load-bearing rules, that's a no-go.**

### Fallacy 2 — "Have the LLM write the code, then run it deterministically"

Replacing `code-filter` with `LLM-author-filter-code → code-run-it`. Tempting because you get determinism from the run step while keeping LLM flexibility in authoring. The failure: **LLMs are quirky and non-deterministic in their authoring**, and a wrong-but-runnable program is worse than no program. Variants:

- The code looks right but has a subtle off-by-one or wrong column name
- The LLM imports a library the runtime doesn't have, or imports the wrong one with the same name
- The LLM writes idiomatic but semantically wrong code for an edge case the description didn't mention
- The same description produces slightly different code each time — version skew across runs

These are *mitigated*, not eliminated, by a deterministic verification step (`code-typecheck`, `code-test`). What you actually want is a triple: **`LLM-author + code-verify + LLM-revise-on-verify-failure`** — the agentic-codegen-with-verification pattern. The verification step is what makes the LLM-as-author primitive trustworthy enough to land in production. Without it, you're hoping.

### What hybrid loops do differently

Hybrid loops don't collapse the algebra. They make explicit that **the deterministic half provides the consistency the LLM lacks, and the LLM half provides the fluency the deterministic side lacks**. The cycle is the answer to: how do you keep a flexible-but-quirky actor productive over many turns? *You don't make it act unilaterally on the parts where consistency matters; you alternate it with a deterministic actor that handles those parts.*

Calibration logs (the `cal_log` MCP server in this repo) make this concrete: they record per-evaluator hit-rate and roll up over time. If an LLM block is reliable enough on its slice of work, fine — keep it. If not, you have data, not opinions, and you decide what to do (tighten the gate downstream, replace the LLM with code, replace the prompt with a deterministic spec). **Every LLM block in the algebra is a candidate for replacement-by-code if its calibration says it should be.** The algebra doesn't claim LLMs are *necessary* anywhere — it claims that if you put one in, you also need the deterministic half nearby to catch what the LLM half misses.

## Pairs: where primitives snap together

A pair is two primitives chained where the first's output type matches the second's input type. The connection is the load-bearing claim — *"data → data" connects to anything; "data → code" requires the next block to be a code-runner; "code → data" follows naturally from any code-on-code analyzer.*

### Pairs starting from "code → data"

- **`code-query` + `LLM-summarize`** — RAG's cheap variant: SQL query → LLM-written report
- **`code-filter` + `LLM-decide`** — narrow candidates by hard rules first, then have the LLM pick from the survivors
- **`code-aggregate` + `LLM-explain`** — compute monthly metrics, LLM writes the narrative for the dashboard
- **`code-instrument` + `LLM-audit`** — capture runtime traces, LLM finds patterns in them (a development-time block)

### Pairs starting from "code → code"

- **`code-codegen` + `code-typecheck`** — generate from spec, immediately verify it compiles
- **`code-refactor` + `LLM-explain`** — transform code, LLM writes commit message describing the change

### Pairs starting from "LLM → data"

- **`LLM-extract` + `code-persist`** — read transcripts, write typed records to substrate (the canonical lens path)
- **`LLM-classify` + `code-dispatch`** — read intent, route to specialist handler (the canonical agentic shape)
- **`LLM-score` + `code-rank`** — assign scores, sort by them
- **`LLM-summarize` + `LLM-summarize`** — multi-stage summarization: per-doc → overall

### Pairs starting from "LLM → code"

- **`LLM-author-SQL` + `code-run`** — natural-language-to-SQL with deterministic execution (text-to-SQL)
- **`LLM-author-regex` + `code-apply`** — LLM writes the pattern, code runs it on the corpus
- **`LLM-author-function` + `code-typecheck`** — LLM proposes a function, type-checker verifies (autocomplete with verification)
- **`LLM-design-schema` + `code-validate`** — LLM proposes JSON Schema, code validates samples against it
- **`LLM-design-notation` + `code-compress`** — schemaforge's pattern: design a notation, deterministically apply it

### Pairs starting from "LLM → data" (LLM read code)

- **`LLM-review-code` + `code-persist`** — code review findings stored as typed records (a calibration substrate for review quality)
- **`LLM-explain-code` + `code-render-docs`** — LLM produces explanations, code stitches them into a doc site
- **`LLM-audit-code` + `code-prioritize`** — LLM emits findings, code prioritizes by severity / cooldown / age

### Pairs starting from "LLM → code" (LLM read code)

- **`LLM-refactor` + `code-test`** — LLM rewrites, test suite verifies
- **`LLM-port` + `code-typecheck`** — LLM translates Python → Go, type-checker enforces
- **`LLM-fix-bug` + `code-bisect`** — LLM proposes patch, code-bisect confirms it resolves the regression

These pairs are not exhaustive, but they cover the well-trodden compositions. **Anywhere two primitives have type-compatible I/O, they snap together — that's the entire compositional rule.**

## Triples: recognizable patterns

Three blocks chained close enough to a loop that you recognize them as named patterns.

### RAG (Retrieval-Augmented Generation)

```
code-query  →  data  →  LLM-reason-with-context  →  data
[fetch docs]      [answer using fetched docs]
```

The simplest RAG is a pair (`code-query` + `LLM-reason`). It becomes a triple when the answer is post-processed by code (citation linking, fact-checking, persisting). The full RAG cycle adds a calibration block: store the answer + later signal of whether it was correct.

### ReAct (Reason + Act, with verification)

```
LLM-decide-next-action  →  code  →  code-run-tool  →  data  →  LLM-reason-on-result  →  ...
                               [the loop continues]
```

Agent decides, deterministic tool runs, LLM reads the result and decides again. Three blocks min; loops indefinitely. The deterministic tool-run is what makes the agent grounded.

### Codegen with verification

```
LLM-author-code  →  code  →  code-typecheck-or-test  →  data  →  LLM-revise-on-failures
                                                    [if errors, loop back]
```

LLM proposes, type checker (or test suite) verifies, LLM iterates on failures. The deterministic verification is the gate that prevents drift; the LLM closes the loop by rewriting in response to typed feedback.

### Canonical hybrid loop (the five-role shape)

```
soft-data  →  LLM-extract  →  data  →  code-filter+score+rank  →  data  →  LLM-reason  →  data  →  code-action
                                                                                                       │
                                                                                                       ▼
                                                                                             new soft-data
                                                                                              (loops back)
```

What `SKILL.md` describes as the default arrangement. Each arrow is one primitive; the loop closes when the action produces new content the lens reads next turn.

### Development-time critique loop (wraps around the runtime)

```
code-instrument-runtime  →  data  →  LLM-critic-panel  →  data  →  code-prioritize  →  data  →  LLM-write-patch  →  code  →  code-apply
                                                                                                                                       │
                                                                                                                                       ▼
                                                                                                                              modified runtime layers
```

Five blocks. Wraps around the runtime hybrid loop. **Most real systems live here** — the runtime is small, the development loop iterates many times across runs, and it's the dev loop that earns the system's reliability over time.

## How this generalizes

The above are named patterns because someone already named them. Every other useful chain is just **more of the same algebra**: pick the next primitive whose input type matches the previous primitive's output type, and you have a valid hybrid-loop fragment.

Branching, joining, dynamic dispatch, recursion (a block calling a sub-loop), parallel fan-out / fan-in — all of these are higher-order composition operators on the same primitive set. They don't change the primitive vocabulary; they change the topology of how primitives connect.

The key disciplines that make composition not collapse:

1. **Schema versioning at every typed edge.** When the upstream block's output schema changes, every downstream block sees the change. Without versioning, composition is fragile.
2. **Calibration at every LLM block** (or at least every load-bearing one). Without it, you can't tell which block is dropping the ball when a chain misbehaves.
3. **Restraint where the LLM block isn't naturally restrained.** The deterministic block immediately downstream of an LLM block is where most projects' opinionated policy lives — filtering, thresholding, normalizing, ranking.
4. **Loop-closure visibility.** Which arrows feed back? When? This file says "the action loops back as new soft-data" but in a real system there might be three different feedback paths with different cadences. Diagram them.

The combinatorial space looks daunting until you realize you're never visiting all of it for any one project. A typical hybrid-loop app is 5-15 blocks arranged in 1-3 cycles. The skill in `../SKILL.md` exists because *most programmers were trained on three or four of the eight primitives, not the whole eight, and certainly not the compositional algebra that connects them.*
