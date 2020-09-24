package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/eden-framework/eden-framework/pkg/enumeration"
)

var InvalidDroneCiStepFailure = errors.New("invalid DroneCiStepFailure")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("DroneCiStepFailure", map[string]string{
		"ignore": "ignore",
	})
}

func ParseDroneCiStepFailureFromString(s string) (DroneCiStepFailure, error) {
	switch s {
	case "":
		return DRONE_CI_STEP_FAILURE_UNKNOWN, nil
	case "ignore":
		return DRONE_CI_STEP_FAILURE__ignore, nil
	}
	return DRONE_CI_STEP_FAILURE_UNKNOWN, InvalidDroneCiStepFailure
}

func ParseDroneCiStepFailureFromLabelString(s string) (DroneCiStepFailure, error) {
	switch s {
	case "":
		return DRONE_CI_STEP_FAILURE_UNKNOWN, nil
	case "ignore":
		return DRONE_CI_STEP_FAILURE__ignore, nil
	}
	return DRONE_CI_STEP_FAILURE_UNKNOWN, InvalidDroneCiStepFailure
}

func (DroneCiStepFailure) EnumType() string {
	return "DroneCiStepFailure"
}

func (DroneCiStepFailure) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_STEP_FAILURE__ignore): {"ignore", "ignore"},
	}
}

func (v DroneCiStepFailure) String() string {
	switch v {
	case DRONE_CI_STEP_FAILURE_UNKNOWN:
		return ""
	case DRONE_CI_STEP_FAILURE__ignore:
		return "ignore"
	}
	return "UNKNOWN"
}

func (v DroneCiStepFailure) Label() string {
	switch v {
	case DRONE_CI_STEP_FAILURE_UNKNOWN:
		return ""
	case DRONE_CI_STEP_FAILURE__ignore:
		return "ignore"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DroneCiStepFailure)(nil)

func (v DroneCiStepFailure) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDroneCiStepFailure
	}
	return []byte(str), nil
}

func (v *DroneCiStepFailure) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDroneCiStepFailureFromString(string(bytes.ToUpper(data)))
	return
}
