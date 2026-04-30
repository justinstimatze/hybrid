// cal_log — calibration logger MCP server for hybrid-loop projects.
//
// Per-evaluator prediction logging with rolling-window hit-rate aggregation.
// Append-only event log at $CAL_LOG_PATH (default: ~/.cal_log/calibration.jsonl).
//
// Six tools:
//   predict(loop, input_hash, prediction, model_id, ...) -> {prediction_id, verdict_due_by}
//   resolve(prediction_id, verdict, verdict_source) -> resolved record
//   hit_rate(loop, window_days) -> {hit_rate, total_resolved, verdict_breakdown}
//   list_pending(loop?, limit) -> unresolved predictions ordered by due date
//   list_recent(loop?, limit) -> recent predictions, most recent first
//   stats() -> per-loop summary across all predictions
//
// Stdio transport. The hybrid claim: an evaluator that can't show its hit-rate is theater.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func resolveDBPath() string {
	if p := os.Getenv("CAL_LOG_PATH"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cannot determine home dir: %v", err)
	}
	return filepath.Join(home, ".cal_log", "calibration.jsonl")
}

func main() {
	dbPath := resolveDBPath()
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		log.Fatalf("cannot create db dir: %v", err)
	}
	store := &Store{Path: dbPath}

	s := server.NewMCPServer(
		"cal_log",
		"0.1.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	registerTools(s, store)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func registerTools(s *server.MCPServer, store *Store) {
	s.AddTool(mcp.NewTool("predict",
		mcp.WithDescription("Log a typed LLM evaluator's prediction. Returns prediction_id for later resolve()."),
		mcp.WithString("loop", mcp.Required(),
			mcp.Description("Stable identifier for the project + evaluator (e.g. 'recruiter-fit-scorer')")),
		mcp.WithString("input_hash", mcp.Required(),
			mcp.Description("Stable hash of the input the prediction was made on")),
		mcp.WithObject("prediction", mcp.Required(),
			mcp.Description("The structured record the evaluator emitted")),
		mcp.WithString("model_id", mcp.Required(),
			mcp.Description("Model identifier (e.g. 'claude-sonnet-4-6')")),
		mcp.WithNumber("schema_version",
			mcp.Description("Schema version (default 1)")),
		mcp.WithNumber("verdict_due_in_days",
			mcp.Description("Days until verdict is expected (default 7)")),
		mcp.WithString("lens_or_reasoner",
			mcp.Description("'lens' or 'reasoner' (default 'reasoner')")),
	), handlePredict(store))

	s.AddTool(mcp.NewTool("resolve",
		mcp.WithDescription("Mark a prediction's verdict. Common verdicts: 'confirmed', 'refuted', 'partial', 'irrelevant'."),
		mcp.WithString("prediction_id", mcp.Required(),
			mcp.Description("ID returned by predict()")),
		mcp.WithString("verdict", mcp.Required(),
			mcp.Description("Verdict category")),
		mcp.WithString("verdict_source",
			mcp.Description("How the verdict was determined (default 'manual')")),
	), handleResolve(store))

	s.AddTool(mcp.NewTool("hit_rate",
		mcp.WithDescription("Hit-rate (confirmed / total resolved) for one loop in the past window_days."),
		mcp.WithString("loop", mcp.Required(),
			mcp.Description("Loop identifier")),
		mcp.WithNumber("window_days",
			mcp.Description("Default 30")),
	), handleHitRate(store))

	s.AddTool(mcp.NewTool("list_pending",
		mcp.WithDescription("Unresolved predictions, ordered by verdict_due_by (oldest first)."),
		mcp.WithString("loop",
			mcp.Description("Optional loop filter")),
		mcp.WithNumber("limit",
			mcp.Description("Default 50")),
	), handleListPending(store))

	s.AddTool(mcp.NewTool("list_recent",
		mcp.WithDescription("Recent predictions (resolved or not), most recent first."),
		mcp.WithString("loop",
			mcp.Description("Optional loop filter")),
		mcp.WithNumber("limit",
			mcp.Description("Default 50")),
	), handleListRecent(store))

	s.AddTool(mcp.NewTool("stats",
		mcp.WithDescription("Top-level summary across all loops: counts and hit-rate where computable."),
	), handleStats(store))
}

func handlePredict(store *Store) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		loop, _ := args["loop"].(string)
		inputHash, _ := args["input_hash"].(string)
		prediction, _ := args["prediction"].(map[string]any)
		modelID, _ := args["model_id"].(string)
		if loop == "" || inputHash == "" || prediction == nil || modelID == "" {
			return mcp.NewToolResultError("loop, input_hash, prediction, and model_id are required"), nil
		}
		schemaVersion := 1
		if v, ok := args["schema_version"].(float64); ok {
			schemaVersion = int(v)
		}
		dueDays := 7
		if v, ok := args["verdict_due_in_days"].(float64); ok {
			dueDays = int(v)
		}
		role := "reasoner"
		if v, ok := args["lens_or_reasoner"].(string); ok && v != "" {
			role = v
		}

		now := time.Now()
		id := uuid.NewString()
		dueBy := now.Add(time.Duration(dueDays) * 24 * time.Hour).Unix()
		event := PredictEvent{
			Event:          "predict",
			PredictionID:   id,
			TS:             now.Unix(),
			Loop:           loop,
			LensOrReasoner: role,
			InputHash:      inputHash,
			Prediction:     prediction,
			ModelID:        modelID,
			SchemaVersion:  schemaVersion,
			VerdictDueBy:   dueBy,
		}
		if err := store.Append(event); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return jsonResult(map[string]any{
			"prediction_id":  id,
			"verdict_due_by": dueBy,
		})
	}
}

