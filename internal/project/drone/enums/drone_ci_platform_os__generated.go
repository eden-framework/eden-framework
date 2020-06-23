package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/profzone/eden-framework/pkg/enumeration"
)

var InvalidDroneCiPlatformOs = errors.New("invalid DroneCiPlatformOs")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("DroneCiPlatformOs", map[string]string{
		"linux":   "linux",
		"windows": "windows",
	})
}

func ParseDroneCiPlatformOsFromString(s string) (DroneCiPlatformOs, error) {
	switch s {
	case "":
		return DRONE_CI_PLATFORM_OS_UNKNOWN, nil
	case "linux":
		return DRONE_CI_PLATFORM_OS__linux, nil
	case "windows":
		return DRONE_CI_PLATFORM_OS__windows, nil
	}
	return DRONE_CI_PLATFORM_OS_UNKNOWN, InvalidDroneCiPlatformOs
}

func ParseDroneCiPlatformOsFromLabelString(s string) (DroneCiPlatformOs, error) {
	switch s {
	case "":
		return DRONE_CI_PLATFORM_OS_UNKNOWN, nil
	case "linux":
		return DRONE_CI_PLATFORM_OS__linux, nil
	case "windows":
		return DRONE_CI_PLATFORM_OS__windows, nil
	}
	return DRONE_CI_PLATFORM_OS_UNKNOWN, InvalidDroneCiPlatformOs
}

func (DroneCiPlatformOs) EnumType() string {
	return "DroneCiPlatformOs"
}

func (DroneCiPlatformOs) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_PLATFORM_OS__linux):   {"linux", "linux"},
		int(DRONE_CI_PLATFORM_OS__windows): {"windows", "windows"},
	}
}

func (v DroneCiPlatformOs) String() string {
	switch v {
	case DRONE_CI_PLATFORM_OS_UNKNOWN:
		return ""
	case DRONE_CI_PLATFORM_OS__linux:
		return "linux"
	case DRONE_CI_PLATFORM_OS__windows:
		return "windows"
	}
	return "UNKNOWN"
}

func (v DroneCiPlatformOs) Label() string {
	switch v {
	case DRONE_CI_PLATFORM_OS_UNKNOWN:
		return ""
	case DRONE_CI_PLATFORM_OS__linux:
		return "linux"
	case DRONE_CI_PLATFORM_OS__windows:
		return "windows"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DroneCiPlatformOs)(nil)

func (v DroneCiPlatformOs) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDroneCiPlatformOs
	}
	return []byte(str), nil
}

func (v *DroneCiPlatformOs) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDroneCiPlatformOsFromString(string(bytes.ToUpper(data)))
	return
}
