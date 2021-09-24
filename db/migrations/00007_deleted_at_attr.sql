-- +goose Up
-- +goose StatementBegin
ALTER TABLE servers ADD COLUMN deleted_at TIMESTAMPTZ NULL;
CREATE INDEX idx_deleted_at ON servers (deleted_at) where deleted_at is null;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE servers DROP COLUMN deleted_at;
DROP INDEX idx_deleted_at;
-- +goose StatementEnd
