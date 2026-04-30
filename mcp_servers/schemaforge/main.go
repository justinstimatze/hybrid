// schemaforge — schema-discovery (compress+verify) MCP server.
//
// Pattern: design a dense notation for a corpus, translate items into it,
// expand back, score the roundtrip, evolve the notation based on per-round
// metrics. Domain-agnostic: caller supplies a corpus (JSONL of {id, spec_text})
// plus a free-form `target` (what the notation expands to) and `rubric`
// (how roundtrip fidelity is scored).
//
// The fine-grained tools (design_notation, compress, expand, evaluate_roundtrip)
// are one LLM round-trip each — bounded for MCP. score_round is pure compute.
// run_round is the convenience that drives a full round across a corpus and
// writes round{N}/ to disk.
//
// Tools:
//
//	design_notation(corpus_path, target, model, previous_notation?, previous_metrics_summary?) -> notation_spec
//	compress(spec_text, notation_spec, target, model) -> {notation_text, notation_tokens}
//	expand(notation_spec, item_notation, target, model) -> {expanded_text, expanded_tokens}
//	evaluate_roundtrip(original_spec, expanded, rubric, model) -> {score, reasoning}
//	score_round(round_dir) -> RoundMetrics
//	run_round(corpus_path, target, rubric, model, output_dir, round_number, ...) -> RoundMetrics
//
// Requires ANTHROPIC_API_KEY in env (except score_round, which is pure compute).
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const defaultModel = "claude-sonnet-4-6"

func main() {
	s := server.NewMCPServer(
		"schemaforge",
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
	s.AddTool(mcp.NewTool("design_notation",
		mcp.WithDescription("Ask the model to design (round 1) or evolve (round N+1) a dense notation for a corpus. One LLM call. Returns the notation spec text. Reads up to 2 longest items from corpus_path as design seeds."),
		mcp.WithString("corpus_path", mcp.Required(), mcp.Description("Path to corpus JSONL")),
		mcp.WithString("target", mcp.Required(), mcp.Description("What the notation should expand to (free-form, e.g. 'TypeScript Hono+Drizzle CRUD app')")),
		mcp.WithString("model", mcp.Description("Anthropic model. Default: "+defaultModel)),
		mcp.WithString("previous_notation", mcp.Description("Prior round's notation (empty for round 1)")),
		mcp.WithString("previous_metrics_summary", mcp.Description("Prior round's metrics summary string (empty for round 1)")),
	), handleDesignNotation)

	s.AddTool(mcp.NewTool("compress",
		mcp.WithDescription("Translate one item's spec into the notation. One LLM call."),
		mcp.WithString("spec_text", mcp.Required(), mcp.Description("Original spec to translate")),
		mcp.WithString("notation_spec", mcp.Required(), mcp.Description("The notation specification produced by design_notation")),
		mcp.WithString("target", mcp.Required(), mcp.Description("What the notation expands to")),
		mcp.WithString("model", mcp.Description("Anthropic model. Default: "+defaultModel)),
	), handleCompress)

	s.AddTool(mcp.NewTool("expand",
		mcp.WithDescription("Expand a notation translation back to the full target form. One LLM call."),
		mcp.WithString("notation_spec", mcp.Required(), mcp.Description("The notation specification")),
		mcp.WithString("item_notation", mcp.Required(), mcp.Description("The compressed notation for this item")),
		mcp.WithString("target", mcp.Required(), mcp.Description("What the notation expands to")),
		mcp.WithString("model", mcp.Description("Anthropic model. Default: "+defaultModel)),
	), handleExpand)

	s.AddTool(mcp.NewTool("evaluate_roundtrip",
		mcp.WithDescription("Score how well an expansion preserves the original spec under the rubric. One LLM call. Returns 0-1 score plus one-sentence reasoning."),
		mcp.WithString("original_spec", mcp.Required(), mcp.Description("Original spec text")),
		mcp.WithString("expanded", mcp.Required(), mcp.Description("Expanded text from expand()")),
		mcp.WithString("rubric", mcp.Description("Free-form scoring rubric. Default: generic semantic-content rubric")),
		mcp.WithString("model", mcp.Description("Anthropic model. Default: "+defaultModel)),
	), handleEvaluate)

	s.AddTool(mcp.NewTool("score_round",
		mcp.WithDescription("Pure compute (no LLM). Read a round directory (.../roundN/items/*.json) and return aggregate metrics."),
		mcp.WithString("round_dir", mcp.Required(), mcp.Description("Path to a round directory like '.../run-X/round2'")),
	), handleScoreRound)

	s.AddTool(mcp.NewTool("run_round",
		mcp.WithDescription("Convenience: drive one full round (design → compress all → expand all → evaluate all → score). Writes round{N}/notation.txt, round{N}/items/{id}.json incrementally (resumable), round{N}/metrics.json on completion. 1 + 3N LLM calls for N items. Caller invokes per round and decides when to stop based on the metrics trend."),
		mcp.WithString("corpus_path", mcp.Required(), mcp.Description("Path to corpus JSONL")),
		mcp.WithString("target", mcp.Required(), mcp.Description("What the notation expands to")),
		mcp.WithString("rubric", mcp.Description("Free-form scoring rubric. Default: generic")),
		mcp.WithString("model", mcp.Description("Anthropic model. Default: "+defaultModel)),
		mcp.WithString("output_dir", mcp.Required(), mcp.Description("Run-level output dir; round subdir created inside")),
		mcp.WithNumber("round_number", mcp.Required(), mcp.Description("1-based round number")),
		mcp.WithString("previous_notation", mcp.Description("Prior round's notation (omit for round 1)")),
		mcp.WithString("previous_metrics_path", mcp.Description("Path to prior round's metrics.json (omit for round 1)")),
	), handleRunRound)
}

// --- handlers ---

func handleDesignNotation(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	llm, err := NewAnthropicClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	corpusPath, _ := args["corpus_path"].(string)
	items, err := LoadCorpus(corpusPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	spec, err := DesignNotation(
		ctx, llm,
		SeedItems(items, 2),
		stringArg(args, "target", ""),
		stringArg(args, "model", defaultModel),
		stringArg(args, "previous_notation", ""),
		stringArg(args, "previous_metrics_summary", ""),
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"notation_spec": spec, "tokens": EstimateTokens(spec)})
}

func handleCompress(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	llm, err := NewAnthropicClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	notation, tokens, err := Compress(
		ctx, llm,
		stringArg(args, "spec_text", ""),
		stringArg(args, "notation_spec", ""),
		stringArg(args, "target", ""),
		stringArg(args, "model", defaultModel),
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"notation_text": notation, "notation_tokens": tokens})
}

