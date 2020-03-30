package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhere(t *testing.T) {
	tt := assert.New(t)
	w := (*where)(nil)
	tt.Nil(w.Expr())
}

func TestStmt(t *testing.T) {
	s := Stmt{}
	if s.Type() != STMT_RAW {
		panic("Stmt should be STMT_RAW")
	}
}

func TestStmtCommon_Order(t *testing.T) {
	o := order{}
	table := T(DB("db"), "t")

	o2 := order{
		expr: Col(table, "F_id").In(),
	}

	exprCases{
		Case(
			"empty groupBy",
			o.Expr(),
			nil,
		),
		Case(
			"empty groupBy",
			o2.Expr(),
			nil,
		),
	}.Run(t, "Stmt common order")
}

func TestStmtCommon_GroupBy(t *testing.T) {
	o := groupBy{}
	table := T(DB("db"), "t")

	exprCases{
		Case(
			"empty groupBy",
			o.Expr(),
			nil,
		),
		Case(
			"simple groupBy",
			o.setBy(Col(table, "F_id"), "").Expr(),
			Expr("GROUP BY (`F_id`)"),
		),
		Case(
			"simple groupBy with order",
			o.setBy(Col(table, "F_id"), ORDER_DESC).Expr(),
			Expr("GROUP BY (`F_id`) DESC"),
		),
		Case(
			"simple groupBy with rollup",
			o.setBy(Col(table, "F_id"), ORDER_DESC).rollup().Expr(),
			Expr("GROUP BY (`F_id`) DESC WITH ROLLUP"),
		),
		Case(
			"multi groupBy",
			o.addBy(Col(table, "F_id"), "").
				addBy(Col(table, "F_b"), "").Expr(),
			Expr("GROUP BY (`F_id`), (`F_b`)"),
		),
		Case(
			"multi groupBy with order",
			o.addBy(Col(table, "F_id"), ORDER_DESC).
				addBy(Col(table, "F_b"), ORDER_ASC).Expr(),
			Expr("GROUP BY (`F_id`) DESC, (`F_b`) ASC"),
		),
		Case(
			"multi groupBy with having",
			o.addBy(Col(table, "F_id"), "").
				having(Col(table, "F_id").Eq(1)).Expr(),
			Expr(
				"GROUP BY (`F_id`) HAVING `F_id` = ?",
				1,
			),
		),
	}.Run(t, "Stmt common group by")
}

func TestStmtCommon_OrderBy(t *testing.T) {
	o := orderBy{}
	table := T(DB("db"), "t")

	exprCases{
		Case(
			"empty orderBy",
			o.Expr(),
			nil,
		),
		Case(
			"simple order by",
			o.setBy(Col(table, "F_id"), ORDER_DESC).Expr(),
			Expr("ORDER BY (`F_id`) DESC"),
		),
		Case(
			"multi order by",
			o.addBy(Col(table, "F_id"), ORDER_DESC).
				addBy(Col(table, "F_b"), ORDER_ASC).Expr(),
			Expr("ORDER BY (`F_id`) DESC, (`F_b`) ASC"),
		),
	}.Run(t, "Stmt common order by")
}

func TestStmtCommon_Limit(t *testing.T) {
	limit := limit{}

	exprCases{
		Case(
			"empty limit",
			limit.Expr(),
			nil,
		),
		Case(
			"limit only with offset",
			limit.offset(1).Expr(),
			nil,
		),
		Case(
			"limit with size",
			limit.limit(1).Expr(),
			Expr(
				"LIMIT 1",
			),
		),
		Case(
			"limit with size and offset",
			limit.limit(1).offset(1).Expr(),
			Expr(
				"LIMIT 1 OFFSET 1",
			),
		),
	}.Run(t, "Stmt common limit")
}
