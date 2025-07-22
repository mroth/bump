package main

import (
	"cmp"

	"github.com/earthboundkid/versioninfo/v2"
)

// build info set by goreleaser during production builds
var (
	buildVersion = ""
	buildCommit  = ""
	buildDate    = ""
	builtBy      = ""
	treeState    = ""
)

func Version() string {
	return cmp.Or(buildVersion, versioninfo.Version)
}
