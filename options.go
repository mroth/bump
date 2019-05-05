package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const tmpusage = `
Usage: bump <owner> <repo> [major|minor|patch|interactive*]

If you are in a git repository that has been cloned from GitHub, owner and
repo args can be omitted, in which case they will be inferred from the remote
origin.

Flags:
    --no-open           Do not automatically open publish URL in browser.
    --verbose, -v       Verbose output.

Environment:
    $BUMP_NO_OPEN       Global default for --no-open
    $BUMP_VERBOSE       Global default for --verbose
`

func usage() {
	fmt.Fprintf(os.Stderr, tmpusage)
}

// Options defines all optional settings understood by the program.
//
// The zero value represents the program defaults.
type Options struct {
	NoOpen  bool // dont auto-open the final URL in browser
	Verbose bool // verbose output requested
}

// Environment variable "key" constants used to map to Options settings.
const (
	EnvKeyNoOpen  = "BUMP_NO_OPEN"
	EnvKeyVerbose = "BUMP_VERBOSE"
)

// NewOptionsFromEnv will return a populated Options struct with any settings
// defined via environment variables applied.
func NewOptionsFromEnv() *Options {
	return &Options{
		NoOpen:  getBoolEnv(EnvKeyNoOpen),
		Verbose: getBoolEnv(EnvKeyVerbose),
	}
}

// dont bother with ok confirmation like in LookupEnv since we just want
// to get a zero value then anyhow
func getBoolEnv(key string) bool {
	val := os.Getenv(key)
	switch strings.ToLower(val) {
	case "true", "yes", "1":
		return true
	default:
		return false
	}
}

// ParseFlags takes Options to use as a starting template -- likely populated
// from NewOptionsFromEnv() -- and parses flags contained in args into a new
// Options and returns that along with the FlagSet which was used so one can
// call Args on it.
func ParseFlags(opts *Options, args []string) (Options, *flag.FlagSet) {
	var newOpts Options
	var flags flag.FlagSet

	flags.BoolVar(&newOpts.NoOpen, "no-open", opts.NoOpen,
		fmt.Sprintf("don't auto-open in browser [$%v]", EnvKeyNoOpen))

	flags.BoolVar(&newOpts.Verbose, "verbose", opts.Verbose,
		fmt.Sprintf("verbose output [$%v]", EnvKeyVerbose))

	flags.BoolVar(&newOpts.Verbose, "v", opts.Verbose,
		fmt.Sprintf("verbose output [$%v]", EnvKeyVerbose))

	flags.Usage = usage
	flags.Parse(args)
	return newOpts, &flags
}
