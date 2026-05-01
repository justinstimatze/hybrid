package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// AuditConfirmationBias — corroboration rate among resolved verdicts.
// Triggered if corroborated/(corroborated+challenged) > 0.75.
// Lifted from winze cmd/metabolism/dreamaudit.go.
func AuditConfirmationBias(s Substrate) AuditorResult {
	r := AuditorResult{
		Bias: "confirmation_bias", BiasName: "Confirmation bias",
		Metric: "corroboration_rate", Threshold: 0.75,
	}
	corroborated, challenged := 0, 0
	for _, rec := range s.Records() {
		switch verdictCategory(rec.Verdict) {
		case "corroborated":
			corroborated++
		case "challenged":
			challenged++
		}
	}
	total := corroborated + challenged
	if total == 0 {
		r.Skipped = true
		r.SkipReason = "no resolved records (need verdicts in {confirmed, refuted, ...})"
		return r
	}
	rate := float64(corroborated) / float64(total)
	r.Value = rate
	r.Detail = fmt.Sprintf("corroborated=%d, challenged=%d, rate=%.3f", corroborated, challenged, rate)
	if rate > r.Threshold {
		r.Triggered = true
		r.Severity = "info"
		r.Conclusion = fmt.Sprintf("Corroboration rate is %.0f%% — the substrate may be over-confirming. "+
			"Consider whether the lens is biased toward 'this looks valid' over 'this looks contradicted.'",
			rate*100)
	} else {
		r.Conclusion = "Corroboration rate is healthy."
	}
	return r
}

// AuditAnchoring — Spearman rank correlation between record age (older first) and
// claim density (records per group). High correlation suggests the substrate
// anchors on early records.
func AuditAnchoring(s Substrate) AuditorResult {
	r := AuditorResult{
		Bias: "anchoring", BiasName: "Anchoring",
		Metric: "spearman_age_density", Threshold: 0.5,
	}
	// Group records by Group field if present, otherwise by year-month of CreatedAt.
	groups := map[string][]Record{}
	for _, rec := range s.Records() {
		key := rec.Group
		if key == "" {
			if rec.CreatedAt.IsZero() {
				continue
			}
			key = rec.CreatedAt.Format("2006-01")
		}
		groups[key] = append(groups[key], rec)
	}
	if len(groups) < 3 {
		r.Skipped = true
		r.SkipReason = "need 3+ groups (Group field or per-month buckets) for rank correlation"
		return r
	}
	// Order groups by oldest creation time within group.
	type ge struct {
		key     string
		oldest  float64 // unix seconds
		density float64 // records in group
	}
	entries := make([]ge, 0, len(groups))
	for k, recs := range groups {
		var oldest float64 = math.MaxFloat64
		for _, rec := range recs {
			if !rec.CreatedAt.IsZero() {
				ts := float64(rec.CreatedAt.Unix())
				if ts < oldest {
					oldest = ts
				}
			}
		}
		if oldest == math.MaxFloat64 {
			oldest = 0
		}
		entries = append(entries, ge{key: k, oldest: oldest, density: float64(len(recs))})
	}
	x := make([]float64, len(entries))
	y := make([]float64, len(entries))
	for i, e := range entries {
		x[i] = e.oldest
		y[i] = e.density
	}
	// Anchoring direction: older groups (lower x) with higher density (higher y) is anchoring.
	// So we expect negative correlation between age-rank and density. We'll negate to make it positive.
	rho := -spearmanRho(x, y)
	r.Value = rho
	r.Detail = fmt.Sprintf("spearman(group_age, density) reversed: %.3f, n=%d groups", rho, len(entries))
	if rho > r.Threshold {
		r.Triggered = true
		r.Severity = "info"
		r.Conclusion = "Older groups have disproportionately higher density. The substrate may be anchored on early records; new records may be under-represented."
	} else {
		r.Conclusion = "No strong anchoring signal."
	}
	return r
}

