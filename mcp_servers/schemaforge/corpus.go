package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
)

// CorpusItem is one record in a corpus JSONL file.
//
// `id` is a unique identifier (used as filename in the per-round output dir).
// `spec_text` is the freeform input the notation should compress and the
// expansion is scored against. Domain is opaque to schemaforge: a CRUD app
// description, a behavioral mechanism narrative, a humor-template instance —
// any text the caller wants to find a dense schema for.
type CorpusItem struct {
	ID       string `json:"id"`
	SpecText string `json:"spec_text"`
}

// LoadCorpus reads a JSONL file. Lines with empty id are rejected. Duplicate ids
// are rejected (would clobber per-item output files).
func LoadCorpus(path string) ([]CorpusItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open corpus: %w", err)
	}
	defer f.Close()

	var items []CorpusItem
	seen := map[string]bool{}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1<<20), 1<<24) // up to 16MB per line
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var item CorpusItem
		if err := json.Unmarshal(line, &item); err != nil {
			return nil, fmt.Errorf("parse line %d: %w", lineNo, err)
		}
		if item.ID == "" {
			return nil, fmt.Errorf("line %d: id is required", lineNo)
		}
		if seen[item.ID] {
			return nil, fmt.Errorf("line %d: duplicate id %q", lineNo, item.ID)
		}
		seen[item.ID] = true
		items = append(items, item)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan corpus: %w", err)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("corpus is empty")
	}
	return items, nil
}

// SeedItems returns the (up to) k items most likely to stress-test a notation.
// Heuristic: longest spec_text. A notation that handles long specs handles
// short ones; not vice versa.
func SeedItems(items []CorpusItem, k int) []CorpusItem {
	if k <= 0 || len(items) == 0 {
		return nil
	}
	if k >= len(items) {
		k = len(items)
	}
	cp := make([]CorpusItem, len(items))
	copy(cp, items)
	sort.SliceStable(cp, func(i, j int) bool {
		return len(cp[i].SpecText) > len(cp[j].SpecText)
	})
	return cp[:k]
}

var wordRe = regexp.MustCompile(`\S+`)

// EstimateTokens is a cheap, deterministic token-count proxy. Exact counts
// come from API usage stats when available; this is for ER metrics and
// per-round comparisons where consistency matters more than precision.
func EstimateTokens(text string) int {
	matches := wordRe.FindAllString(text, -1)
	if len(matches) == 0 {
		return 0
	}
	n := int(float64(len(matches)) * 1.3)
	if n < 1 {
		return 1
	}
	return n
}
