package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_eden_framework_enumeration "gitee.com/eden-framework/enumeration"
)

var InvalidDroneCiKind = errors.New("invalid DroneCiKind")

func init() {
	github_com_eden_framework_enumeration.RegisterEnums("DroneCiKind", map[string]string{
		"secret":    "secret",
		"signature": "signature",
		"pipeline":  "pipeline",
	})
}

func ParseDroneCiKindFromString(s string) (DroneCiKind, error) {
	switch s {
	case "":
		return DRONE_CI_KIND_UNKNOWN, nil
	case "secret":
		return DRONE_CI_KIND__secret, nil
	case "signature":
		return DRONE_CI_KIND__signature, nil
	case "pipeline":
		return DRONE_CI_KIND__pipeline, nil
	}
	return DRONE_CI_KIND_UNKNOWN, InvalidDroneCiKind
}

func ParseDroneCiKindFromLabelString(s string) (DroneCiKind, error) {
	switch s {
	case "":
		return DRONE_CI_KIND_UNKNOWN, nil
	case "secret":
		return DRONE_CI_KIND__secret, nil
	case "signature":
		return DRONE_CI_KIND__signature, nil
	case "pipeline":
		return DRONE_CI_KIND__pipeline, nil
	}
	return DRONE_CI_KIND_UNKNOWN, InvalidDroneCiKind
}

func (DroneCiKind) EnumType() string {
	return "DroneCiKind"
}

func (DroneCiKind) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_KIND__secret):    {"secret", "secret"},
		int(DRONE_CI_KIND__signature): {"signature", "signature"},
		int(DRONE_CI_KIND__pipeline):  {"pipeline", "pipeline"},
	}
}

func (v DroneCiKind) String() string {
	switch v {
	case DRONE_CI_KIND_UNKNOWN:
		return ""
	case DRONE_CI_KIND__secret:
		return "secret"
	case DRONE_CI_KIND__signature:
		return "signature"
	case DRONE_CI_KIND__pipeline:
		return "pipeline"
	}
	return "UNKNOWN"
}

func (v DroneCiKind) Label() string {
	switch v {
	case DRONE_CI_KIND_UNKNOWN:
		return ""
	case DRONE_CI_KIND__secret:
		return "secret"
	case DRONE_CI_KIND__signature:
		return "signature"
	case DRONE_CI_KIND__pipeline:
		return "pipeline"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DroneCiKind)(nil)

func (v DroneCiKind) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDroneCiKind
	}
	return []byte(str), nil
}

func (v *DroneCiKind) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDroneCiKindFromString(string(bytes.ToUpper(data)))
	return
}
