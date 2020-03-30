package mysql_logger_driver

import (
	"database/sql"
	"database/sql/driver"
	"strings"

	"github.com/fatih/color"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type LoggingDriver struct {
	Driver string
}

func (d LoggingDriver) Open(dsn string) (driver.Conn, error) {
	mysqlDriver := &mysql.MySQLDriver{}
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic(err)
	}
	cfg.Passwd = strings.Repeat("*", len(cfg.Passwd))

	conn, err := mysqlDriver.Open(dsn)
	if err != nil {
		logrus.Errorf("failed to open connection: %s %s", cfg.FormatDSN(), err)
		return nil, err
	}
	logrus.Debugf(color.YellowString("connected %s", cfg.FormatDSN()))

	return &loggerConn{cfg: cfg, conn: conn}, nil
}

func init() {
	sql.Register("logger:mysql", &LoggingDriver{"mysql"})
}
