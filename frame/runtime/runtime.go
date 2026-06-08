// Package runtime is the generic hook dispatcher — the single command every
// registered hook runs. One invocation, per hook event, that:
//
//  1. loads the per-session context (state + fuel + cell cache),
//  2. selects the first transition whose trigger matches and whose guard holds,
//     running any cells the guard depends on (lazily, then validating them),
//  3. applies the transition's effects, decrements fuel, advances state, persists,
//  4. emits the hook protocol (additionalContext to inject, exit 2 to block).
//
// Two safety properties are enforced here, not hoped for:
//
//   - fail-open: SafeDispatch recovers from any panic and yields a no-op, so a
//     broken hook can never brick the session (hindcast's Guard discipline).
//   - totality: once fuel reaches 0 the loop can no longer Block, so a Stop-loop
//     terminates regardless of what the oracle says.
//
// The model call lives in Oracle and is intentionally unbound: wire it to your
// provider. Everything else is deterministic and is exercised by package sim
// with a scripted oracle.
package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/justinstimatze/frame/check"
	"github.com/justinstimatze/frame/spec"
)

// Oracle runs a cell and returns the raw model completion.
type Oracle func(spec.Cell, *spec.Context) string

// Output is the result of one dispatch step, before it is rendered to the hook
// protocol.
type Output struct {
	Inject     []string
	Block      *string
	Emits      []Emission
	BudgetHalt bool
}

type Emission struct{ Target, Message string }

func (o Output) Kind() string {
	switch {
	case o.Block != nil:
		return "block"
	case o.BudgetHalt:
		return "budget"
	case len(o.Inject) > 0 || len(o.Emits) > 0:
		return "inject"
	default:
		return "noop"
	}
}

func runCell(c spec.Cell, ctx *spec.Context, oracle Oracle) *spec.CellResult {
	if r, ok := ctx.Cells[c.Name]; ok && r.Ran {
		return r
	}
	raw := oracle(c, ctx)
	res := c.Check(raw) // formal-language membership, then safety
	r := &res
	if ctx.Cells == nil {
		ctx.Cells = map[string]*spec.CellResult{}
	}
	ctx.Cells[c.Name] = r
	return r
}

// Dispatch advances the machine one step for event. It mutates ctx and is pure
// given oracle.
func Dispatch(m spec.Machine, event map[string]any, ctx *spec.Context, oracle Oracle) Output {
	byName := map[string]spec.State{}
	for _, s := range m.States {
		byName[s.Name] = s
	}
	cellMap := map[string]spec.Cell{}
	for _, c := range m.Cells {
		cellMap[c.Name] = c
	}

	// Cell results are ephemeral per hook event: the transcript changes every
	// turn, so a cell must re-run each event. Within a single dispatch the
	// result is cached (one model call even if several guards read it).
	ctx.Cells = map[string]*spec.CellResult{}

	state := byName[ctx.State]
	if state.Terminal {
		return Output{}
	}
	trigger := spec.Trigger(asString(event["hook_event_name"]))

	// Totality: out of fuel -> never block again; release the loop.
	if ctx.Fuel <= 0 {
		out := Output{BudgetHalt: true}
		if spec.Gating(trigger) {
			out.Inject = append(out.Inject, fmt.Sprintf("[%s] step budget (%d) exhausted; halting.", m.Name, m.Fuel))
		}
		return out
	}

	for _, t := range state.On {
		if t.On != trigger {
			continue
		}
		if t.Guard != nil {
			for _, cr := range check.CellReads(t.Guard.Reads) {
				if c, ok := cellMap[cr.Cell]; ok {
					runCell(c, ctx, oracle)
				}
			}
			if !t.Guard.When(ctx) {
				continue
			}
		}
		var out Output
		for _, eff := range t.Do {
			switch e := eff.(type) {
			case spec.Inject:
				out.Inject = append(out.Inject, spec.Resolve(e.Text, ctx))
			case spec.Block:
				r := spec.Resolve(e.Reason, ctx)
				out.Block = &r
			case spec.Run:
				runCell(e.Cell, ctx, oracle)
			case spec.Emit:
				out.Emits = append(out.Emits, Emission{Target: e.Target, Message: spec.Resolve(e.Message, ctx)})
			}
		}
		ctx.Fuel--
		ctx.State = t.To
		return out
	}
	return Output{} // no matching transition: a no-op for this event
}

// SafeDispatch wraps Dispatch with fail-open recovery.
func SafeDispatch(m spec.Machine, event map[string]any, ctx *spec.Context, oracle Oracle) (out Output) {
	defer func() {
		if r := recover(); r != nil {
			out = Output{}
		}
	}()
	return Dispatch(m, event, ctx, oracle)
}

// --- persistence + hook protocol --------------------------------------------

// persisted is the cross-turn state. Cell results are deliberately absent: they
// are ephemeral per event (see Dispatch). Only state, fuel, and machine vars
// survive between hook fires.
type persisted struct {
	State string         `json:"state"`
	Fuel  int            `json:"fuel"`
	Vars  map[string]any `json:"vars"`
}

func statePath(m spec.Machine, session string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".cache", "frame", m.Name)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, session+".json"), nil
}

func LoadContext(m spec.Machine, event map[string]any) *spec.Context {
	session := asString(event["session_id"])
	if session == "" {
		session = "default"
	}
	ctx := &spec.Context{Event: event, State: m.Initial, Fuel: m.Fuel,
		Vars: map[string]any{}, Cells: map[string]*spec.CellResult{}}
	p, err := statePath(m, session)
	if err != nil {
		return ctx
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return ctx
	}
	var ps persisted
	if json.Unmarshal(data, &ps) != nil {
		return ctx
	}
	ctx.State, ctx.Fuel = ps.State, ps.Fuel
	if ps.Vars != nil {
		ctx.Vars = ps.Vars
	}
	return ctx
}

func SaveContext(m spec.Machine, ctx *spec.Context) error {
	session := asString(ctx.Event["session_id"])
	if session == "" {
		session = "default"
	}
	p, err := statePath(m, session)
	if err != nil {
		return err
	}
	data, err := json.Marshal(persisted{State: ctx.State, Fuel: ctx.Fuel, Vars: ctx.Vars})
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p) // atomic
}

// Emit writes the Claude Code hook protocol for out and returns the exit code.
func Emit(trigger spec.Trigger, out Output) int {
	if out.Block != nil && spec.Gating(trigger) {
		fmt.Fprintln(os.Stderr, *out.Block)
		return 2
	}
	if len(out.Inject) > 0 {
		payload := map[string]any{"hookSpecificOutput": map[string]any{
			"hookEventName":     string(trigger),
			"additionalContext": strings.Join(out.Inject, "\n"),
		}}
		b, _ := json.Marshal(payload)
		fmt.Println(string(b))
	}
	return 0
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
