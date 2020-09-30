package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_eden_framework_enumeration "github.com/eden-framework/enumeration"
)

var InvalidDroneCiPlatformArch = errors.New("invalid DroneCiPlatformArch")

func init() {
	github_com_eden_framework_enumeration.RegisterEnums("DroneCiPlatformArch", map[string]string{
		"arm":   "arm32",
		"arm64": "arm64",
		"amd64": "amd64",
	})
}

func ParseDroneCiPlatformArchFromString(s string) (DroneCiPlatformArch, error) {
	switch s {
	case "":
		return DRONE_CI_PLATFORM_ARCH_UNKNOWN, nil
	case "arm":
		return DRONE_CI_PLATFORM_ARCH__arm, nil
	case "arm64":
		return DRONE_CI_PLATFORM_ARCH__arm64, nil
	case "amd64":
		return DRONE_CI_PLATFORM_ARCH__amd64, nil
	}
	return DRONE_CI_PLATFORM_ARCH_UNKNOWN, InvalidDroneCiPlatformArch
}

func ParseDroneCiPlatformArchFromLabelString(s string) (DroneCiPlatformArch, error) {
	switch s {
	case "":
		return DRONE_CI_PLATFORM_ARCH_UNKNOWN, nil
	case "arm32":
		return DRONE_CI_PLATFORM_ARCH__arm, nil
	case "arm64":
		return DRONE_CI_PLATFORM_ARCH__arm64, nil
	case "amd64":
		return DRONE_CI_PLATFORM_ARCH__amd64, nil
	}
	return DRONE_CI_PLATFORM_ARCH_UNKNOWN, InvalidDroneCiPlatformArch
}

func (DroneCiPlatformArch) EnumType() string {
	return "DroneCiPlatformArch"
}

func (DroneCiPlatformArch) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_PLATFORM_ARCH__arm):   {"arm", "arm32"},
		int(DRONE_CI_PLATFORM_ARCH__arm64): {"arm64", "arm64"},
		int(DRONE_CI_PLATFORM_ARCH__amd64): {"amd64", "amd64"},
	}
}

func (v DroneCiPlatformArch) String() string {
	switch v {
	case DRONE_CI_PLATFORM_ARCH_UNKNOWN:
		return ""
	case DRONE_CI_PLATFORM_ARCH__arm:
		return "arm"
	case DRONE_CI_PLATFORM_ARCH__arm64:
		return "arm64"
	case DRONE_CI_PLATFORM_ARCH__amd64:
		return "amd64"
	}
	return "UNKNOWN"
}

func (v DroneCiPlatformArch) Label() string {
	switch v {
	case DRONE_CI_PLATFORM_ARCH_UNKNOWN:
		return ""
	case DRONE_CI_PLATFORM_ARCH__arm:
		return "arm32"
	case DRONE_CI_PLATFORM_ARCH__arm64:
		return "arm64"
	case DRONE_CI_PLATFORM_ARCH__amd64:
		return "amd64"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DroneCiPlatformArch)(nil)

func (v DroneCiPlatformArch) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDroneCiPlatformArch
	}
	return []byte(str), nil
}

func (v *DroneCiPlatformArch) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDroneCiPlatformArchFromString(string(bytes.ToUpper(data)))
	return
}
