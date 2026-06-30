package searchhistory

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/knowledge-search-system/search-engine/internal/model"
)

func TestBuildInsertQuery(t *testing.T) {
	entry := model.SearchHistoryEntry{
		ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Query:     "машинное обучение",
		CreatedAt: time.Date(2026, 6, 25, 10, 0, 0, 0, time.UTC),
	}

	query, args, err := buildInsertQuery(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantQuery := `INSERT INTO search_history (id,query,created_at) VALUES ($1,$2,$3)`
	if query != wantQuery {
		t.Errorf("query = %q, want %q", query, wantQuery)
	}

	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}
	if args[0] != entry.ID {
		t.Errorf("args[0] = %v, want %v", args[0], entry.ID)
	}
	if args[1] != entry.Query {
		t.Errorf("args[1] = %v, want %v", args[1], entry.Query)
	}
	if args[2] != entry.CreatedAt {
		t.Errorf("args[2] = %v, want %v", args[2], entry.CreatedAt)
	}
}
