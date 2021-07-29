-- +goose Up
-- +goose StatementBegin

ALTER TABLE versioned_attributes ADD COLUMN server_component_id UUID REFERENCES server_components(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SET sql_safe_updates = false;

ALTER TABLE versioned_attributes DROP COLUMN server_component_id CASCADE;

-- +goose StatementEnd
