package gen

import (
	"fmt"
	"go/types"
	"io"
	"strings"
	"text/template"

	"golang.org/x/tools/go/loader"

	"github.com/profzone/eden-framework/pkg/codegen"
	"github.com/profzone/eden-framework/pkg/codegen/loaderx"
	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

func NewModel(prog *loader.Program, typeName *types.TypeName, comments string, cfg *Config) *Model {
	m := Model{}
	m.Config = cfg
	m.TypeName = typeName

	m.Table = builder.T(nil, cfg.TableName)

	p := prog.Package(typeName.Pkg().Path())

	forEachStructField(typeName.Type().Underlying().Underlying().(*types.Struct), func(structVal *types.Var, columnName string, tpe string) {
		col := builder.Col(m.Table, columnName).Field(structVal.Name()).Type(tpe)

		for id, o := range p.Defs {
			if o == structVal {
				doc := loaderx.CommentsOf(prog.Fset, id, p.Files...)
				col.Comment = strings.Split(doc, "\n")[0]
			}
		}

		m.AddColumn(col, structVal)
	})

	m.HasSoftDelete = m.Table.F(m.FieldSoftDelete) != nil
	m.HasCreatedAt = m.Table.F(m.FieldCreatedAt) != nil
	m.HasUpdatedAt = m.Table.F(m.FieldUpdatedAt) != nil

	m.Keys = parseKeysFromDoc(comments)
	if m.HasSoftDelete {
		m.Keys.PatchUniqueIndexesWithSoftDelete(m.FieldSoftDelete)
	}
	m.Keys.Bind(m.Table)

	if autoIncrementCol := m.Table.AutoIncrement(); autoIncrementCol != nil {
		m.HasAutoIncrement = true
		m.FieldAutoIncrement = autoIncrementCol.FieldName
	}

	m.Importer = &codegen.Importer{}

	m.Template = T().Funcs(template.FuncMap{
		"use":  m.Importer.Use,
		"dump": m.Importer.Sdump,
		"var":  codegen.ToLowerCamelCase,
	})

	return &m
}

type Model struct {
	*types.TypeName
	*codegen.Importer
	*Config
	*Template
	*Keys
	*builder.Table
	Fields             map[string]*types.Var
	FieldAutoIncrement string
	HasSoftDelete      bool
	HasCreatedAt       bool
	HasUpdatedAt       bool
	HasAutoIncrement   bool
}

func (m *Model) AddColumn(col *builder.Column, tpe *types.Var) {
	m.Table.Columns.Add(col)
	if m.Fields == nil {
		m.Fields = map[string]*types.Var{}
	}
	m.Fields[col.FieldName] = tpe
}

func (m *Model) Render() string {
	blocks := strings.Join(
		[]string{
			m.dataAndTable(),
			m.interfaces(),
			m.methodsForCRUD(),
			m.methodsForList(),
		},
		"\n",
	)

	return fmt.Sprintf(`
	package %s

	%s

	%s
	`,
		m.TypeName.Pkg().Name(),
		m.Importer.String(),
		blocks,
	)
}

func (m *Model) ParseTo(writer io.Writer, tpl string) {
	m.Template.Parse(tpl).Execute(writer, m)
}
