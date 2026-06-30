package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knowledge-search-system/search-engine/config"
	"github.com/knowledge-search-system/search-engine/pkg/postgres"
	"go.uber.org/fx"
)

var Module = fx.Module("repository",
	fx.Provide(
		newPostgresPool,
		NewTransactor,
	),
)

func newPostgresPool(cfg *config.Config) (*pgxpool.Pool, error) {
	return postgres.NewPool(context.Background(), postgres.Config{
		DSN:      cfg.Postgres.DSN,
		MaxConns: cfg.Postgres.MaxConns,
		MinConns: cfg.Postgres.MinConns,
	})
}