// AuditClusteringIllusion — Jaccard overlap between Group sets and Cluster sets.
// High overlap means the system's clusters mirror its file-grouping (i.e. the
// clusters reflect storage layout rather than discovered topology).
func AuditClusteringIllusion(s Substrate) AuditorResult {
	r := AuditorResult{
		Bias: "clustering_illusion", BiasName: "Clustering illusion",
		Metric: "group_cluster_jaccard", Threshold: 0.7,
	}
	// Need both Group and Cluster.
	gToCs := map[string]map[string]bool{}
	cToGs := map[string]map[string]bool{}
	any := false
	for _, rec := range s.Records() {
		if rec.Group == "" || rec.Cluster == "" {
			continue
		}
		any = true
		if gToCs[rec.Group] == nil {
			gToCs[rec.Group] = map[string]bool{}
		}
		if cToGs[rec.Cluster] == nil {
			cToGs[rec.Cluster] = map[string]bool{}
		}
		gToCs[rec.Group][rec.Cluster] = true
		cToGs[rec.Cluster][rec.Group] = true
	}
	if !any {
		r.Skipped = true
		r.SkipReason = "no records have both Group and Cluster hints"
		return r
	}
	// Compute pairwise Jaccard between groups based on their cluster sets.
	keys := make([]string, 0, len(gToCs))
	for k := range gToCs {
		keys = append(keys, k)
	}
	if len(keys) < 2 {
		r.Skipped = true
		r.SkipReason = "need 2+ groups"
		return r
	}
	sort.Strings(keys)
	var sum float64
	var n int
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			a := gToCs[keys[i]]
			b := gToCs[keys[j]]
			inter, union := 0, 0
			seen := map[string]bool{}
			for c := range a {
				seen[c] = true
				if b[c] {
					inter++
				}
				union++
			}
			for c := range b {
				if !seen[c] {
					union++
				}
			}
			if union > 0 {
				sum += float64(inter) / float64(union)
				n++
			}
		}
	}
	avg := 0.0
	if n > 0 {
		avg = sum / float64(n)
	}
	r.Value = avg
	r.Detail = fmt.Sprintf("avg Jaccard(group_clusters) = %.3f over %d pairs", avg, n)
	if avg > r.Threshold {
		r.Triggered = true
		r.Severity = "info"
		r.Conclusion = "Group structure mirrors cluster structure. Clusters may be artifacts of storage layout rather than discovered topology."
	} else {
		r.Conclusion = "Clusters and groups appear structurally distinct."
	}
	return r
}

// AuditAvailabilityHeuristic — Herfindahl-Hirschman Index over provenance source types.
// HHI > 0.25 signals that the substrate over-represents one source.
func AuditAvailabilityHeuristic(s Substrate) AuditorResult {
	r := AuditorResult{
		Bias: "availability_heuristic", BiasName: "Availability heuristic",
		Metric: "provenance_hhi", Threshold: 0.25,
	}
	sourceTypes := map[string]int{}
	total := 0
	for _, rec := range s.Records() {
		for _, p := range rec.Provenance {
			t := p.Type
			if t == "" {
				t = classifyOrigin(p.Origin)
			}
			if t == "" {
				continue
			}
			sourceTypes[t]++
			total++
		}
	}
	if total == 0 {
		r.Skipped = true
		r.SkipReason = "no provenance records found"
		return r
	}
	var hhi float64
	parts := make([]string, 0, len(sourceTypes))
	for t, c := range sourceTypes {
		share := float64(c) / float64(total)
		hhi += share * share
		parts = append(parts, fmt.Sprintf("%s: %d (%.0f%%)", t, c, share*100))
	}
	sort.Strings(parts)
	r.Value = hhi
	r.Detail = fmt.Sprintf("HHI=%.3f across %d provenance records: %s", hhi, total, strings.Join(parts, ", "))
	if hhi > r.Threshold {
		var dominant string
		var maxShare float64
		for t, c := range sourceTypes {
			share := float64(c) / float64(total)
			if share > maxShare {
				maxShare = share
				dominant = t
			}
		}
		r.Triggered = true
		r.Severity = "info"
		r.Conclusion = fmt.Sprintf("Provenance is concentrated: %s = %.0f%%. The substrate may reflect source availability rather than domain importance.", dominant, maxShare*100)
	} else {
		r.Conclusion = "Provenance sources are reasonably diverse."
	}
	return r
}

