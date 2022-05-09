-- +goose Up
-- +goose StatementBegin

CREATE TABLE component_firmware_version (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  component STRING NULL,
  vendor STRING NULL,
  model STRING NULL,
  filename STRING NULL,
  version STRING NULL,
  utility STRING NULL,
  sha STRING NULL,
  upstream_url STRING NULL,
  s3_url STRING NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE component_firmware_version;

-- +goose StatementEnd
