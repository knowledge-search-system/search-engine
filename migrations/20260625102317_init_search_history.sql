-- +goose Up
CREATE TABLE search_history (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    query       text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_search_history_created_at ON search_history (created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS search_history;
