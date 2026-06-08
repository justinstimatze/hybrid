# frame — design notes: the formal oracle and the lexicon dual

Background notes for `frame`. The README covers mechanics; this captures *why*
the shape is what it is, and how `frame` is meant to compose with `lexicon`.

## 1. The formal oracle

A **formal oracle** (Stuart Russell's sense): an LLM whose output is confined to
a formal language `L` that a decision procedure can check for safe actions. The
constraint is on the *output channel*, not the model — the oracle can only ever
emit in `L`, and an action is taken only if the checker certifies it.

`frame`'s `Cell` is exactly this, with the contract split into two decidable
stages (both mandatory; an incomplete oracle is inexpressible):

- **Grammar** — formal-language membership: `raw → (term, wellFormed)`. Is the
  output in `L` at all?
- **Safety** — over a well-formed term, is the denoted action safe?

A term reaches control flow only if `WellFormed && Safe`; otherwise the machine
takes its fail-safe path. Generation-time confinement (constrained decoding /
structured output, so the model *can only* emit `L`) is the runtime's job; the
Grammar is the parse-time backstop.

Why two stages, not one: membership and safety are independent. `read foo` and
`write /etc/passwd` can both be well-formed yet differ in safety; `rm -rf /` is
outside the language entirely. Collapsing them into one opaque validator (the
original design) hides that an action was both *parsed* and *judged safe*.

## 2. Why a total guarded statechart, not a Turing-complete mesh

A hook mesh is trivially Turing complete (arbitrary programs as nodes + a `Stop`
hook that refuses to halt = an unbounded loop), but TC brings the halting problem,
so an LLM-in-the-loop mesh that is TC cannot be formally guaranteed to behave
(Rice's theorem). `frame` targets the class one rung down: a **total**
(always-terminating) guarded statechart, where

- the control-flow skeleton is deterministic and analyzable, and
- the LLM is a fenced oracle that can never move the machine into an unverified
  state.

You deliberately give up Turing-completeness (bound the fuel) to *gain* the
guarantees: reachability, a reachable halt, termination, and the safety property
"no transition is gated on unvalidated oracle output" are all decidable on the
spec, before a turn runs. What is **not** guaranteed — and cannot be — is the
correctness of an individual cell's output. That is the irreducible stochastic
part; the frame contains it, it does not eliminate it.

## 3. Two control-theory readings of the same architecture

- **Russell — Oracle AI + assistance games (safety axis).** The cell is a
  *confined* oracle: it answers, it never acts; only deterministic code acts. The
  validator is the *deference / uncertain-about-the-objective* mechanism — the
  system never treats oracle output as authoritative, and on a failed check it
  defers (fails safe). The fuel bound is the *off-switch / corrigibility*
  guarantee: the loop always halts and yields control.
- **Reflective oracles (self-reference axis).** The bus (an agent's output
  becomes its own or a sibling's input), the calibration loop (the system
  measures its own past predictions), and the dev-time loop (a critic rewrites
  the runtime from transcripts of the runtime) all put the oracle in a position
  to reason about a system that contains it. Reflective oracles are the formal
  fix for that self-reference — but they are uncomputable idealizations whose
  consistency is *internal*. `frame` is the computable approximation where the
  fixed point is *externalized* into the deterministic verifier: the gate does,
  from outside and by force, the consistency work the idealized oracle does from
  inside by definition. The fuel bound tames unbounded self-reference by bounding
  the recursion rather than via a probabilistic fixed point.

Synthesis: `frame` is a *confined* oracle (Russell) embedded in a
*self-referential* mesh (reflective oracles), made computable by externalizing
both the fixed point and the deference into deterministic code.

## 4. The lexicon dual

`lexicon` is a typed cross-domain catalog of low-level cognitive primitives —
each a bounded move with a typed `input → output` signature and a lineage
citation — plus a context-aware surfacing function that renders matched
primitives into the user's working vocabulary. It has a dense internal layer
(formal type signatures, type-checked composition) and a surface layer (the
render function).

`lexicon` and `frame` are **duals**:

- `lexicon` types the **content** — what cognitive moves exist, how they compose
  into recognized patterns, with lineage.
- `frame` types the **control** — when a move fires, that the oracle is fenced,
  that the loop is safe and terminates.

Composed, they form a centaur program typed on both axes at once: its content is
a composition of `lexicon` primitives, its control flow is a `frame` statechart.

### The four joints

1. **`lexicon`'s internal layer is `frame`'s formal language `L`.** Its
   type-checked composition is a decidable membership/typing check, so it drops
   into `Grammar`: a cell's output stops being free text filtered by a closure
   and becomes a type-checked composition of curated primitives. This removes
   `frame`'s largest stub.

2. **Composition unifies the two checkers.** `lexicon` type-checks that
   primitives compose into recognized patterns; `frame.check` proves a
   composition of cells/transitions is sound. Same operation, two layers.
   `lexicon`'s composition-recovery rate (recovering known mental models from
   primitive compositions) is the content-side twin of `frame`'s soundness check
   — one could run `frame`'s checker over a lexicon-typed machine to verify a
   claimed pattern is a *sound* composition, not merely a well-typed one.

3. **The surfacing function is the injection-compliance primitive.** The thread
   that motivated `frame` was "make Claude follow instructions," and slimemold's
   finding was that imperative injection is rejected while factual-in-the-right-
   register is absorbed. `lexicon`'s surfacing function —
   `(matched_primitives, substrate, working_vocabulary) → output in register` —
   is the general mechanism. In `frame` it sits between a cell's `Term` and the
   `Inject` effect: never inject a raw term; render it into the recipient's
   working vocabulary first.

4. **The stack is self-consistently reflective.** `lexicon` is maintained as a
   `winze` project with slimemold's detectors auditing the vocabulary itself, so
   the formal language the oracle speaks is a substrate continuously audited by a
   hybrid loop. The oracle is confined to a language (Russell), the language is a
   fixed-point-audited substrate (reflective), and the audit is deterministic
   detectors (the fixed point externalized). Every layer fences the one below it
   with something deterministic.

Net: one object plays four roles — content language, bus protocol, control
vocabulary, and cognitive-primitive lineage.

## 5. Planned: the `frame/lexicon` adapter (future work)

Not built yet (lexicon was out of session scope when these notes were written).
When lexicon is readable, a thin adapter should bind the two:

- `lexicon.ParseComposition → spec.Grammar` — membership + type-check of a
  primitive composition as the cell's formal language.
- `lexicon.IsRecognized → spec.Safety` — the composition resolves to a recognized
  pattern (not nonsense) as the safety stage.
- a surfacing-render wrapper around `spec.Inject`, so injected context is rendered
  into the recipient's working vocabulary rather than emitted as a raw term.

Open question to resolve against the real repo: does lexicon expose a
deterministic membership/typing check suitable for `Grammar` (the README's
"type-checked composition" suggests yes). If so, `frame`'s per-cell hand-rolled
grammars are replaced wholesale by lexicon membership.

## Prior art these notes lean on

- mcp-dispatch (the bus / piggyback + Stop-hook fallback), slimemold (the
  read-once contract vs data-only injection; the maintenance detectors),
  hindcast (fail-open `Guard`, calibration), Ralph Wiggum (the Stop-loop clock),
  hookshot / Windsurf Cascade (typed guard/effect split), disler
  multi-agent-observability (the star-topology event tap).
- Russell, *Human Compatible* (Oracle AI, assistance games, corrigibility).
- Reflective oracles (Fallenstein, Taylor, Christiano; Critch-adjacent at CHAI).
