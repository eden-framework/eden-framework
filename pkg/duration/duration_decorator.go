package duration

import (
	"github.com/sirupsen/logrus"
)

// PrintDuration print the time duration function process.
// printParam contains the fields which will be appeared in the log.
func PrintDuration(printParam map[string]interface{}) func() {
	start := NewDuration()
	return func() {
		printParam["request_time"] = start.Get()
		logrus.WithFields(logrus.Fields(printParam)).Info()
	}
}
