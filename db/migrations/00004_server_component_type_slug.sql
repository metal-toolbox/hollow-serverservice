-- +goose Up
-- +goose StatementBegin

ALTER TABLE server_component_types ADD COLUMN slug STRING UNIQUE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SET sql_safe_updates = false;
ALTER TABLE server_component_types DROP COLUMN slug CASCADE;

-- +goose StatementEnd
