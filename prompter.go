package main

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v25/github"
	"github.com/manifoldco/promptui"
)

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

func prompt(owner, repo string, currVersion *semver.Version, release *github.RepositoryRelease) (*semver.Version, error) {
	fmt.Printf("🌻 Current version of %v (released %v)\n",
		promptui.Styler(promptui.FGBold)(fmt.Sprintf("%v/%v: %v",
			owner, repo, currVersion)),
		release.GetPublishedAt(),
	)
	// promptui.IconInitial = "🚀"
	choices := []cliVersionOption{
		{"patch", currVersion.IncPatch()},
		{"minor", currVersion.IncMinor()},
		{"major", currVersion.IncMajor()},
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
		return nil, err
	}
	nextVersion := choices[index].Version
	return &nextVersion, nil
}