// AuditSurvivorshipBias — irrelevant:challenged ratio.
// Triggered if the lens classifies many sources as irrelevant but very few as
// challenges; the pipeline may be filtering out valid dissent.
func AuditSurvivorshipBias(s Substrate) AuditorResult {
	r := AuditorResult{
		Bias: "survivorship_bias", BiasName: "Survivorship bias",
		Metric: "irrelevant_to_challenged_ratio", Threshold: 5.0,
	}
	irr, ch := 0, 0
	for _, rec := range s.Records() {
		switch verdictCategory(rec.Verdict) {
		case "irrelevant":
			irr++
		case "challenged":
			ch++
		}
	}
	if irr == 0 && ch == 0 {
		r.Skipped = true
		r.SkipReason = "no irrelevant or challenged verdicts (need verdict vocabulary)"
		return r
	}
	if irr == 0 {
		r.Skipped = true
		r.SkipReason = "no irrelevant verdicts (need 'irrelevant' verdict to detect)"
		return r
	}
	denom := ch
	if denom == 0 {
		denom = 1
	}
	ratio := float64(irr) / float64(denom)
	r.Value = ratio
	r.Detail = fmt.Sprintf("irrelevant=%d, challenged=%d, ratio=%.2f", irr, ch, ratio)
	if ratio > r.Threshold {
		r.Triggered = true
		r.Severity = "warning"
		r.Conclusion = "The pipeline classifies many sources as irrelevant but very few as challenges. May be filtering out valid dissent."
	} else {
		r.Conclusion = "Irrelevant:challenged ratio is healthy."
	}
	return r
}

// AuditFramingEffect — fraction of summaries containing evaluative framing words.
// Triggered if > 15%.
func AuditFramingEffect(s Substrate) AuditorResult {
	r := AuditorResult{
		Bias: "framing_effect", BiasName: "Framing effect",
		Metric: "evaluative_summary_fraction", Threshold: 0.15,
	}
	positiveFraming := []string{
		"groundbreaking", "seminal", "landmark", "revolutionary",
		"brilliant", "influential", "pioneering", "definitive",
		"canonical", "masterful", "profound", "celebrated",
	}
	negativeFraming := []string{
		"controversial", "flawed", "debunked", "discredited",
		"pseudoscientific", "misleading", "simplistic",
		"outdated", "questionable",
	}
	total := 0
	pos, neg := 0, 0
	for _, rec := range s.Records() {
		if rec.SummaryText == "" {
			continue
		}
		total++
		lower := strings.ToLower(rec.SummaryText)
		for _, term := range positiveFraming {
			if containsWordSimple(lower, term) {
				pos++
				break
			}
		}
		for _, term := range negativeFraming {
			if containsWordSimple(lower, term) {
				neg++
				break
			}
		}
	}
	if total == 0 {
		r.Skipped = true
		r.SkipReason = "no records with SummaryText"
		return r
	}
	frac := float64(pos+neg) / float64(total)
	r.Value = frac
	r.Detail = fmt.Sprintf("positive=%d, negative=%d, total_with_summary=%d, fraction=%.3f", pos, neg, total, frac)
	if frac > r.Threshold {
		r.Triggered = true
		r.Severity = "info"
		r.Conclusion = fmt.Sprintf("%.0f%% of summaries use evaluative framing. The substrate may be smuggling in judgment via word choice.", frac*100)
	} else {
		r.Conclusion = "Summaries are reasonably descriptive."
	}
	return r
}

