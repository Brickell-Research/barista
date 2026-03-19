package pipeline

import (
	"fmt"
	"strings"
)

// lineDiff returns a unified-style diff of two multi-line strings.
func lineDiff(oldContent, newContent string) string {
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	common := lcs(oldLines, newLines)

	var sb strings.Builder
	i, j, k := 0, 0, 0

	for k < len(common) {
		for i < len(oldLines) && oldLines[i] != common[k] {
			fmt.Fprintf(&sb, "- %s\n", oldLines[i])
			i++
		}
		for j < len(newLines) && newLines[j] != common[k] {
			fmt.Fprintf(&sb, "+ %s\n", newLines[j])
			j++
		}
		fmt.Fprintf(&sb, "  %s\n", common[k])
		i++
		j++
		k++
	}
	for ; i < len(oldLines); i++ {
		fmt.Fprintf(&sb, "- %s\n", oldLines[i])
	}
	for ; j < len(newLines); j++ {
		fmt.Fprintf(&sb, "+ %s\n", newLines[j])
	}

	return sb.String()
}

// lcs returns the longest common subsequence of two string slices.
func lcs(a, b []string) []string {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] > dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}

	result := make([]string, 0, dp[m][n])
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			result = append([]string{a[i-1]}, result...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}
	return result
}
