// Package spec is the statechart language that compiles to a Claude Code hook
// mesh. A Machine is a finite set of States; Transitions fire on hook Triggers,
// are selected by deterministic Guards, and run Effects.
//
// The one stochastic element — an LLM Cell — is fenced two ways:
//
//   - A Cell cannot be built without a validator (NewCell rejects a nil one, and
//     the validator field is unexported, so a Cell literal without one will not
//     compile in another package). An ungated oracle is therefore inexpressible.
//   - A Guard may only read a cell's *validated* output (see package check),
//     never its raw output, so an unvalidated oracle result can never gate
//     control flow.
//
// Nothing here calls a model or touches the filesystem; it is pure data plus
// deterministic predicates.
package spec

// Trigger is a hook event a transition can fire on (the reliable subset).
type Trigger string

const (
	SessionStart     Trigger = "SessionStart"
	UserPromptSubmit Trigger = "UserPromptSubmit"
	PreToolUse       Trigger = "PreToolUse"
	PostToolUse      Trigger = "PostToolUse"
	Stop             Trigger = "Stop"
	SubagentStop     Trigger = "SubagentStop"
)

// injecting: triggers whose hook output can carry text back into Claude's
// context. gating: triggers whose hook can halt/redirect control flow (exit 2).
// These mirror real Claude Code mechanics, so an Inject on PreToolUse — which
// cannot inject — is a compile-time check failure, not a runtime surprise.
var injecting = map[Trigger]bool{
	SessionStart: true, UserPromptSubmit: true, PostToolUse: true,
	Stop: true, SubagentStop: true,
}
var gating = map[Trigger]bool{PreToolUse: true, UserPromptSubmit: true, Stop: true}

func Injecting(t Trigger) bool { return injecting[t] }
func Gating(t Trigger) bool    { return gating[t] }

// A Cell is a formal oracle: an LLM whose output is confined to a formal
// language and checked for safe actions before it can influence anything.
// Its contract has two decidable stages, both mandatory:
//
//	Grammar  formal-language membership — parse raw output into a well-formed
//	         term, reporting whether it is in the language L at all.
//	Safety   over a well-formed term, decide whether the action it denotes is safe.
//
// A term reaches control flow only if it is WellFormed AND Safe. Generation-time
// confinement (constrained decoding / structured output so the model *can only*
// emit L) is the runtime's job; Grammar is the parse-time backstop.
type Grammar func(raw string) (term string, wellFormed bool)

// Safety decides whether a well-formed term denotes a safe action.
type Safety func(term string) bool

// Cell is a single formal-oracle call. Grammar and Safety are mandatory and
// unexported: a Cell whose output could reach a guard unchecked is inexpressible.
type Cell struct {
	Name         string
	Model        string
	Instructions string
	grammar      Grammar
	safe         Safety
}

// NewCell builds a formal-oracle Cell. It panics if either stage is missing:
// there is no valid Cell that emits unconstrained or unsafety-checked output.
func NewCell(name, model, instructions string, g Grammar, s Safety) Cell {
	if g == nil || s == nil {
		panic("frame/spec: a formal-oracle Cell requires both a grammar (formal language) and a safety check")
	}
	return Cell{Name: name, Model: model, Instructions: instructions, grammar: g, safe: s}
}

// Check runs the formal pipeline over a raw completion: membership, then safety.
func (c Cell) Check(raw string) CellResult {
	term, ok := c.grammar(raw)
	if !ok {
		return CellResult{Raw: raw, Ran: true} // not in the language
	}
	return CellResult{Raw: raw, Term: term, WellFormed: true, Safe: c.safe(term), Ran: true}
}

// Text is a string field that may depend on the live context. Use S for a
// constant.
type Text func(*Context) string

// S lifts a constant string into a Text.
func S(s string) Text { return func(*Context) string { return s } }

// Resolve evaluates a Text against a context (nil Text -> empty string).
func Resolve(t Text, c *Context) string {
	if t == nil {
		return ""
	}
	return t(c)
}

// Effect is the action side of a transition.
type Effect interface{ isEffect() }

// Inject surfaces text into Claude's context (additionalContext).
type Inject struct{ Text Text }

// Block halts the triggering action and feeds Reason back to Claude (exit 2).
// On a Stop trigger this is the "keep going" primitive — it refuses the stop.
type Block struct{ Reason Text }

// Run runs a cell for its side effect (result cached on the context).
type Run struct{ Cell Cell }

// Emit writes a message onto the bus (an inbox file) for another agent.
type Emit struct {
	Target  string
	Message Text
}

func (Inject) isEffect() {}
func (Block) isEffect()  {}
func (Run) isEffect()    {}
func (Emit) isEffect()   {}

// Guard is a deterministic transition predicate. Reads declares the context
// paths it depends on so the checker can prove it never gates on unvalidated
// oracle output. Paths look like "cells.<name>.validated" or "vars.<key>".
type Guard struct {
	Reads []string
	When  func(*Context) bool
}

type Transition struct {
	On    Trigger
	To    string
	Guard *Guard // nil == unconditional
	Do    []Effect
}

type State struct {
	Name     string
	On       []Transition
	Terminal bool
}

type Machine struct {
	Name     string
	Fuel     int    // step budget — the bound that makes the machine total
	Contract string // read-once legitimacy text (the slimemold plane)
	Initial  string
	States   []State
	Cells    []Cell // registry every guard.Reads must resolve against
}

// --- runtime context (carried across steps, persisted between hook fires) ---

// CellResult is the outcome of a formal-oracle Check. A term influences control
// flow only when WellFormed && Safe; otherwise the machine takes its fail-safe.
type CellResult struct {
	Raw        string // the model's raw completion
	Term       string // the parsed, well-formed term (empty if not in the language)
	WellFormed bool   // raw output is a member of the formal language
	Safe       bool   // the term denotes a safe action
	Ran        bool
}

type Context struct {
	Event map[string]any
	State string
	Fuel  int
	Vars  map[string]any
	Cells map[string]*CellResult
}

// Cell returns the cached result for a cell, or a zero result if it has not run.
func (c *Context) Cell(name string) *CellResult {
	if r, ok := c.Cells[name]; ok {
		return r
	}
	return &CellResult{}
}
