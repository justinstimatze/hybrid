package check_test

import (
	"strings"
	"testing"

	"github.com/justinstimatze/frame/check"
	"github.com/justinstimatze/frame/examples/reviewloop"
	"github.com/justinstimatze/frame/spec"
)

func hasCode(errs []string, code string) bool {
	for _, e := range errs {
		if strings.HasPrefix(e, code) {
			return true
		}
	}
	return false
}

func TestExampleIsSound(t *testing.T) {
	if errs := check.Check(reviewloop.Machine()); len(errs) != 0 {
		t.Fatalf("review-loop should be sound, got: %v", errs)
	}
}

func TestUngatedOracleIsInexpressible(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewCell with a nil grammar/safety should panic")
		}
	}()
	_ = spec.NewCell("x", "m", "i", nil, nil)
}

func TestGuardOnRawOutputRejected(t *testing.T) {
	c := spec.NewCell("c", "m", "i",
		func(r string) (string, bool) { return r, true },
		func(string) bool { return true })
	m := spec.Machine{
		Name: "bad", Fuel: 2, Initial: "a", Cells: []spec.Cell{c},
		States: []spec.State{
			{Name: "a", On: []spec.Transition{{
				On: spec.Stop, To: "b",
				Guard: &spec.Guard{Reads: []string{"cells.c.raw"}, When: func(*spec.Context) bool { return true }},
			}}},
			{Name: "b", Terminal: true},
		},
	}
	if errs := check.Check(m); !hasCode(errs, "E-ORACLE") {
		t.Fatalf("expected E-ORACLE for a guard reading cells.c.raw, got: %v", errs)
	}
}

func TestInjectChannelMismatch(t *testing.T) {
	// Inject on PreToolUse: that trigger cannot inject.
	m := spec.Machine{
		Name: "bad", Fuel: 2, Initial: "a",
		States: []spec.State{
			{Name: "a", On: []spec.Transition{{
				On: spec.PreToolUse, To: "b",
				Do: []spec.Effect{spec.Inject{Text: spec.S("hi")}},
			}}},
			{Name: "b", Terminal: true},
		},
	}
	if errs := check.Check(m); !hasCode(errs, "E-INJECT") {
		t.Fatalf("expected E-INJECT, got: %v", errs)
	}
}

func TestBlockChannelMismatch(t *testing.T) {
	// Block on PostToolUse: that trigger cannot halt control flow.
	m := spec.Machine{
		Name: "bad", Fuel: 2, Initial: "a",
		States: []spec.State{
			{Name: "a", On: []spec.Transition{{
				On: spec.PostToolUse, To: "b",
				Do: []spec.Effect{spec.Block{Reason: spec.S("no")}},
			}}},
			{Name: "b", Terminal: true},
		},
	}
	if errs := check.Check(m); !hasCode(errs, "E-BLOCK") {
		t.Fatalf("expected E-BLOCK, got: %v", errs)
	}
}

func TestStructuralProblems(t *testing.T) {
	cases := []struct {
		name string
		m    spec.Machine
		code string
	}{
		{"zero fuel", spec.Machine{Name: "m", Fuel: 0, Initial: "a",
			States: []spec.State{{Name: "a", Terminal: true}}}, "E-FUEL"},
		{"bad target", spec.Machine{Name: "m", Fuel: 1, Initial: "a",
			States: []spec.State{{Name: "a", On: []spec.Transition{{On: spec.Stop, To: "nope"}}}}}, "E-TARGET"},
		{"no halt", spec.Machine{Name: "m", Fuel: 1, Initial: "a",
			States: []spec.State{{Name: "a", On: []spec.Transition{{On: spec.Stop, To: "a"}}}}}, "E-HALT"},
		{"orphan", spec.Machine{Name: "m", Fuel: 1, Initial: "a",
			States: []spec.State{{Name: "a", Terminal: true}, {Name: "b", Terminal: true}}}, "E-ORPHAN"},
		{"undefined initial", spec.Machine{Name: "m", Fuel: 1, Initial: "z",
			States: []spec.State{{Name: "a", Terminal: true}}}, "E-INITIAL"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if errs := check.Check(tc.m); !hasCode(errs, tc.code) {
				t.Fatalf("expected %s, got: %v", tc.code, errs)
			}
		})
	}
}
