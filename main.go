package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v25/github"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

// TODO: maybe use a package scope logger var instead?
// var (
// 	Verbose = false
// )

func main() {
	owner, repo, opts := ParseAll()
	if opts.Verbose {
		log.Printf("ParseAll() opts: %+v owner: %v repo: %v", opts, owner, repo)
	}

	// figure out owner and repo
	//  ...if we got it passed to us already, cool cool
	//  ...if not, call githubRepoDetect() to do our git checking magic
	if owner == "" || repo == "" {
		if opts.Verbose {
			log.Println("owner/repo not specified, checking for local git repo")
		}
		wd, err := os.Getwd()
		if err != nil {
			// couldn't get working directory, something really weird going on
			// we should just fatal in this case
			log.Fatal(err)
		}
		owner, repo, err = githubRepoDetect(wd)
		if err != nil {
			// probably just not in a git repo, no biggie
			// just log what happened in verbose mode, and show usage
			if opts.Verbose {
				log.Println(err)
			}
			usage()
		}
		if opts.Verbose {
			log.Printf("workdir detected as git repo with github remote %v/%v", owner, repo)
		}
	}

	if opts.Verbose {
		log.Printf("checking github for latest release of %v/%v", owner, repo)
	}
	release, err := getLatestRelease(owner, repo)
	if err != nil {
		log.Fatal(err)
	}

	// try to parse tag name from current release into a semantic version
	tag := release.GetTagName()
	version, err := semver.NewVersion(tag)
	if err != nil {
		log.Fatal(err)
	}

	nextVersion, err := prompt(owner, repo, version, release)
	if err != nil {
		log.Fatal(err)
	}

	nextURL := releaseURL(owner, repo, nextVersion)
	fmt.Println("Open sesame:", nextURL)
	browser.OpenURL(nextURL)
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

func getLatestRelease(owner, repo string) (*github.RepositoryRelease, error) {
	client := defaultGithubClient()
	ctx := context.Background()
	defer timeTrack(time.Now(), "client.Repositories.GetLatestRelease()")
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	return release, err
}

func releaseURL(owner, repo string, version *semver.Version) string {
	return fmt.Sprintf(
		"https://github.com/%s/%s/releases/new?tag=v%s&title=v%s",
		owner, repo, version.String(), version.String(),
	)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("TIMING: %s took %s", name, elapsed)
}
