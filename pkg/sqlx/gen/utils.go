package gen

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/google/uuid"

	"github.com/profzone/eden-framework/pkg/codegen"
	"github.com/profzone/eden-framework/pkg/godash"
	"github.com/profzone/eden-framework/pkg/sqlx"
	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

var (
	defRegexp = regexp.MustCompile(`@def ([^\n]+)`)
)

type Keys struct {
	Primary       sqlx.FieldNames
	Indexes       sqlx.Indexes
	UniqueIndexes sqlx.Indexes
}

func (ks *Keys) PatchUniqueIndexesWithSoftDelete(softDeleteField string) {
	if len(ks.UniqueIndexes) > 0 {
		for name, fieldNames := range ks.UniqueIndexes {
			ks.UniqueIndexes[name] = godash.StringUniq(append(fieldNames, softDeleteField))
		}
	}
}

func (ks *Keys) Bind(table *builder.Table) {
	if len(ks.Primary) > 0 {
		cols, err := CheckFields(table, ks.Primary...)
		if err != nil {
			panic(fmt.Errorf("%s, please check primary def", err.Error()))
		}
		ks.Primary = cols.FieldNames()
		table.Keys.Add(builder.PrimaryKey().WithCols(cols.List()...))
	}
	if len(ks.Indexes) > 0 {
		for name, fieldNames := range ks.Indexes {
			cols, err := CheckFields(table, fieldNames...)
			if err != nil {
				panic(fmt.Errorf("%s, please check index def", err.Error()))
			}
			ks.Indexes[name] = cols.FieldNames()
			table.Keys.Add(builder.Index(name).WithCols(cols.List()...))
		}
	}

	if len(ks.UniqueIndexes) > 0 {
		for name, fieldNames := range ks.UniqueIndexes {
			cols, err := CheckFields(table, fieldNames...)
			if err != nil {
				panic(fmt.Errorf("%s, please check unique_index def", err.Error()))
			}
			ks.UniqueIndexes[name] = cols.FieldNames()
			table.Keys.Add(builder.UniqueIndex(name).WithCols(cols.List()...))
		}
	}
}

func CheckFields(table *builder.Table, fieldNames ...string) (cols builder.Columns, err error) {
	for _, fieldName := range fieldNames {
		col := table.F(fieldName)
		if col == nil {
			err = fmt.Errorf("table %s has no field %s", table.Name, fieldName)
			return
		}
		cols.Add(col)
	}
	return
}

func parseKeysFromDoc(doc string) *Keys {
	ks := &Keys{}
	matches := defRegexp.FindAllStringSubmatch(doc, -1)

	for _, subMatch := range matches {
		if len(subMatch) == 2 {
			defs := defSplit(subMatch[1])

			switch strings.ToLower(defs[0]) {
			case "primary":
				if len(defs) < 2 {
					panic(fmt.Errorf("primary at lease 1 Field"))
				}
				ks.Primary = sqlx.FieldNames(defs[1:])
			case "index":
				if len(defs) < 3 {
					panic(fmt.Errorf("index at lease 1 Field"))
				}
				if ks.Indexes == nil {
					ks.Indexes = sqlx.Indexes{}
				}
				ks.Indexes[defs[1]] = sqlx.FieldNames(defs[2:])
			case "unique_index":
				if len(defs) < 3 {
					panic(fmt.Errorf("unique Indexes at lease 1 Field"))
				}
				if ks.UniqueIndexes == nil {
					ks.UniqueIndexes = sqlx.Indexes{}
				}
				ks.UniqueIndexes[defs[1]] = sqlx.FieldNames(defs[2:])
			}
		}
	}
	return ks
}

func defSplit(def string) (defs []string) {
	vs := strings.Split(def, " ")
	for _, s := range vs {
		if s != "" {
			defs = append(defs, s)
		}
	}
	return
}

func toDefaultTableName(name string) string {
	return codegen.ToLowerSnakeCase("t_" + name)
}

func forEachStructField(structType *types.Struct, fn func(fieldVar *types.Var, columnName string, tpe string)) {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)
		if field.Exported() {
			structTag := reflect.StructTag(tag)
			fieldName, exists := structTag.Lookup("db")
			if exists {
				if fieldName != "-" {
					fn(field, fieldName, structTag.Get("sql"))
				}
			} else if field.Anonymous() {
				if nextStructType, ok := field.Type().Underlying().(*types.Struct); ok {
					forEachStructField(nextStructType, fn)
				}
				continue
			}
		}
	}
}

func T() *Template {
	return &Template{}
}

type Template struct {
	tpl     string
	funcMap template.FuncMap
}

func (t Template) Funcs(funcMap template.FuncMap) *Template {
	t.funcMap = funcMap
	return &t
}

func (t Template) Parse(tpl string) *Template {
	t.tpl = tpl
	return &t
}

func (t *Template) Execute(wr io.Writer, data interface{}) {
	tpl, parseErr := template.New(uuid.New().String()).Funcs(t.funcMap).Parse(t.tpl)
	if parseErr != nil {
		panic(fmt.Sprintf("template Prase failded: %s", parseErr.Error()))
	}
	err := tpl.Execute(wr, data)
	if err != nil {
		panic(fmt.Sprintf("template Execute failded: %s", err.Error()))
	}
}

func (t *Template) Render(data interface{}) string {
	buf := new(bytes.Buffer)
	t.Execute(buf, data)
	return buf.String()
}
