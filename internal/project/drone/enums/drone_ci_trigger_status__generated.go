package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/eden-framework/eden-framework/pkg/enumeration"
)

var InvalidDroneCiTriggerStatus = errors.New("invalid DroneCiTriggerStatus")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("DroneCiTriggerStatus", map[string]string{
		"failure": "failure",
		"success": "success",
	})
}

func ParseDroneCiTriggerStatusFromString(s string) (DroneCiTriggerStatus, error) {
	switch s {
	case "":
		return DRONE_CI_TRIGGER_STATUS_UNKNOWN, nil
	case "failure":
		return DRONE_CI_TRIGGER_STATUS__failure, nil
	case "success":
		return DRONE_CI_TRIGGER_STATUS__success, nil
	}
	return DRONE_CI_TRIGGER_STATUS_UNKNOWN, InvalidDroneCiTriggerStatus
}

func ParseDroneCiTriggerStatusFromLabelString(s string) (DroneCiTriggerStatus, error) {
	switch s {
	case "":
		return DRONE_CI_TRIGGER_STATUS_UNKNOWN, nil
	case "failure":
		return DRONE_CI_TRIGGER_STATUS__failure, nil
	case "success":
		return DRONE_CI_TRIGGER_STATUS__success, nil
	}
	return DRONE_CI_TRIGGER_STATUS_UNKNOWN, InvalidDroneCiTriggerStatus
}

func (DroneCiTriggerStatus) EnumType() string {
	return "DroneCiTriggerStatus"
}

func (DroneCiTriggerStatus) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_TRIGGER_STATUS__failure): {"failure", "failure"},
		int(DRONE_CI_TRIGGER_STATUS__success): {"success", "success"},
	}
}

func (v DroneCiTriggerStatus) String() string {
	switch v {
	case DRONE_CI_TRIGGER_STATUS_UNKNOWN:
		return ""
	case DRONE_CI_TRIGGER_STATUS__failure:
		return "failure"
	case DRONE_CI_TRIGGER_STATUS__success:
		return "success"
	}
	return "UNKNOWN"
}

func (v DroneCiTriggerStatus) Label() string {
	switch v {
	case DRONE_CI_TRIGGER_STATUS_UNKNOWN:
		return ""
	case DRONE_CI_TRIGGER_STATUS__failure:
		return "failure"
	case DRONE_CI_TRIGGER_STATUS__success:
		return "success"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DroneCiTriggerStatus)(nil)

func (v DroneCiTriggerStatus) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDroneCiTriggerStatus
	}
	return []byte(str), nil
}

func (v *DroneCiTriggerStatus) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDroneCiTriggerStatusFromString(string(bytes.ToUpper(data)))
	return
}
