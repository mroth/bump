package main

import (
	"os"
	"strings"
	"testing"
)

func TestOptionsPrecedence(t *testing.T) {
	testCases := []struct {
		desc     string
		env      []string
		args     []string
		expected Options
	}{
		{
			desc: "just --verbose arg",
			args: []string{"--verbose"},
			expected: Options{
				NoOpen:  false,
				Verbose: true,
			},
		},
		{
			desc: "just -v arg",
			args: []string{"-v"},
			expected: Options{
				NoOpen:  false,
				Verbose: true,
			},
		},
		{
			desc: "just --no-open arg",
			args: []string{"--no-open"},
			expected: Options{
				NoOpen:  true,
				Verbose: false,
			},
		},
		{
			desc: "multiple args",
			args: []string{"-v", "--no-open"},
			expected: Options{
				NoOpen:  true,
				Verbose: true,
			},
		},
		{
			desc: "env arg bool true format",
			env:  []string{EnvKeyVerbose + "=true"},
			expected: Options{
				Verbose: true,
			},
		},
		{
			desc: "env arg bool TRUE format",
			env:  []string{EnvKeyVerbose + "=TRUE"},
			expected: Options{
				Verbose: true,
			},
		},
		{
			desc: "env arg bool 1 format",
			env:  []string{EnvKeyVerbose + "=1"},
			expected: Options{
				Verbose: true,
			},
		},
		{
			desc: "env arg bool yes format",
			env:  []string{EnvKeyVerbose + "=yes"},
			expected: Options{
				Verbose: true,
			},
		},
		{
			desc: "flags beat env if disagree",
			env:  []string{EnvKeyVerbose + "=yes"},
			args: []string{"--verbose=false"},
			expected: Options{
				Verbose: false,
			},
		},
	}

	originalEnv := os.Environ()
	for _, tC := range testCases {
		tC := tC // pin to avoid scope issues (see scopelint)
		t.Run(tC.desc, func(t *testing.T) {
			resetEnviron(tC.env)
			actualOpts, _ := ParseFlags(NewOptionsFromEnv(), tC.args)
			if actualOpts != tC.expected {
				t.Error("opts not as expected")
			}
		})
	}
	resetEnviron(originalEnv)
}

// resetEnviron cleasrs and then sets the environment to match a []string of
// key=value pairs, which happens to be exactly what os.Environ() from the
// standard library provides us, but with no built in way set back using the
// same thing... even though you can do the same in Cmd.
//
// Used to restore os.Environ after messing with it in a test, or to set it to
// exact values all at once.
func resetEnviron(envPairs []string) {
	os.Clearenv()
	for _, envPair := range envPairs {
		kv := strings.Split(envPair, "=")
		k, v := kv[0], kv[1]
		err := os.Setenv(k, v)
		if err != nil {
			panic("could not set env var in resetEnviron helper")
		}
	}
}
