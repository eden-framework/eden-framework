package generator

import (
	"os"
	"testing"
)

func TestSqlFuncGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	for _, name := range []string{"User", "Org"} {
		g := NewSqlFuncGenerator()
		g.WithComments = true
		g.WithTableName = true
		g.WithTableInterfaces = true
		g.WithMethods = true
		g.Database = "DBTest"
		g.StructName = name

		g.Load(cwd)
		g.Pick()
		g.Output(cwd)
	}
}
