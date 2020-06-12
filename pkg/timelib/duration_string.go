package timelib

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"time"
)

type DurationString time.Duration

var (
	DurationStringZero = DurationString(0)
)

func ParseDurationStringFromString(s string) (ds DurationString, err error) {
	var d time.Duration
	d, err = time.ParseDuration(s)
	ds = DurationString(d)
	return
}

var _ interface {
	sql.Scanner
	driver.Valuer
} = (*DurationString)(nil)

func (ds *DurationString) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		d, err := ParseDurationStringFromString(v)
		if err != nil {
			return err
		}
		*ds = d
	case nil:
		*ds = DurationStringZero
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.DurationString from: %#v", v)
	}
	return nil
}

func (ds DurationString) Value() (driver.Value, error) {
	return ds.String(), nil
}

func (ds DurationString) String() string {
	return time.Duration(ds).String()
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DurationString)(nil)

func (ds DurationString) MarshalText() ([]byte, error) {
	if ds.IsZero() {
		return []byte(""), nil
	}
	str := ds.String()
	return []byte(str), nil
}

func (ds *DurationString) UnmarshalText(data []byte) (err error) {
	str := string(data)
	if len(str) == 0 || str == "0" || str == "0s" {
		str = DurationStringZero.String()
	}
	*ds, err = ParseDurationStringFromString(str)
	return
}

func (ds DurationString) IsZero() bool {
	return time.Duration(ds).Nanoseconds() == 0
}
