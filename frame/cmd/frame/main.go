// Command frame drives the statechart mesh: check a machine for soundness,
// compile it to a settings.json hooks fragment, simulate it deterministically,
// or run it as the live hook dispatcher.
//
//	frame check   <machine>           # static soundness report
//	frame compile <machine> [binary]  # emit settings.json hooks fragment
//	frame sim     <machine>           # run the machine's demo scenarios
//	frame run --machine <machine>     # hook dispatcher (reads a hook event on stdin)
//
// `run` is what a compiled hook invokes. It reads one hook event as JSON on
// stdin, advances the machine, persists, and emits the hook protocol. The model
// call (runtime.Oracle) is unbound in this skeleton; until it is wired, cells
// fail validation and the machine takes its fail-safe path.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/justinstimatze/frame/check"
	"github.com/justinstimatze/frame/compile"
	"github.com/justinstimatze/frame/registry"
	"github.com/justinstimatze/frame/runtime"
	"github.com/justinstimatze/frame/sim"
	"github.com/justinstimatze/frame/spec"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "check":
		os.Exit(cmdCheck(os.Args[2:]))
	case "compile":
		os.Exit(cmdCompile(os.Args[2:]))
	case "sim":
		os.Exit(cmdSim(os.Args[2:]))
	case "run":
		os.Exit(cmdRun(os.Args[2:]))
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "frame: a guarded statechart compiled to a Claude Code hook mesh\n\n")
	fmt.Fprintf(os.Stderr, "  frame check   <machine>\n")
	fmt.Fprintf(os.Stderr, "  frame compile <machine> [binary-path]\n")
	fmt.Fprintf(os.Stderr, "  frame sim     <machine>\n")
	fmt.Fprintf(os.Stderr, "  frame run --machine <machine>   (reads a hook event on stdin)\n\n")
	fmt.Fprintf(os.Stderr, "machines: %s\n", strings.Join(registry.Names(), ", "))
	os.Exit(2)
}

func lookup(name string) registry.Entry {
	e, ok := registry.Get(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown machine %q (have: %s)\n", name, strings.Join(registry.Names(), ", "))
		os.Exit(2)
	}
	return e
}

func cmdCheck(args []string) int {
	if len(args) < 1 {
		usage()
	}
	m := lookup(args[0]).Machine()
	errs := check.Check(m)
	if len(errs) == 0 {
		fmt.Printf("OK: %q is sound (%d states, fuel %d)\n", m.Name, len(m.States), m.Fuel)
		return 0
	}
	fmt.Printf("UNSOUND: %q\n", m.Name)
	for _, e := range errs {
		fmt.Println("  " + e)
	}
	return 1
}

func cmdCompile(args []string) int {
	if len(args) < 1 {
		usage()
	}
	m := lookup(args[0]).Machine()
	binary := "frame"
	if len(args) > 1 {
		binary = args[1]
	}
	out, err := compile.SettingsJSON(m, binary)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Println(out)
	return 0
}

func cmdSim(args []string) int {
	if len(args) < 1 {
		usage()
	}
	e := lookup(args[0])
	m := e.Machine()
	if e.Scenarios == nil {
		fmt.Fprintf(os.Stderr, "%q has no demo scenarios\n", m.Name)
		return 1
	}
	for _, sc := range e.Scenarios() {
		steps, ctx := sim.Run(m, sc)
		fmt.Print(sim.Render(sc.Name, steps))
		fmt.Printf("  -> final state %q, fuel %d\n\n", ctx.State, ctx.Fuel)
	}
	return 0
}

func cmdRun(args []string) int {
	// Fail-open: any problem yields exit 0 with no output.
	name := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--machine" && i+1 < len(args) {
			name = args[i+1]
		}
	}
	e, ok := registry.Get(name)
	if !ok {
		return 0
	}
	m := e.Machine()
	var event map[string]any
	if err := json.NewDecoder(os.Stdin).Decode(&event); err != nil {
		return 0
	}
	ctx := runtime.LoadContext(m, event)
	// Oracle is unbound in the skeleton; a nil-returning stub makes cells fail
	// validation, so the machine takes its fail-safe path rather than guessing.
	out := runtime.SafeDispatch(m, event, ctx, func(spec.Cell, *spec.Context) string { return "" })
	_ = runtime.SaveContext(m, ctx)
	return runtime.Emit(spec.Trigger(asString(event["hook_event_name"])), out)
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
