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

// Validator maps a raw model completion to a validated value, reporting whether
// it was acceptable. ok==false means the oracle output is rejected.
type Validator func(raw string) (validated string, ok bool)

// Cell is a single LLM call. The validator is mandatory and unexported.
type Cell struct {
	Name         string
	Model        string
	Instructions string
	validate     Validator
}

// NewCell builds a Cell. It panics on a nil validator: there is no valid Cell
// whose output could reach a guard unchecked.
func NewCell(name, model, instructions string, v Validator) Cell {
	if v == nil {
		panic("frame/spec: Cell requires a validator (an ungated oracle is inexpressible)")
	}
	return Cell{Name: name, Model: model, Instructions: instructions, validate: v}
}

// Validate runs the cell's validator over a raw completion.
func (c Cell) Validate(raw string) (string, bool) { return c.validate(raw) }

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

type CellResult struct {
	Raw       string `json:"raw"`
	Validated string `json:"validated"`
	Valid     bool   `json:"valid"`
	Ran       bool   `json:"ran"`
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
