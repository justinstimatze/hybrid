// Package check is the static checker — the part that makes unsafe meshes
// inexpressible. Check returns a slice of error strings (empty == sound); each
// is prefixed with a stable code so callers and tests can match on it:
//
//	E-FUEL      fuel is not a positive step budget (totality)
//	E-DUP       duplicate state name
//	E-INITIAL   initial state undefined
//	E-TARGET    transition targets an undefined state
//	E-TERMINAL  a terminal state has outgoing transitions
//	E-ORPHAN    a state is unreachable from initial
//	E-HALT      no terminal state is reachable (the loop can never end)
//	E-CELL      a guard/Run references a cell not in the machine registry
//	E-ORACLE    a guard gates control flow on a cell's *raw* output
//	E-INJECT    an Inject effect sits on a trigger that cannot inject
//	E-BLOCK     a Block effect sits on a trigger that cannot block
//
// Cell-has-a-validator is not checked here: it is enforced one level down, by
// construction (spec.NewCell rejects a nil validator and the field is
// unexported), so an ungated oracle cannot be built in the first place.
package check

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/justinstimatze/frame/spec"
)

var cellPath = regexp.MustCompile(`^cells\.(\w+)\.(\w+)$`)

// CellRead is a parsed "cells.<Cell>.<Field>" guard dependency.
type CellRead struct{ Cell, Field string }

// CellReads extracts the cell dependencies declared in a guard's Reads.
func CellReads(reads []string) []CellRead {
	var out []CellRead
	for _, r := range reads {
		if m := cellPath.FindStringSubmatch(r); m != nil {
			out = append(out, CellRead{Cell: m[1], Field: m[2]})
		}
	}
	return out
}

// Check returns every soundness problem with the machine (empty == sound).
func Check(m spec.Machine) []string {
	var errs []string
	byName := map[string]spec.State{}
	for _, s := range m.States {
		byName[s.Name] = s
	}
	cellNames := map[string]bool{}
	for _, c := range m.Cells {
		cellNames[c.Name] = true
	}

	// --- well-formedness ---
	if m.Fuel <= 0 {
		errs = append(errs, fmt.Sprintf("E-FUEL: fuel must be a positive step budget, got %d", m.Fuel))
	}
	seenName := map[string]int{}
	for _, s := range m.States {
		seenName[s.Name]++
	}
	for _, name := range sortedKeys(seenName) {
		if seenName[name] > 1 {
			errs = append(errs, fmt.Sprintf("E-DUP: duplicate state name %q", name))
		}
	}
	if _, ok := byName[m.Initial]; !ok {
		errs = append(errs, fmt.Sprintf("E-INITIAL: initial state %q is not defined", m.Initial))
	}
	for _, s := range m.States {
		for _, t := range s.On {
			if _, ok := byName[t.To]; !ok {
				errs = append(errs, fmt.Sprintf("E-TARGET: %s --%s--> undefined state %q", s.Name, t.On, t.To))
			}
		}
		if s.Terminal && len(s.On) > 0 {
			errs = append(errs, fmt.Sprintf("E-TERMINAL: terminal state %q has outgoing transitions", s.Name))
		}
	}

	// --- reachability + a reachable halt ---
	if _, ok := byName[m.Initial]; ok {
		seen := map[string]bool{}
		stack := []string{m.Initial}
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if seen[cur] {
				continue
			}
			seen[cur] = true
			for _, t := range byName[cur].On {
				if _, ok := byName[t.To]; ok {
					stack = append(stack, t.To)
				}
			}
		}
		for _, s := range m.States {
			if !seen[s.Name] {
				errs = append(errs, fmt.Sprintf("E-ORPHAN: state %q is unreachable from initial", s.Name))
			}
		}
		halt := false
		for name := range seen {
			if byName[name].Terminal {
				halt = true
				break
			}
		}
		if !halt {
			errs = append(errs, "E-HALT: no terminal (halt) state is reachable from initial")
		}
	}

	// --- per-transition invariants ---
	for _, s := range m.States {
		for _, t := range s.On {
			if t.Guard != nil {
				for _, cr := range CellReads(t.Guard.Reads) {
					if !cellNames[cr.Cell] {
						errs = append(errs, fmt.Sprintf("E-CELL: %s guard reads unknown cell %q (not in machine.Cells)", s.Name, cr.Cell))
					}
					switch cr.Field {
					case "term", "wellformed", "safe":
					default:
						errs = append(errs, fmt.Sprintf("E-ORACLE: %s guard reads cells.%s.%s — control flow may only depend on a cell's formal-language output (term/wellformed/safe), never its raw output", s.Name, cr.Cell, cr.Field))
					}
				}
			}
			for _, eff := range t.Do {
				switch e := eff.(type) {
				case spec.Inject:
					if !spec.Injecting(t.On) {
						errs = append(errs, fmt.Sprintf("E-INJECT: Inject on a %s transition (%s); that trigger cannot inject context", t.On, s.Name))
					}
				case spec.Block:
					if !spec.Gating(t.On) {
						errs = append(errs, fmt.Sprintf("E-BLOCK: Block on a %s transition (%s); that trigger cannot halt control flow", t.On, s.Name))
					}
				case spec.Run:
					if !cellNames[e.Cell.Name] {
						errs = append(errs, fmt.Sprintf("E-CELL: %s runs cell %q not in machine.Cells", s.Name, e.Cell.Name))
					}
				}
			}
		}
	}

	return errs
}

// Validate returns an error listing every problem, or nil if the machine is sound.
func Validate(m spec.Machine) error {
	errs := Check(m)
	if len(errs) == 0 {
		return nil
	}
	out := fmt.Sprintf("machine %q is not sound:", m.Name)
	for _, e := range errs {
		out += "\n  " + e
	}
	return fmt.Errorf("%s", out)
}

func sortedKeys(m map[string]int) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
