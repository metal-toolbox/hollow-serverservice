-- +goose Up
-- +goose StatementBegin

CREATE TABLE firmwares (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  component_id UUID NOT NULL REFERENCES server_components(id) ON DELETE CASCADE,
  vendor STRING NULL,
  model STRING NULL,
  filename STRING NULL,
  version STRING NULL,
  utility STRING NULL,
  sha STRING NULL,
  upstream_url STRING NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE firmwares;

-- +goose StatementEnd
