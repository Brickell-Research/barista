package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	model        = "claude-opus-4-6"
	anthropicURL = "https://api.anthropic.com/v1/messages"
)

type Client struct {
	apiKey string
	http   *http.Client
}

func New() *Client {
	return &Client{
		apiKey: os.Getenv("ANTHROPIC_API_KEY"),
		http:   &http.Client{},
	}
}

type requestBody struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
	Messages  []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseBody struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
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
	payload, err := json.Marshal(requestBody{
		Model:     model,
		MaxTokens: 1024,
		System:    system,
		Messages:  []message{{Role: "user", Content: user}},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("api call: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic API error %d: %s", resp.StatusCode, body)
	}

	var apiResp responseBody
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(apiResp.Content) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	text := apiResp.Content[0].Text
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
