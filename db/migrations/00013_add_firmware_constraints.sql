-- +goose Up
-- +goose StatementBegin

ALTER TABLE component_firmware_version ADD CONSTRAINT vendor_model_version_unique UNIQUE (vendor, model, version);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Not possible to drop unique constraint with ALTER TABLE
-- See https://github.com/cockroachdb/cockroach/issues/42840
DROP INDEX vendor_model_version_unique CASCADE;

-- +goose StatementEnd
