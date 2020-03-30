package gen

import (
	"bytes"
	"fmt"
	"go/types"
	"strings"
)

func (m *Model) methodsForList() string {
	return strings.Join([]string{
		m.methodsForFetchList(),
		m.methodsForBatchList(),
	}, "\n\n")
}

func (m *Model) methodsForFetchList() string {
	buf := &bytes.Buffer{}

	m.ParseTo(buf, fmt.Sprintf(
		`
{{ $method := "FetchList" }}
type {{ .StructName }}List []{{ .StructName }}

// deprecated
func ({{ var .StructName }}List *{{ .StructName }}List) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB, size int32, offset int32, conditions ...*{{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Condition) (count int32, err error)	{
	*{{ var .StructName }}List, count, err = (&{{ .StructName }}{}).FetchList(db, size, offset, conditions...)
	return
}

func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB, size int32, offset int32, conditions ...*{{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Condition) ({{ var .StructName }}List {{ .StructName }}List, count int32, err error) {
	{{ var .StructName }}List = {{ .StructName }}List{}

	table := {{ var .StructName }}.T()

	condition := {{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(conditions...)
	{{ if .HasSoftDelete }}
		condition = {{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(condition, table.F("{{ .FieldSoftDelete }}").Eq({{ use .ConstSoftDeleteTrue }}))
	{{ end }}

	stmt := table.Select().
		Comment("{{ .StructName }}.{{ $method }}").
		Where(condition)

	errForCount := db.Do(stmt.For({{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Count({{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Star()))).Scan(&count).Err()
	if errForCount != nil {
		err = errForCount
		return
	}

	stmt = stmt.Limit(size).Offset(offset)
	{{ if .HasCreatedAt }}
		stmt = stmt.OrderDescBy(table.F("{{ .FieldCreatedAt }}"))
	{{ end }}

	err = db.Do(stmt).Scan(&{{ var .StructName }}List).Err()

	return
}
`))

	m.ParseTo(buf, fmt.Sprintf(
		`
{{ $method := "List" }}

func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB, condition *{{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Condition) ({{ var .StructName }}List {{ .StructName }}List, err error) {
	{{ var .StructName }}List = {{ .StructName }}List{}

	table := {{ var .StructName }}.T()

	{{ if .HasSoftDelete }}
		condition = {{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(condition, table.F("{{ .FieldSoftDelete }}").Eq({{ use .ConstSoftDeleteTrue }}))
	{{ end }}

	stmt := table.Select().
		Comment("{{ .StructName }}.{{ $method }}").
		Where(condition)

	err = db.Do(stmt).Scan(&{{ var .StructName }}List).Err()

	return
}
`))

	m.ParseTo(buf, fmt.Sprintf(
		`
{{ $method := "ListByStruct" }}

func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB) ({{ var .StructName }}List {{ .StructName }}List, err error) {
	{{ var .StructName }}List = {{ .StructName }}List{}

	table := {{ var .StructName }}.T()

	condition := {{ var .StructName }}.ConditionByStruct()

	{{ if .HasSoftDelete }}
		condition = {{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(condition, table.F("{{ .FieldSoftDelete }}").Eq({{ use .ConstSoftDeleteTrue }}))
	{{ end }}

	stmt := table.Select().
		Comment("{{ .StructName }}.{{ $method }}").
		Where(condition)

	err = db.Do(stmt).Scan(&{{ var .StructName }}List).Err()

	return
}
`))

	return buf.String()
}

func ForEachField(typeStruct *types.Struct, cb func(field *types.Var)) {
	for i := 0; i < typeStruct.NumFields(); i++ {
		f := typeStruct.Field(i)
		if f.Anonymous() {
			if s, ok := f.Type().Underlying().(*types.Struct); ok {
				ForEachField(s, cb)
			} else {
				cb(f)
			}
		} else {
			cb(f)
		}
	}
}

func (m *Model) FieldType(name string) (typ string) {
	typeStruct := m.TypeName.Type().Underlying().(*types.Struct)

	ForEachField(typeStruct, func(field *types.Var) {
		if field.Name() == name {
			typ = field.Type().String()
			if strings.Contains(typ, ".") {
				typ = m.Use(field.Type().String())
			}
		}
	})

	return
}

func (m *Model) methodsForBatchList() string {
	buf := &bytes.Buffer{}

	indexedFields := m.IndexFieldNames()

	for _, field := range indexedFields {
		m.ParseTo(buf, fmt.Sprintf(`
{{ $field := "%s" }}
{{ $fieldType := "%s" }}

// deprecated
func ({{ var .StructName }}List *{{ .StructName }}List) BatchFetchBy{{ $field }}List(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB, {{ var $field }}List []{{ $fieldType }}) (err error)	{
	*{{ var .StructName }}List, err = (&{{ .StructName }}{}).BatchFetchBy{{ $field }}List(db, {{ var $field }}List)
	return
}

func ({{ var .StructName }} *{{ .StructName }}) BatchFetchBy{{ $field }}List(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB, {{ var $field }}List []{{ $fieldType }}) ({{ var .StructName }}List {{ .StructName }}List, err error) {
	if len({{ var $field }}List) == 0 {
		return {{ .StructName }}List{}, nil
	}

	table :=  {{ var .StructName }}.T()

	condition := table.F("{{ $field }}").In({{ var $field }}List)

	{{ if .HasSoftDelete }}
		condition = condition.And(table.F("{{ .FieldSoftDelete }}").Eq({{ use .ConstSoftDeleteTrue }}))
	{{ end }}

	stmt := table.Select().
		Comment("{{ .StructName }}.BatchFetchBy{{ $field }}List").
		Where(condition)

	err = db.Do(stmt).Scan(&{{ var .StructName }}List).Err()

	return
}
	`,
			field,
			m.FieldType(field),
		))
	}

	return buf.String()
}
