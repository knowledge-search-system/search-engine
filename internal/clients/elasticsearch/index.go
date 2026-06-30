package elasticsearch

import (
	"context"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const documentsIndexMapping = `{
  "settings": {
    "analysis": {
      "filter": {
        "russian_stop": {
          "type": "stop",
          "stopwords": "_russian_"
        },
        "russian_stemmer": {
          "type": "stemmer",
          "language": "russian"
        }
      },
      "analyzer": {
        "analysis-ru": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "russian_stop", "russian_stemmer"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "chunk_id":    { "type": "keyword" },
      "document_id": { "type": "keyword" },
      "file_name":   { "type": "keyword" },
      "page_number": { "type": "integer" },
      "text":        { "type": "text", "analyzer": "analysis-ru" }
    }
  }
}`

func EnsureDocumentsIndex(ctx context.Context, client *elasticsearch.Client, indexName string) error {
	existsResp, err := esapi.IndicesExistsRequest{Index: []string{indexName}}.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("check documents index exists: %w", err)
	}
	defer existsResp.Body.Close()

	if existsResp.StatusCode == 200 {
		return nil
	}

	createResp, err := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(documentsIndexMapping),
	}.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("create documents index: %w", err)
	}
	defer createResp.Body.Close()

	if createResp.IsError() {
		return fmt.Errorf("create documents index: status %s", createResp.Status())
	}

	return nil
}
