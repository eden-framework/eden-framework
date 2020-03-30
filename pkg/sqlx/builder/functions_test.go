package builder

import (
	"testing"
)

func TestFunc(t *testing.T) {
	exprCases{
		Case(
			"Nil",
			Func("", CanExpr(nil)),
			nil,
		),
		Case(
			"COUNT",
			Count(Star()),
			Expr("COUNT(*)"),
		),
		Case(
			"AVG",
			Avg(Star()),
			Expr("AVG(*)"),
		),
		Case(
			"DISTINCT",
			Distinct(Star()),
			Expr("DISTINCT(*)"),
		),
		Case(
			"MIN",
			Min(Star()),
			Expr("MIN(*)"),
		),
		Case(
			"Max",
			Max(Star()),
			Expr("MAX(*)"),
		),
		Case(
			"First",
			First(Star()),
			Expr("FIRST(*)"),
		),
		Case(
			"Last",
			Last(Star()),
			Expr("LAST(*)"),
		),
		Case(
			"Sum",
			Sum(Star()),
			Expr("SUM(*)"),
		),
	}.Run(t, "Function")

}
