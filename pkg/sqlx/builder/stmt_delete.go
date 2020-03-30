package builder

import (
	"fmt"
)

var (
	DeleteNeedLimitByWhere = fmt.Errorf("no where limit for delete")
)

func Delete(table *Table) *StmtDelete {
	return &StmtDelete{
		table: table,
	}
}

type StmtDelete struct {
	table     *Table
	modifiers []string
	*where
	*orderBy
	*limit
	comment string
}

func (s StmtDelete) Comment(comment string) *StmtDelete {
	s.comment = comment
	return &s
}

func (s *StmtDelete) Type() StmtType {
	return STMT_DELETE
}

func (s StmtDelete) Modifier(modifiers ...string) *StmtDelete {
	s.modifiers = append(s.modifiers, modifiers...)
	return &s
}

func (s StmtDelete) Where(cond *Condition) *StmtDelete {
	s.where = (*where)(cond)
	return &s
}

func (s StmtDelete) OrderBy(col *Column, orderType OrderType) *StmtDelete {
	if s.orderBy == nil {
		s.orderBy = &orderBy{}
	}
	s.orderBy = s.orderBy.setBy(col, orderType)
	return &s
}

func (s StmtDelete) AscendBy(col *Column) *StmtDelete {
	return s.OrderBy(col, ORDER_ASC)
}

func (s StmtDelete) DescendBy(col *Column) *StmtDelete {
	return s.OrderBy(col, ORDER_DESC)
}

func (s StmtDelete) Limit(size int32) *StmtDelete {
	if s.limit == nil {
		s.limit = &limit{}
	}
	s.limit = s.limit.limit(size)
	return &s
}

func (s *StmtDelete) Expr() *Expression {
	if s.where == nil {
		return ExprErr(DeleteNeedLimitByWhere)
	}

	expr := Expr(fmt.Sprintf("%s FROM %s", statement(s.comment, "DELETE", s.modifiers...), s.table.FullName()))

	expr = expr.ConcatBy(" ", s.where)
	expr = expr.ConcatBy(" ", s.orderBy)
	expr = expr.ConcatBy(" ", s.limit)

	return expr
}
