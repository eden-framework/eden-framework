package project

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var VersionRegexp = regexp.MustCompile("(\\d+)\\.(\\d+)\\.(\\d+)")

func FromVersionString(s string) (*Version, error) {
	versions := VersionRegexp.FindAllStringSubmatch(s, -1)

	if len(versions) == 0 || len(versions[0]) != 4 {
		return nil, errors.New(s + " is not an available version")
	}

	vs := versions[0]

	major, _ := strconv.ParseInt(vs[1], 10, 10)
	minor, _ := strconv.ParseInt(vs[2], 10, 10)
	patch, _ := strconv.ParseInt(vs[3], 10, 10)

	return &Version{
		Major: int(major),
		Minor: int(minor),
		Patch: int(patch),
	}, nil
}

type Version struct {
	Latest bool
	Suffix string
	Prefix string
	Major  int
	Minor  int
	Patch  int
}

func (v Version) IsStable() bool {
	return !(v.Latest || v.Suffix != "")
}

func (v Version) String() string {
	version := "latest"
	if !v.Latest {
		versions := []string{
			strconv.Itoa(v.Major),
			strconv.Itoa(v.Minor),
			strconv.Itoa(v.Patch),
		}
		version = strings.Join(versions, ".")
	}

	if v.Prefix != "" {
		version = v.Prefix + "-" + version
	}

	if v.Suffix != "" {
		version = version + "-" + v.Suffix
	}

	return strings.ToLower(version)
}

func (v Version) IncreaseMajor() Version {
	v.Major = v.Major + 1
	v.Minor = 0
	v.Patch = 0
	return v
}

func (v Version) IncreaseMinor() Version {
	v.Minor = v.Minor + 1
	v.Patch = 0
	return v
}

func (v Version) IncreasePatch() Version {
	v.Patch = v.Patch + 1
	return v
}

func (v Version) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

func (v *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	err := unmarshal(&value)
	if err != nil {
		return err
	}
	if version, err := FromVersionString(value); err == nil {
		v.Major = version.Major
		v.Minor = version.Minor
		v.Patch = version.Patch
		return nil
	} else {
		return err
	}
	return nil
}
