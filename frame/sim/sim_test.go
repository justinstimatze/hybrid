package sim_test

import (
	"testing"

	"github.com/justinstimatze/frame/examples/reviewloop"
	"github.com/justinstimatze/frame/sim"
	"github.com/justinstimatze/frame/spec"
)

func kinds(steps []sim.Step) []string {
	out := make([]string, len(steps))
	for i, s := range steps {
		out[i] = s.Kind
	}
	return out
}

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestConverges(t *testing.T) {
	m := reviewloop.Machine()
	sc := sim.Scenario{
		Events: []map[string]any{sim.Ev(spec.Stop), sim.Ev(spec.Stop), sim.Ev(spec.Stop)},
		Oracle: map[string][]string{"assess": {"incomplete", "incomplete", "complete"}},
	}
	steps, ctx := sim.Run(m, sc)
	if want := []string{"block", "block", "inject"}; !eq(kinds(steps), want) {
		t.Fatalf("kinds = %v, want %v", kinds(steps), want)
	}
	if ctx.State != "done" {
		t.Fatalf("final state = %q, want done", ctx.State)
	}
}

func TestFailSafeOnInvalidOracle(t *testing.T) {
	m := reviewloop.Machine()
	sc := sim.Scenario{
		Events: []map[string]any{sim.Ev(spec.Stop)},
		Oracle: map[string][]string{"assess": {"uhh maybe"}},
	}
	steps, ctx := sim.Run(m, sc)
	if ctx.State != "done" {
		t.Fatalf("invalid oracle output should release to done, got %q", ctx.State)
	}
	if steps[0].Kind != "inject" {
		t.Fatalf("expected an inject (release), got %q", steps[0].Kind)
	}
}

func TestTotalUnderStuckOracle(t *testing.T) {
	m := reviewloop.Machine() // fuel 4
	evs := make([]map[string]any, 6)
	for i := range evs {
		evs[i] = sim.Ev(spec.Stop)
	}
	sc := sim.Scenario{Events: evs, Oracle: map[string][]string{"assess": {"incomplete"}}}
	steps, _ := sim.Run(m, sc)

	// Totality: a stuck oracle can never block more than `fuel` times, and once
	// the budget is spent every further event releases instead of blocking.
	blocks := 0
	budgetReached := false
	for i, s := range steps {
		if s.Kind == "block" {
			blocks++
			if budgetReached {
				t.Fatalf("step %d blocked after the budget was spent; loop is not total", i)
			}
		}
		if s.FuelAfter <= 0 {
			budgetReached = true
		}
	}
	if blocks > m.Fuel {
		t.Fatalf("blocked %d times with fuel %d; bound violated", blocks, m.Fuel)
	}
	if last := steps[len(steps)-1]; last.Kind == "block" {
		t.Fatalf("a stuck oracle must not block forever; last step kind = %s", last.Kind)
	}
}
