package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/browser"
)

// build info set by goreleaser during production builds
var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

// VerboseLogging sets whether to log debug/timing info to stderr
var VerboseLogging = false

func logVerbose(format string, v ...any) {
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
		// TODO: possibly warn if have commits ahead of remote locally?
		if err != nil {
			// probably just not in a git repo, no biggie
			// just log what happened in verbose mode, and show usage
			logVerbose("%v", err)
			usage()
		}
		logVerbose("detected .git repo with github remote %v/%v", owner, repo)
	}

	// get latest release version from github
	logVerbose("checking github for latest release of %v/%v", owner, repo)
	previousRelease, err := getLatestRelease(owner, repo)
	if err != nil {
		log.Fatal(err)
	}

	// try to parse tag name from current release into a semantic version
	previousVersion, err := semver.NewVersion(previousRelease.GetTagName())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("🌻 Latest release of %v (published %v)\n",
		boldStyler(fmt.Sprintf("%v/%v: %v", owner, repo, previousVersion)),
		previousRelease.GetPublishedAt().Format("2006 Jan 2"),
	)

	// retrieve changes since last release via GitHub API
	comparison, err := compareRelease(owner, repo, previousRelease.GetTagName())
	if err != nil {
		log.Fatal("failed to retrieve commits", err)
	}

	// display abbreviated changelog to user in CLI, to hopefully aide them in
	// making a decision about what the next semver should be.
	changelog := screenChangelog(comparison)
	fmt.Println(changelog)

	// invoke interactive prompt UI allowing user to select next version
	nextVersion, err := prompt(previousVersion)
	if err != nil {
		log.Fatal(err)
	}

	// create draft URL embedding markdown changelog for next version...
	body := strings.Join([]string{
		markdownChangelog(comparison),
		comparisonURL(owner, repo, previousVersion, nextVersion),
	}, "\n")
	draftURL := draftReleaseURL(owner, repo, nextVersion, body)

	// ...then send user to visit in their web browser!
	if opts.NoOpen {
		fmt.Println("To draft release, visit:", draftURL)
	} else {
		fmt.Println("✨ Drafting new release on GitHub!")
		logVerbose("Opening browser to: %s", draftURL)
		err = browser.OpenURL(draftURL)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// draftReleaseURL constructs a URL to open a new draft release on GitHub for
// given owner/repo with a semver compatible tag based on the semver.Version in
// the tag and title fields, and an encoded body payload to prepopulate the
// form.
func draftReleaseURL(owner, repo string, version *semver.Version, body string) string {
	return fmt.Sprintf(
		"https://github.com/%s/%s/releases/new?tag=v%s&title=v%s&body=%s",
		owner, repo, version.String(), version.String(), url.QueryEscape(body),
	)
}

// markdownChangelog formats a CommitsComparison suitable for displaying on the
// screen to the user, abbreviated to try to not overflow a 80x24 terminal.
//
// Because of this, we only display the 10 most recent commits, with a
// comparison URL targeting HEAD (as draft is not released), so user can view
// the full list on GitHub if desired.
func screenChangelog(comparison *github.CommitsComparison) string {
	var buf bytes.Buffer
	buf.WriteString("Changes since previous release:\n\n")
	const max = 10
	for i, c := range comparison.Commits {
		if i >= max {
			break
		}
		buf.WriteString(
			fmt.Sprintf("  - %v\n",
				strings.Split(c.Commit.GetMessage(), "\n")[0],
			),
		)
	}
	numExtraCommits := len(comparison.Commits) - max
	if numExtraCommits > 0 {
		buf.WriteString(
			fmt.Sprintf("\n...%d more commits, %s\n",
				numExtraCommits, comparison.GetHTMLURL()),
		)
	}
	return buf.String()
}

// markdownChangelog formats a CommitsComparison suitable for markdown display
// in a GitHub Flavored Markdown release notes field.
//
// TODO: cap max number of commits to display? API returns <=250
func markdownChangelog(comparison *github.CommitsComparison) string {
	var buf bytes.Buffer
	buf.WriteString("## Changelog\n\n")
	for _, c := range comparison.Commits {
		buf.WriteString(
			fmt.Sprintf("- %v %.7s\n",
				strings.Split(c.Commit.GetMessage(), "\n")[0],
				c.GetSHA(),
			),
		)
	}
	return buf.String()
}

// comparisonURL makes a GitHub web view URL for comparing two tagged semvers.
func comparisonURL(owner, repo string, base, next *semver.Version) string {
	return fmt.Sprintf(
		"https://github.com/%s/%s/compare/v%s...v%s", owner, repo, base, next,
	)
}
