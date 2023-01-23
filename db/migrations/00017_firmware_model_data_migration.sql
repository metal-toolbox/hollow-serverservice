-- +goose Up
-- +goose StatementBegin

-- Before ALTER TABLE in migration 00018, update the existing strings to look like an array.
UPDATE component_firmware_version SET model = regexp_replace(model, '(.*)', '{\1}') WHERE model IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE component_firmware_version SET model = btrim(model, '{}') WHERE model IS NOT NULL;

-- +goose StatementEnd