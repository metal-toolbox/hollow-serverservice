--
-- +goose Up
-- +goose StatementBegin


CREATE TABLE attributes_firmware_set (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  firmware_set_id UUID REFERENCES component_firmware_set(id) ON DELETE CASCADE,
  namespace STRING NOT NULL,
  data JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL,
  INDEX idx_firmware_set_id (firmware_set_id) WHERE firmware_set_id is not null,
  INVERTED INDEX idx_firmware_set_data (firmware_set_id, namespace, data) WHERE firmware_set_id is not null,
  UNIQUE INDEX idx_firmware_set_namespace (firmware_set_id, namespace) WHERE firmware_set_id is not null
);

-- metadata column dropped in favour of attributes
ALTER TABLE component_firmware_set DROP COLUMN IF EXISTS metadata;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


DROP TABLE firmware_set_attributes;
ALTER TABLE component_firmware_set ADD COLUMN metadata JSONB;
-- +goose StatementEnd