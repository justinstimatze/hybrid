// metacog — cognitive-bias auditor MCP server.
//
// Runs nine bias-signature checks against any typed substrate exposed via the
// Substrate interface. The metrics are general (HHI, Spearman, entropy,
// Jaccard, ratios); winze-specific Go-AST traversal has been abstracted away
// behind the Substrate interface.
//
// Tools:
//
//	audit(substrate_path, format) -> AuditReport      run all nine auditors
//	audit_one(substrate_path, format, auditor) -> r   run one auditor
//	list_auditors() -> []name                          enumerate auditors
//
// Currently supported substrate format: "jsonl" (one JSON record per line,
// schema documented in README.md). Future: "winze" reader for parity.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"metacog",
		"0.1.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)
	registerTools(s)
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func registerTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("audit",
		mcp.WithDescription("Run all nine cognitive-bias auditors against the substrate at substrate_path. Returns per-auditor results plus a summary of which triggered."),
		mcp.WithString("substrate_path", mcp.Required(),
			mcp.Description("Path to the substrate file")),
		mcp.WithString("format",
			mcp.Description("Substrate format: 'jsonl' (default). 'winze' reserved for future Go-AST reader.")),
	), handleAudit)

	s.AddTool(mcp.NewTool("audit_one",
		mcp.WithDescription("Run one named auditor against the substrate."),
		mcp.WithString("substrate_path", mcp.Required(),
			mcp.Description("Path to the substrate file")),
		mcp.WithString("auditor", mcp.Required(),
			mcp.Description("Auditor name: confirmation_bias, anchoring, clustering_illusion, availability_heuristic, survivorship_bias, framing_effect, dunning_kruger, base_rate_neglect, premature_closure")),
		mcp.WithString("format",
			mcp.Description("'jsonl' (default)")),
	), handleAuditOne)

	s.AddTool(mcp.NewTool("list_auditors",
		mcp.WithDescription("Enumerate the nine cognitive-bias auditors and what each measures."),
	), handleListAuditors)
}

func loadSubstrate(args map[string]any) (Substrate, error) {
	path, _ := args["substrate_path"].(string)
	if path == "" {
		return nil, fmt.Errorf("substrate_path is required")
	}
	format, _ := args["format"].(string)
	if format == "" {
		format = "jsonl"
	}
	switch format {
	case "jsonl":
		if _, err := os.Stat(path); err != nil {
			return nil, fmt.Errorf("substrate not found: %w", err)
		}
		return &JSONLSubstrate{Path: path}, nil
	default:
		return nil, fmt.Errorf("unknown format: %q (supported: jsonl)", format)
	}
}

func handleAudit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sub, err := loadSubstrate(req.GetArguments())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	results := make([]AuditorResult, 0, 9)
	triggered := []string{}
	skipped := []string{}
	for _, a := range AllAuditors() {
		r := a(sub)
		results = append(results, r)
		if r.Skipped {
			skipped = append(skipped, r.Bias)
		} else if r.Triggered {
			triggered = append(triggered, r.Bias)
		}
	}
	report := map[string]any{
		"results":   results,
		"triggered": triggered,
		"skipped":   skipped,
		"summary": fmt.Sprintf("%d/%d auditors triggered, %d skipped",
			len(triggered), len(results)-len(skipped), len(skipped)),
	}
	return jsonResult(report)
}

func handleAuditOne(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	auditor, _ := args["auditor"].(string)
	if auditor == "" {
		return mcp.NewToolResultError("auditor is required"), nil
	}
	sub, err := loadSubstrate(args)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	auditors := map[string]Auditor{
		"confirmation_bias":      AuditConfirmationBias,
		"anchoring":              AuditAnchoring,
		"clustering_illusion":    AuditClusteringIllusion,
		"availability_heuristic": AuditAvailabilityHeuristic,
		"survivorship_bias":      AuditSurvivorshipBias,
		"framing_effect":         AuditFramingEffect,
		"dunning_kruger":         AuditDunningKruger,
		"base_rate_neglect":      AuditBaseRateNeglect,
		"premature_closure":      AuditPrematureClosure,
	}
	a, ok := auditors[auditor]
	if !ok {
		return mcp.NewToolResultError(fmt.Sprintf("unknown auditor: %q", auditor)), nil
	}
	return jsonResult(a(sub))
}

func handleListAuditors(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	desc := []map[string]string{
		{"bias": "confirmation_bias", "metric": "corroboration_rate", "needs": "verdicts"},
		{"bias": "anchoring", "metric": "spearman_age_density", "needs": "CreatedAt or Group"},
		{"bias": "clustering_illusion", "metric": "group_cluster_jaccard", "needs": "Group + Cluster"},
		{"bias": "availability_heuristic", "metric": "provenance_hhi", "needs": "Provenance.Origin or Type"},
		{"bias": "survivorship_bias", "metric": "irrelevant_to_challenged_ratio", "needs": "verdicts incl. 'irrelevant'"},
		{"bias": "framing_effect", "metric": "evaluative_summary_fraction", "needs": "SummaryText"},
		{"bias": "dunning_kruger", "metric": "low_complexity_zero_edge_rate", "needs": "Complexity hint + Edges"},
		{"bias": "base_rate_neglect", "metric": "predicate_entropy_bits", "needs": "Edges.Predicate"},
		{"bias": "premature_closure", "metric": "closure_findings", "needs": "SummaryText"},
	}
	return jsonResult(desc)
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}
