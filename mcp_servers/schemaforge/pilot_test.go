//go:build pilot

// Integration test that hits the real Anthropic API. Run with:
//
//	go test -tags pilot -timeout 60m -v -run TestPilot
//
// Requires ANTHROPIC_API_KEY in env. All inputs are env-driven so the test
// works for any corpus/target/rubric:
//
//	SCHEMAFORGE_CORPUS  path to JSONL corpus (required)
//	SCHEMAFORGE_OUT     output dir for round{N}/ subdirs (required)
//	SCHEMAFORGE_TARGET  free-form target string (required)
//	SCHEMAFORGE_RUBRIC  free-form scoring rubric (optional; uses default rubric)
//	SCHEMAFORGE_ROUNDS  number of rounds (default: 3)
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPilot(t *testing.T) {
	corpus := os.Getenv("SCHEMAFORGE_CORPUS")
	outDir := os.Getenv("SCHEMAFORGE_OUT")
	target := os.Getenv("SCHEMAFORGE_TARGET")
	if corpus == "" || outDir == "" || target == "" {
		t.Skip("set SCHEMAFORGE_CORPUS, SCHEMAFORGE_OUT, SCHEMAFORGE_TARGET to run the pilot")
	}
	rubric := os.Getenv("SCHEMAFORGE_RUBRIC") // empty → default rubric in EvaluateRoundtrip
	rounds := 3
	if v := os.Getenv("SCHEMAFORGE_ROUNDS"); v != "" {
		fmt.Sscanf(v, "%d", &rounds)
	}

	llm, err := NewAnthropicClient()
	if err != nil {
		t.Fatalf("anthropic client: %v", err)
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}

	var prevNotation string
	var prevMetricsPath string
	for r := 1; r <= rounds; r++ {
		t.Logf("=== round %d/%d starting ===", r, rounds)
		m, err := RunRound(context.Background(), llm, RunRoundParams{
			CorpusPath:          corpus,
			Target:              target,
			Rubric:              rubric,
			Model:               defaultModel,
			OutputDir:           outDir,
			RoundNumber:         r,
			PreviousNotation:    prevNotation,
			PreviousMetricsPath: prevMetricsPath,
		})
		if err != nil {
			t.Fatalf("round %d: %v", r, err)
		}
		t.Logf("round %d: %s", r, m.Summary)
		t.Logf("  per-item:")
		for _, it := range m.Items {
			if it.SkippedReason != "" {
				t.Logf("    %s  SKIPPED  %s", it.ID, it.SkippedReason)
			} else {
				t.Logf("    %-22s ER=%5.1fx  score=%.2f  (%d→%d toks)",
					it.ID, it.ExpansionRatio, it.RoundtripScore,
					it.NotationTokens, it.ExpandedTokens)
			}
		}
		// Read back the notation we just wrote, for round N+1.
		notationPath := filepath.Join(outDir, fmt.Sprintf("round%d", r), "notation.txt")
		nb, err := os.ReadFile(notationPath)
		if err != nil {
			t.Fatalf("read notation.txt: %v", err)
		}
		prevNotation = string(nb)
		prevMetricsPath = filepath.Join(outDir, fmt.Sprintf("round%d", r), "metrics.json")
	}
}
