package sqlx

import (
	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

type EnumTypeDescriber interface {
	EnumType() string
	Enums() map[int][]string
}

type SqlMetaEnum struct {
	TName string `db:"F_table_name" sql:"varchar(64) NOT NULL"`
	CName string `db:"F_column_name" sql:"varchar(64) NOT NULL"`
	Value int    `db:"F_value" sql:"int NOT NULL"`
	Type  string `db:"F_type" sql:"varchar(255) NOT NULL"`
	Key   string `db:"F_key"  sql:"varchar(255) NOT NULL"`
	Label string `db:"F_label" sql:"varchar(255) NOT NULL"`
}

func (*SqlMetaEnum) TableName() string {
	return "t_sql_meta_enum"
}

func (*SqlMetaEnum) UniqueIndexes() Indexes {
	return Indexes{"I_enum": FieldNames{"TName", "CName", "Value"}}
}

func (database *Database) SyncEnum(db *DB) error {
	task := NewTasks(db)

	metaEnumTable := database.T(&SqlMetaEnum{})

	task = task.With(func(db *DB) error {
		return db.Do(metaEnumTable.Create(true)).Err()
	})

	task = task.With(func(db *DB) error {
		return db.Do(metaEnumTable.Delete().Where(metaEnumTable.F("TName").In(database.TableNames()))).Err()
	})

	stmt := metaEnumTable.Insert()
	hasEnum := false

	for _, table := range database.Tables {
		table.Columns.Range(func(col *builder.Column, idx int) {
			if col.IsEnum() {
				for val := range col.Enums {
					hasEnum = true

					sqlMetaEnum := &SqlMetaEnum{
						TName: table.Name,
						CName: col.Name,
						Type:  col.EnumType,
						Value: val,
					}

					enum := col.Enums[val]

					if len(enum) > 0 {
						sqlMetaEnum.Key = enum[0]
					}

					if len(enum) > 1 {
						sqlMetaEnum.Label = enum[1]
					}

					fieldValues := FieldValuesFromStructByNonZero(sqlMetaEnum)
					cols, vals := metaEnumTable.ColumnsAndValuesByFieldValues(fieldValues)
					stmt = stmt.Columns(cols).Values(vals...)
				}
			}
		})
	}

	if hasEnum {
		task = task.With(func(db *DB) error {
			return db.Do(stmt).Err()
		})
	}

	return task.Do()
}
