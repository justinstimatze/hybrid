package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// LLMClient is the narrow interface every op depends on. Tests pass a fake
// implementation; production wires AnthropicClient. The interface is one
// method on purpose — anything more would couple the ops to SDK shapes.
type LLMClient interface {
	// Call issues one chat completion. system is cached (ephemeral) when
	// the underlying provider supports prompt caching; user is the only
	// per-call variable. Returns the assistant's text response.
	Call(ctx context.Context, system, user, model string, maxTokens int) (string, error)
}

// AnthropicClient wraps anthropic-sdk-go and turns on prompt caching for the
// system prompt. The system prompt in our four ops is reused across many
// per-item calls in a round (notation_spec for compress and expand) — caching
// turns N+N round-trips into "cache hit" cost for the bulk of the prompt.
type AnthropicClient struct {
	c       anthropic.Client
	timeout time.Duration
}

// NewAnthropicClient reads ANTHROPIC_API_KEY from env. Returns a clear error
// if unset rather than failing inside Call — handlers want the error early.
func NewAnthropicClient() (*AnthropicClient, error) {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}
	c := anthropic.NewClient(option.WithAPIKey(key))
	return &AnthropicClient{c: c, timeout: 10 * time.Minute}, nil
}

func (a *AnthropicClient) Call(ctx context.Context, system, user, model string, maxTokens int) (string, error) {
	if maxTokens <= 0 {
		maxTokens = 8192
	}
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	resp, err := a.c.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: int64(maxTokens),
		System: []anthropic.TextBlockParam{
			{
				Text:         system,
				CacheControl: anthropic.CacheControlEphemeralParam{Type: "ephemeral", TTL: anthropic.CacheControlEphemeralTTLTTL1h},
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(user)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("anthropic call: %w", err)
	}
	if resp.StopReason == "max_tokens" {
		return "", fmt.Errorf("response truncated (max_tokens=%d) — try a larger budget", maxTokens)
	}
	for _, block := range resp.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}
	return "", fmt.Errorf("no text block in response")
}
