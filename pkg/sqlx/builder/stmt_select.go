package builder

import (
	"fmt"
)

func SelectFrom(table *Table) *StmtSelect {
	return &StmtSelect{
		table: table,
	}
}

type StmtSelect struct {
	table     *Table
	expr      CanExpr
	modifiers []string
	*where
	*groupBy
	*orderBy
	*limit
	forUpdate bool
	comment   string
}

func (s StmtSelect) Comment(comment string) *StmtSelect {
	s.comment = comment
	return &s
}

func (s *StmtSelect) Type() StmtType {
	return STMT_SELECT
}

func (s StmtSelect) For(expr CanExpr) *StmtSelect {
	s.expr = expr
	return &s
}

func (s StmtSelect) Modifier(modifiers ...string) *StmtSelect {
	s.modifiers = append(s.modifiers, modifiers...)
	return &s
}

func (s StmtSelect) ForUpdate() *StmtSelect {
	s.forUpdate = true
	return &s
}

func (s StmtSelect) Where(cond *Condition) *StmtSelect {
	s.where = (*where)(cond)
	return &s
}

func (s StmtSelect) GroupBy(canExpr CanExpr) *StmtSelect {
	if s.groupBy == nil {
		s.groupBy = &groupBy{}
	}
	s.groupBy = s.groupBy.addBy(canExpr, "")
	return &s
}

func (s StmtSelect) GroupAscBy(canExpr CanExpr) *StmtSelect {
	if s.groupBy == nil {
		s.groupBy = &groupBy{}
	}
	s.groupBy = s.groupBy.addBy(canExpr, ORDER_ASC)
	return &s
}

func (s StmtSelect) GroupDescBy(canExpr CanExpr) *StmtSelect {
	if s.groupBy == nil {
		s.groupBy = &groupBy{}
	}
	s.groupBy = s.groupBy.addBy(canExpr, ORDER_DESC)
	return &s
}

func (s StmtSelect) WithRollup() *StmtSelect {
	s.groupBy = s.groupBy.rollup()
	return &s
}

func (s StmtSelect) Having(cond *Condition) *StmtSelect {
	if s.groupBy == nil {
		s.groupBy = &groupBy{}
	}
	s.groupBy = s.groupBy.having(cond)
	return &s
}

func (s StmtSelect) OrderBy(canExpr CanExpr, orderType OrderType) *StmtSelect {
	if s.orderBy == nil {
		s.orderBy = &orderBy{}
	}
	s.orderBy = s.orderBy.addBy(canExpr, orderType)
	return &s
}

func (s StmtSelect) OrderAscBy(canExpr CanExpr) *StmtSelect {
	return s.OrderBy(canExpr, ORDER_ASC)
}

func (s StmtSelect) OrderDescBy(canExpr CanExpr) *StmtSelect {
	return s.OrderBy(canExpr, ORDER_DESC)
}

func (s StmtSelect) Limit(size int32) *StmtSelect {
	if s.limit == nil {
		s.limit = &limit{}
	}
	s.limit = s.limit.limit(size)
	return &s
}

func (s StmtSelect) Offset(offset int32) *StmtSelect {
	if s.limit == nil {
		s.limit = &limit{}
	}
	s.limit = s.limit.offset(offset)
	return &s
}

func (s *StmtSelect) Expr() *Expression {
	selectExpr := Expr("*")
	if s.expr != nil {
		selectExpr = s.expr.Expr()
	}

	expr := Expr(
		fmt.Sprintf("%s %s FROM %s", statement(s.comment, "SELECT", s.modifiers...), selectExpr.Query, s.table.FullName()),
		selectExpr.Args...,
	)

	expr = expr.ConcatBy(" ", s.where)
	expr = expr.ConcatBy(" ", s.groupBy)
	expr = expr.ConcatBy(" ", s.orderBy)
	expr = expr.ConcatBy(" ", s.limit)

	if s.forUpdate {
		expr = expr.ConcatBy(" ", Expr("FOR UPDATE"))
	}

	return expr
}
