package test

import (
	"github.com/profzone/eden-framework/pkg/sqlx/presets"
	"github.com/profzone/eden-framework/pkg/timelib"
)

// @def primary ID
// @def index I_nickname Nickname
// @def index I_username Username
// @def unique_index I_name Name
type User struct {
	// 姓名
	Name     string                `db:"F_name" json:"name" sql:"varchar(255) NOT NULL DEFAULT ''"`
	Username string                `db:"F_username" json:"username" sql:"varchar(255) NOT NULL DEFAULT ''"`
	Nickname string                `db:"F_nickname" json:"nickname" sql:"varchar(255) NOT NULL DEFAULT ''"`
	Gender   Gender                `db:"F_gender" json:"gender" sql:"int NOT NULL DEFAULT '0'"`
	Birthday timelib.MySQLDatetime `db:"F_birthday" json:"birthday" sql:"datetime NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	Boolean  bool                  `db:"F_boolean" json:"boolean" sql:"boolean NOT NULL DEFAULT '0'"`
	presets.OperateTime
	presets.PrimaryID
	presets.SoftDelete
}

type User2 struct {
	Name string `db:"F_name" json:"name" sql:"varchar(255) NOT NULL DEFAULT ''"`
}
