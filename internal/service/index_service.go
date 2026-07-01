package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/knowledge-search-system/search-engine/config"
	"github.com/knowledge-search-system/search-engine/internal/apperrors"
	"github.com/knowledge-search-system/search-engine/internal/model"
	"go.uber.org/zap"
)

type IndexService struct {
	esClient  *elasticsearch.Client
	indexName string
	logger    *zap.Logger
}

func NewIndexService(esClient *elasticsearch.Client, cfg *config.Config, logger *zap.Logger) *IndexService {
	return &IndexService{
		esClient:  esClient,
		indexName: cfg.Elasticsearch.IndexName,
		logger:    logger,
	}
}

type chunkDocument struct {
	ChunkID    string `json:"chunk_id"`
	DocumentID string `json:"document_id"`
	FileName   string `json:"file_name"`
	PageNumber int32  `json:"page_number"`
	Text       string `json:"text"`
}

func (s *IndexService) IndexChunks(ctx context.Context, chunks []model.Chunk) (int, error) {
	if len(chunks) == 0 {
		return 0, apperrors.ErrNoChunksToIndex
	}

	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: s.esClient,
		Index:  s.indexName,
	})
	if err != nil {
		return 0, apperrors.ErrElasticsearch.WithErr(fmt.Errorf("create bulk indexer: %w", err))
	}

	var failed atomic.Int32

	for _, chunk := range chunks {
		body, err := json.Marshal(chunkDocument{
			ChunkID:    chunk.ChunkID,
			DocumentID: chunk.DocumentID,
			FileName:   chunk.FileName,
			PageNumber: chunk.PageNumber,
			Text:       chunk.Text,
		})
		if err != nil {
			return 0, apperrors.ErrInternal.WithErr(fmt.Errorf("marshal chunk document: %w", err))
		}

		item := esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: chunk.ChunkID,
			Body:       bytes.NewReader(body),
			OnFailure: func(_ context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				failed.Add(1)
				s.logger.Error("failed to index chunk",
					zap.String("chunk_id", item.DocumentID),
					zap.String("es_error", res.Error.Reason),
					zap.Error(err),
				)
			},
		}

		if err := indexer.Add(ctx, item); err != nil {
			return 0, apperrors.ErrElasticsearch.WithErr(fmt.Errorf("add chunk to bulk indexer: %w", err))
		}
	}

	if err := indexer.Close(ctx); err != nil {
		return 0, apperrors.ErrElasticsearch.WithErr(fmt.Errorf("close bulk indexer: %w", err))
	}

	stats := indexer.Stats()
	if failed.Load() > 0 {
		return int(stats.NumIndexed), apperrors.ErrElasticsearch.WithErr(
			fmt.Errorf("%d of %d chunks failed to index", failed.Load(), len(chunks)),
		)
	}

	return int(stats.NumIndexed), nil
}
