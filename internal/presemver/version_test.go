package presemver

import (
	"reflect"
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestType(t *testing.T) {
	tests := []struct {
		verStr string
		want   VersionType
	}{
		{"1.2.3", Patch},
		{"1.2.0", Minor},
		{"1.0.0", Major},
		{"0.0.1", Patch},
		{"1.2.3-rc.1", PrePatch},
		{"1.2.0-foo", PreMinor},
		{"2.0.0-dev1", PreMajor},
	}
	for _, tt := range tests {
		v := semver.MustParse(tt.verStr)
		if got := Type(*v); got != tt.want {
			t.Errorf("Version.Type(%v) = %v, want %v", tt.verStr, got, tt.want)
		}
	}
}

func TestSuggestNext(t *testing.T) {
	type args struct {
		v                  semver.Version
		initialPrereleases bool
	}
	tests := []struct {
		name    string
		args    args
		want    semver.Collection
		wantErr bool
	}{
		{
			name: "typical",
			args: args{
				v:                  *semver.MustParse("1.2.3"),
				initialPrereleases: false,
			},
			want: semver.Collection{
				semver.MustParse("1.2.4"),
				semver.MustParse("1.3.0"),
				semver.MustParse("2.0.0"),
			},
			wantErr: false,
		},
		{
			name: "include initial prereleases",
			args: args{
				v:                  *semver.MustParse("1.2.3"),
				initialPrereleases: true,
			},
			want: semver.Collection{
				semver.MustParse("1.2.4"),
				semver.MustParse("1.3.0"),
				semver.MustParse("2.0.0"),
				semver.MustParse("1.2.4-rc.1"),
				semver.MustParse("1.3.0-rc.1"),
				semver.MustParse("2.0.0-rc.1"),
			},
			wantErr: false,
		},
		{
			name: "existing prerelease",
			args: args{
				v: *semver.MustParse("1.3.0-rc.2"),
			},
			want: semver.Collection{
				semver.MustParse("1.3.0-rc.3"),
				semver.MustParse("1.3.0"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SuggestNext(tt.args.v, tt.args.initialPrereleases)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuggestNext(%v) error = %v, wantErr %v", tt.args.v, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SuggestNext(%v) = %v, want %v", tt.args.v, got, tt.want)
			}
		})
	}
}

func Test_parsePreStr(t *testing.T) {
	tests := []struct {
		pre         string
		wantPrefix  string
		wantCounter uint
		wantErr     bool
	}{
		{"pre1", "pre", 1, false},
		{"pre.1", "pre.", 1, false},
		{"rc32", "rc", 32, false},
		{"", "", 0, true},
		{"foobar", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.pre, func(t *testing.T) {
			gotPrefix, gotCounter, err := parsePreStr(tt.pre)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePreStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPrefix != tt.wantPrefix {
				t.Errorf("parsePreStr() gotPrefix = %v, want %v", gotPrefix, tt.wantPrefix)
			}
			if gotCounter != tt.wantCounter {
				t.Errorf("parsePreStr() gotCounter = %v, want %v", gotCounter, tt.wantCounter)
			}
		})
	}
}

func Test_incPreStr(t *testing.T) {
	tests := []struct {
		pre     string
		want    string
		wantErr bool
	}{
		{"pre1", "pre2", false},
		{"pre.1", "pre.2", false},
		{"dev01", "dev02", false}, // preserve leading zeros
		{"dev.01", "dev.02", false},
		{"rc32", "rc33", false}, // no sep
		{"", "", true},
		{"foobar", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.pre, func(t *testing.T) {
			got, err := incPreStr(tt.pre)
			if (err != nil) != tt.wantErr {
				t.Errorf("incPreStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("incPreStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
