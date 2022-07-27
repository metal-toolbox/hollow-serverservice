-- +goose Up
-- +goose StatementBegin

CREATE TABLE server_secret_types (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  name STRING NOT NULL,
  slug STRING NOT NULL UNIQUE,
  builtin BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

INSERT INTO server_secret_types(name, slug, builtin, created_at, updated_at)
  VALUES ('BMC', 'bmc', true, now(), now());

CREATE TABLE server_secrets (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  server_secret_type_id UUID NOT NULL REFERENCES server_secret_types(id),
  value STRING NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE INDEX idx_server_secrets_by_type (server_id, server_secret_type_id)
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE server_secrets;
DROP TABLE server_secret_types;

-- +goose StatementEnd
