package enums

//go:generate eden generate enum --type-name=DroneCiStepPull
// api:enum
type DroneCiStepPull uint8

//
const (
	DRONE_CI_STEP_PULL_UNKNOWN        DroneCiStepPull = iota
	DRONE_CI_STEP_PULL__if_not_exists                 // if not exists
	DRONE_CI_STEP_PULL__always                        // always
	DRONE_CI_STEP_PULL__never                         // never
)
