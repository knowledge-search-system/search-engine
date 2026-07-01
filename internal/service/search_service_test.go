package service

import (
	"errors"
	"testing"

	"github.com/knowledge-search-system/search-engine/internal/apperrors"
)

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name         string
		page         int32
		pageSize     int32
		wantPage     int32
		wantPageSize int32
		wantErr      error
	}{
		{name: "defaults when zero", page: 0, pageSize: 0, wantPage: defaultPage, wantPageSize: defaultPageSize},
		{name: "explicit values kept", page: 2, pageSize: 25, wantPage: 2, wantPageSize: 25},
		{name: "negative page rejected", page: -1, pageSize: 10, wantErr: apperrors.ErrInvalidPagination},
		{name: "negative page size rejected", page: 1, pageSize: -10, wantErr: apperrors.ErrInvalidPagination},
		{name: "page size over max rejected", page: 1, pageSize: maxPageSize + 1, wantErr: apperrors.ErrInvalidPagination},
		{name: "page size at max accepted", page: 1, pageSize: maxPageSize, wantPage: 1, wantPageSize: maxPageSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPageSize, err := normalizePagination(tt.page, tt.pageSize)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotPage != tt.wantPage {
				t.Errorf("page = %d, want %d", gotPage, tt.wantPage)
			}
			if gotPageSize != tt.wantPageSize {
				t.Errorf("pageSize = %d, want %d", gotPageSize, tt.wantPageSize)
			}
		})
	}
}

func TestBuildCacheKey(t *testing.T) {
	keyA := buildCacheKey("машинное обучение", 1, 10)
	keyB := buildCacheKey("машинное обучение", 1, 10)
	keyC := buildCacheKey("машинное обучение", 2, 10)
	keyD := buildCacheKey("нейронные сети", 1, 10)

	if keyA != keyB {
		t.Errorf("expected identical keys for identical inputs, got %q vs %q", keyA, keyB)
	}
	if keyA == keyC {
		t.Errorf("expected different keys for different page, got the same key %q", keyA)
	}
	if keyA == keyD {
		t.Errorf("expected different keys for different query, got the same key %q", keyA)
	}
}
