package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_eden_framework_enumeration "gitee.com/eden-framework/enumeration"
)

var InvalidDroneCiType = errors.New("invalid DroneCiType")

func init() {
	github_com_eden_framework_enumeration.RegisterEnums("DroneCiType", map[string]string{
		"ssh":        "which executes shell commands on remote servers using the ssh protocol",
		"exec":       "which executes pipeline steps directly on the host machine, with zero isolation",
		"kubernetes": "which executes pipeline steps as containers inside of Kubernetes pods",
		"docker":     "which executes each pipeline steps inside isolated Docker containers",
	})
}

func ParseDroneCiTypeFromString(s string) (DroneCiType, error) {
	switch s {
	case "":
		return DRONE_CI_TYPE_UNKNOWN, nil
	case "ssh":
		return DRONE_CI_TYPE__ssh, nil
	case "exec":
		return DRONE_CI_TYPE__exec, nil
	case "kubernetes":
		return DRONE_CI_TYPE__kubernetes, nil
	case "docker":
		return DRONE_CI_TYPE__docker, nil
	}
	return DRONE_CI_TYPE_UNKNOWN, InvalidDroneCiType
}

func ParseDroneCiTypeFromLabelString(s string) (DroneCiType, error) {
	switch s {
	case "":
		return DRONE_CI_TYPE_UNKNOWN, nil
	case "which executes shell commands on remote servers using the ssh protocol":
		return DRONE_CI_TYPE__ssh, nil
	case "which executes pipeline steps directly on the host machine, with zero isolation":
		return DRONE_CI_TYPE__exec, nil
	case "which executes pipeline steps as containers inside of Kubernetes pods":
		return DRONE_CI_TYPE__kubernetes, nil
	case "which executes each pipeline steps inside isolated Docker containers":
		return DRONE_CI_TYPE__docker, nil
	}
	return DRONE_CI_TYPE_UNKNOWN, InvalidDroneCiType
}

func (DroneCiType) EnumType() string {
	return "DroneCiType"
}

func (DroneCiType) Enums() map[int][]string {
	return map[int][]string{
		int(DRONE_CI_TYPE__ssh):        {"ssh", "which executes shell commands on remote servers using the ssh protocol"},
		int(DRONE_CI_TYPE__exec):       {"exec", "which executes pipeline steps directly on the host machine, with zero isolation"},
		int(DRONE_CI_TYPE__kubernetes): {"kubernetes", "which executes pipeline steps as containers inside of Kubernetes pods"},
		int(DRONE_CI_TYPE__docker):     {"docker", "which executes each pipeline steps inside isolated Docker containers"},
	}
}

func (v DroneCiType) String() string {
	switch v {
	case DRONE_CI_TYPE_UNKNOWN:
		return ""
	case DRONE_CI_TYPE__ssh:
		return "ssh"
	case DRONE_CI_TYPE__exec:
		return "exec"
	case DRONE_CI_TYPE__kubernetes:
		return "kubernetes"
	case DRONE_CI_TYPE__docker:
		return "docker"
	}
	return "UNKNOWN"
}

func (v DroneCiType) Label() string {
	switch v {
	case DRONE_CI_TYPE_UNKNOWN:
		return ""
	case DRONE_CI_TYPE__ssh:
		return "which executes shell commands on remote servers using the ssh protocol"
	case DRONE_CI_TYPE__exec:
		return "which executes pipeline steps directly on the host machine, with zero isolation"
	case DRONE_CI_TYPE__kubernetes:
		return "which executes pipeline steps as containers inside of Kubernetes pods"
	case DRONE_CI_TYPE__docker:
		return "which executes each pipeline steps inside isolated Docker containers"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DroneCiType)(nil)

func (v DroneCiType) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDroneCiType
	}
	return []byte(str), nil
}

func (v *DroneCiType) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDroneCiTypeFromString(string(bytes.ToUpper(data)))
	return
}
