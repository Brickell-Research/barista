package pipeline

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type changelogStatus string

const (
	changelogInitial   changelogStatus = "INITIAL"
	changelogChanged   changelogStatus = "CHANGED"
	changelogBlip      changelogStatus = "BLIP"
	changelogUnchanged changelogStatus = "UNCHANGED"
)

type changelogEntry struct {
	Time      time.Time
	Status    changelogStatus
	SourceURL string
	Note      string
	OldContent string
	NewContent string
}

func appendChangelog(path string, e changelogEntry) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var sb strings.Builder

	fmt.Fprintf(&sb, "\n## %s %s\n", e.Time.Format(time.RFC3339), e.Status)
	if e.SourceURL != "" {
		fmt.Fprintf(&sb, "Source: %s\n", e.SourceURL)
	}
	if e.Note != "" {
		fmt.Fprintf(&sb, "Note: %s\n", e.Note)
	}

	switch e.Status {
	case changelogInitial:
		sb.WriteString("\n")
		for _, line := range strings.Split(e.NewContent, "\n") {
			fmt.Fprintf(&sb, "  %s\n", line)
		}
	case changelogChanged:
		sb.WriteString("\n")
		sb.WriteString(lineDiff(e.OldContent, e.NewContent))
	}

	_, err = f.WriteString(sb.String())
	return err
}
