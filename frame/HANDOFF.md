# Handoff — `frame`

Everything needed to pick this work up cleanly in a local session. `frame` is a
guarded-statechart DSL that compiles to a Claude Code hook mesh; it came out of a
design conversation about using hooks as a deterministic message bus and was
built here as a self-contained Go module.

---

## 1. Status

- **Branch:** `claude/hooks-message-bus-Z1Dw0` on `justinstimatze/hybrid`,
  pushed and in sync with origin.
- **Everything lives under `frame/`** — a self-contained Go module
  (`module github.com/justinstimatze/frame`) nested in the docs repo.
- **Tests green** (`go test ./...`): `spec`, `check`, `sim`.
- **Working tree clean.** Nothing uncommitted, nothing else on the branch.

Commits (newest last):

```
Add frame: a guarded statechart compiled to a Claude Code hook mesh
frame: make Cell a formal oracle (grammar + safety stages)
frame/docs: capture the formal-oracle model and the lexicon dual
(+ this handoff doc)
```

---

## 2. Layout

```
frame/
  README.md                       # mechanics, usage, invariants table
  HANDOFF.md                      # this file
  docs/DESIGN.md                  # the "why": formal oracle, totality, lexicon dual
  go.mod                          # module github.com/justinstimatze/frame (Go 1.24)
  spec/spec.go                    # the DSL: Trigger, Cell (formal oracle), Guard, Effect, Machine, Context
  spec/spec_test.go               # formal-oracle Check semantics; NewCell rejects incomplete oracles
  check/check.go                  # static checker — makes unsafe meshes inexpressible
  check/check_test.go             # one test per invariant (E-ORACLE, E-INJECT, E-BLOCK, E-HALT, ...)
  runtime/runtime.go              # the generic hook dispatcher: fail-open, total, atomic state
  compile/compile.go              # machine -> settings.json hooks fragment
  sim/sim.go                      # deterministic simulator (scripted oracle, no model call)
  sim/sim_test.go                 # converges / fail-safe / total-under-stuck-oracle
  examples/reviewloop/reviewloop.go  # worked Stop-loop machine + demo scenarios
  registry/registry.go            # name -> machine, for the CLI
  cmd/frame/main.go               # CLI: check | compile | sim | run
```

---

## 3. Pull it locally

```bash
cd hybrid                      # your local clone of justinstimatze/hybrid
git fetch origin
git checkout claude/hooks-message-bus-Z1Dw0
cd frame
go test ./...                  # expect ok: spec, check, sim
go vet ./...                   # clean
go run ./cmd/frame sim review-loop      # watch it step / branch / halt
go run ./cmd/frame check   review-loop  # static soundness report
go run ./cmd/frame compile review-loop /usr/local/bin/frame   # settings.json fragment
```

`frame sim review-loop` should print three scenarios: `converges`
(block→block→inject→done), `fail-safe (output outside the language)`, and
`total under a stuck oracle (fuel bound)`.

---

## 4. What it is (one paragraph)

You declare a **machine** — states, transitions, deterministic **guards**, and
fenced LLM **cells** (formal oracles: output confined to a formal language `L`,
checked for safe actions). `check` statically proves the machine is sound
(reachable halt, no orphans, fuel-bounded totality, no oracle on the control
path); `compile` emits the `settings.json` hooks that run it; one generic
`runtime` dispatcher interprets it, fail-open and total. The unsafe shapes are
made *inexpressible* (you can't build a cell without grammar+safety, can't gate a
transition on raw oracle output, can't inject on a trigger that can't inject).
Read `README.md` for mechanics, `docs/DESIGN.md` for the rationale.

---

## 5. Migrate `frame` into its own repo

`frame` was always meant to be standalone. Its module path is already
`github.com/justinstimatze/frame`, so imports and `go install` work unchanged
once it lives at that path.

**Preserving history (recommended):**

```bash
# from the hybrid repo, on this branch:
git subtree split --prefix=frame -b frame-only

git init ../frame && cd ../frame
git pull /absolute/path/to/hybrid frame-only
git remote add origin git@github.com:justinstimatze/frame.git   # after creating the repo
git push -u origin main
```

**Quick and dirty (no history):** `cp -r frame ../frame-new && cd ../frame-new &&
git init && git add . && git commit -m "import frame"`.

After extraction, drop this `HANDOFF.md` if you like — it documents the origin,
not the destination.

---

## 6. Resume point — the `frame/lexicon` adapter

This is the one real build left, and the highest-leverage one: it removes
`frame`'s biggest stub (per-cell hand-rolled grammars) by making `lexicon` the
formal language `L`.

