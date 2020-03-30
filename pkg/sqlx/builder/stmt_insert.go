package builder

import (
	"fmt"
)

var (
	InsertValuesLengthNotMatch = fmt.Errorf("value length is not equal Col length")
)

func Insert(table *Table) *StmtInsert {
	return &StmtInsert{
		table: table,
	}
}

type StmtInsert struct {
	table     *Table
	modifiers []string
	columns   Columns
	Assignments
	values  [][]interface{}
	comment string
}

func (s StmtInsert) Comment(comment string) *StmtInsert {
	s.comment = comment
	return &s
}

func (s *StmtInsert) Type() StmtType {
	return STMT_INSERT
}

func (s StmtInsert) Modifier(modifiers ...string) *StmtInsert {
	s.modifiers = append(s.modifiers, modifiers...)
	return &s
}

func (s StmtInsert) Columns(cols Columns) *StmtInsert {
	s.columns = cols
	return &s
}

func (s StmtInsert) Values(values ...interface{}) *StmtInsert {
	if s.values == nil {
		s.values = [][]interface{}{}
	}
	s.values = append(s.values, values)
	return &s
}

func (s StmtInsert) OnDuplicateKeyUpdate(assigns ...*Assignment) *StmtInsert {
	s.Assignments = assigns
	return &s
}

func (s *StmtInsert) Expr() *Expression {
	expr := Expr(fmt.Sprintf(
		`%s INTO %s`,
		statement(s.comment, "INSERT", s.modifiers...),
		s.table.FullName(),
	))

	expr = expr.ConcatBy(" ", s.columns.Wrap())

	for idx, vals := range s.values {
		if s.columns.Len() != len(vals) {
			return ExprErr(InsertValuesLengthNotMatch)
		}
		joiner := " VALUES "
		if idx > 0 {
			joiner = ","
		}
		expr = expr.ConcatBy(joiner, Expr(
			"("+HolderRepeat(len(vals))+")",
			vals...,
		))
	}

	if len(s.Assignments) > 0 {
		expr = expr.ConcatBy(" ON DUPLICATE KEY UPDATE ", s.Assignments)
	}

	return expr
}