func handleResolve(store *Store) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		pid, _ := args["prediction_id"].(string)
		verdict, _ := args["verdict"].(string)
		if pid == "" || verdict == "" {
			return mcp.NewToolResultError("prediction_id and verdict are required"), nil
		}
		verdictSource := "manual"
		if v, ok := args["verdict_source"].(string); ok && v != "" {
			verdictSource = v
		}
		state, err := store.Fold()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		rec, ok := state[pid]
		if !ok {
			return jsonResult(map[string]any{"error": fmt.Sprintf("prediction_id %s not found", pid)})
		}
		if rec.Verdict != "" {
			return jsonResult(map[string]any{
				"error":    fmt.Sprintf("prediction_id %s already resolved", pid),
				"existing": rec,
			})
		}
		now := time.Now().Unix()
		event := ResolveEvent{
			Event:         "resolve",
			PredictionID:  pid,
			TS:            now,
			Verdict:       verdict,
			VerdictSource: verdictSource,
			VerdictTS:     now,
		}
		if err := store.Append(event); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		rec.Verdict = verdict
		rec.VerdictSource = verdictSource
		rec.VerdictTS = now
		return jsonResult(rec)
	}
}

func handleHitRate(store *Store) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		loop, _ := args["loop"].(string)
		if loop == "" {
			return mcp.NewToolResultError("loop is required"), nil
		}
		windowDays := 30
		if v, ok := args["window_days"].(float64); ok {
			windowDays = int(v)
		}
		state, err := store.Fold()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		cutoff := time.Now().Add(-time.Duration(windowDays) * 24 * time.Hour).Unix()
		breakdown := map[string]int{}
		total := 0
		for _, r := range state {
			if r.Loop != loop || r.Verdict == "" || r.VerdictTS < cutoff {
				continue
			}
			breakdown[r.Verdict]++
			total++
		}
		out := map[string]any{
			"loop":              loop,
			"window_days":       windowDays,
			"total_resolved":    total,
			"verdict_breakdown": breakdown,
		}
		if total == 0 {
			out["hit_rate"] = nil
		} else {
			out["hit_rate"] = float64(breakdown["confirmed"]) / float64(total)
		}
		return jsonResult(out)
	}
}

func handleListPending(store *Store) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		loop, _ := args["loop"].(string)
		limit := 50
		if v, ok := args["limit"].(float64); ok {
			limit = int(v)
		}
		state, err := store.Fold()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		pending := []*Record{}
		for _, r := range state {
			if r.Verdict != "" {
				continue
			}
			if loop != "" && r.Loop != loop {
				continue
			}
			pending = append(pending, r)
		}
		sort.Slice(pending, func(i, j int) bool {
			return pending[i].VerdictDueBy < pending[j].VerdictDueBy
		})
		if len(pending) > limit {
			pending = pending[:limit]
		}
		return jsonResult(pending)
	}
}

func handleListRecent(store *Store) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		loop, _ := args["loop"].(string)
		limit := 50
		if v, ok := args["limit"].(float64); ok {
			limit = int(v)
		}
		state, err := store.Fold()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		records := []*Record{}
		for _, r := range state {
			if loop != "" && r.Loop != loop {
				continue
			}
			records = append(records, r)
		}
		sort.Slice(records, func(i, j int) bool {
			return records[i].TS > records[j].TS
		})
		if len(records) > limit {
			records = records[:limit]
		}
		return jsonResult(records)
	}
}

type loopStats struct {
	Total     int      `json:"total"`
	Resolved  int      `json:"resolved"`
	Confirmed int      `json:"confirmed"`
	HitRate   *float64 `json:"hit_rate"`
}

func handleStats(store *Store) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		state, err := store.Fold()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		byLoop := map[string]*loopStats{}
		for _, r := range state {
			loop := r.Loop
			if loop == "" {
				loop = "unknown"
			}
			b, ok := byLoop[loop]
			if !ok {
				b = &loopStats{}
				byLoop[loop] = b
			}
			b.Total++
			if r.Verdict != "" {
				b.Resolved++
				if r.Verdict == "confirmed" {
					b.Confirmed++
				}
			}
		}
		for _, b := range byLoop {
			if b.Resolved > 0 {
				rate := float64(b.Confirmed) / float64(b.Resolved)
				b.HitRate = &rate
			}
		}
		return jsonResult(map[string]any{
			"db_path":           store.Path,
			"total_predictions": len(state),
			"by_loop":           byLoop,
		})
	}
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}
