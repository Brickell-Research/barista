package pipeline

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type WriteStatus int

const (
	StatusWritten   WriteStatus = iota // new or changed content written
	StatusUnchanged                    // content identical, skipped
	StatusBlip                         // zero guarantees but previous file had some; skipped
)

type WriteResult struct {
	Path   string
	Status WriteStatus
}

// Write writes a translated Caffeine file for the given intermediate. It
// compares against any existing file and:
//   - Skips the write if content is unchanged
//   - Skips the write and logs a BLIP if the new output has no guarantees
//     but the current file on disk does
//   - Writes the new file and appends a changelog entry otherwise
func Write(outputDir string, i *Intermediate, content string) (WriteResult, error) {
	dir := filepath.Join(outputDir, "expectations")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return WriteResult{}, fmt.Errorf("create output dir: %w", err)
	}

	filePath := filepath.Join(dir, i.ServiceName+".caffeine")
	changelogPath := filepath.Join(dir, i.ServiceName+".changelog")

	existing, err := os.ReadFile(filePath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return WriteResult{}, fmt.Errorf("read existing file: %w", err)
	}

	oldContent := string(existing)
	firstWrite := len(existing) == 0
	blip := isNoGuarantees(content) && !isNoGuarantees(oldContent) && !firstWrite
	unchanged := oldContent == content

	switch {
	case blip:
		logChangelog(changelogPath, changelogEntry{
			Time:      time.Now().UTC(),
			Status:    changelogBlip,
			SourceURL: i.SourceURL,
			Note:      "Fetch returned 0 guarantees but previous content exists; not overwriting.",
		})
		return WriteResult{Path: filePath, Status: StatusBlip}, nil

	case unchanged:
		return WriteResult{Path: filePath, Status: StatusUnchanged}, nil
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return WriteResult{}, fmt.Errorf("write file: %w", err)
	}

	status := changelogChanged
	if firstWrite {
		status = changelogInitial
	}
	logChangelog(changelogPath, changelogEntry{
		Time:       time.Now().UTC(),
		Status:     status,
		SourceURL:  i.SourceURL,
		OldContent: oldContent,
		NewContent: content,
	})

	return WriteResult{Path: filePath, Status: StatusWritten}, nil
}

func isNoGuarantees(content string) bool {
	return strings.HasPrefix(strings.TrimSpace(content), "#")
}

// logChangelog appends a changelog entry, logging to stderr on failure
// rather than failing the write itself.
func logChangelog(path string, e changelogEntry) {
	if err := appendChangelog(path, e); err != nil {
		fmt.Fprintf(os.Stderr, "warn: changelog write failed for %s: %v\n", path, err)
	}
}
