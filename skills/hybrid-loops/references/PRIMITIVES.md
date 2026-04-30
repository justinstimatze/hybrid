# Primitives for hybrid loops

Sketches of extractable primitives that recur across hybrid-loop projects. Most don't yet exist as packages. When scaffolding a new project, write a local version of the primitive *and* leave a TODO referencing this file. Once the same primitive is written twice across projects, lift it into a real package.

## cal_log — calibration logger (~50 lines)

Append-only JSONL of predictions and verdicts. The simplest and highest-leverage primitive.

API sketch:
```python
log = CalLog("calibration.jsonl")
log.predict(loop="my_project",
            input_hash=hash(text),
            prediction={"flagged_pattern": "X"},
            model_id="claude-sonnet-4-6",
            verdict_due_in_days=7)
# later:
log.resolve(prediction_id, verdict="confirmed", verdict_source="user_pushback")
log.hit_rate(loop="my_project", window_days=30)
```

Implementation: file append + a small reader for hit-rate. Don't over-engineer. Per-record fields: `ts`, `loop`, `lens_or_reasoner`, `input_hash`, `prediction`, `model_id`, `schema_version`, `verdict_due_by`, `verdict`, `verdict_source`, `verdict_ts`.

## metacog — cognitive-bias auditor

Run cognitive-bias signature checks against any typed substrate. The metrics check the *structure* of the substrate, not its content.

Auditors:
- **confirmation_bias** — corroboration rate among resolved cycles
- **anchoring** — Spearman(file_age, claim_density) (or analogous shape-vs-time correlation)
- **clustering_illusion** — file-grouping vs topology-cluster Jaccard
- **availability_heuristic** — provenance HHI (concentration index over source domains)
- **survivorship_bias** — irrelevant:challenged signal ratio (or analog over the substrate's verdict types)
- **framing_effect** — evaluative-adjective regex frequency in summary fields
- **dunning_kruger** — low-complexity entity centrality vs high-complexity centrality
- **base_rate_neglect** — predicate distribution entropy
- **premature_closure** — thought-terminating cliché detection

API sketch:
```python
audit = Metacog(substrate=load_jsonl("substrate.jsonl"))
results = audit.run_all()
for finding in results.triggered():
    print(finding.bias_name, finding.value, finding.threshold)
```

Generalize the input shape (records with provenance + edges) so it works on any typed substrate. **Per Conjecture 2 in the README, this is the most extractable substrate-level primitive but its substrate-generality is untested.**

## schemaforge — DreamCoder-style schema discovery

Given a corpus and an initial draft schema, refine via compress+verify loop.

API sketch:
```python
forge = Schemaforge(corpus, draft_schema)
result = forge.discover(rounds=4, model="sonnet")
print(result.schema)        # refined schema
print(result.history)       # per-round metrics
```

Implementation: a wake-phase compressor (LLM encodes spec into notation), an adversarial verifier (separate LLM scores fidelity), and a refinement step that proposes notation changes based on compressor failures. Output should include a grammar-size term to guard against bloat. Cite DreamCoder (Ellis et al., 2021) and LILO (Grand et al., 2024) as direct lineage.

## metabolism — phase scheduler

Periodic execution of dream/trip/audit/evolve phases with phase-gate dependencies on a typed substrate.

API sketch:
```python
metabolism = Metabolism(
    substrate=...,
    phases=[
        Phase("dream", run_dream, gate=lambda ctx: ctx.last_dream_age > 1d),
        Phase("trip",  run_trip,  gate=lambda ctx: ctx.bias_audit.passed),
        Phase("audit", run_audit, gate=lambda _: True),
    ]
)
metabolism.run_cycle()
```

Generalize so phases are pluggable. Phase outputs themselves get logged via `cal_log` so phase predictions can be calibrated alongside lens/reasoner predictions.

## mcp_substrate — MCP server template

Wrap a typed substrate as an MCP server with read tools (search, get, stats, query). One template generator that takes a substrate schema + sample queries and emits a working MCP server. Most hybrid-loop projects that grow past v0 want this; once a developer has built two of them by hand, packaging is worth it.

## When to package vs. inline

- v0 prototype: inline a simplified version of `cal_log` directly. Skip the others.
- v1 with > 1 month of expected life: lift `cal_log` into the project as a small file (reuse code, not the package abstraction).
- Once the same primitive is written for the third project: package it. Don't package speculatively; let the third use case define the API.

The bottleneck for most builders isn't "missing packages." It's "Claude doesn't reach for the pattern by default." This skill addresses that. Packaging is downstream.
