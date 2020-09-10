package datatypes

type PrimaryID struct {
	ID uint64 `db:"F_id" sql:"bigint unsigned NOT NULL AUTO_INCREMENT" json:"-"`
}
