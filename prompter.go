package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/manifoldco/promptui"
)

var (
	boldStyler  = promptui.Styler(promptui.FGBold)
	faintStyler = promptui.Styler(promptui.FGFaint)
)

type cliVersionOption struct {
	Name        string
	Version     *semver.Version
	Description string
}

func (o cliVersionOption) String() string {
	return fmt.Sprintf("%v%v",
		o.Name, faintStyler(fmt.Sprintf(" (%v)", o.Version.String())),
	)
}

func prompt(possibilities semver.Collection) (*semver.Version, error) {
	// promptui.IconInitial = "🚀" // default is colored ASCII question mark

	// quick fix: dynamically build cliVersionOption choices from all the
	// Versions in a Collection.  All the information is now contained in the
	// base Version type and helper functions so this is mostly busywork that
	// can be refactored/simplified in the future.
	var choices []cliVersionOption
	for _, c := range possibilities {
		verType := Type(*c)
		choices = append(choices, cliVersionOption{
			strings.ToLower(verType.String()), c, verType.Description(),
		})
	}

	prompt := promptui.Select{
		Label: "Select semver increment to specify next version",
		Items: choices,
		Templates: &promptui.SelectTemplates{
			Details: `{{ .Name }}: {{ .Description }}.`,
		},
		Stdout: &bellSkipper{},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	nextVersion := choices[index].Version
	return nextVersion, nil
}

// bellSkipper implements an io.WriteCloser that skips the terminal bell
// character (ASCII code 7), and writes the rest to os.Stderr. It is used to
// replace readline.Stdout, that is the package used by promptui to display the
// prompts.
//
// This is a workaround for the bell issue documented in
// https://github.com/manifoldco/promptui/issues/49.
type bellSkipper struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal
// bell character.
func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriterCloser over os.Stderr.
func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}
