-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin

-- ALTER COLUMN TYPE can't be used in a transaction, hence the NO TRANSACTION at the top of the file
SET enable_experimental_alter_column_type_general = true;
ALTER TABLE component_firmware_version ALTER model TYPE string array;
SET enable_experimental_alter_column_type_general = false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SET enable_experimental_alter_column_type_general = true;
ALTER TABLE component_firmware_version ALTER model TYPE string;
SET enable_experimental_alter_column_type_general = false;

-- +goose StatementEnd