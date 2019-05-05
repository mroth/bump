package main

import (
	"errors"
	"os/exec"
	"regexp"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

// Returns owner, repo, error
//
// Errors are likely just be a simple "not in git repo" etc and should be
// considered informational rather than fatal.
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

// implementation using go-git
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

// implementation shelling out to local copy of git
// requires
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
