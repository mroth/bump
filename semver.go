package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// standard errors provided by this file
var (
	ErrNotPrerelease     = errors.New("Version does not contain prerelease component")
	ErrInvalidPrerelease = errors.New("Prelease string did not conform to supported format")
)

// Version wraps "Masterminds/semver/v3".Version in order to provide some
// rudimentary support for bumping pre-release versions that follow a specific
// pattern.
// type Version struct {
// 	semver.Version
// }

// IsPrerelease returns whether or not a version has a prerelease version
// component.
// func (v Version) IsPrerelease() bool {
// 	return v.Prerelease() != ""
// }

// func (v Version) IncPrerelease() (Version, error) {
// 	if !v.IsPrerelease() {
// 		return v, ErrNotPrerelease
// 	}

// 	prefix, counter, err := parsePreStr(v.Prerelease())
// 	if err != nil {
// 		return v, err
// 	}
// 	counter++
// 	newPreS := prefix + strconv.Itoa(int(counter))
// 	vNext, err := v.SetPrerelease(newPreS)
// 	return Version{vNext}, err
// }

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
