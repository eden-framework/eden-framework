package datatypes

type OperateTime struct {
	// 创建时间
	CreatedAt MySQLTimestamp `db:"f_created_at,default='0'" json:"createdAt" `
	// 更新时间
	UpdatedAt MySQLTimestamp `db:"f_updated_at,default='0'" json:"updatedAt"`
	// 删除时间
	DeletedAt MySQLTimestamp `db:"f_deleted_at,default='0'" json:"-"`
}
