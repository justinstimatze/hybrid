// Package compile turns a sound machine into a Claude Code settings.json hooks
// fragment. The machine is validated first (a machine that does not pass check
// does not compile). The fragment registers the single generic dispatcher once
// per trigger the machine actually uses — definitions are written once and never
// mutated, so there is no /hooks refresh dance; all dynamism flows through the
// per-session state the runtime reads.
package compile

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/justinstimatze/frame/check"
	"github.com/justinstimatze/frame/spec"
)

// SettingsFragment returns the hooks block for the machine. binary is the path
// to the compiled frame binary the hook should invoke.
func SettingsFragment(m spec.Machine, binary string) (map[string]any, error) {
	if err := check.Validate(m); err != nil {
		return nil, err
	}
	used := map[spec.Trigger]bool{}
	for _, s := range m.States {
		for _, t := range s.On {
			used[t.On] = true
		}
	}
	triggers := make([]string, 0, len(used))
	for t := range used {
		triggers = append(triggers, string(t))
	}
	sort.Strings(triggers)

	command := fmt.Sprintf("%s run --machine %s", binary, m.Name)
	hooks := map[string]any{}
	for _, t := range triggers {
		hooks[t] = []any{map[string]any{
			"matcher": "*",
			"hooks":   []any{map[string]any{"type": "command", "command": command}},
		}}
	}
	return map[string]any{"hooks": hooks}, nil
}

// SettingsJSON renders SettingsFragment as indented JSON.
func SettingsJSON(m spec.Machine, binary string) (string, error) {
	frag, err := SettingsFragment(m, binary)
	if err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(frag, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
