package builder

import (
	"testing"
)

func TestStmtDelete(t *testing.T) {
	table := T(DB("db"), "t")

	if table.Delete().Type() != STMT_DELETE {
		panic("Delete type should be STMT_DELETE")
	}

	exprCases{
		Case(
			"Delete with err",
			table.Delete().Expr(),
			ExprErr(DeleteNeedLimitByWhere),
		),
		Case(
			"Delete with modifier",
			table.Delete().Modifier("IGNORE").Where(
				Col(table, "F_a").Eq(1),
			).Expr(),
			Expr(
				"DELETE IGNORE FROM `db`.`t` WHERE `F_a` = ?",
				1,
			),
		),
		Case(
			"Delete simple",
			table.Delete().Comment("Comment").Where(
				Col(table, "F_a").Eq(1),
			).Expr(),
			Expr(
				"/* Comment */ DELETE FROM `db`.`t` WHERE `F_a` = ?",
				1,
			),
		),
		Case(
			"Delete with limit",
			table.Delete().
				Where(
					Col(table, "F_a").Eq(1),
				).
				Limit(1).
				Expr(),
			Expr(
				"DELETE FROM `db`.`t` WHERE `F_a` = ? LIMIT 1",
				1,
			),
		),
		Case(
			"Delete with order",
			table.Delete().
				Where(Col(table, "F_a").Eq(1)).
				DescendBy(Col(table, "F_b")).
				AscendBy(Col(table, "F_a")).
				Expr(),
			Expr(
				"DELETE FROM `db`.`t` WHERE `F_a` = ? ORDER BY (`F_a`) ASC",
				1,
			),
		),
	}.Run(t, "Stmt delete")
}
