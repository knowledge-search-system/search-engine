package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/knowledge-search-system/search-engine/config"
	"github.com/knowledge-search-system/search-engine/internal/apperrors"
	"github.com/knowledge-search-system/search-engine/internal/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	defaultPage     = int32(1)
	defaultPageSize = int32(10)
	maxPageSize     = int32(100)

	cacheKeyPrefix = "search_engine:search"
)

type SearchService struct {
	esClient    *elasticsearch.Client
	indexName   string
	redisClient *redis.Client
	cacheTTL    time.Duration
	historyRepo model.SearchHistoryRepository
	logger      *zap.Logger
}

func NewSearchService(
	esClient *elasticsearch.Client,
	redisClient *redis.Client,
	historyRepo model.SearchHistoryRepository,
	cfg *config.Config,
	logger *zap.Logger,
) *SearchService {
	return &SearchService{
		esClient:    esClient,
		indexName:   cfg.Elasticsearch.IndexName,
		redisClient: redisClient,
		cacheTTL:    cfg.Redis.CacheTTL,
		historyRepo: historyRepo,
		logger:      logger,
	}
}

type SearchOutcome struct {
	Results  []model.SearchResult
	Total    int32
	Page     int32
	PageSize int32
}

func (s *SearchService) Search(ctx context.Context, query string, page, pageSize int32) (SearchOutcome, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return SearchOutcome{}, apperrors.ErrEmptyQuery
	}

	page, pageSize, err := normalizePagination(page, pageSize)
	if err != nil {
		return SearchOutcome{}, err
	}

	s.saveHistory(ctx, query)

	cacheKey := buildCacheKey(query, page, pageSize)
	if cached, ok := s.getFromCache(ctx, cacheKey); ok {
		return cached, nil
	}

	outcome, err := s.searchElasticsearch(ctx, query, page, pageSize)
	if err != nil {
		return SearchOutcome{}, err
	}

	s.setCache(ctx, cacheKey, outcome)

	return outcome, nil
}

func normalizePagination(page, pageSize int32) (int32, int32, error) {
	if page < 0 || pageSize < 0 {
		return 0, 0, apperrors.ErrInvalidPagination
	}

	if page == 0 {
		page = defaultPage
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		return 0, 0, apperrors.ErrInvalidPagination
	}

	return page, pageSize, nil
}

func (s *SearchService) saveHistory(ctx context.Context, query string) {
	entry := model.SearchHistoryEntry{
		ID:        uuid.New(),
		Query:     query,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.historyRepo.Save(ctx, entry); err != nil {
		s.logger.Error("failed to save search history", zap.Error(err), zap.String("query", query))
	}
}

type esSearchRequestBody struct {
	From  int32 `json:"from"`
	Size  int32 `json:"size"`
	Query struct {
		MultiMatch struct {
			Query  string   `json:"query"`
			Fields []string `json:"fields"`
		} `json:"multi_match"`
	} `json:"query"`
}

type esSearchResponseBody struct {
	Hits struct {
		Total struct {
			Value int32 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Score  float64       `json:"_score"`
			Source chunkDocument `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (s *SearchService) searchElasticsearch(ctx context.Context, query string, page, pageSize int32) (SearchOutcome, error) {
	var reqBody esSearchRequestBody
	reqBody.From = (page - 1) * pageSize
	reqBody.Size = pageSize
	reqBody.Query.MultiMatch.Query = query
	reqBody.Query.MultiMatch.Fields = []string{"text"}

	encoded, err := json.Marshal(reqBody)
	if err != nil {
		return SearchOutcome{}, apperrors.ErrInternal.WithErr(fmt.Errorf("marshal search request: %w", err))
	}

	resp, err := esapi.SearchRequest{
		Index: []string{s.indexName},
		Body:  bytes.NewReader(encoded),
	}.Do(ctx, s.esClient)
	if err != nil {
		return SearchOutcome{}, apperrors.ErrElasticsearch.WithErr(fmt.Errorf("execute search request: %w", err))
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return SearchOutcome{}, apperrors.ErrElasticsearch.WithErr(fmt.Errorf("elasticsearch returned status %s", resp.Status()))
	}

	var parsed esSearchResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return SearchOutcome{}, apperrors.ErrInternal.WithErr(fmt.Errorf("decode search response: %w", err))
	}

	results := make([]model.SearchResult, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		results = append(results, model.SearchResult{
			ChunkID:  hit.Source.ChunkID,
			FileName: hit.Source.FileName,
			Page:     hit.Source.PageNumber,
			Text:     hit.Source.Text,
			Score:    hit.Score,
		})
	}

	return SearchOutcome{
		Results:  results,
		Total:    parsed.Hits.Total.Value,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func buildCacheKey(query string, page, pageSize int32) string {
	hash := sha256.Sum256([]byte(query))
	return fmt.Sprintf("%s:%x:%d:%d", cacheKeyPrefix, hash, page, pageSize)
}

func (s *SearchService) getFromCache(ctx context.Context, key string) (SearchOutcome, bool) {
	raw, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			s.logger.Warn("redis get failed", zap.Error(err))
		}
		return SearchOutcome{}, false
	}

	var outcome SearchOutcome
	if err := json.Unmarshal(raw, &outcome); err != nil {
		s.logger.Warn("failed to unmarshal cached search outcome", zap.Error(err))
		return SearchOutcome{}, false
	}

	return outcome, true
}

func (s *SearchService) setCache(ctx context.Context, key string, outcome SearchOutcome) {
	raw, err := json.Marshal(outcome)
	if err != nil {
		s.logger.Warn("failed to marshal search outcome for cache", zap.Error(err))
		return
	}

	if err := s.redisClient.Set(ctx, key, raw, s.cacheTTL).Err(); err != nil {
		s.logger.Warn("redis set failed", zap.Error(err))
	}
}
