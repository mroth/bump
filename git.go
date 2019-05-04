package main

import (
	"errors"
	"regexp"

	"gopkg.in/src-d/go-git.v4"
)

// Returns owner, repo, error
//
// Errors are likely just be a simple "not in git repo" etc and should be
// considered informational rather than fatal.
func githubRepoDetect(path string) (owner, repo string, err error) {
	gitRepo, err := git.PlainOpen(path)
	if err != nil {
		return
	}

	remote, err := gitRepo.Remote("origin")
	if err != nil {
		return
	}

	fetchURL := remote.Config().URLs[0]
	owner, repo, ok := parseGithubRemote(fetchURL)
	if !ok {
		err = errors.New("non-matching remote url")
	}
	return
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
