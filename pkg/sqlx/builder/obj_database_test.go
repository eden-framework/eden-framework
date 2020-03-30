package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	tt := assert.New(t)

	db := DB("db")

	tt.Nil(db.Table("t"))

	table := T(db, "t")
	db.Register(table)
	tt.NotNil(db.Table("t"))
	tt.Equal([]string{"t"}, db.TableNames())

	exprCases{
		Case(
			"drop Database",
			db.Drop(),
			Expr("DROP DATABASE `db`"),
		),
		Case(
			"create Database if not exists",
			db.Create(true),
			Expr("CREATE DATABASE IF NOT EXISTS `db`"),
		),
		Case(
			"create Database",
			db.Create(false),
			Expr("CREATE DATABASE `db`"),
		),
	}.Run(t, "Database")
}
