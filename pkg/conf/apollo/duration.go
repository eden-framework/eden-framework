package apollo

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Duration struct {
	startTime time.Time
}

func (t *Duration) Reset() {
	t.startTime = time.Now()
}

func (t *Duration) Get() string {
	now := time.Now()
	duration := now.Sub(t.startTime)
	return fmt.Sprintf("%0.3f", float64(duration/time.Millisecond))
}

func (t *Duration) GetAndReset() string {
	defer t.Reset()
	return t.Get()
}

func (t *Duration) ToLogger() *logrus.Entry {
	return logrus.WithField("cost", t.Get())
}

func NewDuration() *Duration {
	return &Duration{
		startTime: time.Now(),
	}
}

// PrintDuration print the time duration function process.
// printParam contains the fields which will be appeared in the log.
func PrintDuration(printParam map[string]interface{}) func() {
	start := NewDuration()
	return func() {
		printParam["request_time"] = start.Get()
		logrus.WithFields(logrus.Fields(printParam)).Debug()
	}
}
