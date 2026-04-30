# Primitives for hybrid loops

Most of these don't yet exist as packages. When scaffolding a new project, write a local version of the primitive *and* leave a TODO referencing this file. Once the same primitive is written twice across projects, lift it into a real package.

## cal_log — calibration logger (~50 lines)

Append-only JSONL of predictions and verdicts. The simplest and highest-leverage primitive.

API sketch:
```python
log = CalLog("calibration.jsonl")
log.predict(loop="slimemold", input_hash=hash(text),
            prediction={"load_bearing_claim": "X"},
            model_id="claude-sonnet-4-6",
            verdict_due_in_days=7)
# later:
log.resolve(prediction_id, verdict="confirmed", verdict_source="user_pushback")
log.hit_rate(loop="slimemold", window_days=30)
```

Implementation: file append + pandas read for hit-rate. Don't over-engineer. Canonical reference: winze's `.metabolism-calibration.jsonl`.

## metacog — cognitive-bias auditor

Run nine cognitive-bias signature checks against any typed substrate.

Auditors (winze-derived):
- confirmation_bias: corroboration rate among resolved cycles
- anchoring: spearman(file_age, claim_density)
- clustering_illusion: file_grouping vs topology_cluster jaccard
- availability_heuristic: provenance HHI
- survivorship_bias: irrelevant:challenged ratio
- framing_effect: evaluative-adjective regex frequency
- dunning_kruger: low-complexity entity centrality vs high-complexity
- base_rate_neglect: predicate distribution entropy
- premature_closure: thought-terminating cliche detection

API sketch:
```python
audit = Metacog(substrate=load_jsonl("substrate.jsonl"))
results = audit.run_all()
for finding in results.triggered():
    print(finding.bias_name, finding.value, finding.threshold)
```

Source the implementations from winze's `cmd/metabolism/dreamaudit.go`. Generalize the input shape (records with provenance + edges) so it works on any typed substrate. **This is the most extractable winze-original primitive.**

## schemaforge — DreamCoder-style schema discovery

Given a corpus + an initial draft schema, refine via compress+verify loop.

API sketch:
```python
forge = Schemaforge(corpus, draft_schema)
result = forge.discover(rounds=4, model="sonnet")
print(result.schema)        # refined schema
print(result.history)       # per-round metrics
```

Implementation: lift from lamina/poc/dense's `discover.py`. Cite DreamCoder and LILO. Output should include grammar-size term to guard against bloat.

## metabolism — phase scheduler

Periodic execution of dream/trip/audit/evolve phases with phase-gate dependencies.

API sketch:
```python
metabolism = Metabolism(
    substrate=...,
    phases=[
        Phase("dream", run_dream, gate=lambda ctx: ctx.last_dream_age > 1d),
        Phase("trip", run_trip, gate=lambda ctx: ctx.bias_audit.passed),
        Phase("audit", run_audit, gate=lambda _: True),
    ]
)
metabolism.run_cycle()
```

Implementation: lift from winze `cmd/metabolism/main.go`. Generalize so phases are pluggable.

## mcp_substrate — MCP server template

Wrap a typed substrate as an MCP server with read tools (search, get, stats, query).

API sketch: see winze `cmd/mcp/main.go` for the canonical shape. The user has done this multiple times; once written for the third project, lift to a template generator.

## When to package vs. inline

For a new hybrid-loop project:

- If it's a v0 prototype: inline a simplified version of cal_log directly. Skip the others.
- If it's a v1 with > 1 month of expected life: lift cal_log into the project as a small file (reuse code, not the package abstraction).
- Once the same primitive is written for the third project: package it. Don't package speculatively; let the third use case define the API.

The user's bottleneck is not "missing packages." It's "Claude doesn't reach for the pattern by default." This skill addresses that. Packaging is downstream.
