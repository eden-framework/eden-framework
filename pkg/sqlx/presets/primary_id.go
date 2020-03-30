package presets

import (
	"database/sql"

	"github.com/profzone/eden-framework/pkg/sqlx"
)

var _ interface {
	sqlx.WithPrimaryKey
} = (*PrimaryID)(nil)

type PrimaryID struct {
	ID uint64 `db:"F_id" sql:"bigint unsigned NOT NULL AUTO_INCREMENT" json:"-"`
}

func (id PrimaryID) PrimaryKey() sqlx.FieldNames {
	return []string{"ID"}
}

func (id *PrimaryID) AfterInsert(result sql.Result) error {
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	id.ID = uint64(lastInsertID)
	return nil
}
