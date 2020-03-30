package builder

import (
	"testing"
)

func TestExpressionTest(t *testing.T) {
	exprCases{
		Case(
			"empty",
			ExprFrom(nil),
			nil,
		),
		Case(
			"expr",
			ExprFrom(Expr("")),
			Expr(""),
		),
	}.Run(t, "Expression")
}
