package main

import (
	"context"
	"os"
	"sort"
	"time"

	"github.com/mroth/bump/internal/presemver"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

type recentReleases struct {
	Full   *github.RepositoryRelease // most recent full release only
	Latest *github.RepositoryRelease // most recent can include prereleases
}

// getRecentReleases will obtain the RecentReleases information for a given
// GitHub repository. It queries the GitHub API and then runs the results
// through our custom precedence ruleset.
func getRecentReleases(owner, repo string) (recentReleases, error) {
	releases, err := githubAPIListReleases(owner, repo)
	rr := extractRecentReleases(releases)
	return rr, err
}

func githubAPIListReleases(owner, repo string) ([]*github.RepositoryRelease, error) {
	client, ctx := defaultGithubClient(), context.Background()
	defer timeTrack(time.Now(), "API call to client.Repositories.ListReleases()")
	opts := &github.ListOptions{}
	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, opts)
	return releases, err
}

func extractRecentReleases(rs []*github.RepositoryRelease) recentReleases {
	res := recentReleases{}

	// Parse all the GitHub releases and filter for valid (parseable) semvers.
	// We will later use this actual semver for sorting rather than timestamps,
	// so we also will maintain a map to easily get back to the original GitHub
	// release data struct.
	col := semver.Collection{}
	releaseMap := make(map[*semver.Version]*github.RepositoryRelease)
	for _, r := range rs {
		rv, err := semver.NewVersion(r.GetTagName())
		if err != nil {
			logVerbose("WARNING: cannot parse %v as valid semver, ignoring...",
				r.GetTagName())
			continue
		}
		col = append(col, rv)
		releaseMap[rv] = r
	}
	logVerbose("releases according to GitHub: %v", col)

	// Sort semantic versions in reverse, which will produce an ordering based
	// on the semver rules, rather than the GitHub release timestamps.
	sort.Sort(sort.Reverse(col))
	logVerbose("reordered releases to semver: %v", col)

	// Assuming we have at least one version, the first one is now the "most
	// recent" release, whether or not it is a prerelease.
	if len(col) >= 1 {
		res.Latest = releaseMap[col[0]]
	}

	// Iterate backwards until we find the first release that's NOT considered a
	// semver pre-release. we need to parse on our own and not trust GitHub's
	// data here, since people manually maintaining releases often dont bother
	// to check the prerelease toggle.
	for _, v := range col {
		if !presemver.HasPrerelease(*v) {
			res.Full = releaseMap[v]
			break
		}
	}

	return res
}

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
