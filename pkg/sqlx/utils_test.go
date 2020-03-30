package sqlx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

func TestFlattenArgs(t *testing.T) {
	tt := assert.New(t)

	{
		q, args := flattenArgs(`#ID IN (?)`, []int{28, 29, 30})
		tt.Equal("#ID IN (?,?,?)", q)
		tt.Equal(args, []interface{}{28, 29, 30})
	}
	{
		q, args := flattenArgs(`#ID = (?)`, []byte(""))
		tt.Equal("#ID = (?)", q)
		tt.Equal(args, []interface{}{[]byte("")})
	}

	{
		q, args := flattenArgs(`#ID = ?`, builder.Expr("#ID + ?", 1))
		tt.Equal("#ID = #ID + ?", q)
		tt.Equal(args, []interface{}{1})
	}
}

func TestStringIndexOf(t *testing.T) {
	tt := assert.New(t)

	tt.Equal(-1, stringIndexOf([]string{"1", "2"}, "3"))
	tt.Equal(1, stringIndexOf([]string{"1", "2"}, "2"))
}
