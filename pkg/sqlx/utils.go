package sqlx

import (
	"database/sql/driver"
	"reflect"
	"regexp"

	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

var queryRegexp = regexp.MustCompile(`(\$\d+)|\?`)

func flattenArgs(query string, args ...interface{}) (finalQuery string, finalArgs []interface{}) {
	index := 0
	finalQuery = queryRegexp.ReplaceAllStringFunc(query, func(i string) string {
		arg := args[index]
		index++

		if canExpr, ok := arg.(builder.CanExpr); ok {
			e := canExpr.Expr()
			resolvedQuery, resolvedArgs := flattenArgs(e.Query, e.Args...)
			finalArgs = append(finalArgs, resolvedArgs...)
			return resolvedQuery
		}

		if _, isValuer := arg.(driver.Valuer); !isValuer {
			if _, isBytes := arg.([]byte); !isBytes {
				if reflect.TypeOf(arg).Kind() == reflect.Slice {
					sliceRv := reflect.ValueOf(arg)
					length := sliceRv.Len()
					for i := 0; i < length; i++ {
						finalArgs = append(finalArgs, sliceRv.Index(i).Interface())
					}
					return builder.HolderRepeat(length)
				}
			}
		}

		finalArgs = append(finalArgs, arg)
		return i
	})
	return
}

func stringIndexOf(slice []string, target string) int {
	for idx, item := range slice {
		if item == target {
			return idx
		}
	}
	return -1
}
