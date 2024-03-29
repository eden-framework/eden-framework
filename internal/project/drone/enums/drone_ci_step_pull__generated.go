package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_eden_framework_enumeration "gitee.com/eden-framework/enumeration"
)

var InvalidDroneCiStepPull = errors.New("invalid DroneCiStepPull")

func init() {
	github_com_eden_framework_enumeration.RegisterEnums("DroneCiStepPull", map[string]string{
		"never":         "never",
		"always":        "always",
		"if_not_exists": "if not exists",
	})
}

func ParseDroneCiStepPullFromString(s string) (DroneCiStepPull, error) {
	switch s {
	case "":
		return DRONE_CI_STEP_PULL_UNKNOWN, nil
	case "never":
		return DRONE_CI_STEP_PULL__never, nil
	case "always":
		return DRONE_CI_STEP_PULL__always, nil
	case "if_not_exists":
		return DRONE_CI_STEP_PULL__if_not_exists, nil
	}
	return DRONE_CI_STEP_PULL_UNKNOWN, InvalidDroneCiStepPull
}

func ParseDroneCiStepPullFromLabelString(s string) (DroneCiStepPull, error) {
	switch s {
	case "":
		return DRONE_CI_STEP_PULL_UNKNOWN, nil
	case "never":
		return DRONE_CI_STEP_PULL__never, nil
	case "always":
		return DRONE_CI_STEP_PULL__always, nil
	case "if not exists":
		return DRONE_CI_STEP_PULL__if_not_exists, nil
	}
	return DRONE_CI_STEP_PULL_UNKNOWN, InvalidDroneCiStepPull
}

func (DroneCiStepPull) EnumType() string {
	return "DroneCiStepPull"
}

func (DroneCiStepPull) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_STEP_PULL__never):         {"never", "never"},
		int(DRONE_CI_STEP_PULL__always):        {"always", "always"},
		int(DRONE_CI_STEP_PULL__if_not_exists): {"if_not_exists", "if not exists"},
	}
}

func (v DroneCiStepPull) String() string {
	switch v {
	case DRONE_CI_STEP_PULL_UNKNOWN:
		return ""
	case DRONE_CI_STEP_PULL__never:
		return "never"
	case DRONE_CI_STEP_PULL__always:
		return "always"
	case DRONE_CI_STEP_PULL__if_not_exists:
		return "if_not_exists"
	}
	return "UNKNOWN"
}

func (v DroneCiStepPull) Label() string {
	switch v {
	case DRONE_CI_STEP_PULL_UNKNOWN:
		return ""
	case DRONE_CI_STEP_PULL__never:
		return "never"
	case DRONE_CI_STEP_PULL__always:
		return "always"
	case DRONE_CI_STEP_PULL__if_not_exists:
		return "if not exists"
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
