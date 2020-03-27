package validatetpl

import (
	"gopkg.in/robfig/cron.v2"
)

const (
	InvalidCrontabType  = "频率串类型错误"
	InvalidCrontabValue = "无效的频率串"
)
const (
	MinimumCrontabLen = 11
	MaximumCrontabLen = 30
)

func ValidateCrontab(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidCrontabType
	}

	if len(s) < MinimumCrontabLen || len(s) > MaximumCrontabLen {
		return false, InvalidCrontabValue
	}

	if _, err := cron.Parse(s); err != nil {
		return false, InvalidCrontabValue
	}
	return true, ""
}

func ValidateCrontabOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidCrontabType
	}

	if len(s) == 0 {
		return true, ""
	}

	if _, err := cron.Parse(s); err != nil {
		return false, InvalidCrontabValue
	}
	return true, ""
}
