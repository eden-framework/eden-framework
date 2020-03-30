package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type exprCases []*exprCase

func (e exprCases) Run(t *testing.T, group string) {
	for _, c := range e {
		t.Log(group + ": " + c.desc)
		assert.Equal(t, ExprFrom(c.expectExpr), ExprFrom(c.expr), c.desc)
	}
}

type exprCase struct {
	desc       string
	expectExpr CanExpr
	expr       CanExpr
}

func Case(desc string, expr CanExpr, expectExpr CanExpr) *exprCase {
	return &exprCase{
		desc:       desc,
		expectExpr: expectExpr,
		expr:       expr,
	}
}

func TestHolderRepeat(t *testing.T) {
	tt := assert.New(t)
	tt.Equal("?,?,?,?,?", HolderRepeat(5))
}

func TestQuote(t *testing.T) {
	tt := assert.New(t)
	tt.Equal("`name`", quote("name"))
}
