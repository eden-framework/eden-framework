package sqlx_test

import (
	"testing"

	"github.com/eden-framework/eden-framework/pkg/sqlx"
)

func BenchmarkDB_DBExecutor(b *testing.B) {
	dbTest := sqlx.NewDatabase("test_for_user")
	db := dbTest.OpenDB(mysqlConnector)

	run := func(db sqlx.DBExecutor) {
		db.D()
	}

	for i := 0; i <= b.N; i++ {
		run(db)
	}
}
