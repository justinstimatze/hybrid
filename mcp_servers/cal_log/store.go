package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// PredictEvent is appended when a typed evaluator records a prediction.
type PredictEvent struct {
	Event          string         `json:"event"`
	PredictionID   string         `json:"prediction_id"`
	TS             int64          `json:"ts"`
	Loop           string         `json:"loop"`
	LensOrReasoner string         `json:"lens_or_reasoner"`
	InputHash      string         `json:"input_hash"`
	Prediction     map[string]any `json:"prediction"`
	ModelID        string         `json:"model_id"`
	SchemaVersion  int            `json:"schema_version"`
	VerdictDueBy   int64          `json:"verdict_due_by"`
}

// ResolveEvent is appended when the verdict for a prediction is determined.
type ResolveEvent struct {
	Event         string `json:"event"`
	PredictionID  string `json:"prediction_id"`
	TS            int64  `json:"ts"`
	Verdict       string `json:"verdict"`
	VerdictSource string `json:"verdict_source"`
	VerdictTS     int64  `json:"verdict_ts"`
}

// Record is the folded current state of a prediction.
type Record struct {
	PredictionID   string         `json:"prediction_id"`
	TS             int64          `json:"ts"`
	Loop           string         `json:"loop"`
	LensOrReasoner string         `json:"lens_or_reasoner"`
	InputHash      string         `json:"input_hash"`
	Prediction     map[string]any `json:"prediction"`
	ModelID        string         `json:"model_id"`
	SchemaVersion  int            `json:"schema_version"`
	VerdictDueBy   int64          `json:"verdict_due_by"`
	Verdict        string         `json:"verdict,omitempty"`
	VerdictSource  string         `json:"verdict_source,omitempty"`
	VerdictTS      int64          `json:"verdict_ts,omitempty"`
}

// Store is an append-only event log on disk.
type Store struct {
	Path string
	mu   sync.Mutex
}

// Append adds one event to the log.
func (s *Store) Append(event any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(parentDir(s.Path), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	f, err := os.OpenFile(s.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %s: %w", s.Path, err)
	}
	defer f.Close()
	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if _, err := f.Write(append(b, '\n')); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// Fold reads the entire log and returns the current state per prediction_id.
// Malformed lines are skipped silently. Resolve events for unknown prediction_ids
// are skipped (orphaned).
func (s *Store) Fold() (map[string]*Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := map[string]*Record{}
	f, err := os.Open(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return state, nil
		}
		return nil, fmt.Errorf("open %s: %w", s.Path, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 64*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var probe struct {
			Event        string `json:"event"`
			PredictionID string `json:"prediction_id"`
		}
		if err := json.Unmarshal(line, &probe); err != nil {
			continue
		}
		switch probe.Event {
		case "predict":
			var p PredictEvent
			if err := json.Unmarshal(line, &p); err != nil {
				continue
			}
			state[p.PredictionID] = &Record{
				PredictionID:   p.PredictionID,
				TS:             p.TS,
				Loop:           p.Loop,
				LensOrReasoner: p.LensOrReasoner,
				InputHash:      p.InputHash,
				Prediction:     p.Prediction,
				ModelID:        p.ModelID,
				SchemaVersion:  p.SchemaVersion,
				VerdictDueBy:   p.VerdictDueBy,
			}
		case "resolve":
			r, ok := state[probe.PredictionID]
			if !ok {
				continue
			}
			var rev ResolveEvent
			if err := json.Unmarshal(line, &rev); err != nil {
				continue
			}
			r.Verdict = rev.Verdict
			r.VerdictSource = rev.VerdictSource
			r.VerdictTS = rev.VerdictTS
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return state, nil
}

func parentDir(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[:i]
		}
	}
	return "."
}
