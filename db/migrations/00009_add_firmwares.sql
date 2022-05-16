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
  s3_url STRING NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE component_firmware_version;

-- +goose StatementEnd
