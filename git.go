package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/google/go-github/v25/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
)

// getLatestRelease is a convenience function wrapping retrieval of latest
// GitHub release for owner and repo
//
// It will automatically use an OAuth scoped token if GITHUB_RELEASE environment
// variable is set, or an unauthed client otherwise.
func getLatestRelease(owner, repo string) (*github.RepositoryRelease, error) {
	client, ctx := defaultGithubClient(), context.Background()
	defer timeTrack(time.Now(), "client.Repositories.GetLatestRelease()")
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
	defer timeTrack(time.Now(), "client.Repositories.CompareCommits()")
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
//
// TODO: actually test me :-)
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

// githubRepoDetect attempts to detect whether a given path is part of a git
// repository that has a GitHub remote as the origin, and if so, returns the
// owner and repo name.
//
// Errors returned are likely just be a simple "not in git repo" etc and should
// be considered informational rather than fatal.
func githubRepoDetect(path string) (owner, repo string, err error) {
	defer timeTrack(time.Now(), "githubRepoDetect()")
	remoteURL, err := _detectRemoteURL_GoGit(path)
	// remoteURL, err := _detectRemoteURL_LocalGit(path)
	if err != nil {
		return
	}
	owner, repo, ok := parseGithubRemote(remoteURL)
	if !ok {
		err = errors.New("non-matching remote url")
	}
	return
}

// detectRemoteURL implementation using go-git
//
// will work even if git is not installed on users machine
// one more dependency to track and keep up to date
//
// go-git adds ~4mb to macOS binary size (from 12MB->16MB ugh)
// benchmarks at 0.1 ms/op
func _detectRemoteURL_GoGit(path string) (string, error) {
	gitRepo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}
	remote, err := gitRepo.Remote("origin")
	if err != nil {
		return "", err
	}
	return remote.Config().URLs[0], nil
}

// detectRemoteURL implementation shelling out to local copy of git
//
// requires git to be installed on machine
// uses os/exec from standard library, does not add a dependency
//
// os/exec adds 242KB to macOS binary size
// benchmarks at 4.2 ms/op
func _detectRemoteURL_LocalGit(path string) (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// parseGithubRemote parses string remoteURL against known patterns matching
// GitHub remotes and returns the owner and repo, along with a boolean ok
// indicating whether a match was found.
//
// Possible GitHub remote formats (HTTPS/SSH):
// 	 https://github.com/mroth/bump.git
//   git@github.com:mroth/bump.git
func parseGithubRemote(remoteURL string) (owner, repo string, ok bool) {
	re := regexp.MustCompile(`^(?:https://|git@)github.com[:/](.*)/(.*)\.git`)
	matches := re.FindStringSubmatch(remoteURL)
	if matches == nil || len(matches) < 3 {
		return
	}
	return matches[1], matches[2], true
}
