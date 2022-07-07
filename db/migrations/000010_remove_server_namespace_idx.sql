-- +goose NO TRANSACTION 
-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS attributes@idx_server_namespace CASCADE;
DROP INDEX IF EXISTS versioned_attributes@idx_server_namespace CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS attributes@idx_server_namespace CASCADE;
-- "pq: unimplemented: cannot drop UNIQUE constraint \"idx_server_namespace\" using ALTER TABLE DROP CONSTRAINT, use DROP INDEX CASCADE instead
DROP INDEX IF EXISTS versioned_attributes@idx_server_namespace CASCADE;

CREATE UNIQUE INDEX idx_server_namespace ON versioned_attributes (server_id, namespace, created_at) where server_id is not null;
CREATE INDEX idx_server_namespace ON attributes (server_id, namespace, created_at) where server_id is not null;

-- +goose StatementEnd