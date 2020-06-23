package enums

//go:generate eden generate enum --type-name=DroneCiType
// api:enum
type DroneCiType uint8

// DroneCI的type类型
const (
	DRONE_CI_TYPE_UNKNOWN     DroneCiType = iota
	DRONE_CI_TYPE__docker                 // which executes each pipeline steps inside isolated Docker containers
	DRONE_CI_TYPE__kubernetes             // which executes pipeline steps as containers inside of Kubernetes pods
	DRONE_CI_TYPE__exec                   // which executes pipeline steps directly on the host machine, with zero isolation
	DRONE_CI_TYPE__ssh                    // which executes shell commands on remote servers using the ssh protocol
)
