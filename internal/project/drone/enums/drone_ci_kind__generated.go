package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_profzone_eden_framework_pkg_enumeration "github.com/profzone/eden-framework/pkg/enumeration"
)

var InvalidDroneCiKind = errors.New("invalid DroneCiKind")

func init() {
	github_com_profzone_eden_framework_pkg_enumeration.RegisterEnums("DroneCiKind", map[string]string{
		"pipeline":  "pipeline",
		"secret":    "secret",
		"signature": "signature",
	})
}

func ParseDroneCiKindFromString(s string) (DroneCiKind, error) {
	switch s {
	case "":
		return DRONE_CI_KIND_UNKNOWN, nil
	case "pipeline":
		return DRONE_CI_KIND__pipeline, nil
	case "secret":
		return DRONE_CI_KIND__secret, nil
	case "signature":
		return DRONE_CI_KIND__signature, nil
	}
	return DRONE_CI_KIND_UNKNOWN, InvalidDroneCiKind
}

func ParseDroneCiKindFromLabelString(s string) (DroneCiKind, error) {
	switch s {
	case "":
		return DRONE_CI_KIND_UNKNOWN, nil
	case "pipeline":
		return DRONE_CI_KIND__pipeline, nil
	case "secret":
		return DRONE_CI_KIND__secret, nil
	case "signature":
		return DRONE_CI_KIND__signature, nil
	}
	return DRONE_CI_KIND_UNKNOWN, InvalidDroneCiKind
}

func (DroneCiKind) EnumType() string {
	return "DroneCiKind"
}

func (DroneCiKind) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_KIND__pipeline):  {"pipeline", "pipeline"},
		int(DRONE_CI_KIND__secret):    {"secret", "secret"},
		int(DRONE_CI_KIND__signature): {"signature", "signature"},
	}
}

func (v DroneCiKind) String() string {
	switch v {
	case DRONE_CI_KIND_UNKNOWN:
		return ""
	case DRONE_CI_KIND__pipeline:
		return "pipeline"
	case DRONE_CI_KIND__secret:
		return "secret"
	case DRONE_CI_KIND__signature:
		return "signature"
	}
	return "UNKNOWN"
}

func (v DroneCiKind) Label() string {
	switch v {
	case DRONE_CI_KIND_UNKNOWN:
		return ""
	case DRONE_CI_KIND__pipeline:
		return "pipeline"
	case DRONE_CI_KIND__secret:
		return "secret"
	case DRONE_CI_KIND__signature:
		return "signature"
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
