package main

import (
	"bufio"
	"encoding/json"
	"os"
	"time"
)

// JSONLSubstrate reads records from a JSONL file. One JSON object per line.
// Empty lines and JSON parse errors are skipped silently (the latter so that
// a malformed line doesn't kill the audit; metacog flags this in stats).
type JSONLSubstrate struct {
	Path    string
	skipped int // count of lines that failed to parse
}

// jsonRecord is the on-disk schema. Times are RFC3339 strings.
type jsonRecord struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Fields      map[string]any `json:"fields"`
	Provenance  []Provenance   `json:"provenance"`
	Edges       []Edge         `json:"edges"`
	CreatedAt   string         `json:"created_at"`
	Verdict     string         `json:"verdict"`
	VerdictTime string         `json:"verdict_time"`
	SummaryText string         `json:"summary_text"`
	Complexity  float64        `json:"complexity"`
	Cluster     string         `json:"cluster"`
	Group       string         `json:"group"`
}

func (j *JSONLSubstrate) Records() []Record {
	f, err := os.Open(j.Path)
	if err != nil {
		return nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 64*1024*1024)
	var out []Record
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var jr jsonRecord
		if err := json.Unmarshal(line, &jr); err != nil {
			j.skipped++
			continue
		}
		r := Record{
			ID:          jr.ID,
			Type:        jr.Type,
			Fields:      jr.Fields,
			Provenance:  jr.Provenance,
			Edges:       jr.Edges,
			Verdict:     jr.Verdict,
			SummaryText: jr.SummaryText,
			Complexity:  jr.Complexity,
			Cluster:     jr.Cluster,
			Group:       jr.Group,
		}
		if jr.CreatedAt != "" {
			if t, err := time.Parse(time.RFC3339, jr.CreatedAt); err == nil {
				r.CreatedAt = t
			}
		}
		if jr.VerdictTime != "" {
			if t, err := time.Parse(time.RFC3339, jr.VerdictTime); err == nil {
				r.VerdictTime = t
			}
		}
		out = append(out, r)
	}
	return out
}

// SkippedLines reports how many lines failed to parse during the last Records() read.
func (j *JSONLSubstrate) SkippedLines() int { return j.skipped }

// inMemorySubstrate is a slice of Records that satisfies Substrate. Used in tests.
type inMemorySubstrate []Record

func (s inMemorySubstrate) Records() []Record { return []Record(s) }
