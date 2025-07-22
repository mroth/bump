package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v29/github"
)

var update = flag.Bool("update", false, "update golden files")

// currentDate represents the current time for test data generation
var currentDate = time.Date(2025, 7, 22, 20, 56, 6, 0, time.UTC)

// timePtr returns a pointer to the given time - helper for creating *time.Time values
func timePtr(t time.Time) *time.Time {
	return &t
}

// testCommitsComparisons contains sample data for testing changelog functions
// NOTE: this test data is AI generated and for testing purposes only.
var testCommitsComparisons = map[string]*github.CommitsComparison{
	"sample": {
		HTMLURL: github.String("https://github.com/owner/repo/compare/v1.0.0...v1.1.0"),
		Commits: []github.RepositoryCommit{
			{
				SHA: github.String("a1b2c3d4e5f6789012345678901234567890abcd"),
				Commit: &github.Commit{
					Message: github.String("feat: add new user authentication system\n\nImplemented OAuth2 integration with Google and GitHub providers.\nAdded user session management and JWT token handling."),
					Author: &github.CommitAuthor{
						Name:  github.String("Alice Johnson"),
						Email: github.String("alice@example.com"),
						Date:  timePtr(currentDate.Add(-5 * time.Minute)),
					},
				},
			},
			{
				SHA: github.String("b2c3d4e5f6789012345678901234567890abcdef"),
				Commit: &github.Commit{
					Message: github.String("fix: resolve memory leak in background worker\n\nFixed goroutine leak that was causing memory usage to grow over time.\nImproved error handling and added proper cleanup."),
					Author: &github.CommitAuthor{
						Name:  github.String("Bob Smith"),
						Email: github.String("bob@example.com"),
						Date:  timePtr(currentDate.Add(-2 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("c3d4e5f6789012345678901234567890abcdef12"),
				Commit: &github.Commit{
					Message: github.String("docs: update API documentation"),
					Author: &github.CommitAuthor{
						Name:  github.String("Carol Williams"),
						Email: github.String("carol@example.com"),
						Date:  timePtr(currentDate.Add(-6 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("d4e5f6789012345678901234567890abcdef1234"),
				Commit: &github.Commit{
					Message: github.String("test: add comprehensive unit tests for auth module"),
					Author: &github.CommitAuthor{
						Name:  github.String("Alice Johnson"),
						Email: github.String("alice@example.com"),
						Date:  timePtr(currentDate.Add(-1 * 24 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("e5f6789012345678901234567890abcdef123456"),
				Commit: &github.Commit{
					Message: github.String("refactor: simplify database connection pooling\n\nReduced complexity and improved performance by 15%."),
					Author: &github.CommitAuthor{
						Name:  github.String("David Brown"),
						Email: github.String("david@example.com"),
						Date:  timePtr(currentDate.Add(-3 * 24 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("f6789012345678901234567890abcdef12345678"),
				Commit: &github.Commit{
					Message: github.String("feat: implement rate limiting middleware"),
					Author: &github.CommitAuthor{
						Name:  github.String("Eva Davis"),
						Email: github.String("eva@example.com"),
						Date:  timePtr(currentDate.Add(-5 * 24 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("789012345678901234567890abcdef123456789a"),
				Commit: &github.Commit{
					Message: github.String("fix: handle edge case in date parsing\n\nFixed panic when parsing malformed timestamps from external APIs."),
					Author: &github.CommitAuthor{
						Name:  github.String("Bob Smith"),
						Email: github.String("bob@example.com"),
						Date:  timePtr(currentDate.Add(-7 * 24 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("89012345678901234567890abcdef123456789ab"),
				Commit: &github.Commit{
					Message: github.String("chore: update dependencies to latest versions"),
					Author: &github.CommitAuthor{
						Name:  github.String("Alice Johnson"),
						Email: github.String("alice@example.com"),
						Date:  timePtr(currentDate.Add(-2 * 7 * 24 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("9012345678901234567890abcdef123456789abc"),
				Commit: &github.Commit{
					Message: github.String("perf: optimize database queries for user lookup\n\nReduced average query time from 150ms to 45ms."),
					Author: &github.CommitAuthor{
						Name:  github.String("Frank Miller"),
						Email: github.String("frank@example.com"),
						Date:  timePtr(currentDate.Add(-3 * 7 * 24 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("012345678901234567890abcdef123456789abcd"),
				Commit: &github.Commit{
					Message: github.String("feat: add webhook support for external integrations"),
					Author: &github.CommitAuthor{
						Name:  github.String("Grace Wilson"),
						Email: github.String("grace@example.com"),
						Date:  timePtr(currentDate.Add(-6 * 7 * 24 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("12345678901234567890abcdef123456789abcde"),
				Commit: &github.Commit{
					Message: github.String("fix: correct timezone handling in scheduled tasks\n\nEnsured all scheduled operations respect user's local timezone."),
					Author: &github.CommitAuthor{
						Name:  github.String("David Brown"),
						Email: github.String("david@example.com"),
						Date:  timePtr(currentDate.Add(-2 * 30 * 24 * time.Hour)),
					},
				},
			},
			{
				SHA: github.String("2345678901234567890abcdef123456789abcdef"),
				Commit: &github.Commit{
					Message: github.String("style: format code according to new linting rules"),
					Author: &github.CommitAuthor{
						Name:  github.String("Eva Davis"),
						Email: github.String("eva@example.com"),
						Date:  timePtr(currentDate.Add(-3 * 30 * 24 * time.Hour)),
					},
				},
			},
		},
	},
}

func TestRenderChangelogScreen(t *testing.T) {
	for name, comparison := range testCommitsComparisons {
		t.Run(name, func(t *testing.T) {
			got := RenderChangelogScreen(comparison)

			goldenFile := filepath.Join("testdata", name+"_screen.golden")

			if *update {
				err := os.WriteFile(goldenFile, []byte(got), 0644)
				if err != nil {
					t.Fatalf("failed to update golden file: %v", err)
				}
				return
			}

			wantBytes, err := os.ReadFile(goldenFile)
			if err != nil {
				t.Fatalf("failed to read golden file %s: %v", goldenFile, err)
			}
			want := string(wantBytes)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("RenderChangelogScreen() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRenderChangelogMarkdown(t *testing.T) {
	for name, comparison := range testCommitsComparisons {
		t.Run(name, func(t *testing.T) {
			got := RenderChangelogMarkdown(comparison)

			goldenFile := filepath.Join("testdata", name+"_markdown.golden")

			if *update {
				err := os.WriteFile(goldenFile, []byte(got), 0644)
				if err != nil {
					t.Fatalf("failed to update golden file: %v", err)
				}
				return
			}

			wantBytes, err := os.ReadFile(goldenFile)
			if err != nil {
				t.Fatalf("failed to read golden file %s: %v", goldenFile, err)
			}
			want := string(wantBytes)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("RenderChangelogMarkdown() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
