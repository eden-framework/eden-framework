package enums

//go:generate eden generate enum --type-name=DroneCiKind
// api:enum
type DroneCiKind uint8

// DroneCI的kind类型
const (
	DRONE_CI_KIND_UNKNOWN    DroneCiKind = iota
	DRONE_CI_KIND__pipeline              // pipeline
	DRONE_CI_KIND__signature             // signature
	DRONE_CI_KIND__secret                // secret
)
