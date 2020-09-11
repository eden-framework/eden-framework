package datatypes

type PrimaryID struct {
	ID uint64 `db:"F_id,autoincrement" json:"-"`
}
