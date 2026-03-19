package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

type Client struct {
	api anthropic.Client
}

func New() *Client {
	return &Client{
		api: anthropic.NewClient(), // reads ANTHROPIC_API_KEY from env
	}
}

type GuaranteeJSON struct {
	Name       string  `json:"name"`
	Threshold  float64 `json:"threshold"`
	WindowDays int     `json:"window_days"`
}

type ExtractResponse struct {
	Guarantees []GuaranteeJSON `json:"guarantees"`
}

func (c *Client) ExtractGuarantees(ctx context.Context, system, user string) (*ExtractResponse, error) {
	msg, err := c.api.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeOpus4_6,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: system},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(user)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("anthropic API: %w", err)
	}

	var text string
	for _, block := range msg.Content {
		if b, ok := block.AsAny().(anthropic.TextBlock); ok {
			text = b.Text
			break
		}
	}

	if text == "" {
		return nil, fmt.Errorf("empty response from API")
	}

	text = strings.TrimPrefix(text, "```json\n")
	text = strings.TrimPrefix(text, "```\n")
	text = strings.TrimSuffix(text, "\n```")
	text = strings.TrimSpace(text)

	var result ExtractResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("parse LLM JSON: %w", err)
	}

	return &result, nil
}
