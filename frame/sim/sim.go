// Package sim runs a machine deterministically against a scripted oracle, so the
// whole control-flow skeleton can be proven out without a model call. It is the
// empirical complement to package check: check proves the shape is sound; sim
// shows a sound machine actually steps, branches, halts, and stays total.
package sim

import (
	"fmt"
	"strings"

	"github.com/justinstimatze/frame/runtime"
	"github.com/justinstimatze/frame/spec"
)

// Scenario scripts a run: a list of hook events and, per cell, the sequence of
// raw completions the oracle should return (the last value repeats if events
// outrun the script).
type Scenario struct {
	Name   string
	Events []map[string]any
	Oracle map[string][]string
}

// Ev is a convenience constructor for a hook event of the given trigger.
func Ev(t spec.Trigger) map[string]any {
	return map[string]any{"hook_event_name": string(t), "session_id": "sim"}
}

// Step records one dispatched event for the trace.
type Step struct {
	Trigger    string
	From, To   string
	FuelBefore int
	FuelAfter  int
	Kind       string
	Detail     string
}

// Run executes the scenario and returns the trace plus the final context.
func Run(m spec.Machine, sc Scenario) ([]Step, *spec.Context) {
	counters := map[string]int{}
	oracle := func(c spec.Cell, _ *spec.Context) string {
		seq := sc.Oracle[c.Name]
		i := counters[c.Name]
		counters[c.Name]++
		if len(seq) == 0 {
			return ""
		}
		if i >= len(seq) {
			i = len(seq) - 1
		}
		return seq[i]
	}

	ctx := &spec.Context{State: m.Initial, Fuel: m.Fuel,
		Vars: map[string]any{}, Cells: map[string]*spec.CellResult{}}
	var steps []Step
	for _, ev := range sc.Events {
		ctx.Event = ev
		from, fuelBefore := ctx.State, ctx.Fuel
		out := runtime.SafeDispatch(m, ev, ctx, oracle)
		steps = append(steps, Step{
			Trigger:    asString(ev["hook_event_name"]),
			From:       from,
			To:         ctx.State,
			FuelBefore: fuelBefore,
			FuelAfter:  ctx.Fuel,
			Kind:       out.Kind(),
			Detail:     detail(out),
		})
	}
	return steps, ctx
}

// Render formats a trace as a table.
func Render(name string, steps []Step) string {
	var b strings.Builder
	fmt.Fprintf(&b, "scenario %q\n", name)
	fmt.Fprintf(&b, "  %-12s %-8s %-8s %-5s %-7s %s\n", "trigger", "from", "to", "fuel", "kind", "detail")
	for _, s := range steps {
		fmt.Fprintf(&b, "  %-12s %-8s %-8s %d->%-2d %-7s %s\n",
			s.Trigger, s.From, s.To, s.FuelBefore, s.FuelAfter, s.Kind, s.Detail)
	}
	return b.String()
}

func detail(o runtime.Output) string {
	switch o.Kind() {
	case "block":
		return "exit2: " + *o.Block
	case "inject":
		return "ctx: " + strings.Join(o.Inject, " / ")
	case "budget":
		if len(o.Inject) > 0 {
			return strings.Join(o.Inject, " / ")
		}
		return "released (out of fuel)"
	default:
		return "-"
	}
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
