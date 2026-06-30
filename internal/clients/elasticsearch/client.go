package elasticsearch

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/knowledge-search-system/search-engine/config"
)

func NewClient(cfg *config.Config) (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}

	return client, nil
}
