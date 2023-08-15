-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin

ALTER TABLE component_firmware_set_map DROP CONSTRAINT fk_firmware_id_ref_component_firmware_version;
ALTER TABLE component_firmware_set_map ADD CONSTRAINT fk_firmware_id_ref_component_firmware_version FOREIGN KEY (firmware_id) REFERENCES component_firmware_version(id) ON DELETE RESTRICT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE component_firmware_set_map DROP CONSTRAINT fk_firmware_id_ref_component_firmware_version;
ALTER TABLE component_firmware_set_map ADD CONSTRAINT fk_firmware_id_ref_component_firmware_version FOREIGN KEY (firmware_id) REFERENCES component_firmware_version(id) ON DELETE CASCADE;

-- +goose StatementEnd