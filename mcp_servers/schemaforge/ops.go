package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// The four LLM ops. Each is one round-trip; each has a tag-extraction step
// that's tested against the FakeLLMClient — that's where the contract between
// system prompts and downstream parsing lives.

// DesignNotation asks the model to design (round 1) or evolve (round N+1) a
// dense notation. Up to 2 longest items in `seeds` anchor the design.
//
// `target` is free-form: "complete TypeScript implementation using Hono+Drizzle",
// "structured behavioral mechanism with Effect/Process/Necessity sections",
// "standalone humor template instance with stage directions", etc. The
// notation will be designed to expand to that shape.
//
// On evolution rounds, `previousNotation` and `previousMetricsSummary` flow
// the prior round's findings back into the design prompt. Empty strings = round 1.
func DesignNotation(
	ctx context.Context,
	llm LLMClient,
	seeds []CorpusItem,
	target string,
	model string,
	previousNotation string,
	previousMetricsSummary string,
) (string, error) {
	if len(seeds) == 0 {
		return "", fmt.Errorf("no seed items provided")
	}
	if target == "" {
		return "", fmt.Errorf("target description is required")
	}

	var user string
	if previousNotation != "" {
		user = fmt.Sprintf(`Here is your previous notation:

<notation_spec>
%s
</notation_spec>

Metrics from last round:
%s

Decide whether to change the notation. The default is to keep it unchanged.
Only modify the notation if there are SPECIFIC failures justifying the change:

- If mean correctness >= 0.85 and no item scored below 0.70: output the notation UNCHANGED.
  Minor terminology drift at this level is rubric noise, not notation failure.
- If items scored below 0.70 due to a specific structural ambiguity in the notation:
  make a TARGETED fix to that ambiguity. Do not touch unrelated parts.
- If a primitive was clearly unused across all items: remove it.
- Do NOT add new primitives unless a specific failure requires them.

A growing notation that doesn't address concrete failures is a regression.

The notation should expand to: %s

Output ONLY the (possibly unchanged) notation specification between <notation_spec> tags.`,
			previousNotation, previousMetricsSummary, target)
	} else {
		var seedsBuf strings.Builder
		for i, s := range seeds {
			if i > 0 {
				seedsBuf.WriteString("\n---\n")
			}
			seedsBuf.WriteString("# ")
			seedsBuf.WriteString(s.ID)
			seedsBuf.WriteString("\n")
			seedsBuf.WriteString(s.SpecText)
		}
		user = fmt.Sprintf(`Design a maximally dense notation for items in this corpus.

The notation should expand to: %s

Seed items (use these to anchor the notation design — pick the hard cases, simple cases follow):

%s

Output ONLY the notation specification between <notation_spec> tags. Do NOT translate the items yet.`,
			target, seedsBuf.String())
	}

	resp, err := llm.Call(ctx, designSystem, user, model, 8192)
	if err != nil {
		return "", fmt.Errorf("design call: %w", err)
	}
	spec := extractTag(resp, "notation_spec")
	if spec == "" {
		return "", fmt.Errorf("no <notation_spec> block in response")
	}
	return spec, nil
}

// Compress translates one corpus item into the notation. One LLM call.
// Returns the notation text and its estimated token count.
func Compress(
	ctx context.Context,
	llm LLMClient,
	specText string,
	notationSpec string,
	target string,
	model string,
) (string, int, error) {
	if specText == "" {
		return "", 0, fmt.Errorf("spec_text is required")
	}
	if notationSpec == "" {
		return "", 0, fmt.Errorf("notation_spec is required")
	}
	user := fmt.Sprintf(`Translate this item into the notation.

The notation expands to: %s

Item:

%s

Output ONLY the notation between <notation> tags.`, target, specText)

	resp, err := llm.Call(ctx, compressSystem+"\n\nNOTATION SPEC:\n"+notationSpec, user, model, 4096)
	if err != nil {
		return "", 0, fmt.Errorf("compress call: %w", err)
	}
	notation := extractTag(resp, "notation")
	if notation == "" {
		return "", 0, fmt.Errorf("no <notation> block in response")
	}
	return notation, EstimateTokens(notation), nil
}

// Expand translates one notation back to the full target form. One LLM call.
// Returns the expanded text and its estimated token count.
func Expand(
	ctx context.Context,
	llm LLMClient,
	notationSpec string,
	itemNotation string,
	target string,
	model string,
) (string, int, error) {
	if notationSpec == "" {
		return "", 0, fmt.Errorf("notation_spec is required")
	}
	if itemNotation == "" {
		return "", 0, fmt.Errorf("item_notation is required")
	}
	user := fmt.Sprintf(`Expand this notation into a complete %s.

Notation:

%s`, target, itemNotation)

	resp, err := llm.Call(ctx, expandSystem+"\n\nNOTATION SPEC:\n"+notationSpec, user, model, 16384)
	if err != nil {
		return "", 0, fmt.Errorf("expand call: %w", err)
	}
	resp = strings.TrimSpace(resp)
	if resp == "" {
		return "", 0, fmt.Errorf("expansion is empty")
	}
	return resp, EstimateTokens(resp), nil
}

// EvaluateRoundtrip scores how well an expansion preserves the original
// content under the supplied rubric. Returns a 0-1 score and the model's
// one-sentence reasoning. On parse failure returns score=0.5 with the parse
// error in reasoning — better than failing the whole round.
func EvaluateRoundtrip(
	ctx context.Context,
	llm LLMClient,
	originalSpec string,
	expanded string,
	rubric string,
	model string,
) (float64, string, error) {
	if rubric == "" {
		rubric = defaultRubric
	}
	user := fmt.Sprintf(`Original specification:
%s

Expanded implementation:
%s

Rubric:
%s

Output ONLY the JSON object.`, originalSpec, expanded, rubric)

	resp, err := llm.Call(ctx, evaluateSystem, user, model, 512)
	if err != nil {
		return 0, "", fmt.Errorf("evaluate call: %w", err)
	}
	score, reasoning, ok := parseScoreJSON(resp)
	if !ok {
		return 0.5, "parse_failed: " + truncate(resp, 120), nil
	}
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	return score, reasoning, nil
}

// --- helpers ---

func extractTag(s, tag string) string {
	re := regexp.MustCompile(`(?s)<` + regexp.QuoteMeta(tag) + `[^>]*>(.*?)</` + regexp.QuoteMeta(tag) + `>`)
	m := re.FindStringSubmatch(s)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}

var jsonObjRe = regexp.MustCompile(`(?s)\{[^{}]*"overall"[^{}]*\}`)

func parseScoreJSON(s string) (float64, string, bool) {
	m := jsonObjRe.FindString(s)
	if m == "" {
		return 0, "", false
	}
	var parsed struct {
		Overall   float64 `json:"overall"`
		Reasoning string  `json:"reasoning"`
	}
	if err := json.Unmarshal([]byte(m), &parsed); err != nil {
		return 0, "", false
	}
	return parsed.Overall, parsed.Reasoning, true
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
