package pipeline

import (
	"fmt"
	"strings"
)

func Translate(i *Intermediate) string {
	if len(i.Guarantees) == 0 {
		return fmt.Sprintf("# No guarantees found for %s/%s", i.ProviderName, i.ServiceName)
	}

	var sb strings.Builder
	sb.WriteString("Expectations")

	for _, g := range i.Guarantees {
		sb.WriteString("\n\n")
		fmt.Fprintf(&sb, "%q:\n", g.Name)
		fmt.Fprintf(&sb, "  Guarantees %s%% over %dd window", formatThreshold(g.Threshold), g.WindowDays)
	}

	return sb.String()
}

func formatThreshold(t float64) string {
	if t == float64(int(t)) {
		return fmt.Sprintf("%d", int(t))
	}
	return fmt.Sprintf("%g", t)
}
