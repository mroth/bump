package main

import (
	"context"
	"os"
	"time"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

// getLatestRelease is a convenience function wrapping retrieval of latest
// GitHub release for owner and repo
//
// It will automatically use an OAuth scoped token if GITHUB_RELEASE environment
// variable is set, or an unauthed client otherwise.
func getLatestRelease(owner, repo string) (*github.RepositoryRelease, error) {
	client, ctx := defaultGithubClient(), context.Background()
	defer timeTrack(time.Now(), "API call to client.Repositories.GetLatestRelease()")
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	return release, err
}

// compareRelease is a convenience function wrapping retrieval of a commits
// comparison between current default branch HEAD and a given release tagName,
// via the GitHub API.
//
// It will automatically use an OAuth scoped token if GITHUB_RELEASE environment
// variable is set, or an unauthed client otherwise.
//
// Unlike git log base..head, which returns commits in reverse chronlogical
// order, the GitHub V3 API returns in chronological order. To do what most
// changelog generators do and show commits in reverse chronological, we just
// mutate the actual CommitsComparison.Commits subfield prior to returning it
// here.
//
// Interestingly enough, this is the exact opposite of what is claimed in the
// GitHub API documentation, which says log is chronological, see:
// https://developer.github.com/v3/repos/commits/#compare-two-commits.
func compareRelease(owner, repo, tagName string) (*github.CommitsComparison, error) {
	client, ctx := defaultGithubClient(), context.Background()
	defer timeTrack(time.Now(), "API call to client.Repositories.CompareCommits()")
	cc, _, err := client.Repositories.CompareCommits(ctx, owner, repo, tagName, "HEAD")
	if cc != nil {
		reverseCommitOrder(cc)
	}
	return cc, err
}

func reverseCommitOrder(cc *github.CommitsComparison) {
	for i := len(cc.Commits)/2 - 1; i >= 0; i-- {
		opp := len(cc.Commits) - 1 - i
		cc.Commits[i], cc.Commits[opp] = cc.Commits[opp], cc.Commits[i]
	}
}

// defaultGithubClient returns a OAuth scoped Github API Client if GITHUB_TOKEN
// is set the local environment, or an unauthorized one otherwise.
func defaultGithubClient() *github.Client {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if ok {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		return github.NewClient(tc)
	}
	return github.NewClient(nil)
}