// AuditDunningKruger — fraction of low-complexity records with zero edges.
// Generalized from winze: if many "simple" records have no relationships, the
// system may be treating simple things as complete-in-isolation when they
// likely have unmodeled connections.
// Skipped if no Complexity hints are present.
func AuditDunningKruger(s Substrate) AuditorResult {
	r := AuditorResult{
		Bias: "dunning_kruger", BiasName: "Dunning-Kruger effect",
		Metric: "low_complexity_zero_edge_rate", Threshold: 0.90,
	}
	hasComplexity := false
	lowTotal, lowZero := 0, 0
	for _, rec := range s.Records() {
		if rec.Complexity > 0 {
			hasComplexity = true
		}
		if rec.Complexity > 0 && rec.Complexity < 0.3 {
			lowTotal++
			if len(rec.Edges) == 0 {
				lowZero++
			}
		}
	}
	if !hasComplexity {
		r.Skipped = true
		r.SkipReason = "no records have Complexity hint set"
		return r
	}
	if lowTotal == 0 {
		r.Skipped = true
		r.SkipReason = "no low-complexity records (Complexity in (0, 0.3))"
		return r
	}
	rate := float64(lowZero) / float64(lowTotal)
	r.Value = rate
	r.Detail = fmt.Sprintf("low_complexity_records=%d, zero_edges=%d, rate=%.3f", lowTotal, lowZero, rate)
	if rate > r.Threshold {
		r.Triggered = true
		r.Severity = "info"
		r.Conclusion = "Most simple records have no edges. The substrate may be treating simple things as complete-in-isolation when they likely have unmodeled connections."
	} else {
		r.Conclusion = "Simple records have reasonable connection density."
	}
	return r
}

// AuditBaseRateNeglect — Shannon entropy of edge predicate distribution.
// Triggered if entropy < 3.0 bits (substrate uses too few predicate types).
func AuditBaseRateNeglect(s Substrate) AuditorResult {
	r := AuditorResult{
		Bias: "base_rate_neglect", BiasName: "Base rate neglect",
		Metric: "predicate_entropy_bits", Threshold: 3.0,
	}
	counts := map[string]int{}
	total := 0
	for _, rec := range s.Records() {
		for _, e := range rec.Edges {
			if e.Predicate == "" {
				continue
			}
			counts[e.Predicate]++
			total++
		}
	}
	if total == 0 {
		r.Skipped = true
		r.SkipReason = "no edges with predicates"
		return r
	}
	var h float64
	for _, c := range counts {
		p := float64(c) / float64(total)
		if p > 0 {
			h -= p * math.Log2(p)
		}
	}
	r.Value = h
	r.Detail = fmt.Sprintf("predicate_entropy=%.3f bits over %d edges, %d distinct predicates", h, total, len(counts))
	// Note: triggered when value is BELOW threshold (low entropy = concentration).
	if h < r.Threshold {
		r.Triggered = true
		r.Severity = "info"
		r.Conclusion = fmt.Sprintf("Predicate entropy is %.2f bits (threshold %.1f). The substrate uses too few relationship types; consider whether new edge predicates are needed.", h, r.Threshold)
	} else {
		r.Conclusion = "Predicate distribution is healthy."
	}
	return r
}

