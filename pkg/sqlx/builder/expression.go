package builder

func ExprFrom(v interface{}) *Expression {
	switch v.(type) {
	case *Expression:
		return v.(*Expression)
	case CanExpr:
		return v.(CanExpr).Expr()
	}
	return nil
}

func Expr(query string, args ...interface{}) *Expression {
	return &Expression{Query: query, Args: args}
}

func ExprErr(err error) *Expression {
	return &Expression{Err: err}
}

type CanExpr interface {
	Expr() *Expression
}

type Expression struct {
	Err   error
	Query string
	Args  []interface{}
}

func (e *Expression) Expr() *Expression {
	return e
}

func (e Expression) ConcatBy(joiner string, canExpr CanExpr) *Expression {
	expr := canExpr.Expr()
	if expr == nil {
		return &e
	}
	return Expr(e.Query+joiner+expr.Query, append(e.Args, expr.Args...)...)
}
