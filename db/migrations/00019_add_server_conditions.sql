-- +goose Up
-- +goose StatementBegin

-- server_condition_types declares the possible conditions that can be associated with server
CREATE TABLE server_condition_types (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  slug STRING NOT NULL,
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL,
  UNIQUE INDEX idx_slug (slug)
);


--- predefine server condition types
INSERT INTO server_condition_types(slug, created_at, updated_at)
  VALUES
    ('firmwareUpdate', now(), now()),
    ('inventoryOutofband', now(), now());


-- server_condition_status_types declares the possible condition statuses the server condition can have.
CREATE TABLE server_condition_status_types (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  slug STRING NOT NULL,
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL,
  UNIQUE INDEX idx_slug (slug)
);

-- predefine server condition status values.
INSERT INTO server_condition_status_types(slug, created_at, updated_at)
  VALUES
    ('pending', now(), now()),
    ('active', now(), now()),
    ('succeeded', now(), now()),
    ('failed', now(), now());


-- server_conditions enables a server_condition_type and server_condition to be associated with a server.
CREATE TABLE server_conditions (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  server_condition_type_id UUID NOT NULL REFERENCES server_condition_types(id),
  server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  server_condition_status_type_id UUID NOT NULL REFERENCES server_condition_status_types(id),
  parameters JSONB NOT NULL,    -- condition input parameters
  status_output JSONB NOT NULL,     -- condition controller output
  created_at TIMESTAMPTZ NULL,
  updated_at TIMESTAMPTZ NULL,
  UNIQUE (server_id, server_condition_type_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

SET sql_safe_updates = false;

DROP TABLE server_condition_types;
DROP TABLE server_condition_status_types;
DROP TABLE server_condtions;

-- +goose StatementEnd
