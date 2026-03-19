package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"barista/internal/config"
)

func TestLoad(t *testing.T) {
	yml := `
providers:
  - name: aws
    services:
      - name: s3
        url: https://aws.amazon.com/s3/sla/
      - name: rds
        url: https://aws.amazon.com/rds/sla/
  - name: stripe
    services:
      - name: payments
        url: https://stripe.com/sla
`
	path := filepath.Join(t.TempDir(), "services.yml")
	if err := os.WriteFile(path, []byte(yml), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if len(cfg.Providers) != 2 {
		t.Errorf("want 2 providers, got %d", len(cfg.Providers))
	}
	if cfg.OutputDir != "output" {
		t.Errorf("want default output dir, got %q", cfg.OutputDir)
	}
}

func TestAllServices(t *testing.T) {
	yml := `
providers:
  - name: aws
    services:
      - name: s3
        url: https://aws.amazon.com/s3/sla/
      - name: rds
        url: https://aws.amazon.com/rds/sla/
  - name: stripe
    services:
      - name: payments
        url: https://stripe.com/sla
`
	path := filepath.Join(t.TempDir(), "services.yml")
	os.WriteFile(path, []byte(yml), 0644)

	cfg, _ := config.Load(path)
	services := cfg.AllServices()

	if len(services) != 3 {
		t.Errorf("want 3 services, got %d", len(services))
	}
}

func TestServiceKey(t *testing.T) {
	yml := `
providers:
  - name: aws
    services:
      - name: s3
        url: https://aws.amazon.com/s3/sla/
`
	path := filepath.Join(t.TempDir(), "services.yml")
	os.WriteFile(path, []byte(yml), 0644)

	cfg, _ := config.Load(path)
	svc := cfg.AllServices()[0]

	if svc.Key() != "aws/s3" {
		t.Errorf("want key aws/s3, got %q", svc.Key())
	}
}

func TestFindService(t *testing.T) {
	yml := `
providers:
  - name: aws
    services:
      - name: s3
        url: https://aws.amazon.com/s3/sla/
`
	path := filepath.Join(t.TempDir(), "services.yml")
	os.WriteFile(path, []byte(yml), 0644)

	cfg, _ := config.Load(path)

	svc, ok := cfg.FindService("aws/s3")
	if !ok {
		t.Fatal("expected to find aws/s3")
	}
	if svc.Name != "s3" {
		t.Errorf("want s3, got %q", svc.Name)
	}

	_, ok = cfg.FindService("aws/unknown")
	if ok {
		t.Error("expected not to find aws/unknown")
	}
}
