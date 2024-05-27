package containers

import (
	"Service/pkg/utils"
	"context"
	"path/filepath"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const StartUpTimeout = 5 * time.Second

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:alpine"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(StartUpTimeout)),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	goModPath, err := utils.FindGoModPath()
	if err != nil {
		return nil, err
	}

	migrationsDir := filepath.Join(filepath.Dir(goModPath), "/db")
	gooseDB, err := goose.OpenDBWithDriver("pgx", connStr)
	if err != nil {
		return nil, err
	}
	defer gooseDB.Close()

	if err = goose.Up(gooseDB, migrationsDir); err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}
