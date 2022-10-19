--
-- +goose NO TRANSACTION
--  Once the predicate eval bug is fixed, this can be run in transaction
--    https://github.com/cockroachdb/cockroach/issues/79613
--
-- +goose Up
-- +goose StatementBegin


-- reference component firmware set id in attributes
ALTER TABLE attributes ADD COLUMN component_firmware_set_id UUID REFERENCES component_firmware_set(id) ON DELETE CASCADE;

--  indexes for component_firmware_set attributes
CREATE INDEX idx_component_firmware_set_id ON attributes (component_firmware_set_id) where component_firmware_set_id is not null;
CREATE INDEX idx_component_firmware_set_namespace ON attributes (component_firmware_set_id, namespace, created_at) where component_firmware_set_id is not null;
CREATE INVERTED INDEX idx_component_firmware_set_data ON attributes (component_firmware_set_id, namespace, data) where component_firmware_set_id is not null;

-- drop older constraint
ALTER TABLE attributes DROP CONSTRAINT check_server_id_server_component_id;

-- add constraint to ensure either one of these foreign keys are present
ALTER TABLE attributes ADD CONSTRAINT check_server_component_firmware_set_id CHECK (
  (
    (server_id is not null)::integer +
    (server_component_id is not null)::integer +
    (component_firmware_set_id is not null)::integer
  ) = 1
);

-- metadata column dropped in favour of attributes
ALTER TABLE component_firmware_set DROP COLUMN IF EXISTS metadata;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE attributes ADD CONSTRAINT check_server_id_server_component_id CHECK (
  (
    (server_id is not null)::integer +
    (server_component_id is not null)::integer
  ) = 1
);

ALTER TABLE attributes DROP CONSTRAINT IF EXISTS check_server_component_firmware_set_id;

DROP INDEX IF EXISTS attributes@idx_component_firmware_set_data;
DROP INDEX IF EXISTS attributes@idx_component_firmware_set_namespace;
DROP INDEX IF EXISTS attributes@idx_component_firmware_set_id;

SET sql_safe_updates = false;

ALTER TABLE component_firmware_set ADD COLUMN metadata JSONB;

ALTER TABLE attributes DROP COLUMN component_firmware_set_id CASCADE;

-- +goose StatementEnd