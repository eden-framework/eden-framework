package mysql_logger_driver

import (
	"database/sql/driver"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type loggingTx struct {
	tx driver.Tx
}

func (tx *loggingTx) Commit() error {
	if err := tx.tx.Commit(); err != nil {
		logrus.Debugf("failed to commit transaction: %s", err)
		return err
	}
	logrus.Debugf(color.YellowString("=========== Committed Transaction ==========="))
	return nil
}

func (tx *loggingTx) Rollback() error {
	if err := tx.tx.Rollback(); err != nil {
		logrus.Debugf("failed to rollback transaction: %s", err)
		return err
	}
	logrus.Debugf("=========== Rollback Transaction ===========")
	return nil
}
