-- +goose NO TRANSACTION 
-- +goose Up
-- +goose StatementBegin
ALTER TABLE versioned_attributes DROP CONSTRAINT IF EXISTS check_server_id_server_component_id;
ALTER TABLE versioned_attributes ADD CONSTRAINT check_server_id_server_component_id CHECK (
  (
    (server_id is not null)::integer +
    (server_component_id is not null)::integer
  ) = 2
);

ALTER TABLE attributes DROP CONSTRAINT IF EXISTS check_server_id_server_component_id;
ALTER TABLE attributes ADD CONSTRAINT check_server_id_server_component_id CHECK (
  (
    (server_id is not null)::integer +
    (server_component_id is not null)::integer
  ) = 2
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE versioned_attributes DROP CONSTRAINT IF EXISTS check_server_id_server_component_id;
ALTER TABLE versioned_attributes ADD CONSTRAINT check_server_id_server_component_id CHECK (
  (
    (server_id is not null)::integer +
    (server_component_id is not null)::integer
  ) = 1
);


ALTER TABLE attributes DROP CONSTRAINT IF EXISTS check_server_id_server_component_id;
ALTER TABLE attributes ADD CONSTRAINT check_server_id_server_component_id CHECK (
  (
    (server_id is not null)::integer +
    (server_component_id is not null)::integer
  ) = 1
);


-- +goose StatementEnd