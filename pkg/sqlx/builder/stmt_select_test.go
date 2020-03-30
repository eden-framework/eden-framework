package builder

import (
	"testing"
)

func TestStmtSelect(t *testing.T) {
	table := T(DB("db"), "t")

	if table.Select().Type() != STMT_SELECT {
		panic("select type should be STMT_SELECT")
	}

	exprCases{
		Case(
			"Select with modifier",
			table.Select().Modifier("DISTINCT").Where(
				Col(table, "F_a").Eq(1),
			).Expr(),
			Expr(
				"SELECT DISTINCT * FROM `db`.`t` WHERE `F_a` = ?",
				1,
			),
		),
		Case(
			"SELECT simple",
			table.Select().Comment("Comment").Where(
				Col(table, "F_a").Eq(1),
			).Expr(),
			Expr(
				"/* Comment */ SELECT * FROM `db`.`t` WHERE `F_a` = ?",
				1,
			),
		),
		Case(
			"SELECT with select expr",
			table.Select().For(Col(table, "F_a")).Where(
				Col(table, "F_a").Eq(1),
			).Expr(),
			Expr(
				"SELECT `F_a` FROM `db`.`t` WHERE `F_a` = ?",
				1,
			),
		),
		Case(
			"Select for update",
			table.Select().
				Where(Col(table, "F_a").Eq(1)).
				ForUpdate().
				Expr(),
			Expr(
				"SELECT * FROM `db`.`t` WHERE `F_a` = ? FOR UPDATE",
				1,
			),
		),
		Case(
			"Select with group by",
			table.Select().
				Where(Col(table, "F_a").Eq(1)).
				GroupBy(Col(table, "F_a")).
				Having(Col(table, "F_a").Eq(1)).
				Expr(),
			Expr(
				"SELECT * FROM `db`.`t` WHERE `F_a` = ? GROUP BY (`F_a`) HAVING `F_a` = ?",
				1, 1,
			),
		),
		Case(
			"Select with group by with rollup",
			table.Select().
				Where(Col(table, "F_a").Eq(1)).
				GroupBy(Col(table, "F_a")).
				WithRollup().
				Having(Col(table, "F_a").Eq(1)).
				Expr(),
			Expr(
				"SELECT * FROM `db`.`t` WHERE `F_a` = ? GROUP BY (`F_a`) WITH ROLLUP HAVING `F_a` = ?",
				1, 1,
			),
		),
		Case(
			"Select with desc group by",
			table.Select().
				Where(Col(table, "F_a").Eq(1)).
				GroupDescBy(Col(table, "F_b")).
				Expr(),
			Expr(
				"SELECT * FROM `db`.`t` WHERE `F_a` = ? GROUP BY (`F_b`) DESC",
				1,
			),
		),
		Case(
			"Select with combined ordered group by ",
			table.Select().
				Where(Col(table, "F_a").Eq(1)).
				GroupAscBy(Col(table, "F_a")).
				GroupDescBy(Col(table, "F_b")).
				Expr(),
			Expr(
				"SELECT * FROM `db`.`t` WHERE `F_a` = ? GROUP BY (`F_a`) ASC, (`F_b`) DESC",
				1,
			),
		),
		Case(
			"Select with having only should skip",
			table.Select().
				Where(Col(table, "F_a").Eq(1)).
				Having(Col(table, "F_a").Eq(1)).
				Expr(),
			Expr(
				"SELECT * FROM `db`.`t` WHERE `F_a` = ?",
				1,
			),
		),
		Case(
			"Select with limit",
			table.Select().
				Where(
					Col(table, "F_a").Eq(1),
				).
				Limit(1).
				Expr(),
			Expr(
				"SELECT * FROM `db`.`t` WHERE `F_a` = ? LIMIT 1",
				1,
			),
		),
		Case(
			"Select with limit with offset",
			table.Select().
				Where(
					Col(table, "F_a").Eq(1),
				).
				Offset(200).
				Limit(1).
				Expr(),
			Expr(
				"SELECT * FROM `db`.`t` WHERE `F_a` = ? LIMIT 1 OFFSET 200",
				1,
			),
		),
		Case(
			"Select with order",
			table.Select().
				Where(Col(table, "F_a").Eq(1)).
				OrderAscBy(Col(table, "F_a")).
				OrderDescBy(Col(table, "F_b")).
				Expr(),
			Expr(
				"SELECT * FROM `db`.`t` WHERE `F_a` = ? ORDER BY (`F_a`) ASC, (`F_b`) DESC",
				1,
			),
		),
	}.Run(t, "Stmt select")
}
