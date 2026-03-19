package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"barista/internal/config"
	"barista/internal/fetcher"
	"barista/internal/gitops"
	"barista/internal/llm"
	"barista/internal/pipeline"
)

func main() {
	cfg, err := config.Load("config/services.yml")
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	gitops.EnsureRepo(cfg.OutputDir, cfg.OutputRepo)

	llmClient := llm.New()
	synth := pipeline.NewSynthesizer(llmClient)
	ctx := context.Background()

	var services []config.Service
	if len(os.Args) > 1 {
		svc, ok := cfg.FindService(os.Args[1])
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown service: %s\n", os.Args[1])
			os.Exit(1)
		}
		services = []config.Service{*svc}
	} else {
		services = cfg.AllServices()
	}

	for _, svc := range services {
		slog.Info("fetching", "service", svc.Key(), "url", svc.URL)
		content, err := fetcher.Fetch(svc.URL)
		if err != nil {
			slog.Error("fetch failed", "service", svc.Key(), "err", err)
			continue
		}

		slog.Info("synthesizing", "service", svc.Key())
		intermediate, err := synth.Synthesize(ctx, svc, content)
		if err != nil {
			slog.Error("synthesize failed", "service", svc.Key(), "err", err)
			continue
		}

		caffeine := pipeline.Translate(intermediate)
		result, err := pipeline.Write(cfg.OutputDir, intermediate, caffeine)
		if err != nil {
			slog.Error("write failed", "service", svc.Key(), "err", err)
			continue
		}

		switch result.Status {
		case pipeline.StatusWritten:
			slog.Info("written", "service", svc.Key(), "path", result.Path)
		case pipeline.StatusUnchanged:
			slog.Info("unchanged", "service", svc.Key())
		case pipeline.StatusBlip:
			slog.Warn("blip: zero guarantees returned, previous file preserved", "service", svc.Key())
		}
	}

	gitops.CommitAndPush(cfg.OutputDir, "update expectations")
}