func handleExpand(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	llm, err := NewAnthropicClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	expanded, tokens, err := Expand(
		ctx, llm,
		stringArg(args, "notation_spec", ""),
		stringArg(args, "item_notation", ""),
		stringArg(args, "target", ""),
		stringArg(args, "model", defaultModel),
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"expanded_text": expanded, "expanded_tokens": tokens})
}

func handleEvaluate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	llm, err := NewAnthropicClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	score, reasoning, err := EvaluateRoundtrip(
		ctx, llm,
		stringArg(args, "original_spec", ""),
		stringArg(args, "expanded", ""),
		stringArg(args, "rubric", ""),
		stringArg(args, "model", defaultModel),
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"score": score, "reasoning": reasoning})
}

func handleScoreRound(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	dir := stringArg(args, "round_dir", "")
	if dir == "" {
		return mcp.NewToolResultError("round_dir is required"), nil
	}
	m, err := ScoreRoundDir(dir)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(m)
}

func handleRunRound(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	llm, err := NewAnthropicClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	roundNum := 0
	if v, ok := args["round_number"].(float64); ok {
		roundNum = int(v)
	}
	p := RunRoundParams{
		CorpusPath:          stringArg(args, "corpus_path", ""),
		Target:              stringArg(args, "target", ""),
		Rubric:              stringArg(args, "rubric", ""),
		Model:               stringArg(args, "model", defaultModel),
		OutputDir:           stringArg(args, "output_dir", ""),
		RoundNumber:         roundNum,
		PreviousNotation:    stringArg(args, "previous_notation", ""),
		PreviousMetricsPath: stringArg(args, "previous_metrics_path", ""),
	}
	m, err := RunRound(ctx, llm, p)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(m)
}

// --- helpers ---

func stringArg(args map[string]any, key, dflt string) string {
	v, _ := args[key].(string)
	if v == "" {
		return dflt
	}
	return v
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}
