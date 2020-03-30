package presets

type ReserveFields struct {
	// 预留整数1
	ReserveInt1 int32 `db:"F_reserve_int1" json:"-" sql:"int(32) NOT NULL DEFAULT '0'"`
	// 预留整数2
	ReserveInt2 int32 `db:"F_reserve_int2" json:"-" sql:"int(32) NOT NULL DEFAULT '0'"`
	// 预留字符串1
	ReserveString1 string `db:"F_reserve_string1" json:"-" sql:"varchar(255) NOT NULL DEFAULT ''"`
	// 预留字符串2
	ReserveString2 string `db:"F_reserve_string2" json:"-" sql:"varchar(255) NOT NULL DEFAULT ''"`
}
