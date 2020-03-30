package sqlx

import (
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

func Do(db *DB, stmt builder.Statement) (result *Result) {
	result = &Result{}

	e := stmt.Expr()
	if e == nil {
		result.err = NewSqlError(sqlErrTypeInvalidSql, "")
		return
	}
	if e.Err != nil {
		result.err = NewSqlError(sqlErrTypeInvalidSql, e.Err.Error())
		logrus.Errorf("%s", result.err)
		return
	}

	result.stmtType = stmt.Type()

	switch result.stmtType {
	case builder.STMT_SELECT:
		rows, queryErr := db.Query(e.Query, e.Args...)
		if queryErr != nil {
			result.err = queryErr
			return
		}
		result.Rows = rows
	case builder.STMT_INSERT, builder.STMT_UPDATE:
		sqlResult, execErr := db.Exec(e.Query, e.Args...)
		if execErr != nil {
			if mysqlErr, ok := execErr.(*mysql.MySQLError); ok && mysqlErr.Number == DuplicateEntryErrNumber {
				result.err = NewSqlError(sqlErrTypeConflict, mysqlErr.Error())
			} else {
				result.err = execErr
			}
			return
		}
		result.Result = sqlResult
	case builder.STMT_DELETE, builder.STMT_RAW:
		sqlResult, execErr := db.Exec(e.Query, e.Args...)
		if execErr != nil {
			result.err = execErr
			return
		}
		result.Result = sqlResult
	}
	return
}

type Result struct {
	stmtType builder.StmtType
	err      error
	*sql.Rows
	sql.Result
}

func (r *Result) Err() error {
	return r.err
}

func (r *Result) Scan(v interface{}) *Result {
	if r.err != nil {
		return r
	}

	if r.Rows != nil {
		defer r.Rows.Close()

		if scanner, ok := v.(sql.Scanner); ok {
			for r.Rows.Next() {
				if scanErr := r.Rows.Scan(scanner); scanErr != nil {
					r.err = scanErr
					return r
				}
			}
		} else {

			modelType := reflect.TypeOf(v)
			if modelType.Kind() != reflect.Ptr {
				r.err = NewSqlError(sqlErrTypeInvalidScanTarget, "can not scan to a none pointer variable")
				return r
			}

			modelType = modelType.Elem()

			isSlice := false
			if modelType.Kind() == reflect.Slice {
				modelType = modelType.Elem()
				isSlice = true
			}

			if modelType.Kind() == reflect.Struct || isSlice {
				columns, getErr := r.Rows.Columns()
				if getErr != nil {
					r.err = getErr
					return r
				}

				rv := reflect.Indirect(reflect.ValueOf(v))

				rowLength := 0

				for r.Rows.Next() {
					if !isSlice && rowLength > 1 {
						r.err = NewSqlError(sqlErrTypeSelectShouldOne, "more than one records found, but only one")
						return r
					}

					rowLength++
					length := len(columns)
					dest := make([]interface{}, length)
					itemRv := rv

					if isSlice {
						itemRv = reflect.New(modelType).Elem()
					}

					destIndexes := make(map[int]bool, length)

					ForEachStructFieldValue(itemRv, func(structFieldValue reflect.Value, structField reflect.StructField, columnName string) {
						idx := stringIndexOf(columns, columnName)
						if idx >= 0 {
							dest[idx] = structFieldValue.Addr().Interface()
							destIndexes[idx] = true
						}
					})

					for index := range dest {
						if !destIndexes[index] {
							placeholder := emptyScanner(0)
							dest[index] = &placeholder
						} else {
							// todo null ignore
							dest[index] = newNullableScanner(dest[index])
						}
					}

					if scanErr := r.Rows.Scan(dest...); scanErr != nil {
						r.err = scanErr
						return r
					}

					if isSlice {
						rv.Set(reflect.Append(rv, itemRv))
					}
				}

				if !isSlice && rowLength == 0 {
					r.err = NewSqlError(sqlErrTypeNotFound, "record is not found")
					return r
				}
			} else {
				for r.Rows.Next() {
					if scanErr := r.Rows.Scan(v); scanErr != nil {
						r.err = scanErr
						return r
					}
				}
			}
		}
		if err := r.Rows.Err(); err != nil {
			r.err = err
			return r
		}

		// Make sure the query can be processed to completion with no errors.
		if err := r.Rows.Close(); err != nil {
			r.err = err
			return r
		}
	}

	return r
}

type emptyScanner int

var _ interface {
	sql.Scanner
	driver.Valuer
} = (*emptyScanner)(nil)

func (e *emptyScanner) Scan(value interface{}) error {
	return nil
}

func (e emptyScanner) Value() (driver.Value, error) {
	return 0, nil
}

func newNullableScanner(dest interface{}) *nullableScanner {
	return &nullableScanner{
		dest: dest,
	}
}

type nullableScanner struct {
	dest interface{}
}

var _ interface {
	sql.Scanner
} = (*nullableScanner)(nil)

func (scanner *nullableScanner) Scan(src interface{}) error {
	if scanner, ok := scanner.dest.(sql.Scanner); ok {
		return scanner.Scan(src)
	}
	if src == nil {
		if zeroSetter, ok := scanner.dest.(ZeroSetter); ok {
			zeroSetter.SetToZero()
			return nil
		}
		return nil
	}
	return convertAssign(scanner.dest, src)
}
