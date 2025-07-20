package postgres

import (
	"context"

	"github.com/Sidney-Bernardin/Chat-Alloy/internal"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var (
	ErrNoRows = errors.New("no rows")
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *internal.Config) (*Repository, error) {

	pool, err := pgxpool.New(ctx, cfg.POSTGRES_URL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create connection pool")
	}

	m, err := migrate.New("file://internal/repos/postgres/Migrations", cfg.POSTGRES_URL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create migration instance")
	}

	if err := m.Migrate(20250718202520); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, errors.Wrap(err, "cannot migrate")
	}

	return &Repository{pool}, nil
}

func row[T any](ctx context.Context, repo *Repository, query string, args ...any) (*T, error) {

	row, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query failed")
	}

	model, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[T])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNoRows
		default:
			return nil, errors.Wrap(err, "cannot collect row")
		}
	}

	return model, nil
}
