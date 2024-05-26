package postgresStorage

import (
	"fmt"
	// required to register the PostgreSQL driver for the database/sql.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	ApplicationName string
	PgDriver        string
}

func New(cfg Config) (*sqlx.DB, error) {
	return sqlx.Connect(
		cfg.PgDriver, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s application_name=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.DBName,
			cfg.SSLMode,
			cfg.ApplicationName,
		))
}
