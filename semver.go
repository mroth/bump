// Additional functions for working with github.com.Masterminds/semver/v3 in
// order to provide some rudimentary support for bumping pre-release versions
// that follow a specific pattern, and semi-intelligent suggestion of what
// possible versions would come next, taking into account prereleases.
//
// Originally, I planned to wrap "Masterminds/semver/v3".Version entirely, but
// that results in enough boilerplate code that this is simpler, even if
// slightly less ergonomic as an API iuser.  In the future, consider whether any
// of this can be generalized enough to be upstreamed into a PR for semver
// itself.
package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Masterminds/semver/v3"
)

var (
	// InitialPrereleasePrefix is the prefix that will utilized for initial
	// prerelease versions when there is no previous prerelease to base it off
	// of.
	InitialPrereleasePrefix = "rc."
)

type VersionType int8

// Possible VersionType values
const (
	Patch VersionType = iota
	Minor
	Major
	PrePatch
	PreMinor
	PreMajor
)

func (vt VersionType) String() string {
	switch vt {
	case Patch:
		return "Patch"
	case Minor:
		return "Minor"
	case Major:
		return "Major"
	case PrePatch:
		return "PrePatch"
	case PreMinor:
		return "PreMinor"
	case PreMajor:
		return "PreMajor"
	default:
		return strconv.Itoa(int(vt))
	}
}

// standard errors provided by this file
var (
	ErrNotPrerelease     = errors.New("Version does not contain prerelease component")
	ErrInvalidPrerelease = errors.New("Prelease string did not conform to supported format")
)

// Type returns the VersionType of v.
func Type(v semver.Version) VersionType {
	switch {
	case HasPrerelease(v) && v.Minor() == 0 && v.Patch() == 0:
		return PreMajor
	case HasPrerelease(v) && v.Minor() != 0 && v.Patch() == 0:
		return PreMinor
	case HasPrerelease(v) && v.Patch() != 0:
		return PrePatch
	case v.Minor() == 0 && v.Patch() == 0:
		return Major
	case v.Minor() != 0 && v.Patch() == 0:
		return Minor
	case v.Patch() != 0:
		return Patch
	default:
		panic("unreachable! (fatal logic error in identifying version type)")
	}
}

// SuggestNext will return reasonable suggestions for a Version that would
// come after the current Version.
//
// For a typical release version, ***
//
// For an existing pre-release version, this should be the incremented prelease
// version, and the following version.
func SuggestNext(v semver.Version) semver.Collection {
	if HasValidPrerelease(v) {
		nextPre, _ := IncPrerelease(v) // safety: can't err due to guard
		final := Finalize(v)
		return semver.Collection{&nextPre, &final}
	}

	patch := v.IncPatch()
	minor := v.IncMinor()
	major := v.IncMajor()
	return semver.Collection{
		&patch,
		&minor,
		&major,
		// TODO: PrePatch
		// TODO: PreMinor
		// TODO: PreMajor
	}
}

// HasPrerelease returns whether or not a version has a prerelease version
// component.
//
// To determine whether the prerelease version component conforms to the
// standard format used by bump, use ValidPrelease() instead.
func HasPrerelease(v semver.Version) bool {
	return v.Prerelease() != ""
}

// HasValidPrerelease whether or not the version hash a prerelease version
// component which conforms to the standardized format used by bump.
func HasValidPrerelease(v semver.Version) bool {
	_, _, err := parsePreStr(v.Prerelease())
	return err != nil
}

// IncPrerelease returns an incremented prerelease Version from a Version that
// contains an existing valid prerelease.
func IncPrerelease(v semver.Version) (semver.Version, error) {
	if !HasPrerelease(v) {
		return v, ErrNotPrerelease
	}
	prefix, counter, err := parsePreStr(v.Prerelease())
	if err != nil {
		return v, err
	}
	counter++
	newPreS := prefix + strconv.Itoa(int(counter))
	vNext, err := v.SetPrerelease(newPreS)
	return vNext, err
}

// Finalize returns the semver.Version that v is a prerelease of.
func Finalize(v semver.Version) semver.Version {
	nv, _ := v.SetPrerelease("")
	// Safety: look at implementation of SetPrerelease, should be okay.
	return nv
}

var rePre = regexp.MustCompile(`(.*?)(\d+)$`)

func parsePreStr(pre string) (prefix string, counter uint, err error) {
	m := rePre.FindStringSubmatch(pre)
	if len(m) != 3 {
		return "", 0, ErrInvalidPrerelease
	}
	c, err := strconv.Atoi(m[2])
	if err != nil {
		// should be impossible based on regex, so panic to cause a lot of noise
		panic(err)
	}
	return m[1], uint(c), nil
}

// incPreStr will attempt to increment the numeric suffix in a prerelease
// string, preserving any leading zeroes.
func incPreStr(pre string) (string, error) {
	prefix, counter, err := parsePreStr(pre)
	if err != nil {
		return "", err
	}
	counter++
	origSuffixLen := len(pre) - len(prefix)
	return fmt.Sprintf("%s%0*d", prefix, origSuffixLen, counter), nil
}
