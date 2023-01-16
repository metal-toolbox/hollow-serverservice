-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin

-- Not possible to drop unique constraint with ALTER TABLE
-- See https://github.com/cockroachdb/cockroach/issues/42840
DROP INDEX vendor_model_version_unique CASCADE;

-- Before ALTER TABLE update the existing strings to look like an array.
UPDATE component_firmware_version SET model = regexp_replace(model, '(.*)', '{\1}') WHERE model IS NOT NULL;

-- ALTER COLUMN TYPE can't be used in a transaction, hence the NO TRANSACTION at the top of the file
SET enable_experimental_alter_column_type_general = true;
ALTER TABLE component_firmware_version ALTER model TYPE string array;
SET enable_experimental_alter_column_type_general = false;

ALTER TABLE component_firmware_version ADD CONSTRAINT vendor_component_version_unique UNIQUE (vendor, component, version);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX vendor_component_version_unique CASCADE;
SET enable_experimental_alter_column_type_general = true;
ALTER TABLE component_firmware_version ALTER model TYPE string;
SET enable_experimental_alter_column_type_general = false;
ALTER TABLE component_firmware_version ADD CONSTRAINT vendor_model_version_unique UNIQUE (vendor, model, version);

-- +goose StatementEnd