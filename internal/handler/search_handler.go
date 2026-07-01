package handler

import (
	"context"

	searchenginev1 "github.com/knowledge-search-system/search-engine/proto/searchengine/v1"
	"github.com/knowledge-search-system/search-engine/internal/model"
	"github.com/knowledge-search-system/search-engine/internal/service"
)

type SearchHandler struct {
	searchenginev1.UnimplementedSearchServiceServer

	searchService *service.SearchService
	indexService  *service.IndexService
}

func NewSearchHandler(searchService *service.SearchService, indexService *service.IndexService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		indexService:  indexService,
	}
}

func (h *SearchHandler) Search(ctx context.Context, req *searchenginev1.SearchRequest) (*searchenginev1.SearchResponse, error) {
	outcome, err := h.searchService.Search(ctx, req.GetQ(), req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, err
	}

	results := make([]*searchenginev1.SearchResult, 0, len(outcome.Results))
	for _, r := range outcome.Results {
		results = append(results, &searchenginev1.SearchResult{
			ChunkId:  r.ChunkID,
			FileName: r.FileName,
			Page:     r.Page,
			Text:     r.Text,
			Score:    r.Score,
		})
	}

	return &searchenginev1.SearchResponse{
		Results:  results,
		Total:    outcome.Total,
		Page:     outcome.Page,
		PageSize: outcome.PageSize,
	}, nil
}

func (h *SearchHandler) IndexChunks(ctx context.Context, req *searchenginev1.IndexChunksRequest) (*searchenginev1.IndexChunksResponse, error) {
	chunks := make([]model.Chunk, 0, len(req.GetChunks()))
	for _, c := range req.GetChunks() {
		chunks = append(chunks, model.Chunk{
			ChunkID:    c.GetChunkId(),
			DocumentID: c.GetDocumentId(),
			FileName:   c.GetFileName(),
			PageNumber: c.GetPageNumber(),
			Text:       c.GetText(),
		})
	}

	indexed, err := h.indexService.IndexChunks(ctx, chunks)
	if err != nil {
		return nil, err
	}

	return &searchenginev1.IndexChunksResponse{
		IndexedCount: int32(indexed),
	}, nil
}
