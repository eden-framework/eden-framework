package builder

import (
	"testing"
)

func TestStmtUpdate(t *testing.T) {
	table := T(DB("db"), "t")

	if table.Update().Type() != STMT_UPDATE {
		panic("update type should be STMT_UPDATE")
	}

	exprCases{
		Case(
			"Update with err",
			table.Update().Set(
				Col(table, "F_a").By(1),
			).Expr(),
			ExprErr(UpdateNeedLimitByWhere),
		),
		Case(
			"Update with modifier",
			table.Update().Modifier("IGNORE").Set(
				Col(table, "F_a").By(1),
			).Where(
				Col(table, "F_a").Eq(1),
			).Expr(),
			Expr(
				"UPDATE IGNORE `db`.`t` SET `F_a` = ? WHERE `F_a` = ?",
				1, 1,
			),
		),
		Case(
			"Update simple",
			table.Update().Comment("Comment").Set(
				Col(table, "F_a").By(1),
				Col(table, "F_b").By(2),
			).Where(
				Col(table, "F_a").Eq(1),
			).Expr(),
			Expr(
				"/* Comment */ UPDATE `db`.`t` SET `F_a` = ?, `F_b` = ? WHERE `F_a` = ?",
				1, 2, 1,
			),
		),
		Case(
			"Update with limit",
			table.Update().Set(
				Col(table, "F_a").By(1),
			).
				Where(
					Col(table, "F_a").Eq(1),
				).
				Limit(1).
				Expr(),
			Expr(
				"UPDATE `db`.`t` SET `F_a` = ? WHERE `F_a` = ? LIMIT 1",
				1, 1,
			),
		),
		Case(
			"Update with order",
			table.Update().Set(
				Col(table, "F_a").By(Col(table, "F_a").Incr(1)),
				Col(table, "F_b").By(Col(table, "F_b").Desc(2)),
			).Where(
				Col(table, "F_a").Eq(3),
			).OrderDescBy(
				Col(table, "F_b"),
			).OrderAscBy(
				Col(table, "F_a"),
			).Expr(),
			Expr(
				"UPDATE `db`.`t` SET `F_a` = `F_a` + ?, `F_b` = `F_b` - ? WHERE `F_a` = ? ORDER BY (`F_a`) ASC",
				1, 2, 3,
			),
		),
	}.Run(t, "Stmt update")
}
