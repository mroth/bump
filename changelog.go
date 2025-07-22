package main

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v29/github"
)

// RenderChangelogScreen formats a CommitsComparison suitable for displaying on the
// screen to the user, abbreviated to try to not overflow a 80x24 terminal.
//
// Because of this, we only display the 10 most recent commits, with a
// comparison URL targeting HEAD (as draft is not released), so user can view
// the full list on GitHub if desired.
func RenderChangelogScreen(comparison *github.CommitsComparison) string {
	const maxDisplayCommits = 10
	var buf strings.Builder
	buf.WriteString("Changes since previous release:\n\n")

	for _, c := range comparison.Commits[:min(maxDisplayCommits, len(comparison.Commits))] {
		fmt.Fprintf(&buf, "  - %v\n", firstCommitMsgLine(c))
	}

	if numExtraCommits := len(comparison.Commits) - maxDisplayCommits; numExtraCommits > 0 {
		fmt.Fprintf(&buf, "\n...%d more commits, %s\n", numExtraCommits, comparison.GetHTMLURL())
	}

	return buf.String()
}

// RenderChangelogMarkdown formats a CommitsComparison suitable for markdown display
// in a GitHub Flavored Markdown release notes field.
//
// TODO: cap max number of commits to display? API returns <=250
func RenderChangelogMarkdown(comparison *github.CommitsComparison) string {
	var buf strings.Builder
	buf.WriteString("## Changelog\n\n")

	for _, c := range comparison.Commits {
		fmt.Fprintf(&buf, "- %v %.7s\n", firstCommitMsgLine(c), c.GetSHA())
	}

	return buf.String()
}

func firstCommitMsgLine(c github.RepositoryCommit) string {
	msg := c.Commit.GetMessage()
	lines := strings.SplitN(msg, "\n", 2)
	return lines[0]
}

// comparisonURL makes a GitHub web view URL for comparing two tagged semvers.
func comparisonURL(owner, repo string, base, next *semver.Version) string {
	return fmt.Sprintf(
		"https://github.com/%s/%s/compare/v%s...v%s", owner, repo, base, next,
	)
}
