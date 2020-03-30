package mysql_logger_driver

import (
	"database/sql/driver"

	"github.com/fatih/color"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type loggerConn struct {
	cfg  *mysql.Config
	conn driver.Conn
}

func (c *loggerConn) Begin() (driver.Tx, error) {
	logrus.Debugf(color.YellowString("=========== Beginning Transaction ==========="))
	tx, err := c.conn.Begin()
	if err != nil {
		logrus.Errorf("failed to begin transaction: %s", err)
		return nil, err
	}
	return &loggingTx{tx: tx}, nil
}

func (c *loggerConn) Close() error {
	if err := c.conn.Close(); err != nil {
		logrus.Errorf("failed to close connection: %s", err)
		return err
	}
	return nil
}

func (c *loggerConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.conn.Prepare(query)
	if err != nil {
		logrus.Errorf("failed to prepare query: %s, err: %s", query, err)
		return nil, err
	}
	return &loggerStmt{cfg: c.cfg, query: query, stmt: stmt}, nil
}
