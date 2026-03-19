package pipeline_test

import (
	"strings"
	"testing"

	"barista/internal/pipeline"
)

func TestTranslate_WithGuarantees(t *testing.T) {
	i := &pipeline.Intermediate{
		ServiceName:  "s3",
		ProviderName: "aws",
		SourceURL:    "https://aws.amazon.com/s3/sla/",
		Guarantees: []pipeline.Guarantee{
			{Name: "monthly_uptime_percentage", Threshold: 99.9, WindowDays: 30},
		},
	}

	got := pipeline.Translate(i)
	want := "# === Expectations ===\n" +
		"Unmeasured Expectations\n" +
		"  # Source: https://aws.amazon.com/s3/sla/\n" +
		"\n" +
		"  * \"Monthly Uptime Percentage\":\n" +
		"    Provides {\n" +
		"      threshold: 99.9%,\n" +
		"      window_in_days: 30\n" +
		"    }"

	if got != want {
		t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
	}
}

func TestTranslate_IntegerThreshold(t *testing.T) {
	i := &pipeline.Intermediate{
		ServiceName:  "rds",
		ProviderName: "aws",
		SourceURL:    "https://aws.amazon.com/rds/sla/",
		Guarantees: []pipeline.Guarantee{
			{Name: "monthly_uptime", Threshold: 99.0, WindowDays: 30},
		},
	}

	got := pipeline.Translate(i)
	if !contains(got, "threshold: 99%,") {
		t.Errorf("expected integer threshold formatting, got:\n%s", got)
	}
}

func TestTranslate_TitleCase(t *testing.T) {
	i := &pipeline.Intermediate{
		ServiceName:  "s3",
		ProviderName: "aws",
		SourceURL:    "https://aws.amazon.com/s3/sla/",
		Guarantees: []pipeline.Guarantee{
			{Name: "monthly_uptime_percentage", Threshold: 99.9, WindowDays: 30},
		},
	}

	got := pipeline.Translate(i)
	if !contains(got, "\"Monthly Uptime Percentage\"") {
		t.Errorf("expected title case name, got:\n%s", got)
	}
}

func TestTranslate_MultipleGuarantees(t *testing.T) {
	i := &pipeline.Intermediate{
		ServiceName:  "s3",
		ProviderName: "aws",
		SourceURL:    "https://aws.amazon.com/s3/sla/",
		Guarantees: []pipeline.Guarantee{
			{Name: "monthly_uptime", Threshold: 99.9, WindowDays: 30},
			{Name: "weekly_uptime", Threshold: 99.99, WindowDays: 7},
		},
	}

	got := pipeline.Translate(i)
	for _, s := range []string{"Monthly Uptime", "Weekly Uptime", "99.9%", "99.99%"} {
		if !contains(got, s) {
			t.Errorf("expected output to contain %q, got:\n%s", s, got)
		}
	}
}

func TestTranslate_NoGuarantees(t *testing.T) {
	i := &pipeline.Intermediate{
		ServiceName:  "lambda",
		ProviderName: "aws",
		SourceURL:    "https://aws.amazon.com/lambda/sla/",
		Guarantees:   []pipeline.Guarantee{},
	}

	got := pipeline.Translate(i)
	if !contains(got, "# No guarantees found for aws/lambda") {
		t.Errorf("unexpected output: %q", got)
	}
	if !contains(got, "https://aws.amazon.com/lambda/sla/") {
		t.Errorf("expected source URL in no-guarantees output: %q", got)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
