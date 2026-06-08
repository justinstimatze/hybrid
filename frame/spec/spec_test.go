package spec_test

import (
	"testing"

	"github.com/justinstimatze/frame/spec"
)

// A formal oracle over a tiny command language: L = {"read <x>", "write <x>"};
// the safety check forbids writes. Demonstrates the two independent stages —
// membership in the language, then safety of the denoted action.
func TestFormalOracleCheck(t *testing.T) {
	cell := spec.NewCell("cmd", "m", "emit `read X` or `write X`",
		func(raw string) (string, bool) { // Grammar: membership in L
			if len(raw) > 5 && (raw[:5] == "read " || raw[:6] == "write ") {
				return raw, true
			}
			return "", false
		},
		func(term string) bool { return term[:4] == "read" }, // Safety: reads only
	)

	cases := []struct {
		name             string
		raw              string
		wellFormed, safe bool
	}{
		{"in-language, safe", "read notes.md", true, true},
		{"in-language, unsafe", "write /etc/passwd", true, false},
		{"outside the language", "rm -rf /", false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := cell.Check(tc.raw)
			if r.WellFormed != tc.wellFormed || r.Safe != tc.safe {
				t.Fatalf("Check(%q) = {WellFormed:%v Safe:%v}, want {%v %v}",
					tc.raw, r.WellFormed, r.Safe, tc.wellFormed, tc.safe)
			}
		})
	}
}

func TestNewCellRejectsMissingStages(t *testing.T) {
	g := func(r string) (string, bool) { return r, true }
	s := func(string) bool { return true }
	for _, tc := range []struct {
		name string
		g    spec.Grammar
		s    spec.Safety
	}{
		{"nil grammar", nil, s},
		{"nil safety", g, nil},
		{"both nil", nil, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic for an incomplete formal oracle")
				}
			}()
			_ = spec.NewCell("x", "m", "i", tc.g, tc.s)
		})
	}
}
