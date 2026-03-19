package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"

	"barista/internal/config"
	"barista/internal/fetcher"
	"barista/internal/gitops"
	"barista/internal/llm"
	"barista/internal/pipeline"
)

const (
	TaskDiscover = "barista:discover"
	TaskExplore  = "barista:explore"
)

type explorePayload struct {
	ServiceKey string `json:"service_key"`
}

type Handler struct {
	cfg    *config.Config
	synth  *pipeline.Synthesizer
	client *asynq.Client
}

func New(cfg *config.Config, llmClient *llm.Client, redisOpt asynq.RedisClientOpt) *Handler {
	return &Handler{
		cfg:    cfg,
		synth:  pipeline.NewSynthesizer(llmClient),
		client: asynq.NewClient(redisOpt),
	}
}

func (h *Handler) HandleDiscover(ctx context.Context, _ *asynq.Task) error {
	services := h.cfg.AllServices()
	slog.Info("discovering services", "count", len(services))

	for _, svc := range services {
		payload, err := json.Marshal(explorePayload{ServiceKey: svc.Key()})
		if err != nil {
			return fmt.Errorf("marshal payload for %s: %w", svc.Key(), err)
		}

		task := asynq.NewTask(TaskExplore, payload, asynq.Queue("exploration"))
		if _, err := h.client.EnqueueContext(ctx, task); err != nil {
			slog.Error("failed to enqueue", "service", svc.Key(), "err", err)
		}
	}

	return nil
}

func (h *Handler) HandleExplore(ctx context.Context, t *asynq.Task) error {
	var payload explorePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	svc, ok := h.cfg.FindService(payload.ServiceKey)
	if !ok {
		return fmt.Errorf("unknown service: %s", payload.ServiceKey)
	}

	slog.Info("fetching", "service", svc.Key(), "url", svc.URL)
	content, err := fetcher.Fetch(svc.URL)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	slog.Info("synthesizing", "service", svc.Key())
	intermediate, err := h.synth.Synthesize(ctx, *svc, content)
	if err != nil {
		return fmt.Errorf("synthesize: %w", err)
	}

	caffeine := pipeline.Translate(intermediate)

	result, err := pipeline.Write(h.cfg.OutputDir, intermediate, caffeine)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	switch result.Status {
	case pipeline.StatusWritten:
		slog.Info("written", "service", svc.Key(), "path", result.Path)
		gitops.CommitAndPush(h.cfg.OutputDir, fmt.Sprintf("update %s expectations", svc.Key()))
	case pipeline.StatusUnchanged:
		slog.Info("unchanged", "service", svc.Key())
	case pipeline.StatusBlip:
		slog.Warn("blip: zero guarantees returned, previous file preserved", "service", svc.Key())
	}
	return nil
}
