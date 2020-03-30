package builder

import (
	"fmt"
	"strings"
)

type StmtType int16

const (
	STMT_UNKNOWN StmtType = iota
	STMT_INSERT
	STMT_DELETE
	STMT_UPDATE
	STMT_SELECT
	STMT_RAW
)

type Stmt Expression

func (s *Stmt) Type() StmtType {
	return STMT_RAW
}

func (s *Stmt) Expr() *Expression {
	return (*Expression)(s)
}

type Statement interface {
	CanExpr
	Type() StmtType
}

func comment(c string) string {
	if c == "" {
		return ""
	}
	return fmt.Sprintf("/* %s */ ", c)
}

func statement(c string, tpe string, modifiers ...string) string {
	if modifiers == nil {
		return comment(c) + tpe
	}
	return comment(c) + tpe + " " + strings.Join(modifiers, " ")
}

type where Condition

func (w *where) Expr() *Expression {
	if w == nil || w.Query == "" {
		return nil
	}
	return Expr("WHERE "+w.Query, w.Args...)
}

type OrderType string

const (
	ORDER_NO   OrderType = ""
	ORDER_ASC  OrderType = "ASC"
	ORDER_DESC OrderType = "DESC"
)

type order struct {
	expr CanExpr
	tpe  OrderType
}

func (o order) Expr() *Expression {
	if o.expr == nil {
		return nil
	}
	expr := o.expr.Expr()
	if expr == nil {
		return nil
	}
	suffix := ""
	if o.tpe != "" {
		suffix = " " + string(o.tpe)
	}
	return Expr(
		fmt.Sprintf("(%s)%s", expr.Query, suffix),
		expr.Args...,
	)
}

type groupBy struct {
	groups     []order
	withRollup bool
	havingCond *Condition
}

func (g groupBy) setBy(expr CanExpr, orderType OrderType) *groupBy {
	g.groups = []order{{
		expr: expr,
		tpe:  orderType,
	}}
	return &g
}

func (g groupBy) addBy(expr CanExpr, orderType OrderType) *groupBy {
	g.groups = append(g.groups, order{
		expr: expr,
		tpe:  orderType,
	})
	return &g
}

func (g groupBy) rollup() *groupBy {
	g.withRollup = true
	return &g
}

func (g groupBy) having(cond *Condition) *groupBy {
	g.havingCond = cond
	return &g
}

func (g *groupBy) Expr() *Expression {
	if g == nil {
		return nil
	}
	if len(g.groups) > 0 {
		expr := Expr("GROUP BY")
		for i, order := range g.groups {
			if i == 0 {
				expr = expr.ConcatBy(" ", order)
			} else {
				expr = expr.ConcatBy(", ", order)
			}
		}

		if g.withRollup {
			expr = expr.ConcatBy(" ", Expr("WITH ROLLUP"))
		}

		expr = expr.ConcatBy(" HAVING ", g.havingCond)
		return expr
	}
	return nil
}

type orderBy struct {
	orders []order
}

func (o orderBy) setBy(expr CanExpr, orderType OrderType) *orderBy {
	o.orders = []order{{
		expr: expr,
		tpe:  orderType,
	}}
	return &o
}

func (o orderBy) addBy(expr CanExpr, orderType OrderType) *orderBy {
	o.orders = append(o.orders, order{
		expr: expr,
		tpe:  orderType,
	})
	return &o
}

func (o *orderBy) Expr() *Expression {
	if o == nil {
		return nil
	}
	if len(o.orders) > 0 {
		expr := Expr("ORDER BY")
		for i, order := range o.orders {
			if i == 0 {
				expr = expr.ConcatBy(" ", order)
			} else {
				expr = expr.ConcatBy(", ", order)
			}
		}
		return expr
	}
	return nil
}

type limit struct {
	// limit
	rowCount int32
	// offset
	offsetCount int32
}

func (l limit) limit(size int32) *limit {
	l.rowCount = size
	return &l
}

func (l limit) offset(offset int32) *limit {
	l.offsetCount = offset
	return &l
}

func (l *limit) Expr() *Expression {
	if l == nil {
		return nil
	}
	if l.rowCount > 0 {
		if l.offsetCount > 0 {
			return Expr(fmt.Sprintf("LIMIT %d OFFSET %d", l.rowCount, l.offsetCount))
		}
		return Expr(fmt.Sprintf(fmt.Sprintf("LIMIT %d", l.rowCount)))
	}
	return nil
}
