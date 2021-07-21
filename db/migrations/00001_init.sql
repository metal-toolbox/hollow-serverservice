-- +goose Up
-- +goose StatementBegin

CREATE TABLE servers (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  name STRING NULL,
  facility_code STRING NULL,
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL,
  INDEX idx_facility (facility_code)
);

CREATE TABLE server_component_types (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  name STRING NOT NULL,
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL,
  UNIQUE INDEX idx_name (name)
);

CREATE TABLE server_components (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  name STRING NULL,
  vendor STRING NULL,
  model STRING NULL,
  serial STRING NULL,
  server_component_type_id UUID NOT NULL REFERENCES server_component_types(id),
  server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL,
  INDEX idx_server_component_type_id (server_component_type_id),
  INDEX idx_server_id (server_id)
);

CREATE TABLE versioned_attributes (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  server_id UUID NULL REFERENCES servers(id) ON DELETE CASCADE,
  namespace STRING NOT NULL,
  data JSONB NULL,
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL,
  -- ensure exactly one relationship is set
  CHECK (
    (
      (server_id is not null)::integer +
      0
    ) = 1
  ),
  INDEX idx_server_id (server_id) where server_id is not null,
  INDEX idx_server_namespace (server_id, namespace, created_at) where server_id is not null,
  INVERTED INDEX idx_server_data (server_id, namespace, data) where server_id is not null
);

CREATE TABLE attributes (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  server_id UUID NULL REFERENCES servers(id) ON DELETE CASCADE,
  server_component_id UUID REFERENCES server_components(id) ON DELETE CASCADE,
  namespace STRING NOT NULL,
  data JSONB NULL,
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL,
  -- ensure exactly one relationship is set
  CHECK (
    (
      (server_id is not null)::integer +
      (server_component_id is not null)::integer
    ) = 1
  ),
  INDEX idx_server_id (server_id) where server_id is not null,
  UNIQUE INDEX idx_server_namespace (server_id, namespace) WHERE server_id is not null,
  INVERTED INDEX idx_server_data (server_id, namespace, data) WHERE server_id is not null,
  INDEX idx_server_component_id (server_component_id) WHERE server_component_id is not null,
  UNIQUE INDEX idx_server_component_namespace (server_component_id, namespace) WHERE server_component_id is not null,
  INVERTED INDEX idx_server_component_data (server_component_id, namespace, data) WHERE server_component_id is not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE attributes;
DROP TABLE versioned_attributes;
DROP TABLE server_components;
DROP TABLE server_component_types;
DROP TABLE servers;

-- +goose StatementEnd
