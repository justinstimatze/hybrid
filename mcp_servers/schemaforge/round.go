package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ItemResult is the per-item record written to round{N}/items/{id}.json.
// It's everything score_round needs and everything a caller wants to inspect
// when a particular item scored poorly.
type ItemResult struct {
	ID              string  `json:"id"`
	NotationText    string  `json:"notation_text"`
	NotationTokens  int     `json:"notation_tokens"`
	ExpandedText    string  `json:"expanded_text"`
	ExpandedTokens  int     `json:"expanded_tokens"`
	RoundtripScore  float64 `json:"roundtrip_score"`
	ScoreReasoning  string  `json:"score_reasoning"`
	ExpansionRatio  float64 `json:"expansion_ratio"`
	SkippedReason   string  `json:"skipped_reason,omitempty"`
}

// RoundMetrics is the aggregate written to round{N}/metrics.json. The Summary
// field is a human-readable string that's also fed back into the next round's
// design call as previousMetricsSummary.
type RoundMetrics struct {
	RoundNumber          int          `json:"round_number"`
	NotationSpecTokens   int          `json:"notation_spec_tokens"`
	ItemCount            int          `json:"item_count"`
	SkippedCount         int          `json:"skipped_count"`
	MeanExpansionRatio   float64      `json:"mean_expansion_ratio"`
	MeanCorrectness      float64      `json:"mean_correctness"`
	TotalNotationTokens  int          `json:"total_notation_tokens"`
	TotalExpandedTokens  int          `json:"total_expanded_tokens"`
	Items                []ItemResult `json:"items"`
	Summary              string       `json:"summary"`
}

// RunRoundParams are the inputs to one full round.
type RunRoundParams struct {
	CorpusPath           string
	Target               string
	Rubric               string
	Model                string
	OutputDir            string // run-level dir; round subdir created inside
	RoundNumber          int
	PreviousNotation     string // empty for round 1
	PreviousMetricsPath  string // path to prior round's metrics.json; empty for round 1
}

// RunRound drives one full round: design → compress all → expand all →
// evaluate all → score. Writes round{N}/notation.txt, round{N}/items/{id}.json
// per item (incrementally — caller can ctrl-C and resume), round{N}/metrics.json
// at the end.
//
// Resume semantics: if round{N}/items/{id}.json already exists, that item is
// loaded from disk and not re-processed. The notation.txt similarly is reused
// if present (so a partially-completed round picks up where it left off).
func RunRound(ctx context.Context, llm LLMClient, p RunRoundParams) (*RoundMetrics, error) {
	if p.RoundNumber < 1 {
		return nil, fmt.Errorf("round_number must be >= 1")
	}
	if p.OutputDir == "" {
		return nil, fmt.Errorf("output_dir is required")
	}
	items, err := LoadCorpus(p.CorpusPath)
	if err != nil {
		return nil, err
	}

	roundDir := filepath.Join(p.OutputDir, fmt.Sprintf("round%d", p.RoundNumber))
	itemsDir := filepath.Join(roundDir, "items")
	if err := os.MkdirAll(itemsDir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir round dir: %w", err)
	}

	// Notation: reuse if already on disk, else design.
	notationPath := filepath.Join(roundDir, "notation.txt")
	notationSpec, err := os.ReadFile(notationPath)
	if err != nil || len(notationSpec) == 0 {
		prevSummary := ""
		if p.PreviousMetricsPath != "" {
			prev, perr := LoadRoundMetrics(p.PreviousMetricsPath)
			if perr == nil {
				prevSummary = prev.Summary
			}
		}
		seeds := SeedItems(items, 2)
		spec, err := DesignNotation(ctx, llm, seeds, p.Target, p.Model, p.PreviousNotation, prevSummary)
		if err != nil {
			return nil, fmt.Errorf("design_notation: %w", err)
		}
		if err := os.WriteFile(notationPath, []byte(spec), 0o644); err != nil {
			return nil, fmt.Errorf("write notation.txt: %w", err)
		}
		notationSpec = []byte(spec)
	}
	specStr := string(notationSpec)

	// Per-item: compress → expand → evaluate.
	for _, item := range items {
		itemPath := filepath.Join(itemsDir, item.ID+".json")
		if _, err := os.Stat(itemPath); err == nil {
			continue // already done; resume
		}

		result := ItemResult{ID: item.ID}

		notation, nTokens, err := Compress(ctx, llm, item.SpecText, specStr, p.Target, p.Model)
		if err != nil {
			result.SkippedReason = "compress: " + err.Error()
			_ = writeItem(itemPath, result)
			continue
		}
		result.NotationText = notation
		result.NotationTokens = nTokens

		expanded, eTokens, err := Expand(ctx, llm, specStr, notation, p.Target, p.Model)
		if err != nil {
			result.SkippedReason = "expand: " + err.Error()
			_ = writeItem(itemPath, result)
			continue
		}
		result.ExpandedText = expanded
		result.ExpandedTokens = eTokens
		if nTokens > 0 {
			result.ExpansionRatio = float64(eTokens) / float64(nTokens)
		}

		score, reasoning, err := EvaluateRoundtrip(ctx, llm, item.SpecText, expanded, p.Rubric, p.Model)
		if err != nil {
			result.SkippedReason = "evaluate: " + err.Error()
			_ = writeItem(itemPath, result)
			continue
		}
		result.RoundtripScore = score
		result.ScoreReasoning = reasoning

		if err := writeItem(itemPath, result); err != nil {
			return nil, fmt.Errorf("write item %s: %w", item.ID, err)
		}
	}

	// Aggregate.
	metrics, err := scoreRoundDir(roundDir, p.RoundNumber, EstimateTokens(specStr))
	if err != nil {
		return nil, fmt.Errorf("score round: %w", err)
	}
	metricsPath := filepath.Join(roundDir, "metrics.json")
	if err := writeJSON(metricsPath, metrics); err != nil {
		return nil, fmt.Errorf("write metrics.json: %w", err)
	}
	return metrics, nil
}

