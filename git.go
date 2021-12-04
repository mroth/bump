package main

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
)

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
		err = errors.New("cannot pattern match remote url: " + remoteURL)
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
// bytes adds 218kb
// benchmarks at 5.1 ms/op
//
// NOTE: not currently used at all, was only here for benchmarking purposes.
// FIXME: Does not respect path (known issue, would need address if using this in future).
func _detectRemoteURL_LocalGit(path string) (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(output)), nil
}

// parseGithubRemote parses string remoteURL against known patterns matching
// GitHub remotes and returns the owner and repo, along with a boolean ok
// indicating whether a match was found.
//
// Possible GitHub remote formats (HTTPS/SSH):
// 	 https://github.com/mroth/bump.git
//   git@github.com:mroth/bump.git
func parseGithubRemote(remoteURL string) (owner, repo string, ok bool) {
	re := regexp.MustCompile(`^(?:https://|git@)github.com[:/](.*)/(.*?)(?:\.git$|$)`)
	matches := re.FindStringSubmatch(remoteURL)
	if matches == nil || len(matches) < 3 {
		return
	}
	return matches[1], matches[2], true
}
