package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/go-github/v29/github"
)

func Test_extractRecentReleases(t *testing.T) {
	tests := []struct {
		sampleKey            string
		wantLatest, wantFull string
	}{
		// My forked version of scalafmt to provide a good sample of messy data.
		// Has a more recent v2.3.3-RC2 tag which is intentionally not marked as
		// a release on GitHub (emulating a mistake in the actual repo). The
		// current decision is that bump should intentionally only look at GH
		// releases, but not use their prerelease flag as source of truth, do
		// both ordering and prerelease status based on semver rules and
		// precedence.
		{
			sampleKey:  "mroth-scalafmt-20200121",
			wantLatest: "v2.3.3-RC1",
			wantFull:   "v2.3.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sampleKey, func(t *testing.T) {
			// load sample data
			dat, err := ioutil.ReadFile("testdata/releases/" + tt.sampleKey + ".sample.json")
			if err != nil {
				t.Fatal(err)
			}
			var rs []*github.RepositoryRelease
			err = json.Unmarshal(dat, &rs)
			if err != nil {
				t.Fatal(err)
			}

			got := extractRecentReleases(rs)
			gotFull, gotPre := got.Full.GetTagName(), got.Latest.GetTagName()
			if (tt.wantFull != gotFull) || (tt.wantLatest != gotPre) {
				t.Errorf("want [%v %v] got [%v %v]", tt.wantFull, tt.wantLatest, gotFull, gotPre)
			}
		})
	}
}
