package confpostgres

import (
	"context"
	"fmt"
	"github.com/profzone/eden-framework/pkg/sqlx"
	"github.com/profzone/eden-framework/pkg/sqlx/postgresqlconnector"
	"github.com/profzone/envconfig"
	"time"
)

type Postgres struct {
	Host            string
	SlaveHost       string
	Port            int
	User            string
	Password        envconfig.Password
	Extra           string
	Extensions      []string
	PoolSize        int
	ConnMaxLifetime envconfig.Duration
	Database        *sqlx.Database `ignored:"true"`

	*sqlx.DB `ignored:"true"`
	slaveDB  *sqlx.DB `ignored:"true"`
}

func (m *Postgres) LivenessCheck() map[string]string {
	s := map[string]string{}

	_, err := m.DB.ExecContext(context.Background(), "SELECT 1")
	if err != nil {
		s[m.Host] = err.Error()
	} else {
		s[m.Host] = "ok"
	}

	if m.slaveDB != nil {
		_, err := m.slaveDB.ExecContext(context.Background(), "SELECT 1")
		if err != nil {
			s[m.SlaveHost] = err.Error()
		} else {
			s[m.SlaveHost] = "ok"
		}
	}

	return s
}

func (m *Postgres) SetDefaults() {
	if m.Host == "" {
		m.Host = "127.0.0.1"
	}

	if m.Port == 0 {
		m.Port = 5432
	}

	if m.PoolSize == 0 {
		m.PoolSize = 10
	}

	if m.ConnMaxLifetime == 0 {
		m.ConnMaxLifetime = envconfig.Duration(1 * time.Hour)
	}

	if m.Extra == "" {
		m.Extra = "sslmode=disable"
	}
}

func (m *Postgres) url(host string) string {
	password := m.Password
	if password != "" {
		password = ":" + password
	}
	return fmt.Sprintf("postgres://%s%s@%s:%d", m.User, password, host, m.Port)
}

func (m *Postgres) conn(host string) (*sqlx.DB, error) {
	db := m.Database.OpenDB(&postgresqlconnector.PostgreSQLConnector{
		Host:       m.url(host),
		Extra:      m.Extra,
		Extensions: m.Extensions,
	})

	db.SetMaxOpenConns(m.PoolSize)
	db.SetMaxIdleConns(m.PoolSize / 2)
	db.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))

	_, err := db.ExecContext(context.Background(), "SELECT 1")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m *Postgres) UseSlave() sqlx.DBExecutor {
	if m.slaveDB != nil {
		return m.slaveDB
	}
	return m.DB
}

func (m *Postgres) Init() {
	r := Retry{Repeats: 5, Interval: envconfig.Duration(1 * time.Second)}

	err := r.Do(func() error {
		db, err := m.conn(m.Host)
		if err != nil {
			return err
		}
		m.DB = db
		return nil
	})

	if err != nil {
		panic(err)
	}

	if m.SlaveHost != "" {
		err := r.Do(func() error {
			db, err := m.conn(m.Host)
			if err != nil {
				return err
			}
			m.slaveDB = db
			return nil
		})

		if err != nil {
			panic(err)
		}
	}
}

func SwitchSlave(executor sqlx.DBExecutor) sqlx.DBExecutor {
	if canSlave, ok := executor.(CanSlave); !ok {
		return canSlave.UseSlave()
	}
	return executor
}

type CanSlave interface {
	UseSlave() sqlx.DBExecutor
}
