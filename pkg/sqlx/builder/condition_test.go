package builder

import (
	"testing"
)

func TestBuilderCond(t *testing.T) {
	table := T(DB("db"), "t")

	exprCases{
		Case("CondRules",
			NewCondRules().
				When(true, Col(table, "a").Eq(1)).
				When(true, Col(table, "b").Like(`g`)).
				When(false, Col(table, "b").Like(`g`)).
				ToCond(),
			Expr(
				"(`a` = ?) AND (`b` LIKE ?)",
				1, "%g%",
			),
		),
		Case(
			"Chain Condition",
			Col(table, "a").Eq(1).
				And(Col(table, "b").LeftLike("c")).
				Or(Col(table, "a").Eq(2)).
				Xor(Col(table, "b").RightLike("g")).Expr(),
			Expr(
				"(((`a` = ?) AND (`b` LIKE ?)) OR (`a` = ?)) XOR (`b` LIKE ?)",
				1, "%c", 2, "g%",
			),
		),
		Case(
			"Compose Condition",
			Xor(
				Or(
					And(
						Col(table, "a").Eq(1),
						Col(table, "b").Like("c"),
					),
					Col(table, "a").Eq(2),
				),
				Col(table, "b").Like("g"),
			).Expr(),
			Expr(
				"(((`a` = ?) AND (`b` LIKE ?)) OR (`a` = ?)) XOR (`b` LIKE ?)",
				1, "%c%", 2, "%g%",
			),
		),
		Case(
			"Skip nil",
			Xor(
				Col(table, "a").In(),
				Or(
					Col(table, "a").NotIn(),
					And(
						nil,
						Col(table, "a").Eq(1),
						Col(table, "b").Like("c"),
					),
					Col(table, "a").Eq(2),
				),
				Col(table, "b").Like("g"),
			).Expr(),
			Expr(
				"(((`a` = ?) AND (`b` LIKE ?)) OR (`a` = ?)) XOR (`b` LIKE ?)",
				1, "%c%", 2, "%g%",
			),
		),
		Case(
			"XOR",
			Xor(
				Col(table, "a").In(),
				Or(
					Col(table, "a").NotIn(),
					And(
						nil,
						Col(table, "a").Eq(1),
						Col(table, "b").Like("c"),
					),
					Col(table, "a").Eq(2),
				),
				Col(table, "b").Like("g"),
			).Expr(),
			Expr(
				"(((`a` = ?) AND (`b` LIKE ?)) OR (`a` = ?)) XOR (`b` LIKE ?)",
				1, "%c%", 2, "%g%",
			),
		),
		Case(
			"XOR",
			Xor(
				Col(table, "a").Eq(1),
				Col(table, "b").Like("g"),
			).Expr(),
			Expr(
				"(`a` = ?) XOR (`b` LIKE ?)",
				1, "%g%",
			),
		),
		Case(
			"Like",
			Col(table, "d").Like("e").Expr(),
			Expr(
				"`d` LIKE ?",
				"%e%",
			),
		),
		Case(
			"Not like",
			Col(table, "d").NotLike("e").Expr(),
			Expr(
				"`d` NOT LIKE ?",
				"%e%",
			),
		),
		Case(
			"Equal",
			Col(table, "d").Eq("e").Expr(),
			Expr(
				"`d` = ?",
				"e",
			),
		),
		Case(
			"Not Equal",
			Col(table, "d").Neq("e").Expr(),
			Expr(
				"`d` <> ?",
				"e",
			),
		),
		Case(
			"In",
			Col(table, "d").In("e", "f").Expr(),
			Expr(
				"`d` IN (?,?)",
				"e", "f",
			),
		),
		Case(
			"In With Select",
			Col(table, "d").In(SelectFrom(table).Where(Col(table, "d").Eq(1))).Expr(),
			Expr(
				"`d` IN (SELECT * FROM `db`.`t` WHERE `d` = ?)",
				1,
			),
		),
		Case(
			"NotIn",
			Col(table, "d").NotIn("e", "f").Expr(),
			Expr(
				"`d` NOT IN (?,?)",
				"e", "f",
			),
		),
		Case(
			"Not In With Select",
			Col(table, "d").NotIn(SelectFrom(table).Where(Col(table, "d").Eq(1))).Expr(),
			Expr(
				"`d` NOT IN (SELECT * FROM `db`.`t` WHERE `d` = ?)",
				1,
			),
		),
		Case(
			"Less than",
			Col(table, "d").Lt(3).Expr(),
			Expr(
				"`d` < ?",
				3,
			),
		),
		Case(
			"Less or equal than",
			Col(table, "d").Lte(3).Expr(),
			Expr(
				"`d` <= ?",
				3,
			),
		),
		Case(
			"Greater than",
			Col(table, "d").Gt(3).Expr(),
			Expr(
				"`d` > ?",
				3,
			),
		),
		Case(
			"Greater or equal than",
			Col(table, "d").Gte(3).Expr(),
			Expr(
				"`d` >= ?",
				3,
			),
		),
		Case(
			"Between",
			Col(table, "d").Between(0, 2).Expr(),
			Expr(
				"`d` BETWEEN ? AND ?",
				0, 2,
			),
		),
		Case(
			"Not between",
			Col(table, "d").NotBetween(0, 2).Expr(),
			Expr(
				"`d` NOT BETWEEN ? AND ?",
				0, 2,
			),
		),
		Case(
			"Is null",
			Col(table, "d").IsNull().Expr(),
			Expr(
				"`d` IS NULL",
			),
		),
		Case(
			"Is not null",
			Col(table, "d").IsNotNull().Expr(),
			Expr(
				"`d` IS NOT NULL",
			),
		),
	}.Run(t, "Condition")
}
