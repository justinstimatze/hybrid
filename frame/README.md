# frame

A guarded statechart that compiles to a Claude Code **hook mesh**.

You declare a machine — states, transitions, deterministic guards, and fenced
LLM cells. `frame` statically proves it is sound (reachable halt, no orphans,
bounded, no oracle on the control path) and emits the `settings.json` hooks that
run it. One generic dispatcher interprets the machine at runtime; all dynamism
flows through per-session state, so hook *definitions* are written once and never
need a `/hooks` refresh.

This is a **skeleton** — the architecture, the checker, the runtime, and a
deterministic simulator, with the one model call left unbound. It is meant to be
lifted into its own repo.

## Why this shape

A hook mesh is trivially Turing complete (arbitrary programs as nodes + a `Stop`
hook that refuses to halt = an unbounded loop). But Turing-completeness brings
the halting problem, so an LLM-in-the-loop mesh that is TC **cannot** be formally
guaranteed to behave. `frame` deliberately targets the class one rung below: a
**total** (always-terminating) guarded statechart, where

- the **control-flow skeleton is deterministic and analyzable**, and
- the **LLM is a fenced oracle** whose output can never move the machine into an
  unverified state.

A `Cell` is a **formal oracle** (Russell's sense): an LLM whose output is
confined to a formal language and checked for safe actions before it can
influence anything. Its contract is two decidable stages, both mandatory:

- **Grammar** — formal-language membership: parse raw output into a well-formed
  term, or reject it as outside the language `L`.
- **Safety** — over a well-formed term, decide whether the action it denotes is
  safe.

A term reaches control flow only if it is `WellFormed && Safe`; otherwise the
machine takes its fail-safe path. (Generation-time confinement — constrained
decoding so the model *can only* emit `L` — is the runtime's job; the Grammar is
the parse-time backstop. See *What's stubbed*.)

The load-bearing invariant: *an LLM output may never gate a transition without a
deterministic check between it and the transition.* `frame` makes the unsafe
shapes **inexpressible**:

| Unsafe shape | How it is made impossible |
|---|---|
| An oracle that emits unconstrained or unchecked output | `spec.NewCell` requires both a grammar and a safety check; the fields are unexported, so a `Cell` literal without them won't compile elsewhere |
| A guard branching on raw oracle output | `check` rejects any guard reading `cells.X.raw` (only `.term`/`.wellformed`/`.safe`) — code `E-ORACLE` |
| Injecting on a trigger that can't inject | `E-INJECT` (e.g. `Inject` on `PreToolUse`) |
| Blocking on a trigger that can't block | `E-BLOCK` (e.g. `Block` on `PostToolUse`) |
| A loop that can never end | `E-HALT` (no reachable terminal) + a mandatory `Fuel` bound (totality) |
| A hook crash bricking the session | `runtime.SafeDispatch` recovers and fails open (exit 0, no output) |

## Mapping to prior art

`frame` is the type system over primitives that already exist separately:

| Layer | Prior art it abstracts |
|---|---|
| typed guard/effect split | hookshot (`Pre*`=decision, `Post*`=fire-and-forget), Windsurf Cascade |
| the bus / inbox transport (`Emit`) | mcp-dispatch (piggyback delivery + Stop-hook fallback) |
| the read-once legitimacy `Contract` | slimemold (behavioral contract vs data-only injection) |
| fail-open dispatcher, calibration tap | hindcast (`Guard`, bounded logging) |
| the `Stop`-loop clock + `Fuel` bound | Ralph Wiggum loop, made total |
| whole-mesh observation | disler multi-agent-observability (star tap) |

## Layout

```
spec/      the statechart language (pure data + deterministic predicates)
check/     the static checker — makes unsafe meshes inexpressible
runtime/   the generic hook dispatcher (fail-open, total)
compile/   machine -> settings.json hooks fragment
sim/       deterministic simulator (scripted oracle, no model call)
registry/  name -> machine, for the CLI
examples/  reviewloop: a worked Stop-loop
cmd/frame/ the CLI
```

## Use

```bash
go build ./...

frame check   review-loop                 # static soundness report
frame compile review-loop /path/to/frame   # settings.json hooks fragment
frame sim     review-loop                  # run the demo scenarios
frame run --machine review-loop            # hook dispatcher (hook event JSON on stdin)
```

`frame sim review-loop` shows the worked example converging, failing safe on
invalid oracle output, and staying total under a stuck oracle:

```
scenario "converges"
  trigger   from   to     fuel  kind    detail
  Stop      loop   loop   4->3  block   exit2: Task not yet complete — continue ...
  Stop      loop   loop   3->2  block   exit2: Task not yet complete — continue ...
  Stop      loop   done   2->1  inject  ctx: Verified: the task is complete.
```

## Defining a machine

```go
assess := spec.NewCell("assess", "claude-sonnet-4-6",
    "Output exactly 'complete' or 'incomplete'.",
    func(raw string) (string, bool) {        // Grammar: membership in L
        v := strings.ToLower(strings.TrimSpace(raw))
        return v, v == "complete" || v == "incomplete"
    },
    func(string) bool { return true })       // Safety: a classification is safe

loop := spec.State{Name: "loop", On: []spec.Transition{
    {On: spec.Stop, To: "done",
        Guard: &spec.Guard{Reads: []string{"cells.assess.term"},   // formal output only
            When: func(c *spec.Context) bool { return c.Cell("assess").WellFormed && c.Cell("assess").Term == "complete" }},
        Do: []spec.Effect{spec.Inject{Text: spec.S("Verified complete.")}}},
    {On: spec.Stop, To: "loop",
        Guard: &spec.Guard{Reads: []string{"cells.assess.term"},
            When: func(c *spec.Context) bool { return c.Cell("assess").WellFormed && c.Cell("assess").Term == "incomplete" }},
        Do: []spec.Effect{spec.Block{Reason: spec.S("Not done — continue.")}}},
}}

m := spec.Machine{Name: "review-loop", Fuel: 4, Initial: "loop",
    Contract: "You installed review-loop ...",
    States: []spec.State{loop, {Name: "done", Terminal: true}},
    Cells:  []spec.Cell{assess}}
```

Register it in `registry`, then `check` / `compile` / `sim` / `run` it.

## What's stubbed

- `runtime.Oracle` — the model call, and with it **generation-time confinement**:
  a real binding should constrain decoding (grammar / structured output / tool
  schema) so the model *can only* emit the cell's language `L`, making the
  Grammar stage a backstop rather than the sole guarantee. Until wired, `frame
  run` uses a nil-returning stub, so cells fall outside `L` and machines take
  their fail-safe path (rather than guessing).
- `Emit` records bus messages but does not yet write inbox files — drop in the
  mcp-dispatch relay here.
- No calibration tap yet — the hindcast-style predict/verdict log per cell.

## Properties the checker buys you

Because the skeleton is deterministic, these are decidable on the spec, before a
single turn runs: **reachability**, **a reachable halt**, **termination** (the
fuel bound), and the safety property *no transition is gated on unvalidated
oracle output*. What is **not** guaranteed — and cannot be — is the correctness
of an individual cell's output. That is the irreducible stochastic part; the
frame's job is to contain it, not eliminate it.
