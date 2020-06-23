package enums

//go:generate eden generate enum --type-name=DroneCiTriggerEvent
// api:enum
type DroneCiTriggerEvent uint8

//
const (
	DRONE_CI_TRIGGER_EVENT_UNKNOWN       DroneCiTriggerEvent = iota
	DRONE_CI_TRIGGER_EVENT__cron                             // cron
	DRONE_CI_TRIGGER_EVENT__custom                           // custom
	DRONE_CI_TRIGGER_EVENT__push                             // push
	DRONE_CI_TRIGGER_EVENT__pull_request                     // pull request
	DRONE_CI_TRIGGER_EVENT__tag                              // tag
	DRONE_CI_TRIGGER_EVENT__promote                          // promote
	DRONE_CI_TRIGGER_EVENT__rollback                         // rollback
)
