package apperrors

var (
	ErrEmptyQuery            = New(CodeInvalidArgument, "search.empty_query")
	ErrInvalidPagination     = New(CodeInvalidArgument, "search.invalid_pagination")
	ErrNoChunksToIndex       = New(CodeInvalidArgument, "index.no_chunks")
	ErrElasticsearch         = New(CodeUnavailable, "search.elasticsearch_unavailable")
	ErrSearchHistoryNotFound = New(CodeNotFound, "search.history_not_found")
	ErrInternal              = New(CodeInternal, "common.internal_error")
)
