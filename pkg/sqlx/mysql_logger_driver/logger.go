package mysql_logger_driver

import (
	"github.com/go-sql-driver/mysql"
)

func init() {
	mysql.SetLogger(&logger{})
}

type logger struct{}

func (l *logger) Print(args ...interface{}) {
}
