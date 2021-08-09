-- +goose Up
-- +goose StatementBegin

ALTER TABLE versioned_attributes ADD COLUMN tally INTEGER;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SET sql_safe_updates = false;
ALTER TABLE versioned_attributes DROP COLUMN tally CASCADE;

-- +goose StatementEnd
