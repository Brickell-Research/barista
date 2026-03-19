package pipeline

import (
	"fmt"
	"strings"
)

func Translate(i *Intermediate) string {
	if len(i.Guarantees) == 0 {
		return fmt.Sprintf("# No guarantees found for %s/%s\n# Source: %s", i.ProviderName, i.ServiceName, i.SourceURL)
	}

	var sb strings.Builder

	sb.WriteString("# === Expectations ===\n")
	fmt.Fprintf(&sb, "Unmeasured Expectations\n")
	fmt.Fprintf(&sb, "  # Source: %s\n", i.SourceURL)

	for _, g := range i.Guarantees {
		sb.WriteString("\n")
		fmt.Fprintf(&sb, "  * %q:\n", toTitleCase(g.Name))
		sb.WriteString("    Provides {\n")
		fmt.Fprintf(&sb, "      threshold: %s%%,\n", formatThreshold(g.Threshold))
		fmt.Fprintf(&sb, "      window_in_days: %d\n", g.WindowDays)
		sb.WriteString("    }")
	}

	return sb.String()
}

func toTitleCase(s string) string {
	words := strings.Split(s, "_")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func formatThreshold(t float64) string {
	if t == float64(int(t)) {
		return fmt.Sprintf("%d", int(t))
	}
	return fmt.Sprintf("%g", t)
}
