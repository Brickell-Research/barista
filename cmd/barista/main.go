package main

import (
	"log/slog"
	"net/url"
	"os"

	"github.com/hibiken/asynq"

	"barista/internal/config"
	"barista/internal/gitops"
	"barista/internal/llm"
	"barista/internal/worker"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	cfg, err := config.Load("config/services.yml")
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	gitops.EnsureRepo(cfg.OutputDir, cfg.OutputRepo)

	redisOpt := asynq.RedisClientOpt{Addr: parseRedisAddr(redisURL)}
	llmClient := llm.New()
	h := worker.New(cfg, llmClient, redisOpt)

	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 5,
		Queues: map[string]int{
			"exploration": 1,
			"default":     1,
		},
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TaskDiscover, h.HandleDiscover)
	mux.HandleFunc(worker.TaskExplore, h.HandleExplore)

	scheduler := asynq.NewScheduler(redisOpt, nil)
	if _, err := scheduler.Register("0 * * * *", asynq.NewTask(worker.TaskDiscover, nil)); err != nil {
		slog.Error("failed to register schedule", "err", err)
		os.Exit(1)
	}

	go func() {
		if err := scheduler.Run(); err != nil {
			slog.Error("scheduler stopped", "err", err)
		}
	}()

	slog.Info("starting barista worker")
	if err := srv.Run(mux); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func parseRedisAddr(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "localhost:6379"
	}
	return u.Host
}
