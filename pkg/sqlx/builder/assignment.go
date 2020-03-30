package builder

type Assignment Expression

func (a Assignment) Expr() *Expression {
	return (*Expression)(&a)
}

type Assignments []*Assignment

func (assigns Assignments) Expr() (e *Expression) {
	e = Expr("")
	for i, assignment := range assigns {
		joiner := ", "
		if i == 0 {
			joiner = ""
		}
		e = e.ConcatBy(joiner, assignment)
	}
	return e
}
