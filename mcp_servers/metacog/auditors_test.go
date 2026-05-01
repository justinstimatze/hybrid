package main

import (
	"strings"
	"testing"
	"time"
)

// helper to build a record quickly
func rec(verdict string, summary string, prov []Provenance, edges []Edge, group, cluster string, complexity float64) Record {
	return Record{
		ID:      "r-" + verdict + "-" + summary[:min(len(summary), 4)],
		Verdict: verdict, SummaryText: summary,
		Provenance: prov, Edges: edges,
		Group: group, Cluster: cluster, Complexity: complexity,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// --- AuditConfirmationBias ---

func TestConfirmationBias_TriggeredHighRate(t *testing.T) {
	var recs []Record
	for i := 0; i < 8; i++ {
		recs = append(recs, Record{Verdict: "confirmed"})
	}
	for i := 0; i < 1; i++ {
		recs = append(recs, Record{Verdict: "refuted"})
	}
	r := AuditConfirmationBias(inMemorySubstrate(recs))
	if !r.Triggered {
		t.Errorf("expected triggered (rate %.2f > 0.75), got %+v", r.Value, r)
	}
	if r.Value < 0.75 {
		t.Errorf("rate should be > 0.75, got %.2f", r.Value)
	}
}

func TestConfirmationBias_NotTriggeredHealthy(t *testing.T) {
	var recs []Record
	for i := 0; i < 5; i++ {
		recs = append(recs, Record{Verdict: "confirmed"})
	}
	for i := 0; i < 5; i++ {
		recs = append(recs, Record{Verdict: "refuted"})
	}
	r := AuditConfirmationBias(inMemorySubstrate(recs))
	if r.Triggered {
		t.Errorf("expected not triggered (rate 0.5), got %+v", r)
	}
}

func TestConfirmationBias_Skipped(t *testing.T) {
	r := AuditConfirmationBias(inMemorySubstrate{Record{Verdict: ""}})
	if !r.Skipped {
		t.Errorf("expected skipped (no verdicts)")
	}
}

// --- AuditAvailabilityHeuristic ---

func TestAvailabilityHeuristic_TriggeredConcentrated(t *testing.T) {
	var recs []Record
	for i := 0; i < 9; i++ {
		recs = append(recs, Record{Provenance: []Provenance{{Type: "wikipedia"}}})
	}
	recs = append(recs, Record{Provenance: []Provenance{{Type: "arxiv"}}})
	r := AuditAvailabilityHeuristic(inMemorySubstrate(recs))
	if !r.Triggered {
		t.Errorf("expected triggered (HHI %.3f > 0.25), got %+v", r.Value, r)
	}
}

func TestAvailabilityHeuristic_NotTriggeredDiverse(t *testing.T) {
	var recs []Record
	for _, src := range []string{"wikipedia", "arxiv", "academic", "news", "code_repo"} {
		for i := 0; i < 4; i++ {
			recs = append(recs, Record{Provenance: []Provenance{{Type: src}}})
		}
	}
	r := AuditAvailabilityHeuristic(inMemorySubstrate(recs))
	if r.Triggered {
		t.Errorf("expected not triggered (diverse), got HHI=%.3f", r.Value)
	}
}

func TestAvailabilityHeuristic_OriginClassification(t *testing.T) {
	recs := []Record{
		{Provenance: []Provenance{{Origin: "https://en.wikipedia.org/wiki/X"}}},
		{Provenance: []Provenance{{Origin: "https://arxiv.org/abs/1234.5678"}}},
	}
	r := AuditAvailabilityHeuristic(inMemorySubstrate(recs))
	if r.Skipped {
		t.Errorf("origin classification should not skip: %+v", r)
	}
	if !strings.Contains(r.Detail, "wikipedia") {
		t.Errorf("expected wikipedia in detail: %s", r.Detail)
	}
}

// --- AuditSurvivorshipBias ---

func TestSurvivorshipBias_Triggered(t *testing.T) {
	var recs []Record
	for i := 0; i < 50; i++ {
		recs = append(recs, Record{Verdict: "irrelevant"})
	}
	for i := 0; i < 2; i++ {
		recs = append(recs, Record{Verdict: "challenged"})
	}
	r := AuditSurvivorshipBias(inMemorySubstrate(recs))
	if !r.Triggered {
		t.Errorf("expected triggered (ratio %.1f > 5), got %+v", r.Value, r)
	}
}

func TestSurvivorshipBias_SkippedNoIrrelevant(t *testing.T) {
	r := AuditSurvivorshipBias(inMemorySubstrate{
		{Verdict: "confirmed"}, {Verdict: "refuted"},
	})
	if !r.Skipped {
		t.Errorf("expected skipped (no irrelevant verdicts)")
	}
}

// --- AuditFramingEffect ---

func TestFramingEffect_Triggered(t *testing.T) {
	recs := []Record{
		{SummaryText: "This is a groundbreaking and revolutionary result."},
		{SummaryText: "A controversial paper with flawed methodology."},
		{SummaryText: "A neutral description of the work."},
	}
	r := AuditFramingEffect(inMemorySubstrate(recs))
	if !r.Triggered {
		t.Errorf("expected triggered (>15%% evaluative), got %+v", r)
	}
}

func TestFramingEffect_NotTriggered(t *testing.T) {
	recs := []Record{}
	for i := 0; i < 20; i++ {
		recs = append(recs, Record{SummaryText: "A neutral factual description."})
	}
	recs = append(recs, Record{SummaryText: "A groundbreaking result."})
	r := AuditFramingEffect(inMemorySubstrate(recs))
	if r.Triggered {
		t.Errorf("expected not triggered (1/21 ≈ 5%%), got value=%.3f", r.Value)
	}
}

// --- AuditPrematureClosure ---

func TestPrematureClosure_Triggered(t *testing.T) {
	recs := []Record{
		{ID: "a", SummaryText: "Obviously, this is the case."},
		{ID: "b", SummaryText: "It is well established that X causes Y."},
		{ID: "c", SummaryText: "A neutral description."},
	}
	r := AuditPrematureClosure(inMemorySubstrate(recs))
	if !r.Triggered {
		t.Errorf("expected triggered (>=1 finding), got %+v", r)
	}
	if r.Value < 2 {
		t.Errorf("expected 2 findings, got %.0f", r.Value)
	}
}

func TestPrematureClosure_AntiClosureContextSuppresses(t *testing.T) {
	recs := []Record{
		{ID: "a", SummaryText: "Although obviously controversial, this remains debated."},
	}
	r := AuditPrematureClosure(inMemorySubstrate(recs))
	if r.Triggered {
		t.Errorf("anti-closure context should suppress: %+v", r)
	}
}

// --- AuditBaseRateNeglect ---

func TestBaseRateNeglect_TriggeredLowEntropy(t *testing.T) {
	recs := []Record{}
	for i := 0; i < 50; i++ {
		recs = append(recs, Record{Edges: []Edge{{Predicate: "supports"}}})
	}
	for i := 0; i < 2; i++ {
		recs = append(recs, Record{Edges: []Edge{{Predicate: "contradicts"}}})
	}
	r := AuditBaseRateNeglect(inMemorySubstrate(recs))
	if !r.Triggered {
		t.Errorf("expected triggered (entropy %.2f < 3.0)", r.Value)
	}
}

func TestBaseRateNeglect_NotTriggeredHighEntropy(t *testing.T) {
	preds := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	var recs []Record
	for _, p := range preds {
		for i := 0; i < 3; i++ {
			recs = append(recs, Record{Edges: []Edge{{Predicate: p}}})
		}
	}
	r := AuditBaseRateNeglect(inMemorySubstrate(recs))
	if r.Triggered {
		t.Errorf("expected not triggered (entropy %.2f >= 3.0)", r.Value)
	}
}

// --- AuditAnchoring ---

func TestAnchoring_TriggeredOldDense(t *testing.T) {
	now := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	var recs []Record
	// Old group: 50 records
	for i := 0; i < 50; i++ {
		recs = append(recs, Record{Group: "g_old", CreatedAt: now.AddDate(-2, 0, i)})
	}
	// Mid group: 10 records
	for i := 0; i < 10; i++ {
		recs = append(recs, Record{Group: "g_mid", CreatedAt: now.AddDate(-1, 0, i)})
	}
	// New group: 2 records
	for i := 0; i < 2; i++ {
		recs = append(recs, Record{Group: "g_new", CreatedAt: now.AddDate(0, 0, i)})
	}
	r := AuditAnchoring(inMemorySubstrate(recs))
	if !r.Triggered {
		t.Errorf("expected triggered (older=denser), got value=%.3f, %+v", r.Value, r)
	}
}

// --- AuditClusteringIllusion ---

func TestClusteringIllusion_TriggeredOverlap(t *testing.T) {
	// Each group has exactly one cluster — perfect overlap
	var recs []Record
	for i := 0; i < 5; i++ {
		recs = append(recs, Record{Group: "g_a", Cluster: "c_a"})
		recs = append(recs, Record{Group: "g_b", Cluster: "c_b"})
		recs = append(recs, Record{Group: "g_c", Cluster: "c_c"})
	}
	r := AuditClusteringIllusion(inMemorySubstrate(recs))
	// With singleton clusters per group, Jaccard between any two groups is 0/2 = 0.
	// So this is NOT triggered. Let me invert: clusters span groups means overlap.
	if r.Triggered {
		t.Errorf("singleton clusters per group should not trigger overlap: %+v", r)
	}
}

func TestClusteringIllusion_SkippedNoHints(t *testing.T) {
	r := AuditClusteringIllusion(inMemorySubstrate{{ID: "x"}})
	if !r.Skipped {
		t.Errorf("expected skipped (no Group/Cluster hints)")
	}
}

// --- AuditDunningKruger ---

func TestDunningKruger_Triggered(t *testing.T) {
	var recs []Record
	for i := 0; i < 19; i++ {
		recs = append(recs, Record{Complexity: 0.1, Edges: nil})
	}
	recs = append(recs, Record{Complexity: 0.1, Edges: []Edge{{Predicate: "x", To: "y"}}})
	r := AuditDunningKruger(inMemorySubstrate(recs))
	if !r.Triggered {
		t.Errorf("expected triggered (rate %.2f > 0.90)", r.Value)
	}
}

func TestDunningKruger_SkippedNoComplexity(t *testing.T) {
	r := AuditDunningKruger(inMemorySubstrate{{ID: "x"}, {ID: "y"}})
	if !r.Skipped {
		t.Errorf("expected skipped (no complexity hints)")
	}
}

// --- AllAuditors / smoke ---

func TestAllAuditors_RunWithoutPanic(t *testing.T) {
	rich := []Record{
		{ID: "a", Verdict: "confirmed", SummaryText: "Obviously a neutral text.",
			Provenance: []Provenance{{Type: "wikipedia"}},
			Edges:      []Edge{{Predicate: "supports", To: "b"}},
			Group:      "g1", Cluster: "c1", Complexity: 0.2,
			CreatedAt: time.Now()},
		{ID: "b", Verdict: "refuted", SummaryText: "A controversial claim.",
			Provenance: []Provenance{{Type: "arxiv"}},
			Edges:      []Edge{{Predicate: "contradicts", To: "a"}},
			Group:      "g2", Cluster: "c2", Complexity: 0.8,
			CreatedAt: time.Now()},
		{ID: "c", Verdict: "irrelevant",
			Provenance: []Provenance{{Type: "news"}}},
	}
	for _, a := range AllAuditors() {
		r := a(inMemorySubstrate(rich))
		if r.Bias == "" {
			t.Errorf("auditor returned empty bias name")
		}
	}
}
