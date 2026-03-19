package pipeline

import (
	"context"
	"fmt"

	"barista/internal/config"
	"barista/internal/llm"
	"barista/internal/prompts"
)

type Synthesizer struct {
	client *llm.Client
}

func NewSynthesizer(client *llm.Client) *Synthesizer {
	return &Synthesizer{client: client}
}

func (s *Synthesizer) Synthesize(ctx context.Context, svc config.Service, content string) (*Intermediate, error) {
	userPrompt := prompts.SLAUser(svc.ProviderName, svc.Name, content)

	resp, err := s.client.ExtractGuarantees(ctx, prompts.SLASystem, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("LLM extraction failed: %w", err)
	}

	guarantees := make([]Guarantee, 0, len(resp.Guarantees))
	for _, g := range resp.Guarantees {
		guarantees = append(guarantees, Guarantee{
			Name:       g.Name,
			Threshold:  g.Threshold,
			WindowDays: g.WindowDays,
		})
	}

	return &Intermediate{
		ServiceName:  svc.Name,
		ProviderName: svc.ProviderName,
		SourceURL:    svc.URL,
		Guarantees:   guarantees,
	}, nil
}
