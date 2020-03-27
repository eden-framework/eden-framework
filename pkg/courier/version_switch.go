package courier

import (
	"fmt"
	"regexp"
)

var VersionSwitchKey = "X-Version"
var reVersionSwitch = regexp.MustCompile("@x-version\\(([a-zA-Z_]+)\\)")

type WithVersionSwitch struct {
	XVersion string `name:"X-Version,omitempty" in:"header"`
}

func MetadataWithVersionSwitch(xVersion string) Metadata {
	return Metadata{
		VersionSwitchKey: []string{xVersion},
	}
}

func ModifyRequestIDWithVersionSwitch(requestID string, version string) string {
	others, _, exists := ParseVersionSwitch(requestID)
	if exists {
		requestID = others
	}
	return requestID + fmt.Sprintf("@x-version(%s)", version)
}

func ParseVersionSwitch(s string) (string, string, bool) {
	version := reVersionSwitch.FindString(s)
	if version == "" {
		return s, version, false
	}
	return reVersionSwitch.ReplaceAllStringFunc(s, func(s string) string {
		version = reVersionSwitch.FindStringSubmatch(s)[1]
		return ""
	}), version, true
}
