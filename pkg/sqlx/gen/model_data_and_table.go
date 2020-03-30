package gen

import (
	"bytes"
	"sort"

	"github.com/profzone/eden-framework/pkg/godash"
	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

func (m *Model) dataAndTable() string {
	buf := &bytes.Buffer{}

	if m.WithTableInterfaces {
		m.ParseTo(buf, `
var {{ .StructName }}Table *{{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Table

func init() {
	{{ .StructName }}Table = {{ .Database }}.Register(&{{ .StructName }}{})
}

func ({{ var .StructName }} *{{ .StructName }}) D() *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.Database {
	return {{ .Database }}
}


func ({{ var .StructName }} *{{ .StructName }}) T() *{{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Table {
	return {{ .StructName }}Table
}

func ({{ var .StructName }} *{{ .StructName }}) TableName() string {
	return "{{ .TableName }}"
}
`)
	}

	m.ParseTo(buf, `
	{{ $structName := .StructName }}

	type {{ .StructName }}Fields struct {
		{{ range $k, $field := ( .FieldNames ) }}{{ print $field }} *{{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Column
		{{ end }}
	}

	var {{ $structName }}Field = struct {
		{{ range $k, $field := ( .FieldNames ) }}{{ print $field }} string
		{{ end }}
	}{
		{{ range $k, $field := ( .FieldNames ) }}{{ print $field }}: "{{ print $field }}",
		{{ end }}
	}

	func ({{ var .StructName }} *{{ .StructName }}) Fields() *{{ .StructName }}Fields {
		table := {{ var .StructName }}.T()

		return &{{ .StructName }}Fields{
			{{ range $k, $field := ( .FieldNames ) }}{{ print $field }}: table.F({{ $structName }}Field.{{ print $field }}),
			{{ end }}
		}
	}

	func ({{ var .StructName }} *{{ .StructName }}) IndexFieldNames() []string {
		return {{ dump .IndexFieldNames }}
	}

	func ({{ var .StructName }} *{{ .StructName }}) ConditionByStruct() *{{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Condition  {
		table := {{ var .StructName }}.T()

		fieldValues := {{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.FieldValuesFromStructByNonZero({{ var .StructName }})

		conditions := []*{{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.Condition{}

		for _, fieldName := range {{ var .StructName }}.IndexFieldNames() {
			if v, exists := fieldValues[fieldName]; exists {
				conditions = append(conditions, table.F(fieldName).Eq(v))
				delete(fieldValues, fieldName)
			}
		}

		if len(conditions) == 0 {
			panic(fmt.Errorf("at least one of field for indexes has value"))
		}

		for fieldName, v := range fieldValues {
			conditions = append(conditions, table.F(fieldName).Eq(v))
		}

		condition := {{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(conditions...)

		{{ if .HasSoftDelete }}
			condition = {{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(condition, table.F("{{ .FieldSoftDelete }}").Eq({{ use .ConstSoftDeleteTrue }}))
		{{ end }}
		return condition
	}
	`)

	return buf.String()
}

func (m *Model) IndexFieldNames() []string {
	indexedFields := []string{}

	m.Table.Keys.Range(func(key *builder.Key, idx int) {
		fieldNames := key.Columns.FieldNames()
		indexedFields = append(indexedFields, fieldNames...)
	})

	indexedFields = godash.StringUniq(indexedFields)

	indexedFields = godash.StringFilter(indexedFields, func(item string, i int) bool {
		if m.HasSoftDelete {
			return item != m.FieldSoftDelete
		}
		return true
	})

	sort.Strings(indexedFields)
	return indexedFields
}
