-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_deleted_at;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE INDEX idx_deleted_at ON servers (deleted_at) where deleted_at is null;
-- +goose StatementEnd
