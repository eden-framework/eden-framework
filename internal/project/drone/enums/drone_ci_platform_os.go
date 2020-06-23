package enums

//go:generate eden generate enum --type-name=DroneCiPlatformOs
// api:enum
type DroneCiPlatformOs uint8

//
const (
	DRONE_CI_PLATFORM_OS_UNKNOWN  DroneCiPlatformOs = iota
	DRONE_CI_PLATFORM_OS__linux                     // linux
	DRONE_CI_PLATFORM_OS__windows                   // windows
)
