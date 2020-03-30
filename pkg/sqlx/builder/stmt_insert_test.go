package builder

import (
	"testing"
)

func TestStmtInsert(t *testing.T) {
	table := T(DB("db"), "t")

	if table.Insert().Type() != STMT_INSERT {
		panic("Insert type should be STMT_INSERT")
	}

	exprCases{
		Case(
			"Insert sql failed",
			table.
				Insert().
				Columns(Cols(table, "F_a", "F_b")).
				Values(1).
				Expr(),
			ExprErr(InsertValuesLengthNotMatch),
		),
		Case(
			"Insert simple",
			table.
				Insert().
				Comment("Comment").
				Columns(Cols(table, "F_a", "F_b")).
				Values(1, 2).
				Expr(),
			Expr(
				"/* Comment */ INSERT INTO `db`.`t` (`F_a`,`F_b`) VALUES (?,?)",
				1, 2,
			),
		),
		Case(
			"Insert with modifier",
			table.
				Insert().
				Modifier("IGNORE").
				Columns(Cols(table, "F_a", "F_b")).
				Values(1, 2).
				Expr(),
			Expr(
				"INSERT IGNORE INTO `db`.`t` (`F_a`,`F_b`) VALUES (?,?)",
				1, 2,
			),
		),
		Case(
			"Insert on on duplicate key update",
			Insert(table).
				Columns(Cols(table, "F_a", "F_b")).
				Values(1, 2).
				OnDuplicateKeyUpdate(Col(table, "F_b").By(2)).
				Expr(),
			Expr(
				"INSERT INTO `db`.`t` (`F_a`,`F_b`) VALUES (?,?) ON DUPLICATE KEY UPDATE `F_b` = ?",
				1, 2, 2,
			),
		),
		Case(
			"Insert multiple",
			Insert(table).
				Columns(Cols(table, "F_a", "F_b")).
				Values(1, 2).
				Values(1, 2).
				Values(1, 2).
				Expr(),
			Expr(
				"INSERT INTO `db`.`t` (`F_a`,`F_b`) VALUES (?,?),(?,?),(?,?)",
				1, 2, 1, 2, 1, 2,
			),
		),
	}.Run(t, "Stmt insert")
}
