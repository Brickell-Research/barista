package pipeline_test

import (
	"os"
	"path/filepath"
	"testing"

	"barista/internal/pipeline"
)

func newIntermediate(provider, service, url string, guarantees []pipeline.Guarantee) *pipeline.Intermediate {
	return &pipeline.Intermediate{
		ProviderName: provider,
		ServiceName:  service,
		SourceURL:    url,
		Guarantees:   guarantees,
	}
}

func TestWrite_FirstWrite(t *testing.T) {
	dir := t.TempDir()
	i := newIntermediate("aws", "s3", "https://aws.amazon.com/s3/sla/", []pipeline.Guarantee{
		{Name: "monthly_uptime", Threshold: 99.9, WindowDays: 30},
	})
	content := pipeline.Translate(i)

	result, err := pipeline.Write(dir, i, content)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if result.Status != pipeline.StatusWritten {
		t.Errorf("want StatusWritten, got %v", result.Status)
	}

	wantPath := filepath.Join(dir, "expectations", "s3.caffeine")
	if result.Path != wantPath {
		t.Errorf("got path %q, want %q", result.Path, wantPath)
	}

	data, _ := os.ReadFile(result.Path)
	if len(data) == 0 {
		t.Error("expected non-empty file")
	}

	// changelog should exist
	changelogPath := filepath.Join(dir, "expectations", "s3.changelog")
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		t.Error("expected changelog file to be created")
	}
}

func TestWrite_Unchanged(t *testing.T) {
	dir := t.TempDir()
	i := newIntermediate("aws", "s3", "https://aws.amazon.com/s3/sla/", []pipeline.Guarantee{
		{Name: "monthly_uptime", Threshold: 99.9, WindowDays: 30},
	})
	content := pipeline.Translate(i)

	pipeline.Write(dir, i, content)
	result, err := pipeline.Write(dir, i, content)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if result.Status != pipeline.StatusUnchanged {
		t.Errorf("want StatusUnchanged, got %v", result.Status)
	}
}

func TestWrite_Changed(t *testing.T) {
	dir := t.TempDir()

	i1 := newIntermediate("aws", "s3", "https://aws.amazon.com/s3/sla/", []pipeline.Guarantee{
		{Name: "monthly_uptime", Threshold: 99.9, WindowDays: 30},
	})
	pipeline.Write(dir, i1, pipeline.Translate(i1))

	i2 := newIntermediate("aws", "s3", "https://aws.amazon.com/s3/sla/", []pipeline.Guarantee{
		{Name: "monthly_uptime", Threshold: 99.95, WindowDays: 30},
	})
	result, err := pipeline.Write(dir, i2, pipeline.Translate(i2))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if result.Status != pipeline.StatusWritten {
		t.Errorf("want StatusWritten on change, got %v", result.Status)
	}
}

func TestWrite_Blip(t *testing.T) {
	dir := t.TempDir()

	// Write a valid file first
	i := newIntermediate("aws", "s3", "https://aws.amazon.com/s3/sla/", []pipeline.Guarantee{
		{Name: "monthly_uptime", Threshold: 99.9, WindowDays: 30},
	})
	pipeline.Write(dir, i, pipeline.Translate(i))

	// Now simulate a blip: empty guarantees
	iEmpty := newIntermediate("aws", "s3", "https://aws.amazon.com/s3/sla/", nil)
	blipContent := pipeline.Translate(iEmpty) // starts with "#"

	result, err := pipeline.Write(dir, iEmpty, blipContent)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if result.Status != pipeline.StatusBlip {
		t.Errorf("want StatusBlip, got %v", result.Status)
	}

	// Original file should be untouched (non-empty, still has valid content)
	data, _ := os.ReadFile(result.Path)
	if len(data) == 0 {
		t.Error("blip should not overwrite existing file")
	}
}
