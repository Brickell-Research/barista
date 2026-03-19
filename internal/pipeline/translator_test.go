package pipeline_test

import (
	"testing"

	"barista/internal/pipeline"
)

func TestTranslate_WithGuarantees(t *testing.T) {
	i := &pipeline.Intermediate{
		ServiceName:  "s3",
		ProviderName: "aws",
		Guarantees: []pipeline.Guarantee{
			{Name: "monthly_uptime_percentage", Threshold: 99.9, WindowDays: 30},
		},
	}

	got := pipeline.Translate(i)
	want := "Expectations\n\n\"monthly_uptime_percentage\":\n  Guarantees 99.9% over 30d window"

	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestTranslate_MultipleGuarantees(t *testing.T) {
	i := &pipeline.Intermediate{
		ServiceName:  "s3",
		ProviderName: "aws",
		Guarantees: []pipeline.Guarantee{
			{Name: "monthly_uptime", Threshold: 99.9, WindowDays: 30},
			{Name: "weekly_uptime", Threshold: 99.99, WindowDays: 7},
		},
	}

	got := pipeline.Translate(i)
	if got == "" {
		t.Error("expected non-empty output")
	}
	// Both guarantees should appear
	for _, name := range []string{"monthly_uptime", "weekly_uptime"} {
		if !contains(got, name) {
			t.Errorf("expected output to contain %q", name)
		}
	}
}

func TestTranslate_IntegerThreshold(t *testing.T) {
	i := &pipeline.Intermediate{
		ServiceName:  "rds",
		ProviderName: "aws",
		Guarantees: []pipeline.Guarantee{
			{Name: "monthly_uptime", Threshold: 99.0, WindowDays: 30},
		},
	}

	got := pipeline.Translate(i)
	want := "Expectations\n\n\"monthly_uptime\":\n  Guarantees 99% over 30d window"

	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestTranslate_NoGuarantees(t *testing.T) {
	i := &pipeline.Intermediate{
		ServiceName:  "lambda",
		ProviderName: "aws",
		Guarantees:   []pipeline.Guarantee{},
	}

	got := pipeline.Translate(i)
	want := "# No guarantees found for aws/lambda"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsRune(s, substr))
}

func containsRune(s, substr string) bool {
	for i := range s {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
