package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	tt := assert.New(t)

	db := DB("db")
	table := T(db, "t")

	tableNext := table.Define(
		Col(table, "F_id").Type("bigint(64) unsigned NOT NULL AUTO_INCREMENT"),
		Col(table, "F_name").Type("varchar(128) NOT NULL DEFAULT ''"),
		Col(table, "F_created_at").Type("bigint(64) NOT NULL DEFAULT '0'"),
		Col(table, "F_updated_at").Type("bigint(64) NOT NULL DEFAULT '0'"),
		Index("I_name").WithCols(Col(table, "F_name")),
		UniqueIndex("I_username").WithCols(Col(table, "F_username"), Col(table, "F_id")),
	)

	table = table.Define(
		Col(table, "F_id").Field("ID"), // skip without type
		Col(table, "F_id").Field("ID").Type("bigint(64) unsigned NOT NULL AUTO_INCREMENT"),
		Col(table, "F_name").Field("Name").Type("varchar(255) NOT NULL DEFAULT ''"),
		Col(table, "F_username").Field("Username").Type("varchar(255) NOT NULL DEFAULT ''"),
		Col(table, "F_created_at").Field("CreatedAt").Type("bigint(64) NOT NULL DEFAULT '0'"),
		Col(table, "F_updated_at").Field("UpdatedAt").Type("bigint(64) NOT NULL DEFAULT '0'"),
		PrimaryKey(), // skip without Columns
		PrimaryKey().WithCols(Col(table, "F_id")),
		Index("I_name").WithCols(Col(table, "F_name")),
		UniqueIndex("I_username").WithCols(Col(table, "F_name"), Col(table, "F_id")),
	)

	{
		cols, values := table.ColumnsAndValuesByFieldValues(FieldValues{
			"ID":   1,
			"Name": "1",
		})

		tt.Equal(table.Fields("ID", "Name").Len(), cols.Len())
		if values[0] == 1 {
			tt.Equal([]interface{}{1, "1"}, values)
		} else {
			tt.Equal([]interface{}{"1", 1}, values)
		}
	}

	{
		assignments := table.AssignsByFieldValues(FieldValues{
			"ID":   1,
			"Name": "1",
		})

		expr := assignments.Expr()

		if expr.Args[0] == 1 {
			tt.Equal(Expr("`F_id` = ?, `F_name` = ?", 1, "1"), assignments.Expr())
		} else {
			tt.Equal(Expr("`F_name` = ?, `F_id` = ?", "1", 1), assignments.Expr())
		}

	}

	exprCases{
		Case(
			"create Table",
			table.Create(true),
			Expr("CREATE TABLE IF NOT EXISTS `db`.`t` ("+
				"`F_id` bigint(64) unsigned NOT NULL AUTO_INCREMENT, "+
				"`F_name` varchar(255) NOT NULL DEFAULT '', "+
				"`F_username` varchar(255) NOT NULL DEFAULT '', "+
				"`F_created_at` bigint(64) NOT NULL DEFAULT '0', "+
				"`F_updated_at` bigint(64) NOT NULL DEFAULT '0', "+
				"PRIMARY KEY (`F_id`), "+
				"INDEX `I_name` (`F_name`), "+
				"UNIQUE INDEX `I_username` (`F_name`,`F_id`)"+
				") ENGINE=InnoDB CHARSET=utf8"),
		),
		Case(
			"create Table",
			table.Create(false),
			Expr("CREATE TABLE `db`.`t` ("+
				"`F_id` bigint(64) unsigned NOT NULL AUTO_INCREMENT, "+
				"`F_name` varchar(255) NOT NULL DEFAULT '', "+
				"`F_username` varchar(255) NOT NULL DEFAULT '', "+
				"`F_created_at` bigint(64) NOT NULL DEFAULT '0', "+
				"`F_updated_at` bigint(64) NOT NULL DEFAULT '0', "+
				"PRIMARY KEY (`F_id`), "+
				"INDEX `I_name` (`F_name`), "+
				"UNIQUE INDEX `I_username` (`F_name`,`F_id`)"+
				") ENGINE=InnoDB CHARSET=utf8"),
		),
		Case(
			"cond",
			table.Cond("#ID = ? AND #Username = ?"),
			Expr("`F_id` = ? AND `F_username` = ?"),
		),
		Case(
			"cond with unregister col field",
			table.Cond("#ID = ? AND #Usernames = ?"),
			Expr("`F_id` = ? AND #Usernames = ?"),
		),
		Case(
			"diff for migrate",
			table.Diff(tableNext),
			Expr("ALTER TABLE `db`.`t` "+
				"DROP COLUMN `F_username`, "+
				"MODIFY COLUMN `F_name` varchar(128) NOT NULL DEFAULT '', "+
				"DROP PRIMARY KEY, "+
				"DROP INDEX `I_username`, ADD UNIQUE INDEX `I_username` (`F_username`,`F_id`)",
			),
		),
		Case(
			"revert diff for migrate",
			tableNext.Diff(table),
			Expr("ALTER TABLE `db`.`t` "+
				"MODIFY COLUMN `F_name` varchar(255) NOT NULL DEFAULT '', "+
				"ADD COLUMN `F_username` varchar(255) NOT NULL DEFAULT '', "+
				"DROP INDEX `I_username`, ADD UNIQUE INDEX `I_username` (`F_name`,`F_id`), "+
				"ADD PRIMARY KEY (`F_id`)",
			),
		),
		Case(
			"diff without change",
			table.Diff(table),
			nil,
		),
		Case(
			"drop Table",
			table.Drop(),
			Expr("DROP TABLE `db`.`t`"),
		),
		Case(
			"truncate Table",
			table.Truncate(),
			Expr("TRUNCATE TABLE `db`.`t`"),
		),
	}.Run(t, "Table")
}
