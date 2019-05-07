package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const usageText = `Usage: bump <owner> <repo>

If you are in a git repository that has been cloned from GitHub, owner and
repo args can be omitted, in which case they will be inferred from the remote
origin.

Flags:
    --no-open           Do not automatically open publish URL in browser.
    --verbose, -v       Verbose output.
    --version           Print version and exit.
    --help              Print help and exit.

Environment:
    $BUMP_NO_OPEN       Global default for --no-open
    $BUMP_VERBOSE       Global default for --verbose
    $GITHUB_TOKEN       Optional, will use if present to access private repos
`

func usage() {
	fmt.Fprintf(os.Stderr, usageText)
	os.Exit(1)
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
//
// If version was requested, we just output and shortcircuit exit, since same
// thing would happen if --help was requested as per normal FlagSet behavior.
//
// Note none of the usage text here actually shows up in help output since just
// manually overriding that page for now.
func ParseFlags(opts *Options, args []string) (Options, *flag.FlagSet) {
	var newOpts Options
	var flags flag.FlagSet

	flags.BoolVar(&newOpts.NoOpen, "no-open", opts.NoOpen, "")
	flags.BoolVar(&newOpts.Verbose, "verbose", opts.Verbose, "")
	flags.BoolVar(&newOpts.Verbose, "v", opts.Verbose, "")
	version := flags.Bool("version", false, "")
	flags.Usage = usage

	// explicitly swallow error to appease errcheck
	// (FlagSet.Parse returns "ErrHelp if -help or -h were set but not defined")
	_ = flags.Parse(args)
	if *version {
		fmt.Println(buildVersion)
		os.Exit(0)
	}

	return newOpts, &flags
}

// ParseAll rolls up all CLI option parsing curently needed for main()
func ParseAll() (owner, repo string, opts Options) {
	opts, flags := ParseFlags(NewOptionsFromEnv(), os.Args[1:])
	owner = flags.Arg(0)
	repo = flags.Arg(1)
	return
}

// For future possible expansion...
// Usage: bump <owner> <repo> [major|minor|patch|interactive*]
// type Strategy int

// const (
// 	Patch Strategy = iota
// 	Minor
// 	Major
// 	Interactive
// )
