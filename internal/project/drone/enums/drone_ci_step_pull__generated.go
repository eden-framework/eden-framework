package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/profzone/eden-framework/pkg/enumeration"
)

var InvalidDroneCiStepPull = errors.New("invalid DroneCiStepPull")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("DroneCiStepPull", map[string]string{
		"always":        "always",
		"if_not_exists": "if not exists",
		"never":         "never",
	})
}

func ParseDroneCiStepPullFromString(s string) (DroneCiStepPull, error) {
	switch s {
	case "":
		return DRONE_CI_STEP_PULL_UNKNOWN, nil
	case "always":
		return DRONE_CI_STEP_PULL__always, nil
	case "if-not-exists":
		return DRONE_CI_STEP_PULL__if_not_exists, nil
	case "never":
		return DRONE_CI_STEP_PULL__never, nil
	}
	return DRONE_CI_STEP_PULL_UNKNOWN, InvalidDroneCiStepPull
}

func ParseDroneCiStepPullFromLabelString(s string) (DroneCiStepPull, error) {
	switch s {
	case "":
		return DRONE_CI_STEP_PULL_UNKNOWN, nil
	case "always":
		return DRONE_CI_STEP_PULL__always, nil
	case "if not exists":
		return DRONE_CI_STEP_PULL__if_not_exists, nil
	case "never":
		return DRONE_CI_STEP_PULL__never, nil
	}
	return DRONE_CI_STEP_PULL_UNKNOWN, InvalidDroneCiStepPull
}

func (DroneCiStepPull) EnumType() string {
	return "DroneCiStepPull"
}

func (DroneCiStepPull) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_STEP_PULL__always):        {"always", "always"},
		int(DRONE_CI_STEP_PULL__if_not_exists): {"if_not_exists", "if not exists"},
		int(DRONE_CI_STEP_PULL__never):         {"never", "never"},
	}
}

func (v DroneCiStepPull) String() string {
	switch v {
	case DRONE_CI_STEP_PULL_UNKNOWN:
		return ""
	case DRONE_CI_STEP_PULL__always:
		return "always"
	case DRONE_CI_STEP_PULL__if_not_exists:
		return "if-not-exists"
	case DRONE_CI_STEP_PULL__never:
		return "never"
	}
	return "UNKNOWN"
}

func (v DroneCiStepPull) Label() string {
	switch v {
	case DRONE_CI_STEP_PULL_UNKNOWN:
		return ""
	case DRONE_CI_STEP_PULL__always:
		return "always"
	case DRONE_CI_STEP_PULL__if_not_exists:
		return "if not exists"
	case DRONE_CI_STEP_PULL__never:
		return "never"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DroneCiStepPull)(nil)

func (v DroneCiStepPull) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDroneCiStepPull
	}
	return []byte(str), nil
}

func (v *DroneCiStepPull) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDroneCiStepPullFromString(string(bytes.ToUpper(data)))
	return
}
