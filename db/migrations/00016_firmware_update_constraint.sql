-- +goose Up
-- +goose StatementBegin

-- Not possible to drop unique constraint with ALTER TABLE
-- See https://github.com/cockroachdb/cockroach/issues/42840
DROP INDEX vendor_model_version_unique CASCADE;
ALTER TABLE component_firmware_version ADD CONSTRAINT vendor_component_version_unique UNIQUE (vendor, component, version);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX vendor_component_version_unique CASCADE;
ALTER TABLE component_firmware_version ADD CONSTRAINT vendor_model_version_unique UNIQUE (vendor, model, version);

-- +goose StatementEnd