package main

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver"
	"github.com/chzyer/readline"
	"github.com/google/go-github/v25/github"
	"github.com/manifoldco/promptui"
)

type cliVersionOption struct {
	Name    string
	Version semver.Version
}

func (o cliVersionOption) String() string {
	return fmt.Sprintf(
		"%v%v",
		o.Name,
		promptui.Styler(promptui.FGFaint)(
			fmt.Sprintf(" (%v)", o.Version.String()),
		),
	)
}

func prompt(
	owner, repo string, currVersion *semver.Version,
	release *github.RepositoryRelease) (*semver.Version, error) {

	fmt.Printf("🌻 Current version of %v (released %v)\n",
		promptui.Styler(promptui.FGBold)(fmt.Sprintf("%v/%v: %v",
			owner, repo, currVersion)),
		release.GetPublishedAt().Format("2006 Jan 2"),
	)

	// promptui.IconInitial = "🚀" // default is colored ASCII question mark
	choices := []cliVersionOption{
		{"patch", currVersion.IncPatch()},
		{"minor", currVersion.IncMinor()},
		{"major", currVersion.IncMajor()},
	}

	prompt := promptui.Select{
		Label: "Select semver increment to specify new version",
		Items: choices,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	nextVersion := choices[index].Version
	return &nextVersion, nil
}

// below is all boilerplate copy and pasted to workaround bell issue documented
// in https://github.com/manifoldco/promptui/issues/49. :-(

// stderr implements an io.WriteCloser that skips the terminal bell character
// (ASCII code 7), and writes the rest to os.Stderr. It's used to replace
// readline.Stdout, that is the package used by promptui to display the prompts.
type stderr struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal
// bell character.
func (s *stderr) Write(b []byte) (int, error) {
	if len(b) == 1 && b[0] == readline.CharBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriterCloser over os.Stderr.
func (s *stderr) Close() error {
	return os.Stderr.Close()
}

func init() {
	readline.Stdout = &stderr{}
}
