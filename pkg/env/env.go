package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	UNKNOWN = iota
	ONLINE
	PRE
	DEMO
	TEST
	STAGING
	DEV
	LOCAL
)

type GoEnv int

func (goEnv GoEnv) String() string {
	switch goEnv {
	case ONLINE:
		return "ONLINE"
	case PRE:
		return "PRE"
	case DEMO:
		return "DEMO"
	case TEST:
		return "TEST"
	case STAGING:
		return "STAGING"
	case DEV:
		return "DEV"
	default:
		panic(fmt.Sprintf("invalid go env %d", goEnv))
	}
}

func GetRuntimeEnv() GoEnv {
	goEnv := os.Getenv("GOENV")
	if goEnv == "" {
		goEnv = "DEV"
	}
	switch strings.ToUpper(goEnv) {
	case "ONLINE":
		return ONLINE
	case "PRE":
		return PRE
	case "DEMO":
		return DEMO
	case "TEST":
		return TEST
	case "STAGING":
		return STAGING
	case "DEV":
		return DEV
	case "LOCAL":
		return LOCAL
	default:
		panic("invalid go env " + goEnv)
	}
}

func IsOnline() bool {
	return GetRuntimeEnv() == ONLINE
}
