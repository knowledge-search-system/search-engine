package searchhistory

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knowledge-search-system/search-engine/internal/model"
	"github.com/knowledge-search-system/search-engine/internal/repository/sql"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) model.SearchHistoryRepository {
	return &Repository{pool: pool}
}

func (r *Repository) Save(ctx context.Context, entry model.SearchHistoryEntry) error {
	query, args, err := buildInsertQuery(entry)
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	executor := sql.ExecutorFromContext(ctx, r.pool)
	if _, err := executor.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insert search history: %w", err)
	}

	return nil
}