`docs/DESIGN.md` §4–5 has the full reasoning. The binding:

```go
// frame/lexicon/adapter.go (to be written)
func LexiconCell(name, model, instr string, lex *lexicon.Library) spec.Cell {
    return spec.NewCell(name, model, instr,
        lex.ParseComposition, // Grammar: membership + type-check of a primitive composition
        lex.IsRecognized)     // Safety:  composes to a recognized pattern (not nonsense)
}

// and a surfacing-render wrapper so injected context is rendered into the
// recipient's working vocabulary rather than emitted as a raw term:
//   render(term, substrate, workingVocab) -> spec.Inject{Text: ...}
```

**Resolve first (against the real lexicon repo):** does `lexicon` expose a
*deterministic* membership/typing check usable as `Grammar`? The README's
"type-checked composition" in the dense internal layer suggests yes. If so,
per-cell grammars in `frame` are replaced wholesale by lexicon membership.

Locally you have auth to `lexicon` (the session-scope wall we hit is gone), so
clone it alongside and wire the adapter.

---

## 7. Known stubs / gaps (the TODO surface)

- **`runtime.Oracle` is unbound** — the actual model call. A real binding should
  also do *generation-time confinement* (constrained decoding / structured output
  / tool schema) so the model *can only* emit the cell's language `L`, making
  `Grammar` a backstop rather than the sole guarantee. Until wired, `frame run`
  uses a nil-returning stub, so cells fall outside `L` and machines fail safe.
- **`spec.Emit` records bus messages but does not write inbox files.** Drop in the
  mcp-dispatch relay here (atomic write to `{dispatch_dir}/{target}/`), so the
  statechart's `Emit` effect actually delivers onto the bus.
- **No calibration tap yet** — the hindcast-style predict/verdict log per cell.
  This is what makes the stack self-reflective (the system measuring its own
  oracle); it's also where the "audit the vocabulary" loop from lexicon lands.
- **`Safety` is trivial in the example** (a classification is always safe). The
  stage earns its keep when a cell's term denotes an *action*; see
  `spec/spec_test.go`'s read/write command-language oracle for the real shape.
- **`run` reads/writes per-session state under `~/.cache/frame/`** — fine for a
  single host; revisit if you ever want cross-host runs.

---

## 8. Design context (so the "why" isn't lost)

`docs/DESIGN.md` is the canonical record. In brief:

- **Formal oracle (Russell's sense):** an LLM whose output is confined to a
  formal language checked for safe actions. `Cell` = `Grammar` (membership) +
  `Safety` (safe action); a term reaches control flow only if `WellFormed &&
  Safe`.
- **Total guarded statechart, not a Turing-complete mesh:** a hook mesh is
  trivially TC (Stop-hook that refuses to halt = unbounded loop), but TC forfeits
  the guarantees (halting problem). `frame` bounds fuel to stay *total*, buying
  decidable reachability / termination / "no unvalidated oracle on the control
  path."
- **Two readings of the same architecture:** Russell (confinement / deference /
  corrigibility — the safety axis) and reflective oracles (self-reference — the
  bus + calibration + dev-time loop; `frame` externalizes the fixed point into
  the deterministic verifier).
- **The `lexicon` dual:** `lexicon` types the *content* (cognitive primitives,
  composition, lineage); `frame` types the *control* (when moves fire, fencing,
  termination). Composed, a centaur program typed on both axes. lexicon's
  surfacing function is also the injection-compliance primitive (render into the
  recipient's working vocabulary, the generalized slimemold "factual not
  imperative" finding).

---

## 9. Reference repos used (NOT on this branch)

To reason about the design, these were cloned into the build container's
`/home/user/_scope` and read, but **nothing from them is on the branch** — they
are your own repos, untouched, and they vanish with the container:

- `mcp-dispatch` — the bus (piggyback delivery + Stop-hook fallback). The `Emit`
  target.
- `slimemold` — the read-once contract vs data-only injection; the maintenance
  detectors that audit a substrate.
- `hindcast` — fail-open `Guard`, calibration substrate. The calibration-tap
  target.
- prior art: `disler/claude-code-hooks-multi-agent-observability`,
  `thedotmack/claude-mem`, `GowayLee/cchooks`, `syou6162/cchook`,
  `CorridorSecurity/hookshot`.

`frame`'s README maps each layer back to these.

---

## 10. Optional

No PR was opened (none was requested). If you'd rather review via PR than the raw
branch, open one against `main` from `claude/hooks-message-bus-Z1Dw0`.
