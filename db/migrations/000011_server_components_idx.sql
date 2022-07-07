-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_server_components ON server_components (server_id, serial, server_component_type_id) WHERE server_id IS NOT NULL AND serial IS NOT NULL AND server_component_type_id IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS server_components@idx_server_components CASCADE;

-- +goose StatementEnd