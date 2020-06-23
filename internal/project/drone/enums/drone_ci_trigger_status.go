package enums

//go:generate eden generate enum --type-name=DroneCiTriggerStatus
// api:enum
type DroneCiTriggerStatus uint8

//
const (
	DRONE_CI_TRIGGER_STATUS_UNKNOWN  DroneCiTriggerStatus = iota
	DRONE_CI_TRIGGER_STATUS__success                      // success
	DRONE_CI_TRIGGER_STATUS__failure                      // failure
)
