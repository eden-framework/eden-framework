package gen

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/profzone/eden-framework/pkg/godash"
	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

func (m *Model) methodsForCRUD() string {
	return strings.Join([]string{
		m.methodsForBasic(),
		m.methodsForKeys(),
	}, "\n\n")
}

func (m *Model) EnableIfNeed() string {
	buf := &bytes.Buffer{}

	if m.HasSoftDelete {
		m.ParseTo(buf, `{{ var .StructName }}.{{ .FieldSoftDelete }} = {{ use .ConstSoftDeleteTrue }}
		`)
	}

	return buf.String()
}

func (m *Model) TypeTime(fieldName string) string {
	if field, ok := m.Fields[fieldName]; ok {
		return field.Type().String()
	}
	return ""
}

func (m *Model) SetCreatedAtIfNeed() string {
	buf := &bytes.Buffer{}

	if m.HasCreatedAt {
		m.ParseTo(buf, `if {{ var .StructName }}.{{ .FieldCreatedAt }}.IsZero() {
			{{ var .StructName }}.{{ .FieldCreatedAt }} = {{ use ( .TypeTime .FieldCreatedAt ) }}({{ use "time" }}.Now())
		}
		`)

		if m.HasUpdatedAt {
			m.ParseTo(buf, `{{ var .StructName }}.{{ .FieldUpdatedAt }} = {{ var .StructName }}.{{ .FieldCreatedAt }}
			`)
		}
	}

	return buf.String()
}

func (m *Model) SetUpdatedForFieldValuesAtIfNeed() string {
	buf := &bytes.Buffer{}

	if m.HasCreatedAt {
		if m.HasUpdatedAt {
			m.ParseTo(buf, `
			if _, ok := fieldValues["{{ .FieldUpdatedAt }}"]; !ok {
				fieldValues["{{ .FieldUpdatedAt }}"] = {{ use ( .TypeTime .FieldUpdatedAt ) }}({{ use "time" }}.Now())
			}
			`)
		}
	}

	return buf.String()
}

func (m *Model) SetEnabledForFieldValuesAtIfNeed() string {
	buf := &bytes.Buffer{}

	if m.HasSoftDelete {
		m.ParseTo(buf, `
		if _, ok := fieldValues["{{ .FieldSoftDelete }}"]; !ok {
			fieldValues["{{ .FieldSoftDelete }}"] = {{ use .ConstSoftDeleteTrue }}
		}
		`)
	}

	return buf.String()
}

func (m *Model) SetUpdatedAtIfNeed() string {
	buf := &bytes.Buffer{}

	if m.HasCreatedAt {
		if m.HasUpdatedAt {
			m.ParseTo(buf, `{{ var .StructName }}.{{ .FieldUpdatedAt }} = {{ use ( .TypeTime .FieldUpdatedAt ) }}({{ use "time" }}.Now())
			`)
		}
	}

	return buf.String()
}

