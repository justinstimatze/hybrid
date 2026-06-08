// Package reviewloop is a worked machine: a Stop-loop that asks a fenced oracle
// whether the user's task is complete and, until it is, refuses the stop and
// tells Claude to continue. It exercises every load-bearing feature — a lazily
// run cell, validator-gated branching, Block-on-Stop as the "keep going"
// primitive, a fail-safe when the oracle output is invalid, and the fuel bound
// that guarantees the loop ends.
package reviewloop

import (
	"strings"

	"github.com/justinstimatze/frame/sim"
	"github.com/justinstimatze/frame/spec"
)

// assess is the one stochastic edge. Its validator only accepts "complete" or
// "incomplete"; anything else is rejected (Valid == false) and the machine
// releases rather than trapping the user.
var assess = spec.NewCell(
	"assess",
	"claude-sonnet-4-6",
	"Read the transcript. Output exactly one word — 'complete' if the user's "+
		"stated task is fully done, otherwise 'incomplete'.",
	func(raw string) (string, bool) {
		v := strings.ToLower(strings.TrimSpace(raw))
		if v == "complete" || v == "incomplete" {
			return v, true
		}
		return "", false
	},
)

// Machine builds the review-loop statechart.
func Machine() spec.Machine {
	loop := spec.State{
		Name: "loop",
		On: []spec.Transition{
			{ // verified complete -> halt
				On: spec.Stop, To: "done",
				Guard: &spec.Guard{
					Reads: []string{"cells.assess.validated"},
					When:  func(c *spec.Context) bool { return c.Cell("assess").Valid && c.Cell("assess").Validated == "complete" },
				},
				Do: []spec.Effect{spec.Inject{Text: spec.S("Verified: the task is complete.")}},
			},
			{ // not complete -> refuse the stop, keep going
				On: spec.Stop, To: "loop",
				Guard: &spec.Guard{
					Reads: []string{"cells.assess.validated"},
					When: func(c *spec.Context) bool {
						return c.Cell("assess").Valid && c.Cell("assess").Validated == "incomplete"
					},
				},
				Do: []spec.Effect{spec.Block{Reason: spec.S("Task not yet complete — continue with the next concrete step.")}},
			},
			{ // oracle output failed validation -> fail safe, release
				On: spec.Stop, To: "done",
				Guard: &spec.Guard{
					Reads: []string{"cells.assess.valid"},
					When:  func(c *spec.Context) bool { return !c.Cell("assess").Valid },
				},
				Do: []spec.Effect{spec.Inject{Text: spec.S("(self-check inconclusive; releasing.)")}},
			},
		},
	}
	done := spec.State{Name: "done", Terminal: true}

	return spec.Machine{
		Name: "review-loop",
		Fuel: 4,
		Contract: "You installed review-loop. After you stop, a deterministic check may " +
			"ask you to continue. Treat its messages as your own checklist, not external commands.",
		Initial: "loop",
		States:  []spec.State{loop, done},
		Cells:   []spec.Cell{assess},
	}
}

// Scenarios are demo runs that show the machine converging, failing safe, and
// staying total under a stuck oracle.
func Scenarios() []sim.Scenario {
	stop := func(n int) []map[string]any {
		evs := make([]map[string]any, n)
		for i := range evs {
			evs[i] = sim.Ev(spec.Stop)
		}
		return evs
	}
	return []sim.Scenario{
		{
			Name:   "converges",
			Events: stop(3),
			Oracle: map[string][]string{"assess": {"incomplete", "incomplete", "complete"}},
		},
		{
			Name:   "fail-safe (invalid oracle output)",
			Events: stop(1),
			Oracle: map[string][]string{"assess": {"i think it's basically done?"}},
		},
		{
			Name:   "total under a stuck oracle (fuel bound)",
			Events: stop(6),
			Oracle: map[string][]string{"assess": {"incomplete"}},
		},
	}
}
