package builder

func Count(canExpr CanExpr) *Function {
	return Func("COUNT", canExpr)
}

func Avg(canExpr CanExpr) *Function {
	return Func("AVG", canExpr)
}

func Distinct(canExpr CanExpr) *Function {
	return Func("DISTINCT", canExpr)
}

func Min(canExpr CanExpr) *Function {
	return Func("MIN", canExpr)
}

func Max(canExpr CanExpr) *Function {
	return Func("MAX", canExpr)
}

func First(canExpr CanExpr) *Function {
	return Func("FIRST", canExpr)
}

func Last(canExpr CanExpr) *Function {
	return Func("LAST", canExpr)
}

func Sum(canExpr CanExpr) *Function {
	return Func("SUM", canExpr)
}

func Func(name string, canExpr CanExpr) *Function {
	return &Function{
		Name:    name,
		canExpr: canExpr,
	}
}

type Function struct {
	Name    string
	canExpr CanExpr
}

func (f *Function) Expr() *Expression {
	if f.canExpr != nil {
		e := f.canExpr.Expr()
		if e != nil {
			return Expr(f.Name+"("+e.Query+")", e.Args...)
		}
	}
	return nil
}
