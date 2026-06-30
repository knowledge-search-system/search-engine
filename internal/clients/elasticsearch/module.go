package elasticsearch

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/knowledge-search-system/search-engine/config"
	"go.uber.org/fx"
)

var Module = fx.Module("elasticsearch_client",
	fx.Provide(NewClient),
	fx.Invoke(registerIndexHook),
)

func registerIndexHook(lc fx.Lifecycle, client *elasticsearch.Client, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return EnsureDocumentsIndex(ctx, client, cfg.Elasticsearch.IndexName)
		},
	})
}
