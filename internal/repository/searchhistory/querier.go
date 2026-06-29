package searchhistory

import (
	"github.com/Masterminds/squirrel"
	"github.com/knowledge-search-system/search-engine/internal/model"
	"github.com/knowledge-search-system/search-engine/internal/repository/dbconsts"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func buildInsertQuery(entry model.SearchHistoryEntry) (string, []any, error) {
	return psql.Insert(dbconsts.SearchHistoryTable).
		Columns("id", "query", "created_at").
		Values(entry.ID, entry.Query, entry.CreatedAt).
		ToSql()
}
