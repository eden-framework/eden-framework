package enums

//go:generate eden generate enum --type-name=DroneCiPlatformArch
// api:enum
type DroneCiPlatformArch uint8

//
const (
	DRONE_CI_PLATFORM_ARCH_UNKNOWN DroneCiPlatformArch = iota
	DRONE_CI_PLATFORM_ARCH__amd64                      // amd64
	DRONE_CI_PLATFORM_ARCH__arm64                      // arm64
	DRONE_CI_PLATFORM_ARCH__arm                        // arm32
)
