package sqlx

import (
	"fmt"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

type Task func(db *DB) error

func (task Task) Run(db *DB) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %s; calltrace:%s", fmt.Sprint(e), string(debug.Stack()))
		}
	}()
	return task(db)
}

func NewTasks(db *DB) *Tasks {
	return &Tasks{
		db: db,
	}
}

type Tasks struct {
	db    *DB
	tasks []Task
}

func (tasks Tasks) With(task ...Task) *Tasks {
	tasks.tasks = append(tasks.tasks, task...)
	return &tasks
}

func (tasks *Tasks) Do() (err error) {
	if len(tasks.tasks) == 0 {
		return nil
	}

	db := tasks.db
	inTxScope := false

	if !db.IsTx() {
		db, err = db.Begin()
		if err != nil {
			return err
		}
		inTxScope = true
	}

	for _, task := range tasks.tasks {
		if runErr := task.Run(db); runErr != nil {
			if inTxScope {
				// err will bubble up，just handle and rollback in outermost layer
				logrus.Errorf("SQL FAILED: %s", runErr.Error())
				if rollBackErr := db.Rollback(); rollBackErr != nil {
					logrus.Errorf("ROLLBACK FAILED: %s", rollBackErr.Error())
					err = rollBackErr
					return
				}
			}
			return runErr
		}
	}

	if inTxScope {
		if commitErr := db.Commit(); commitErr != nil {
			logrus.Errorf("TRANSACTION COMMIT FAILED: %s", commitErr.Error())
			return commitErr
		}
	}

	return nil
}
