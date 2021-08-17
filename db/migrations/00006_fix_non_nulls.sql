-- +goose Up
-- +goose StatementBegin

ALTER TABLE server_component_types ALTER COLUMN slug SET NOT NULL;
ALTER TABLE server_components ALTER COLUMN server_component_type_id SET NOT NULL;
ALTER TABLE server_components ALTER COLUMN server_id SET NOT NULL;
ALTER TABLE versioned_attributes ALTER COLUMN data SET NOT NULL;
ALTER TABLE versioned_attributes ALTER COLUMN tally SET DEFAULT 0;
ALTER TABLE versioned_attributes ALTER COLUMN tally SET NOT NULL;
ALTER TABLE attributes ALTER COLUMN data SET NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE server_component_types ALTER COLUMN slug DROP NOT NULL;
ALTER TABLE server_components ALTER COLUMN server_component_type_id DROP NOT NULL;
ALTER TABLE server_components ALTER COLUMN server_id DROP NOT NULL;
ALTER TABLE versioned_attributes ALTER COLUMN data DROP NOT NULL;
ALTER TABLE versioned_attributes ALTER COLUMN tally DROP NOT NULL;
ALTER TABLE attributes ALTER COLUMN data DROP NOT NULL;

-- +goose StatementEnd
