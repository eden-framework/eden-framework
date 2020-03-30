package gen

import (
	"bytes"

	"github.com/profzone/eden-framework/pkg/sqlx/builder"
)

func (m *Model) GetComments() map[string]string {
	comments := map[string]string{}
	m.Columns.Range(func(col *builder.Column, idx int) {
		comments[col.FieldName] = col.Comment
	})
	return comments
}

func (m *Model) interfaces() string {
	buf := &bytes.Buffer{}

	if len(m.Keys.Primary) > 0 {
		m.ParseTo(buf, `func ({{ var .StructName }} *{{ .StructName }}) PrimaryKey() {{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.FieldNames {
		return {{ dump .Keys.Primary }}
	}
	`)
	}

	if len(m.Keys.Indexes) > 0 {
		m.ParseTo(buf, `func ({{ var .StructName }} *{{ .StructName }}) Indexes() {{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.Indexes {
		return {{ dump .Keys.Indexes }}
	}
	`)
	}

	if len(m.Keys.UniqueIndexes) > 0 {
		m.ParseTo(buf, `func ({{ var .StructName }} *{{ .StructName }}) UniqueIndexes() {{ use "github.com/profzone/eden-framework/pkg/sqlx" }}.Indexes {
		return {{ dump .Keys.UniqueIndexes }}
	}
	`)
	}

	if m.WithComments {
		m.ParseTo(buf, `func ({{ var .StructName }} *{{ .StructName }}) Comments() map[string]string {
		return {{ dump ( .GetComments ) }}
	}`)
	}

	return buf.String()
}