func (m *Model) methodsForBasic() string {
	buf := &bytes.Buffer{}

	m.ParseTo(buf, `
	{{ $method := "Create"}}
	func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB) error {
	{{ ( .EnableIfNeed ) }}
	{{ ( .SetCreatedAtIfNeed ) }}

	stmt := {{ var .StructName }}.D().
		Insert({{ var .StructName }}).
		Comment("{{ .StructName }}.{{ $method }}")

	dbRet := db.Do(stmt)
	err := dbRet.Err()

	{{ if .HasAutoIncrement }}
		if err == nil {
			lastInsertID, _ := dbRet.LastInsertId()
			{{ var .StructName }}.{{ .FieldAutoIncrement }} = {{ .FieldType .FieldAutoIncrement }}(lastInsertID)
		}
	{{ end }}

	return err
	}
	`)

	m.ParseTo(buf, fmt.Sprintf(
		`
{{ $method := "DeleteByStruct" }}
func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB) (err error) {
	table :=  {{ var .StructName }}.T()

	stmt := table.Delete().
		Comment("{{ .StructName }}.{{ $method }}").
		Where({{ var .StructName }}.ConditionByStruct())

	err = db.Do(stmt).Err()
	return
}
`))

	if len(m.Keys.UniqueIndexes) > 0 {
		m.ParseTo(buf, `
{{ $method := "CreateOnDuplicateWithUpdateFields"}}
func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB, updateFields []string) error {
	if len(updateFields) == 0 {
		panic({{ use "fmt"}}.Errorf("must have update fields"))
	}

	{{ ( .EnableIfNeed ) }}
	{{ ( .SetCreatedAtIfNeed ) }}

	table := {{ var .StructName }}.T()

	fieldValues := {{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.FieldValuesFromStructByNonZero({{ var .StructName }}, updateFields...)

	{{ if .HasAutoIncrement }}
		delete(fieldValues, "{{ .FieldAutoIncrement }}")
	{{ end }}

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	m := make(map[string]bool, len(updateFields))
	for _, field := range updateFields {
		m[field] = true
	}

	// fields of unique index can not update
	{{ if .HasCreatedAt }} delete(m, "{{ .FieldCreatedAt }}")
	{{ end }}

	for _, fieldNames := range {{ var .StructName }}.UniqueIndexes() {
		for _, field := range fieldNames {
			delete(m, field)
		}
	}

	if len(m) == 0 {
		panic(fmt.Errorf("no fields for updates"))
	}

	for field := range fieldValues {
		if !m[field] {
			delete(fieldValues, field)
		}
	}

	stmt := table.
		Insert().Columns(cols).Values(vals...).
		OnDuplicateKeyUpdate(table.AssignsByFieldValues(fieldValues)...).
		Comment("{{ .StructName }}.{{ $method }}")

	return db.Do(stmt).Err()
}
	`)
	}

	return buf.String()
}

func toExactlyConditionFrom(fieldNames ...string) string {
	buf := &bytes.Buffer{}
	for _, fieldName := range fieldNames {
		buf.WriteString(fmt.Sprintf(`table.F("%s").Eq({{ var .StructName }}.%s),
		`, fieldName, fieldName))
	}
	return buf.String()
}

func createMethod(method string, fieldNames ...string) string {
	return fmt.Sprintf(method, strings.Join(fieldNames, "And"))
}

