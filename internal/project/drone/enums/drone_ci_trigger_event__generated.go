package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_eden_framework_enumeration "github.com/eden-framework/enumeration"
)

var InvalidDroneCiTriggerEvent = errors.New("invalid DroneCiTriggerEvent")

func init() {
	github_com_eden_framework_enumeration.RegisterEnums("DroneCiTriggerEvent", map[string]string{
		"rollback":     "rollback",
		"promote":      "promote",
		"tag":          "tag",
		"pull_request": "pull request",
		"push":         "push",
		"custom":       "custom",
		"cron":         "cron",
	})
}

func ParseDroneCiTriggerEventFromString(s string) (DroneCiTriggerEvent, error) {
	switch s {
	case "":
		return DRONE_CI_TRIGGER_EVENT_UNKNOWN, nil
	case "rollback":
		return DRONE_CI_TRIGGER_EVENT__rollback, nil
	case "promote":
		return DRONE_CI_TRIGGER_EVENT__promote, nil
	case "tag":
		return DRONE_CI_TRIGGER_EVENT__tag, nil
	case "pull_request":
		return DRONE_CI_TRIGGER_EVENT__pull_request, nil
	case "push":
		return DRONE_CI_TRIGGER_EVENT__push, nil
	case "custom":
		return DRONE_CI_TRIGGER_EVENT__custom, nil
	case "cron":
		return DRONE_CI_TRIGGER_EVENT__cron, nil
	}
	return DRONE_CI_TRIGGER_EVENT_UNKNOWN, InvalidDroneCiTriggerEvent
}

func ParseDroneCiTriggerEventFromLabelString(s string) (DroneCiTriggerEvent, error) {
	switch s {
	case "":
		return DRONE_CI_TRIGGER_EVENT_UNKNOWN, nil
	case "rollback":
		return DRONE_CI_TRIGGER_EVENT__rollback, nil
	case "promote":
		return DRONE_CI_TRIGGER_EVENT__promote, nil
	case "tag":
		return DRONE_CI_TRIGGER_EVENT__tag, nil
	case "pull request":
		return DRONE_CI_TRIGGER_EVENT__pull_request, nil
	case "push":
		return DRONE_CI_TRIGGER_EVENT__push, nil
	case "custom":
		return DRONE_CI_TRIGGER_EVENT__custom, nil
	case "cron":
		return DRONE_CI_TRIGGER_EVENT__cron, nil
	}
	return DRONE_CI_TRIGGER_EVENT_UNKNOWN, InvalidDroneCiTriggerEvent
}

func (DroneCiTriggerEvent) EnumType() string {
	return "DroneCiTriggerEvent"
}

func (DroneCiTriggerEvent) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_TRIGGER_EVENT__rollback):     {"rollback", "rollback"},
		int(DRONE_CI_TRIGGER_EVENT__promote):      {"promote", "promote"},
		int(DRONE_CI_TRIGGER_EVENT__tag):          {"tag", "tag"},
		int(DRONE_CI_TRIGGER_EVENT__pull_request): {"pull_request", "pull request"},
		int(DRONE_CI_TRIGGER_EVENT__push):         {"push", "push"},
		int(DRONE_CI_TRIGGER_EVENT__custom):       {"custom", "custom"},
		int(DRONE_CI_TRIGGER_EVENT__cron):         {"cron", "cron"},
	}
}

func (v DroneCiTriggerEvent) String() string {
	switch v {
	case DRONE_CI_TRIGGER_EVENT_UNKNOWN:
		return ""
	case DRONE_CI_TRIGGER_EVENT__rollback:
		return "rollback"
	case DRONE_CI_TRIGGER_EVENT__promote:
		return "promote"
	case DRONE_CI_TRIGGER_EVENT__tag:
		return "tag"
	case DRONE_CI_TRIGGER_EVENT__pull_request:
		return "pull_request"
	case DRONE_CI_TRIGGER_EVENT__push:
		return "push"
	case DRONE_CI_TRIGGER_EVENT__custom:
		return "custom"
	case DRONE_CI_TRIGGER_EVENT__cron:
		return "cron"
	}
	return "UNKNOWN"
}

func (v DroneCiTriggerEvent) Label() string {
	switch v {
	case DRONE_CI_TRIGGER_EVENT_UNKNOWN:
		return ""
	case DRONE_CI_TRIGGER_EVENT__rollback:
		return "rollback"
	case DRONE_CI_TRIGGER_EVENT__promote:
		return "promote"
	case DRONE_CI_TRIGGER_EVENT__tag:
		return "tag"
	case DRONE_CI_TRIGGER_EVENT__pull_request:
		return "pull request"
	case DRONE_CI_TRIGGER_EVENT__push:
		return "push"
	case DRONE_CI_TRIGGER_EVENT__custom:
		return "custom"
	case DRONE_CI_TRIGGER_EVENT__cron:
		return "cron"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DroneCiTriggerEvent)(nil)

func (v DroneCiTriggerEvent) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDroneCiTriggerEvent
	}
	return []byte(str), nil
}

func (v *DroneCiTriggerEvent) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDroneCiTriggerEventFromString(string(bytes.ToUpper(data)))
	return
}