// AuditPrematureClosure — count of summaries containing thought-terminating
// clichés without anti-closure context. Any finding triggers (threshold 1.0).
func AuditPrematureClosure(s Substrate) AuditorResult {
	r := AuditorResult{
		Bias: "premature_closure", BiasName: "Premature closure",
		Metric: "closure_findings", Threshold: 1.0,
	}
	closurePhrases := []string{
		"it goes without saying",
		"obviously",
		"everyone knows",
		"it is well known",
		"it is well established",
		"there is no doubt",
		"undeniably",
		"unquestionably",
		"it is clear that",
		"needless to say",
		"of course",
		"self-evident",
		"beyond dispute",
		"universally accepted",
		"widely recognized",
		"generally agreed",
	}
	antiClosure := []string{
		"remains debated", "remains contested",
		"still debated", "still contested",
		"not universally", "despite", "although", "however",
	}
	total := 0
	findings := []string{}
	for _, rec := range s.Records() {
		if rec.SummaryText == "" {
			continue
		}
		total++
		lower := strings.ToLower(rec.SummaryText)
		// skip if anti-closure context is also present
		hasAnti := false
		for _, ac := range antiClosure {
			if strings.Contains(lower, ac) {
				hasAnti = true
				break
			}
		}
		if hasAnti {
			continue
		}
		for _, p := range closurePhrases {
			if strings.Contains(lower, p) {
				if rec.ID != "" {
					findings = append(findings, fmt.Sprintf("%s: %q", rec.ID, p))
				} else {
					findings = append(findings, fmt.Sprintf("(no id): %q", p))
				}
				break
			}
		}
	}
	if total == 0 {
		r.Skipped = true
		r.SkipReason = "no records with SummaryText"
		return r
	}
	r.Value = float64(len(findings))
	exampleN := len(findings)
	if exampleN > 5 {
		exampleN = 5
	}
	r.Detail = fmt.Sprintf("findings=%d across %d summaries; examples: %s",
		len(findings), total, strings.Join(findings[:exampleN], "; "))
	if r.Value >= r.Threshold {
		r.Triggered = true
		r.Severity = "info"
		r.Conclusion = "Found thought-terminating clichés in record summaries. Consider whether these reflect genuine consensus or premature closure."
	} else {
		r.Conclusion = "No closure clichés found."
	}
	return r
}

// classifyOrigin maps a provenance origin URL/string into a coarse source-type bucket.
func classifyOrigin(origin string) string {
	o := strings.ToLower(origin)
	switch {
	case strings.Contains(o, "wikipedia.org"):
		return "wikipedia"
	case strings.Contains(o, "arxiv.org"):
		return "arxiv"
	case strings.Contains(o, "doi.org") ||
		strings.Contains(o, ".edu/") ||
		strings.Contains(o, "scholar."):
		return "academic"
	case strings.Contains(o, "github.com") ||
		strings.Contains(o, "gitlab.com"):
		return "code_repo"
	case strings.Contains(o, "news.") ||
		strings.Contains(o, "nytimes") ||
		strings.Contains(o, "guardian") ||
		strings.Contains(o, "bbc.") ||
		strings.Contains(o, "reuters"):
		return "news"
	case strings.HasPrefix(o, "http://") || strings.HasPrefix(o, "https://"):
		return "web"
	case o != "":
		return "other"
	default:
		return ""
	}
}

// containsWordSimple checks if `term` appears in `text` with non-letter boundaries.
// Cheap regex-free approximation (good enough for evaluative-word detection).
func containsWordSimple(text, term string) bool {
	idx := strings.Index(text, term)
	for idx >= 0 {
		startOK := idx == 0 || !isLetter(text[idx-1])
		end := idx + len(term)
		endOK := end >= len(text) || !isLetter(text[end])
		if startOK && endOK {
			return true
		}
		next := strings.Index(text[idx+1:], term)
		if next < 0 {
			break
		}
		idx = idx + 1 + next
	}
	return false
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// spearmanRho computes Spearman's rank correlation between x and y.
func spearmanRho(x, y []float64) float64 {
	if len(x) != len(y) || len(x) < 2 {
		return 0
	}
	rx := assignRanks(x)
	ry := assignRanks(y)
	n := float64(len(x))
	var sum float64
	for i := range rx {
		d := rx[i] - ry[i]
		sum += d * d
	}
	return 1 - (6*sum)/(n*(n*n-1))
}

func assignRanks(vals []float64) []float64 {
	n := len(vals)
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return vals[idx[i]] < vals[idx[j]] })
	ranks := make([]float64, n)
	i := 0
	for i < n {
		j := i + 1
		for j < n && vals[idx[j]] == vals[idx[i]] {
			j++
		}
		// Tie correction: assign avg rank
		avg := float64(i+j+1) / 2
		for k := i; k < j; k++ {
			ranks[idx[k]] = avg
		}
		i = j
	}
	return ranks
}
