package timelib

import (
	"time"
)

const (
	CST_ZONE_OFFSET = 8 * 60 * 60
)

func Now() time.Time {
	return time.Now().In(CST)
}

var (
	UTC *time.Location
	CST *time.Location
)

func init() {
	UTC = time.UTC
	CST = time.FixedZone("CST", CST_ZONE_OFFSET)
}
