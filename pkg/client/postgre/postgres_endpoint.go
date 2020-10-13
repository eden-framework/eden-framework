package confpostgres

import (
	"context"
	"fmt"
	"github.com/eden-framework/sqlx"
	"github.com/eden-framework/sqlx/postgresqlconnector"
	"github.com/profzone/envconfig"
	"time"
)

type PostgresEndpoint struct {
	Endpoint envconfig.Endpoint `env:""`
	Database *sqlx.Database     `env:"-"`

	Extensions      []string
	PoolSize        int
	ConnMaxLifetime envconfig.Duration

	*sqlx.DB `env:"-"`
}

func (m *PostgresEndpoint) LivenessCheck() map[string]string {
	s := map[string]string{}

	_, err := m.DB.ExecContext(context.Background(), "SELECT 1")
	if err != nil {
		s[m.Endpoint.Host()] = err.Error()
	} else {
		s[m.Endpoint.Host()] = "ok"
	}

	return s
}

func (m *PostgresEndpoint) SetDefaults() {
	if m.PoolSize == 0 {
		m.PoolSize = 10
	}

	if m.ConnMaxLifetime == 0 {
		m.ConnMaxLifetime = envconfig.Duration(1 * time.Hour)
	}

	if m.Endpoint.IsZero() {
		m.Endpoint.Hostname = "127.0.0.1"
		m.Endpoint.Port = 5432
	}

	if m.Database.Name == "" {
		if len(m.Endpoint.Base) > 0 {
			m.Database.Name = m.Endpoint.Base
		}
	}
}

func (m *PostgresEndpoint) url(host string) string {
	password := m.Endpoint.Password
	if password != "" {
		password = ":" + password
	}
	return fmt.Sprintf("postgres://%s%s@%s", m.Endpoint.Username, password, host)
}

func (m *PostgresEndpoint) conn(host string) (*sqlx.DB, error) {
	db := m.Database.OpenDB(&postgresqlconnector.PostgreSQLConnector{
		Host:       m.url(host),
		Extra:      m.Endpoint.Extra.Encode(),
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

func (m *PostgresEndpoint) Init() {
	r := Retry{Repeats: 5, Interval: envconfig.Duration(1 * time.Second)}

	err := r.Do(func() error {
		db, err := m.conn(m.Endpoint.Host())
		if err != nil {
			return err
		}
		m.DB = db
		return nil
	})

	if err != nil {
		panic(err)
	}
}