func (m *Model) methodsForKeys() string {
	buf := &bytes.Buffer{}

	m.Table.Keys.Range(func(key *builder.Key, idx int) {
		fieldNames := key.Columns.FieldNames()
		fieldNamesWithoutEnabled := godash.StringFilter(fieldNames, func(item string, i int) bool {
			if m.HasSoftDelete {
				return item != m.FieldSoftDelete
			}
			return true
		})
		if m.HasSoftDelete && key.Type == builder.PRIMARY {
			fieldNames = append(fieldNames, m.FieldSoftDelete)
		}
		if key.Type == builder.PRIMARY || key.Type == builder.UNIQUE_INDEX {

			m.ParseTo(buf, fmt.Sprintf(`
{{ $method := "%s" }}
func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB) error {
	{{ ( .EnableIfNeed ) }}

	table :=  {{ var .StructName }}.T()
	stmt := table.Select().
		Comment("{{ .StructName }}.{{ $method }}").
		Where({{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(
			%s
		))

	return db.Do(stmt).Scan({{ var .StructName }}).Err()
}
	`,
				createMethod("FetchBy%s", fieldNamesWithoutEnabled...),
				toExactlyConditionFrom(fieldNames...),
			))

			m.ParseTo(buf, fmt.Sprintf(`
{{ $method := "%s" }}
func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB) error {
	{{ ( .EnableIfNeed ) }}

	table :=  {{ var .StructName }}.T()
	stmt := table.Select().
		Comment("{{ .StructName }}.{{ $method }}").
		Where({{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(
			%s
		)).
		ForUpdate()

	return db.Do(stmt).Scan({{ var .StructName }}).Err()
}
					`,
				createMethod("FetchBy%sForUpdate", fieldNamesWithoutEnabled...),
				toExactlyConditionFrom(fieldNames...),
			))

			m.ParseTo(buf, fmt.Sprintf(`
{{ $method := "%s" }}
func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB) error {
	{{ ( .EnableIfNeed ) }}

	table :=  {{ var .StructName }}.T()
	stmt := table.Delete().
		Comment("{{ .StructName }}.{{ $method }}").
		Where({{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(
			%s
		))

	return db.Do(stmt).Scan({{ var .StructName }}).Err()
}
					`,
				createMethod("DeleteBy%s", fieldNamesWithoutEnabled...),
				toExactlyConditionFrom(fieldNames...),
			))

			m.ParseTo(buf, fmt.Sprintf(`
{{ $method := "%s" }}
{{ $methodForFetch := "%s" }}
func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB, fieldValues {{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.FieldValues) error {
	{{ ( .SetUpdatedForFieldValuesAtIfNeed ) }}
	{{ ( .EnableIfNeed ) }}

	table := {{ var .StructName }}.T()

	{{ if .HasAutoIncrement }}
		delete(fieldValues, "{{ .FieldAutoIncrement }}")
	{{ end }}

	stmt := table.Update().
		Comment("{{ .StructName }}.{{ $method }}").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where({{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(
			%s
		))

	dbRet := db.Do(stmt).Scan({{ var .StructName }})
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return {{ var .StructName }}.{{ $methodForFetch }}(db)
	}
	return nil
}
					`,
				createMethod("UpdateBy%sWithMap", fieldNamesWithoutEnabled...),
				createMethod("FetchBy%s", fieldNamesWithoutEnabled...),
				toExactlyConditionFrom(fieldNames...),
			))

			m.ParseTo(buf, fmt.Sprintf(`
{{ $method := "%s" }}
{{ $methodForUpdateWithMap := "%s" }}
func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB, zeroFields ...string) error {
	fieldValues := {{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.FieldValuesFromStructByNonZero({{ var .StructName }}, zeroFields...)
	return {{ var .StructName }}.{{ $methodForUpdateWithMap }}(db, fieldValues)
}
					`,
				createMethod("UpdateBy%sWithStruct", fieldNamesWithoutEnabled...),
				createMethod("UpdateBy%sWithMap", fieldNamesWithoutEnabled...),
			))

			if m.HasSoftDelete {

				m.ParseTo(buf, fmt.Sprintf(`
{{ $method := "%s" }}
{{ $methodForDelete := "%s" }}
func ({{ var .StructName }} *{{ .StructName }}) {{ $method }}(db *{{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DB) error {
	{{ ( .EnableIfNeed ) }}
	table :=  {{ var .StructName }}.T()

	fieldValues := {{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.FieldValues{}
	fieldValues["{{ .FieldSoftDelete }}"] = {{ use .ConstSoftDeleteFalse }}

	{{ ( .SetUpdatedForFieldValuesAtIfNeed ) }}

	stmt := table.Update().
		Comment("{{ .StructName }}.{{ $method }}").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where({{ use "github.com/profzone/eden-framework/pkg/sqlx/builder" }}.And(
			%s
		))

	dbRet := db.Do(stmt).Scan({{ var .StructName }})
	err := dbRet.Err()
	if err != nil {
		dbErr := {{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.DBErr(err)
		if dbErr.IsConflict() {
			return 	{{ var .StructName }}.{{ $methodForDelete }}(db)
		}
		return err
	}
	return nil
}
					`,
					createMethod("SoftDeleteBy%s", fieldNamesWithoutEnabled...),
					createMethod("DeleteBy%s", fieldNamesWithoutEnabled...),
					toExactlyConditionFrom(fieldNames...),
				))
			}
		}
	})

	return buf.String()
}
