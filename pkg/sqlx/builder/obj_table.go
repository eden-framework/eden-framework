package builder

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/profzone/eden-framework/pkg/env"
)

type TableDef interface {
	IsValidDef() bool
	Def() *Expression
}

func T(db *Database, tableName string) *Table {
	return &Table{
		DB:   db,
		Name: tableName,
	}
}

type Table struct {
	DB   *Database
	Name string
	Columns
	Keys
	Engine  string
	Charset string
}

func (t Table) Define(defs ...TableDef) *Table {
	for _, def := range defs {
		if def.IsValidDef() {
			switch def.(type) {
			case *Column:
				t.Columns.Add(def.(*Column))
			case *Key:
				t.Keys.Add(def.(*Key))
			}
		}
	}
	return &t
}

var (
	fieldNamePlaceholder = regexp.MustCompile("#[A-Z][A-Za-z0-9_]+")
)

func (t *Table) Ex(query string, args ...interface{}) *Expression {
	finalQuery := fieldNamePlaceholder.ReplaceAllStringFunc(query, func(i string) string {
		fieldName := strings.TrimLeft(i, "#")
		if col := t.F(fieldName); col != nil {
			return col.String()
		}
		return i
	})
	return Expr(finalQuery, args...)
}

func (t *Table) Cond(query string, args ...interface{}) *Condition {
	return (*Condition)(t.Ex(query, args...))
}

type FieldValues map[string]interface{}

func (t *Table) ColumnsAndValuesByFieldValues(fieldValues FieldValues) (columns Columns, args []interface{}) {
	fieldNames := make([]string, 0)
	for fieldName := range fieldValues {
		fieldNames = append(fieldNames, fieldName)
	}

	sort.Strings(fieldNames)

	for _, fieldName := range fieldNames {
		if col := t.F(fieldName); col != nil {
			columns.Add(col)
			args = append(args, fieldValues[fieldName])
		}
	}
	return
}

func (t *Table) AssignsByFieldValues(fieldValues FieldValues) (assignments Assignments) {
	for fieldName, value := range fieldValues {
		col := t.F(fieldName)
		if col != nil {
			assignments = append(assignments, col.By(value))
		}
	}
	return
}

func (t *Table) String() string {
	return quote(t.Name)
}

func (t *Table) FullName() string {
	return t.DB.String() + "." + t.String()
}

func (t *Table) Insert() *StmtInsert {
	return Insert(t)
}

func (t *Table) Delete() *StmtDelete {
	return Delete(t)
}

func (t *Table) Select() *StmtSelect {
	return SelectFrom(t)
}

func (t *Table) Update() *StmtUpdate {
	return Update(t)
}

func (t *Table) Drop() *Stmt {
	return (*Stmt)(Expr(fmt.Sprintf("DROP TABLE %s", t.FullName())))
}

func (t *Table) Truncate() *Stmt {
	return (*Stmt)(Expr(fmt.Sprintf("TRUNCATE TABLE %s", t.FullName())))
}

func (t *Table) Diff(table *Table) *Stmt {
	colsDiffResult := t.Columns.Diff(table.Columns)
	keysDiffResult := t.Keys.Diff(table.Keys)

	colsChanged := colsDiffResult.IsChanged()
	indexesChanged := keysDiffResult.IsChanged()

	if !colsChanged && !indexesChanged {
		return nil
	}
	expr := Expr(fmt.Sprintf(`ALTER TABLE %s `, t.FullName()))

	joiner := ""

	if colsChanged {
		if Configuration.DropColumnWhenMigration || env.IsOnline() {
			colsDiffResult.colsForDelete.Range(func(col *Column, idx int) {
				expr = expr.ConcatBy(joiner, col.Drop())
				joiner = ", "
			})
		}
		colsDiffResult.colsForUpdate.Range(func(col *Column, idx int) {
			expr = expr.ConcatBy(joiner, col.Modify())
			joiner = ", "
		})
		colsDiffResult.colsForAdd.Range(func(col *Column, idx int) {
			expr = expr.ConcatBy(joiner, col.Add())
			joiner = ", "
		})
	}

	if indexesChanged {
		keysDiffResult.keysForDelete.Range(func(key *Key, idx int) {
			expr = expr.ConcatBy(joiner, key.Drop())
			joiner = ", "
		})
		keysDiffResult.keysForUpdate.Range(func(key *Key, idx int) {
			expr = expr.ConcatBy(joiner, key.Drop())
			joiner = ", "
			expr = expr.ConcatBy(joiner, key.Add())
		})
		keysDiffResult.keysForAdd.Range(func(key *Key, idx int) {
			expr = expr.ConcatBy(joiner, key.Add())
			joiner = ", "
		})
	}

	return (*Stmt)(expr)
}

func (t *Table) Create(ifNotExists bool) *Stmt {
	expr := Expr("CREATE TABLE")
	if ifNotExists {
		expr = expr.ConcatBy(" ", Expr("IF NOT EXISTS"))
	}
	expr.Query = expr.Query + fmt.Sprintf(" %s (", t.FullName())

	if !t.Columns.IsEmpty() {
		isFirstCol := true

		t.Columns.Range(func(col *Column, idx int) {
			joiner := ", "
			if isFirstCol {
				joiner = ""
			}
			def := col.Def()
			if def != nil {
				isFirstCol = false
				expr = expr.ConcatBy(joiner, col.Def())
			}
		})

		t.Keys.Range(func(key *Key, idx int) {
			expr = expr.ConcatBy(", ", key.Def())
		})
	}

	engine := t.Engine
	if engine == "" {
		engine = "InnoDB"
	}

	charset := t.Charset
	if charset == "" {
		charset = "utf8"
	}

	expr.Query = fmt.Sprintf("%s) ENGINE=%s CHARSET=%s", expr.Query, engine, charset)
	return (*Stmt)(expr)
}

type Tables map[string]*Table

func (tables Tables) TableNames() (names []string) {
	for name := range tables {
		names = append(names, name)
	}
	return
}

func (tables Tables) Add(table *Table) {
	tables[table.Name] = table
}
