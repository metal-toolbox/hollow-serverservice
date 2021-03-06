-- +goose Up
-- +goose StatementBegin

CREATE TABLE component_firmware_version (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  component STRING NOT NULL,
  vendor STRING NOT NULL,
  model STRING NOT NULL,
  filename STRING NOT NULL,
  version STRING NOT NULL,
  checksum STRING NOT NULL,
  upstream_url STRING NOT NULL,
  repository_url STRING NOT NULL,
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE component_firmware_version;

-- +goose StatementEnd
