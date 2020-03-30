package presets

import (
	"time"

	"github.com/profzone/eden-framework/pkg/timelib"
)

type OperateTime struct {
	CreateTime timelib.MySQLTimestamp `db:"F_create_time" sql:"bigint(64) NOT NULL DEFAULT '0'" json:"createTime"`
	UpdateTime timelib.MySQLTimestamp `db:"F_update_time" sql:"bigint(64) NOT NULL DEFAULT '0'" json:"updateTime"`
}

func (t *OperateTime) BeforeUpdate() {
	t.UpdateTime = timelib.MySQLTimestamp(time.Now())
}

func (t *OperateTime) BeforeInsert() {
	t.CreateTime = timelib.MySQLTimestamp(time.Now())
	t.UpdateTime = t.CreateTime
}
