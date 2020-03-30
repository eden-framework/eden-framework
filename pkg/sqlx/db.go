package sqlx

import (
	"database/sql"
	"errors"
	"time"

	"github.com/profzone/eden-framework/pkg/sqlx/builder"
	_ "github.com/profzone/eden-framework/pkg/sqlx/mysql_logger_driver"
)

var ErrNotTx = errors.New("db is not *sql.Tx")
var ErrNotDB = errors.New("db is not *sql.DB")

func Open(driverName string, dataSourceName string, openFunc func(string, string) (*sql.DB, error)) (*DB, error) {
	db, err := openFunc(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{
		SqlExecutor: db,
	}, nil
}

func MustOpen(driverName string, dataSourceName string, openFunc func(string, string) (*sql.DB, error)) *DB {
	db, err := Open(driverName, dataSourceName, openFunc)
	if err != nil {
		panic(err)
	}
	return db
}

type DB struct {
	SqlExecutor
}

func (d *DB) Do(stmt builder.Statement) (result *Result) {
	return Do(d, stmt)
}

func (d *DB) IsTx() bool {
	_, ok := d.SqlExecutor.(*sql.Tx)
	return ok
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	query, args = flattenArgs(query, args...)
	return d.SqlExecutor.Query(query, args...)
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	query, args = flattenArgs(query, args...)
	return d.SqlExecutor.Exec(query, args...)
}

func (d *DB) Begin() (*DB, error) {
	if d.IsTx() {
		return nil, ErrNotDB
	}
	db, err := d.SqlExecutor.(*sql.DB).Begin()
	if err != nil {
		return nil, err
	}
	return &DB{
		SqlExecutor: db,
	}, nil
}

func (d *DB) Commit() error {
	if !d.IsTx() {
		return ErrNotTx
	}
	return d.SqlExecutor.(*sql.Tx).Commit()
}

func (d *DB) Rollback() error {
	if !d.IsTx() {
		return ErrNotTx
	}
	return d.SqlExecutor.(*sql.Tx).Rollback()
}

func (d *DB) SetMaxOpenConns(n int) {
	d.SqlExecutor.(*sql.DB).SetMaxOpenConns(n)
}

func (d *DB) SetMaxIdleConns(n int) {
	d.SqlExecutor.(*sql.DB).SetMaxIdleConns(n)
}

func (d *DB) SetConnMaxLifetime(t time.Duration) {
	d.SqlExecutor.(*sql.DB).SetConnMaxLifetime(t)
}
