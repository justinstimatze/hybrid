package main

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func nearly(got, want float64) bool {
	return math.Abs(got-want) < 1e-9
}

// FakeLLMClient lets us drive the ops deterministically. Each Call pops the
// next response off Responses; CallLog records (system, user, model) so tests
// can assert prompt assembly.
type FakeLLMClient struct {
	Responses []string
	CallLog   []FakeCall
}

type FakeCall struct {
	System string
	User   string
	Model  string
}

func (f *FakeLLMClient) Call(_ context.Context, system, user, model string, _ int) (string, error) {
	f.CallLog = append(f.CallLog, FakeCall{System: system, User: user, Model: model})
	if len(f.Responses) == 0 {
		return "", nil
	}
	r := f.Responses[0]
	f.Responses = f.Responses[1:]
	return r, nil
}

// --- corpus ---

func TestLoadCorpus(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "c.jsonl")
	body := `{"id":"a","spec_text":"alpha"}
{"id":"b","spec_text":"beta beta"}
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	items, err := LoadCorpus(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].ID != "a" || items[1].SpecText != "beta beta" {
		t.Fatalf("bad load: %+v", items)
	}
}

func TestLoadCorpus_RejectsDuplicates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "c.jsonl")
	os.WriteFile(path, []byte(`{"id":"a","spec_text":"x"}
{"id":"a","spec_text":"y"}
`), 0o644)
	_, err := LoadCorpus(path)
	if err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestSeedItems_LongestFirst(t *testing.T) {
	items := []CorpusItem{
		{ID: "short", SpecText: "x"},
		{ID: "long", SpecText: strings.Repeat("y ", 100)},
		{ID: "mid", SpecText: strings.Repeat("z ", 20)},
	}
	seeds := SeedItems(items, 2)
	if len(seeds) != 2 {
		t.Fatalf("want 2 seeds, got %d", len(seeds))
	}
	if seeds[0].ID != "long" || seeds[1].ID != "mid" {
		t.Fatalf("seeds not in length order: %+v", seeds)
	}
}

func TestEstimateTokens(t *testing.T) {
	if EstimateTokens("") != 0 {
		t.Errorf("empty should be 0")
	}
	if EstimateTokens("hello world") < 2 {
		t.Errorf("two words should yield >= 2 tokens")
	}
}

// --- ops ---

func TestDesignNotation_Round1(t *testing.T) {
	llm := &FakeLLMClient{Responses: []string{
		"Here's my design.\n<notation_spec>\nE = entity. F = field.\n</notation_spec>\nDone.",
	}}
	seeds := []CorpusItem{{ID: "x", SpecText: "stress test item content"}}
	spec, err := DesignNotation(context.Background(), llm, seeds, "TS code", "claude-test", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if spec != "E = entity. F = field." {
		t.Fatalf("bad spec: %q", spec)
	}
	if len(llm.CallLog) != 1 {
		t.Fatalf("want 1 call, got %d", len(llm.CallLog))
	}
	if !strings.Contains(llm.CallLog[0].User, "TS code") {
		t.Errorf("target should appear in user prompt")
	}
	if strings.Contains(llm.CallLog[0].User, "previous notation") {
		t.Errorf("round 1 should not mention previous")
	}
}

func TestDesignNotation_RoundN_UsesPrevious(t *testing.T) {
	llm := &FakeLLMClient{Responses: []string{
		"<notation_spec>v2 spec</notation_spec>",
	}}
	seeds := []CorpusItem{{ID: "x", SpecText: "anything"}}
	_, err := DesignNotation(context.Background(), llm, seeds, "T", "m",
		"v1 spec text", "Round 1: ER=2.0x, correctness=80%")
	if err != nil {
		t.Fatal(err)
	}
	user := llm.CallLog[0].User
	if !strings.Contains(user, "v1 spec text") {
		t.Errorf("previous notation should appear in user prompt")
	}
	if !strings.Contains(user, "ER=2.0x") {
		t.Errorf("previous metrics summary should appear")
	}
}

func TestDesignNotation_MissingTag(t *testing.T) {
	llm := &FakeLLMClient{Responses: []string{"no tags here, just prose"}}
	_, err := DesignNotation(context.Background(), llm,
		[]CorpusItem{{ID: "x", SpecText: "y"}}, "T", "m", "", "")
	if err == nil || !strings.Contains(err.Error(), "notation_spec") {
		t.Errorf("expected tag error, got %v", err)
	}
}

func TestCompress(t *testing.T) {
	llm := &FakeLLMClient{Responses: []string{"<notation>E[F:s]</notation>"}}
	notation, tokens, err := Compress(context.Background(), llm, "an entity with a string field", "spec", "T", "m")
	if err != nil {
		t.Fatal(err)
	}
	if notation != "E[F:s]" {
		t.Errorf("bad notation: %q", notation)
	}
	if tokens < 1 {
		t.Errorf("expected positive token count")
	}
	if !strings.Contains(llm.CallLog[0].System, "spec") {
		t.Errorf("notation_spec should be in system prompt for caching")
	}
}

func TestExpand(t *testing.T) {
	expansion := "function getThing() { return {}; }"
	llm := &FakeLLMClient{Responses: []string{expansion}}
	out, tokens, err := Expand(context.Background(), llm, "spec", "E[F:s]", "TS code", "m")
	if err != nil {
		t.Fatal(err)
	}
	if out != expansion {
		t.Errorf("bad expansion: %q", out)
	}
	if tokens < 1 {
		t.Errorf("expected positive token count")
	}
}

func TestEvaluateRoundtrip_GoodJSON(t *testing.T) {
	llm := &FakeLLMClient{Responses: []string{`{"overall": 0.85, "reasoning": "captures most fields"}`}}
	score, reasoning, err := EvaluateRoundtrip(context.Background(), llm, "spec", "expanded", "rubric", "m")
	if err != nil {
		t.Fatal(err)
	}
	if score != 0.85 {
		t.Errorf("score: got %v want 0.85", score)
	}
	if reasoning != "captures most fields" {
		t.Errorf("reasoning: got %q", reasoning)
	}
}

func TestEvaluateRoundtrip_ParseFailureFallsBack(t *testing.T) {
	llm := &FakeLLMClient{Responses: []string{"sorry, I cannot produce a JSON score"}}
	score, reasoning, err := EvaluateRoundtrip(context.Background(), llm, "s", "e", "r", "m")
	if err != nil {
		t.Fatalf("parse failure should not error, got %v", err)
	}
	if score != 0.5 {
		t.Errorf("fallback score should be 0.5, got %v", score)
	}
	if !strings.HasPrefix(reasoning, "parse_failed") {
		t.Errorf("reasoning should mark parse failure: %q", reasoning)
	}
}

func TestEvaluateRoundtrip_ClampsRange(t *testing.T) {
	llm := &FakeLLMClient{Responses: []string{`{"overall": 1.5, "reasoning": "too high"}`}}
	score, _, _ := EvaluateRoundtrip(context.Background(), llm, "s", "e", "r", "m")
	if score != 1.0 {
		t.Errorf("score >1 should clamp to 1, got %v", score)
	}
}

// --- score_round ---

func TestScoreRoundDir(t *testing.T) {
	dir := t.TempDir()
	roundDir := filepath.Join(dir, "round1")
	itemsDir := filepath.Join(roundDir, "items")
	os.MkdirAll(itemsDir, 0o755)

	os.WriteFile(filepath.Join(roundDir, "notation.txt"), []byte("E = entity"), 0o644)

	items := []ItemResult{
		{ID: "a", NotationTokens: 10, ExpandedTokens: 100, ExpansionRatio: 10, RoundtripScore: 0.9},
		{ID: "b", NotationTokens: 20, ExpandedTokens: 100, ExpansionRatio: 5, RoundtripScore: 0.7},
		{ID: "c", SkippedReason: "compress: timeout"},
	}
	for _, it := range items {
		raw, _ := json.Marshal(it)
		os.WriteFile(filepath.Join(itemsDir, it.ID+".json"), raw, 0o644)
	}

	m, err := ScoreRoundDir(roundDir)
	if err != nil {
		t.Fatal(err)
	}
	if m.RoundNumber != 1 {
		t.Errorf("round_number: got %d want 1", m.RoundNumber)
	}
	if m.ItemCount != 3 {
		t.Errorf("item_count: got %d want 3", m.ItemCount)
	}
	if m.SkippedCount != 1 {
		t.Errorf("skipped_count: got %d want 1", m.SkippedCount)
	}
	if m.MeanExpansionRatio != 7.5 {
		t.Errorf("mean ER: got %v want 7.5", m.MeanExpansionRatio)
	}
	wantCorr := (0.9 + 0.7) / 2
	if !nearly(m.MeanCorrectness, wantCorr) {
		t.Errorf("mean correctness: got %v want %v", m.MeanCorrectness, wantCorr)
	}
	if !strings.Contains(m.Summary, "Round 1") {
		t.Errorf("summary missing round number: %q", m.Summary)
	}
}

func TestScoreRoundDir_BadPath(t *testing.T) {
	_, err := ScoreRoundDir("/nope/round1")
	if err == nil {
		t.Errorf("expected error for missing dir")
	}
}

func TestScoreRoundDir_BadName(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "not-a-round")
	os.MkdirAll(filepath.Join(bad, "items"), 0o755)
	_, err := ScoreRoundDir(bad)
	if err == nil || !strings.Contains(err.Error(), "round number") {
		t.Errorf("expected round-number error, got %v", err)
	}
}

// --- run_round (full flow with fake) ---

func TestRunRound_FullFlowWithFake(t *testing.T) {
	dir := t.TempDir()
	corpus := filepath.Join(dir, "c.jsonl")
	os.WriteFile(corpus, []byte(`{"id":"a","spec_text":"alpha entity"}
{"id":"b","spec_text":"beta entity with more content"}
`), 0o644)

	// Order of LLM calls in run_round (round 1, 2 items):
	// 1. design  → <notation_spec>
	// 2. compress(a) → <notation>
	// 3. expand(a)   → text
	// 4. evaluate(a) → JSON
	// 5. compress(b) → <notation>
	// 6. expand(b)   → text
	// 7. evaluate(b) → JSON
	llm := &FakeLLMClient{Responses: []string{
		"<notation_spec>E = entity, F = field</notation_spec>",
		"<notation>Ea</notation>",
		"expanded code for a",
		`{"overall": 0.9, "reasoning": "good"}`,
		"<notation>Eb[F:F]</notation>",
		"expanded code for b which is longer",
		`{"overall": 0.8, "reasoning": "ok"}`,
	}}

	out := filepath.Join(dir, "out")
	m, err := RunRound(context.Background(), llm, RunRoundParams{
		CorpusPath:  corpus,
		Target:      "TypeScript implementation",
		Rubric:      "score it",
		Model:       "claude-test",
		OutputDir:   out,
		RoundNumber: 1,
	})
	if err != nil {
		t.Fatalf("run_round: %v", err)
	}
	if m.ItemCount != 2 || m.SkippedCount != 0 {
		t.Errorf("counts: got items=%d skipped=%d", m.ItemCount, m.SkippedCount)
	}
	if !nearly(m.MeanCorrectness, 0.85) {
		t.Errorf("mean correctness: got %v want 0.85", m.MeanCorrectness)
	}
	if len(llm.CallLog) != 7 {
		t.Errorf("expected 7 LLM calls, got %d", len(llm.CallLog))
	}

	// Check the on-disk artifacts.
	if _, err := os.Stat(filepath.Join(out, "round1", "notation.txt")); err != nil {
		t.Errorf("notation.txt not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "round1", "items", "a.json")); err != nil {
		t.Errorf("a.json not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "round1", "metrics.json")); err != nil {
		t.Errorf("metrics.json not written: %v", err)
	}
}

func TestRunRound_ResumesFromPartial(t *testing.T) {
	dir := t.TempDir()
	corpus := filepath.Join(dir, "c.jsonl")
	os.WriteFile(corpus, []byte(`{"id":"a","spec_text":"alpha"}
{"id":"b","spec_text":"beta"}
`), 0o644)

	out := filepath.Join(dir, "out")
	roundDir := filepath.Join(out, "round1")
	os.MkdirAll(filepath.Join(roundDir, "items"), 0o755)

	// Pre-write notation.txt and item a.json; only b should re-run.
	os.WriteFile(filepath.Join(roundDir, "notation.txt"), []byte("PRE-EXISTING SPEC"), 0o644)
	preItem := ItemResult{
		ID: "a", NotationText: "Ea", NotationTokens: 5,
		ExpandedText: "old expansion", ExpandedTokens: 50,
		ExpansionRatio: 10, RoundtripScore: 0.95,
	}
	raw, _ := json.Marshal(preItem)
	os.WriteFile(filepath.Join(roundDir, "items", "a.json"), raw, 0o644)

	// Only 3 calls expected: compress(b), expand(b), evaluate(b).
	llm := &FakeLLMClient{Responses: []string{
		"<notation>Eb</notation>",
		"new expansion for b",
		`{"overall": 0.6, "reasoning": "weaker"}`,
	}}
	m, err := RunRound(context.Background(), llm, RunRoundParams{
		CorpusPath:  corpus,
		Target:      "T",
		OutputDir:   out,
		RoundNumber: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(llm.CallLog) != 3 {
		t.Errorf("expected 3 calls (resume), got %d", len(llm.CallLog))
	}
	if m.ItemCount != 2 {
		t.Errorf("items: got %d want 2", m.ItemCount)
	}
	// a's pre-existing score (0.95) and b's new score (0.6) → mean 0.775
	if !nearly(m.MeanCorrectness, 0.775) {
		t.Errorf("mean correctness: got %v want 0.775", m.MeanCorrectness)
	}
}

// --- training pairs ---

func TestExtractTrainingPairs_FiltersByScore(t *testing.T) {
	dir := t.TempDir()
	roundDir := filepath.Join(dir, "round1")
	itemsDir := filepath.Join(roundDir, "items")
	os.MkdirAll(itemsDir, 0o755)
	os.WriteFile(filepath.Join(roundDir, "notation.txt"), []byte("S"), 0o644)

	for _, it := range []ItemResult{
		{ID: "high", NotationText: "n", ExpandedText: "e", RoundtripScore: 0.9, NotationTokens: 1, ExpandedTokens: 10},
		{ID: "low", NotationText: "n", ExpandedText: "e", RoundtripScore: 0.4, NotationTokens: 1, ExpandedTokens: 10},
		{ID: "skipped", SkippedReason: "x"},
	} {
		raw, _ := json.Marshal(it)
		os.WriteFile(filepath.Join(itemsDir, it.ID+".json"), raw, 0o644)
	}

	pairs, err := ExtractTrainingPairs(roundDir, 0.7)
	if err != nil {
		t.Fatal(err)
	}
	if len(pairs) != 1 || pairs[0].ID != "high" {
		t.Fatalf("expected just 'high' above threshold, got %+v", pairs)
	}
}
