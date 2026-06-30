package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Chunk struct {
	ChunkID    string
	DocumentID string
	FileName   string
	PageNumber int32
	Text       string
}

type SearchResult struct {
	ChunkID  string
	FileName string
	Page     int32
	Text     string
	Score    float64
}

type SearchHistoryEntry struct {
	ID        uuid.UUID
	Query     string
	CreatedAt time.Time
}

type SearchHistoryRepository interface {
	Save(ctx context.Context, entry SearchHistoryEntry) error
}
