-- +goose Up
-- +goose StatementBegin

-- metadata column dropped in favour of attributes
ALTER TABLE component_firmware_set DROP COLUMN IF EXISTS metadata;

ALTER TABLE attributes ADD COLUMN component_firmware_set_id UUID REFERENCES component_firmware_set(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SET sql_safe_updates = false;

ALTER TABLE component_firmware_set ADD COLUMN metadata JSONB;

ALTER TABLE attributes DROP COLUMN component_firmware_set_id CASCADE;

-- +goose StatementEnd