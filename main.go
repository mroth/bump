package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Masterminds/semver"
	"github.com/pkg/browser"
)

// build set by goreleaser on production builds
var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

// VerboseLogging sets whether to log debug/timing info to stderr
var VerboseLogging = false

func logVerbose(format string, v ...interface{}) {
	if VerboseLogging {
		log.Printf(format, v...)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	logVerbose("TIMING: %s took %s", name, elapsed)
}

func main() {
	owner, repo, opts := ParseAll()
	VerboseLogging = opts.Verbose
	logVerbose("BUILD INFO: %v %v %v", buildVersion, buildCommit, buildDate)
	logVerbose("ParseAll() opts: %+v owner: %v repo: %v", opts, owner, repo)

	// figure out owner and repo
	//  ...if we got it passed to us already, cool cool
	//  ...if not, call githubRepoDetect() to do our git checking magic
	if owner == "" || repo == "" {
		logVerbose("owner/repo not specified, checking for local git repo")
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
			logVerbose("%v", err)
			usage()
		}
		logVerbose("wd detected as git repo with github remote %v/%v", owner, repo)
	}

	// get latest release version from github
	logVerbose("checking github for latest release of %v/%v", owner, repo)
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

	// invoke interactive prompt UI helping user select next version
	nextVersion, err := prompt(owner, repo, version, release)
	if err != nil {
		log.Fatal(err)
	}

	// create draft URL for next version, send user to visit it
	nextURL := releaseURL(owner, repo, nextVersion)
	if opts.NoOpen {
		fmt.Println("To draft release:", nextURL)
	} else {
		fmt.Println("✨ Opening", nextURL)
		err = browser.OpenURL(nextURL)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func releaseURL(owner, repo string, version *semver.Version) string {
	return fmt.Sprintf(
		"https://github.com/%s/%s/releases/new?tag=v%s&title=v%s",
		owner, repo, version.String(), version.String(),
	)
}
