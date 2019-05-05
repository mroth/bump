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

// TODO: maybe use a package scope logger var instead?
// var (
// 	Verbose = false
// )

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
	opts, flags := ParseFlags(NewOptionsFromEnv(), os.Args[1:])
	log.Printf("ParseFlags() new options: %+v", opts)

	// figure out owner and repo
	//  ...if we got it passed to us, cool cool
	//  ...if not, call githubRepoDetect() to do our git checking magic
	var owner, repo string
	if len(flags.Args()) < 2 {
		log.Println("owner/repo not specified, checking for local git repo")
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err) // actually something weird going on
		}
		owner, repo, err = githubRepoDetect(wd)
		if err != nil {
			// probably just not in a git repo, no biggie
			if opts.Verbose {
				log.Println(err)
			}
			usage()
			os.Exit(1)
		}
		if opts.Verbose {
			log.Printf("workdir detected as git repo with github remote %v/%v", owner, repo)
		}
	} else {
		owner, repo = flags.Arg(0), flags.Arg(1)
	}

	// TODO: deal with possible passed auto-action
	// if len(os.Args) >= 4 {
	// 	switch os.Args[3] {
	// 	case ""
	// 	}
	// }

	if opts.Verbose {
		log.Printf("checking github for latest release of %v/%v", owner, repo)
	}
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
	choices := []cliVersionOption{
		{"patch", version.IncPatch()},
		{"minor", version.IncMinor()},
		{"major", version.IncMajor()},
	}
	prompt := promptui.Select{
		Label: "Select semver increment to specify new version",
		Items: choices,
		// Templates: &promptui.SelectTemplates{
		// Active: `🚀 {{ . | red }}`,
		// Help: `{{ "Use the arrow (or vim) keys to navigate: ↓ ↑ → ←" | faint }}`,
		// },
	}

	index, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	nextVersion := choices[index].Version
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