// ScoreRoundDir reads a round directory and computes RoundMetrics without
// touching the LLM. Useful for re-aggregating after manual edits, and as the
// underlying impl of the score_round tool.
func ScoreRoundDir(roundDir string) (*RoundMetrics, error) {
	// Recover round number from path: ".../round{N}".
	base := filepath.Base(strings.TrimRight(roundDir, "/"))
	var roundNum int
	_, err := fmt.Sscanf(base, "round%d", &roundNum)
	if err != nil || roundNum < 1 {
		return nil, fmt.Errorf("could not infer round number from %q (expected '.../roundN')", roundDir)
	}
	specBytes, _ := os.ReadFile(filepath.Join(roundDir, "notation.txt"))
	return scoreRoundDir(roundDir, roundNum, EstimateTokens(string(specBytes)))
}

func scoreRoundDir(roundDir string, roundNum, specTokens int) (*RoundMetrics, error) {
	itemsDir := filepath.Join(roundDir, "items")
	entries, err := os.ReadDir(itemsDir)
	if err != nil {
		return nil, fmt.Errorf("read items dir: %w", err)
	}
	var items []ItemResult
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		var it ItemResult
		raw, err := os.ReadFile(filepath.Join(itemsDir, e.Name()))
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, &it); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		items = append(items, it)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })

	m := &RoundMetrics{
		RoundNumber:        roundNum,
		NotationSpecTokens: specTokens,
		ItemCount:          len(items),
		Items:              items,
	}
	var erSum, scoreSum float64
	var erCount, scoredCount int
	for _, it := range items {
		if it.SkippedReason != "" {
			m.SkippedCount++
			continue
		}
		m.TotalNotationTokens += it.NotationTokens
		m.TotalExpandedTokens += it.ExpandedTokens
		if it.ExpansionRatio > 0 {
			erSum += it.ExpansionRatio
			erCount++
		}
		scoreSum += it.RoundtripScore
		scoredCount++
	}
	if erCount > 0 {
		m.MeanExpansionRatio = erSum / float64(erCount)
	}
	if scoredCount > 0 {
		m.MeanCorrectness = scoreSum / float64(scoredCount)
	}
	m.Summary = fmt.Sprintf(
		"Round %d: spec=%d toks, %d items (%d skipped), mean ER=%.1fx, mean correctness=%.0f%%",
		roundNum, specTokens, m.ItemCount, m.SkippedCount,
		m.MeanExpansionRatio, m.MeanCorrectness*100,
	)
	return m, nil
}

func LoadRoundMetrics(path string) (*RoundMetrics, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m RoundMetrics
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func writeItem(path string, r ItemResult) error {
	if r.NotationTokens > 0 && r.ExpandedTokens > 0 {
		r.ExpansionRatio = float64(r.ExpandedTokens) / float64(r.NotationTokens)
	}
	return writeJSON(path, r)
}

func writeJSON(path string, v any) error {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}

// TrainingPair is one (notation, expansion) pair extracted from a round for
// fine-tuning. Only items with score >= minScore are included.
type TrainingPair struct {
	ID             string  `json:"id"`
	Input          string  `json:"input"`
	Output         string  `json:"output"`
	ExpansionRatio float64 `json:"expansion_ratio"`
	Score          float64 `json:"score"`
}

// ExtractTrainingPairs reads a round and returns items above minScore as
// (notation→expansion) pairs. Caller writes them to JSONL or wherever needed.
func ExtractTrainingPairs(roundDir string, minScore float64) ([]TrainingPair, error) {
	m, err := ScoreRoundDir(roundDir)
	if err != nil {
		return nil, err
	}
	var pairs []TrainingPair
	for _, it := range m.Items {
		if it.SkippedReason != "" {
			continue
		}
		if it.RoundtripScore < minScore {
			continue
		}
		pairs = append(pairs, TrainingPair{
			ID:             it.ID,
			Input:          it.NotationText,
			Output:         it.ExpandedText,
			ExpansionRatio: it.ExpansionRatio,
			Score:          it.RoundtripScore,
		})
	}
	return pairs, nil
}
