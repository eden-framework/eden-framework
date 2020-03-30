package builder

type star string

func Star() star {
	return star("")
}

func (s star) Expr() *Expression {
	return Expr("*")
}
