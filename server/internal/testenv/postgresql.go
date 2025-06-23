package testenv

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/speakeasy-api/gram/internal/o11y"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var pgCloneMutex sync.Mutex

func nextRandom() string {
	return fmt.Sprintf("%d", uuid.New().ID())
}

type PostgresDBCloneFunc func(t *testing.T, name string) (*pgxpool.Pool, error)

// NewTestPostgres creates a new Postgres container with a template database built
// from a SQL init script. A reference to the container is returned as well as
// a function to create test databases from the template. All "clone" databases
// are automatically dropped when the test ends using t.Cleanup() hooks.
func NewTestPostgres(ctx context.Context) (*postgres.PostgresContainer, PostgresDBCloneFunc, error) {
	container, err := postgres.Run(
		ctx,
		"postgres:17",
		postgres.WithUsername("gotest"),
		postgres.WithPassword("gotest"),
		postgres.WithDatabase("gotestdb"),
		postgres.WithInitScripts(filepath.Join("..", "..", "database", "schema.sql")),
		postgres.BasicWaitStrategies(),
		// Store the database in-memory for faster tests
		testcontainers.WithTmpfs(map[string]string{"/var/lib/postgresql/data": "rw"}),
		testcontainers.WithEnv(map[string]string{"PGDATA": "/var/lib/postgresql/data"}),
		testcontainers.WithLogger(NewTestcontainersLogger()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("start postgres container: %w", err)
	}

	uri, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("read connection string: %w", err)
	}

	conn, err := pgx.Connect(ctx, uri)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to template database: %w", err)
	}
	defer o11y.NoLogDefer(func() error {
		return conn.Close(ctx)
	})

	_, err = conn.Exec(ctx, "ALTER DATABASE gotestdb WITH is_template = true;")
	if err != nil {
		return nil, nil, fmt.Errorf("mark template database: %w", err)
	}

	return container, newPostgresCloneFunc(container), nil
}

func newPostgresCloneFunc(container *postgres.PostgresContainer) PostgresDBCloneFunc {
	return func(t *testing.T, name string) (*pgxpool.Pool, error) {
		t.Helper()
		ctx := t.Context()
		uri, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			return nil, fmt.Errorf("read connection string: %w", err)
		}

		pgCloneMutex.Lock()
		defer pgCloneMutex.Unlock()

		conn, err := pgx.Connect(ctx, uri)
		if err != nil {
			return nil, fmt.Errorf("connect to template database: %w", err)
		}
		defer o11y.NoLogDefer(func() error {
			return conn.Close(ctx)
		})

		clonename := fmt.Sprintf("%s_%s", name, nextRandom())
		_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE gotestdb;", clonename))
		if err != nil {
			return nil, fmt.Errorf("create test database: %w", err)
		}

		cloneuri := strings.Replace(uri, "gotestdb", clonename, 1)
		pool, err := pgxpool.New(ctx, cloneuri)
		if err != nil {
			return nil, fmt.Errorf("create pgx pool: %w", err)
		}

		t.Cleanup(func() {
			timeoutCtx, cancel := context.WithTimeout(context.WithoutCancel(t.Context()), 60*time.Second)
			defer cancel()

			pool.Close()

			conn, err := pgx.Connect(timeoutCtx, uri)
			if err != nil {
				panic(fmt.Errorf("drop test database: connect: %w", err))
			}
			defer o11y.NoLogDefer(func() error {
				return conn.Close(timeoutCtx)
			})

			_, err = conn.Exec(timeoutCtx, fmt.Sprintf("DROP DATABASE %s;", clonename))
			if err != nil {
				panic(fmt.Errorf("drop test database: exec: %w", err))
			}
		})

		return pool, nil
	}
}
