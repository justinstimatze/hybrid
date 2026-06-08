// Package registry maps machine names to the compiled-in machines the CLI can
// check, compile, simulate, and run. A real deployment registers its own
// machines here (or builds its own binary that imports this package).
package registry

import (
	"sort"

	"github.com/justinstimatze/frame/examples/reviewloop"
	"github.com/justinstimatze/frame/sim"
	"github.com/justinstimatze/frame/spec"
)

type Entry struct {
	Machine   func() spec.Machine
	Scenarios func() []sim.Scenario
}

var machines = map[string]Entry{
	"review-loop": {Machine: reviewloop.Machine, Scenarios: reviewloop.Scenarios},
}

func Get(name string) (Entry, bool) {
	e, ok := machines[name]
	return e, ok
}

func Names() []string {
	ns := make([]string, 0, len(machines))
	for n := range machines {
		ns = append(ns, n)
	}
	sort.Strings(ns)
	return ns
}
