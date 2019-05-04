package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v25/github"
	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

// type Options struct {
// 	NoOpen        bool // dont auto-open the final URL in browser
// 	NoInteractive bool // TODO: disable interactive mode? (implies --no-open)
// 	Verbose       bool // TODO verbose mode?
// }

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <owner> <repo>\n", os.Args[0])
	// fmt.Fprintf(os.Stderr, "Usage: %s <owner> <repo> [major|minor|patch]\n", os.Args[0])
}

type cliVersionOption struct {
	Name    string
	Version semver.Version
}

func (o cliVersionOption) String() string {
	return fmt.Sprintf(
		"%v %v",
		o.Name,
		promptui.Styler(promptui.FGFaint)(
			fmt.Sprintf("(%v)", o.Version.String()),
		),
	)
}

func main() {
	var owner, repo string
	if len(os.Args) < 3 {
		log.Println("owner/repo not specified, checking for local git repo")
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err) // actually something weird going on
		}
		owner, repo, err = githubRepoDetect(wd)
		if err != nil { // probably just not in a git repo
			log.Println(err) // TODO: only verbose
			usage()
			os.Exit(1)
		}
		log.Printf("workdir detected as git repo with github remote %v/%v", owner, repo) // TODO: verbose
	} else {
		owner, repo = os.Args[1], os.Args[2]
	}

	// if len(os.Args) >= 4 {
	// 	switch os.Args[3] {
	// 	case ""
	// 	}
	// }

	log.Printf("checking github for latest release of %v/%v", owner, repo) // TODO: verbose
	release, err := getLatestRelease(owner, repo)
	if err != nil {
		log.Fatal(err)
	}

	tag := release.GetTagName()
	version, err := semver.NewVersion(tag)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🌻 Current version of %v (released %v)\n",
		promptui.Styler(promptui.FGBold)(fmt.Sprintf("%v/%v: %v", owner, repo, version)),
		release.GetPublishedAt(),
	)
	// promptui.IconInitial = "🚀"
	options := []cliVersionOption{
		{"patch", version.IncPatch()},
		{"minor", version.IncMinor()},
		{"major", version.IncMajor()},
	}
	prompt := promptui.Select{
		Label: "Select semver increment to specify new version",
		Items: options,
		// Templates: &promptui.SelectTemplates{
		// Active: `🚀 {{ . | red }}`,
		// Help: `{{ "Use the arrow (or vim) keys to navigate: ↓ ↑ → ←" | faint }}`,
		// },
	}

	index, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	nextVersion := options[index].Version
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

func releaseURL(owner, repo string, version semver.Version) string {
	return fmt.Sprintf(
		"https://github.com/%s/%s/releases/new?tag=v%s&title=v%s",
		owner, repo, version.String(), version.String(),
	)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("TIMING: %s took %s", name, elapsed)
}
