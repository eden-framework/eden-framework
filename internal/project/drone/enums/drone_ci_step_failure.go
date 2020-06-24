package enums

//go:generate eden generate enum --type-name=DroneCiStepFailure
// api:enum
type DroneCiStepFailure uint8

//
const (
	DRONE_CI_STEP_FAILURE_UNKNOWN DroneCiStepFailure = iota
	DRONE_CI_STEP_FAILURE__ignore                    // ignore
)
