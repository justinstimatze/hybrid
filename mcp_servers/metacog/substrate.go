package main

import "time"

// Substrate exposes typed records from any backing store. Implement this
// interface to make any typed substrate auditable by metacog.
type Substrate interface {
	Records() []Record
}

// Record is the abstracted unit of analysis. Auditors operate on slices of
// these. Optional fields (CreatedAt, Verdict, Complexity, etc.) cause
// graceful skips in the auditors that need them; they don't have to be
// populated for the substrate to be auditable on the metrics they don't
// require.
type Record struct {
	ID          string
	Type        string
	Fields      map[string]any
	Provenance  []Provenance
	Edges       []Edge
	CreatedAt   time.Time
	Verdict     string
	VerdictTime time.Time
	SummaryText string

	// Optional hints. Auditors degrade gracefully if absent.
	Complexity float64 // pre-computed complexity score, 0-1 (used by dunning_kruger)
	Cluster    string  // pre-computed cluster id (used by clustering_illusion)
	Group      string  // file or coarse grouping (used by clustering_illusion fallback, anchoring)
}

// Provenance describes a citation or source for a record.
type Provenance struct {
	Origin string // URL, citation, or source label
	Type   string // optional pre-classified type (e.g. "academic", "news", "primary")
	Quote  string // optional verbatim quote
}

// Edge is a relationship to another record.
type Edge struct {
	Predicate string // relationship type
	To        string // target record ID
}

// AuditorResult is the output of one bias auditor running on a substrate.
type AuditorResult struct {
	Bias       string  `json:"bias"`      // canonical identifier (e.g. "confirmation_bias")
	BiasName   string  `json:"bias_name"` // human-readable name
	Metric     string  `json:"metric"`    // what was measured
	Value      float64 `json:"value"`     // measured value
	Threshold  float64 `json:"threshold"` // alerting threshold
	Triggered  bool    `json:"triggered"` // value exceeds threshold
	Severity   string  `json:"severity,omitempty"`
	Detail     string  `json:"detail"`     // numeric breakdown
	Conclusion string  `json:"conclusion"` // what this means; what to do
	Skipped    bool    `json:"skipped,omitempty"`
	SkipReason string  `json:"skip_reason,omitempty"`
}

// Auditor is a function that audits a substrate for one bias signature.
type Auditor func(s Substrate) AuditorResult

// AllAuditors returns the canonical nine.
func AllAuditors() []Auditor {
	return []Auditor{
		AuditConfirmationBias,
		AuditAnchoring,
		AuditClusteringIllusion,
		AuditAvailabilityHeuristic,
		AuditSurvivorshipBias,
		AuditFramingEffect,
		AuditDunningKruger,
		AuditBaseRateNeglect,
		AuditPrematureClosure,
	}
}

// verdictCategory normalizes a verdict string into one of:
// "corroborated", "challenged", "irrelevant", "partial", "unknown".
// This lets metacog work across substrates with different verdict vocabularies.
func verdictCategory(v string) string {
	switch v {
	case "confirmed", "corroborated", "supported", "validated", "verified":
		return "corroborated"
	case "refuted", "challenged", "contradicted", "rejected":
		return "challenged"
	case "irrelevant", "noise", "off_topic":
		return "irrelevant"
	case "partial", "mixed":
		return "partial"
	case "":
		return "unknown"
	default:
		return "unknown"
	}
}
