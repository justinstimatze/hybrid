package main

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return &Store{Path: filepath.Join(t.TempDir(), "calibration.jsonl")}
}

func mustAppend(t *testing.T, s *Store, ev any) {
	t.Helper()
	if err := s.Append(ev); err != nil {
		t.Fatalf("append: %v", err)
	}
}

func TestEmptyFoldReturnsEmptyMap(t *testing.T) {
	s := newTestStore(t)
	state, err := s.Fold()
	if err != nil {
		t.Fatalf("fold: %v", err)
	}
	if len(state) != 0 {
		t.Errorf("expected 0 records, got %d", len(state))
	}
}

func TestPredictThenFold(t *testing.T) {
	s := newTestStore(t)
	mustAppend(t, s, PredictEvent{
		Event: "predict", PredictionID: "p1",
		TS: 1000, Loop: "test", LensOrReasoner: "reasoner",
		InputHash: "h", Prediction: map[string]any{"x": 1.0},
		ModelID: "m", SchemaVersion: 1, VerdictDueBy: 2000,
	})
	state, err := s.Fold()
	if err != nil {
		t.Fatalf("fold: %v", err)
	}
	if len(state) != 1 {
		t.Fatalf("expected 1 record, got %d", len(state))
	}
	r := state["p1"]
	if r.Loop != "test" {
		t.Errorf("unexpected loop: %q", r.Loop)
	}
	if r.Verdict != "" {
		t.Errorf("expected unresolved, got verdict=%q", r.Verdict)
	}
}

func TestPredictThenResolve(t *testing.T) {
	s := newTestStore(t)
	mustAppend(t, s, PredictEvent{Event: "predict", PredictionID: "p1", Loop: "test"})
	mustAppend(t, s, ResolveEvent{
		Event: "resolve", PredictionID: "p1",
		Verdict: "confirmed", VerdictSource: "manual", VerdictTS: 5000,
	})
	state, err := s.Fold()
	if err != nil {
		t.Fatalf("fold: %v", err)
	}
	r := state["p1"]
	if r.Verdict != "confirmed" {
		t.Errorf("expected confirmed, got %q", r.Verdict)
	}
	if r.VerdictTS != 5000 {
		t.Errorf("expected verdict_ts 5000, got %d", r.VerdictTS)
	}
}

func TestResolveOrphanedSkipped(t *testing.T) {
	s := newTestStore(t)
	mustAppend(t, s, ResolveEvent{Event: "resolve", PredictionID: "ghost", Verdict: "confirmed"})
	state, _ := s.Fold()
	if _, ok := state["ghost"]; ok {
		t.Errorf("orphaned resolve should not create a record")
	}
}

func TestMalformedLineSkipped(t *testing.T) {
	s := newTestStore(t)
	mustAppend(t, s, PredictEvent{Event: "predict", PredictionID: "p1", Loop: "test"})
	f, err := os.OpenFile(s.Path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if _, err := f.Write([]byte("not json\n")); err != nil {
		t.Fatalf("write: %v", err)
	}
	f.Close()
	state, err := s.Fold()
	if err != nil {
		t.Fatalf("fold should not error on malformed: %v", err)
	}
	if len(state) != 1 {
		t.Errorf("expected 1, got %d", len(state))
	}
}

func TestMultiplePredictsAndResolves(t *testing.T) {
	s := newTestStore(t)
	for _, id := range []string{"p1", "p2", "p3"} {
		mustAppend(t, s, PredictEvent{Event: "predict", PredictionID: id, Loop: "test"})
	}
	mustAppend(t, s, ResolveEvent{Event: "resolve", PredictionID: "p1", Verdict: "confirmed", VerdictTS: 100})
	mustAppend(t, s, ResolveEvent{Event: "resolve", PredictionID: "p2", Verdict: "refuted", VerdictTS: 200})
	state, _ := s.Fold()
	if len(state) != 3 {
		t.Fatalf("expected 3 records, got %d", len(state))
	}
	if state["p1"].Verdict != "confirmed" {
		t.Errorf("p1: %+v", state["p1"])
	}
	if state["p2"].Verdict != "refuted" {
		t.Errorf("p2: %+v", state["p2"])
	}
	if state["p3"].Verdict != "" {
		t.Errorf("p3 should be unresolved: %+v", state["p3"])
	}
}

func TestMultipleResolveEventsLastWins(t *testing.T) {
	s := newTestStore(t)
	mustAppend(t, s, PredictEvent{Event: "predict", PredictionID: "p1", Loop: "test"})
	mustAppend(t, s, ResolveEvent{Event: "resolve", PredictionID: "p1", Verdict: "confirmed", VerdictTS: 100})
	mustAppend(t, s, ResolveEvent{Event: "resolve", PredictionID: "p1", Verdict: "refuted", VerdictTS: 200})
	state, _ := s.Fold()
	r := state["p1"]
	if r.Verdict != "refuted" || r.VerdictTS != 200 {
		t.Errorf("expected last-write-wins refuted/200, got %s/%d", r.Verdict, r.VerdictTS)
	}
}
