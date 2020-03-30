package builder

import (
	"fmt"
)

var (
	UpdateNeedLimitByWhere = fmt.Errorf("no where limit for update")
)

func Update(table *Table) *StmtUpdate {
	return &StmtUpdate{
		table: table,
	}
}

type StmtUpdate struct {
	table     *Table
	modifiers []string
	Assignments
	*where
	*orderBy
	*limit
	comment string
}

func (s StmtUpdate) Comment(comment string) *StmtUpdate {
	s.comment = comment
	return &s
}

func (s *StmtUpdate) Type() StmtType {
	return STMT_UPDATE
}

func (s StmtUpdate) Modifier(modifiers ...string) *StmtUpdate {
	s.modifiers = append(s.modifiers, modifiers...)
	return &s
}

func (s StmtUpdate) Set(assignments ...*Assignment) *StmtUpdate {
	if s.Assignments == nil {
		s.Assignments = []*Assignment{}
	}
	s.Assignments = append(s.Assignments, assignments...)
	return &s
}

func (s StmtUpdate) Where(cond *Condition) *StmtUpdate {
	s.where = (*where)(cond)
	return &s
}

func (s StmtUpdate) OrderBy(col *Column, orderType OrderType) *StmtUpdate {
	if s.orderBy == nil {
		s.orderBy = &orderBy{}
	}
	s.orderBy = s.orderBy.setBy(col, orderType)
	return &s
}

func (s StmtUpdate) OrderAscBy(col *Column) *StmtUpdate {
	return s.OrderBy(col, ORDER_ASC)
}

func (s StmtUpdate) OrderDescBy(col *Column) *StmtUpdate {
	return s.OrderBy(col, ORDER_DESC)
}

func (s StmtUpdate) Limit(size int32) *StmtUpdate {
	if s.limit == nil {
		s.limit = &limit{}
	}
	s.limit = s.limit.limit(size)
	return &s
}

func (s *StmtUpdate) Expr() *Expression {
	if s.where == nil {
		return ExprErr(UpdateNeedLimitByWhere)
	}

	expr := Expr(fmt.Sprintf("%s %s SET", statement(s.comment, "UPDATE", s.modifiers...), s.table.FullName()))
	expr = expr.ConcatBy(" ", s.Assignments)
	expr = expr.ConcatBy(" ", s.where)
	expr = expr.ConcatBy(" ", s.orderBy)
	expr = expr.ConcatBy(" ", s.limit)

	return expr
}
