# frame — context for Claude

`frame` is a guarded-statechart DSL that compiles to a Claude Code hook mesh.
You declare a machine (states, transitions, deterministic guards, fenced LLM
"cells"); a static checker proves it sound; a compiler emits the `settings.json`
hooks; one generic runtime dispatcher interprets it, fail-open and total.

## Start here

1. `HANDOFF.md` — current state, how to pull/verify, how to migrate this into its
   own repo, the resume point, and known stubs. Read it first.
2. `README.md` — mechanics, usage, the invariants table.
3. `docs/DESIGN.md` — the rationale (formal oracle, why total-not-Turing-complete,
   Russell / reflective-oracle readings, the `lexicon` dual).

## State

- Self-contained Go module: `github.com/justinstimatze/frame`, Go 1.24, **stdlib
  only** (no third-party deps — keep it that way unless there's a strong reason).
- `go test ./...` and `go vet ./...` are green; keep them green.
- It was built on branch `claude/hooks-message-bus-Z1Dw0` of a docs repo and is
  meant to be lifted into its own project (see `HANDOFF.md` §5). Sorting out
  branches/commits and the extraction is expected — do it however is cleanest.

## Likely task

The one real build left is the **`frame/lexicon` adapter** (`HANDOFF.md` §6,
`docs/DESIGN.md` §5): bind `lexicon.ParseComposition → spec.Grammar` and
`lexicon.IsRecognized → spec.Safety`, plus a surfacing-render wrapper around
`spec.Inject`. First answer: does `lexicon` expose a deterministic
membership/typing check usable as `Grammar`?

## Invariants — do not weaken these

The whole point is that unsafe meshes are *inexpressible*. When changing code:

- A `Cell` is a **formal oracle**: it must keep both stages (`Grammar` =
  language membership, `Safety` = safe action). `NewCell` rejects either being
  nil. Never let an LLM output reach a guard unchecked.
- **Guards may only read `cells.X.{term,wellformed,safe}`, never `.raw`.** The
  checker enforces this (`E-ORACLE`); don't add an escape hatch.
- **Keep machines total:** every machine needs a positive `Fuel` bound and a
  reachable terminal. The runtime must stop blocking once fuel hits 0.
- **Hooks fail open:** any error/panic in the dispatcher yields exit 0 with no
  output (`runtime.SafeDispatch`). A broken hook must never brick a session.
- `Inject` only on triggers that can inject; `Block` only on triggers that can
  block (the checker enforces `E-INJECT` / `E-BLOCK`).

## Verify

```bash
go test ./... && go vet ./...
go run ./cmd/frame sim review-loop      # converges / fail-safe / total
```
