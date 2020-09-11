package enummeta

import (
	"github.com/profzone/eden-framework/pkg/enumeration"
	"reflect"

	"github.com/profzone/eden-framework/pkg/sqlx"
	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

func SyncEnum(db sqlx.DBExecutor) error {
	metaEnumTable := builder.T((&SqlMetaEnum{}).TableName())
	builder.ScanDefToTable(reflect.ValueOf(&SqlMetaEnum{}), metaEnumTable)

	dialect := db.Dialect()

	task := sqlx.NewTasks(db.WithSchema(""))

	task = task.With(func(db sqlx.DBExecutor) error {
		_, err := db.ExecExpr(dialect.DropTable(metaEnumTable))
		return err
	})

	exprs := dialect.CreateTableIsNotExists(metaEnumTable)

	for i := range exprs {
		expr := exprs[i]
		task = task.With(func(db sqlx.DBExecutor) error {
			_, err := db.ExecExpr(expr)
			return err
		})
	}

	{
		// insert values
		stmtForInsert := builder.Insert().Into(metaEnumTable)
		vals := make([]interface{}, 0)

		columns := &builder.Columns{}

		db.D().Tables.Range(func(table *builder.Table, idx int) {
			table.Columns.Range(func(col *builder.Column, idx int) {
				v := reflect.New(col.ColumnType.Type).Interface()
				if enumValue, ok := v.(enumeration.EnumTypeDescriber); ok {
					for value, enum := range enumValue.Enums() {
						sqlMetaEnum := &SqlMetaEnum{
							TName: table.Name,
							CName: col.Name,
							Value: value,
							Type:  enumValue.EnumType(),
						}

						if len(enum) > 0 {
							sqlMetaEnum.Key = enum[0]
						}

						if len(enum) > 1 {
							sqlMetaEnum.Label = enum[1]
						}

						fieldValues := builder.FieldValuesFromStructByNonZero(sqlMetaEnum, "Value")
						cols, values := metaEnumTable.ColumnsAndValuesByFieldValues(fieldValues)
						vals = append(vals, values...)
						columns = cols
					}
				}
			})
		})

		if len(vals) > 0 {
			stmtForInsert = stmtForInsert.Values(columns, vals...)

			task = task.With(func(db sqlx.DBExecutor) error {
				_, err := db.ExecExpr(stmtForInsert)
				return err
			})
		}
	}

	return task.Do()
}
