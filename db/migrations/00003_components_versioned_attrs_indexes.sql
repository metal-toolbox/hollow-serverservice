-- The inverted index would fail because the transaction adding the new column needed to complete before they could be added
-- so moved all the other changes outside of the table change into it's own migration

-- +goose Up
-- +goose StatementBegin

CREATE INDEX idx_server_component_id ON versioned_attributes (server_component_id) where server_component_id is not null;
CREATE INDEX idx_server_component_namespace ON versioned_attributes (server_component_id, namespace, created_at) where server_component_id is not null;
CREATE INVERTED INDEX idx_server_component_data ON versioned_attributes (server_component_id, namespace, data) where server_component_id is not null;

ALTER TABLE versioned_attributes ADD CONSTRAINT check_server_id_server_component_id CHECK (
  (
    (server_id is not null)::integer +
    (server_component_id is not null)::integer
  ) = 1
);

ALTER TABLE versioned_attributes DROP CONSTRAINT check_server_id;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE versioned_attributes ADD CONSTRAINT check_server_id CHECK (
  (
    (server_id is not null)::integer +
    0
  ) = 1
);

ALTER TABLE versioned_attributes DROP CONSTRAINT check_server_id_server_component_id;

DROP INDEX versioned_attributes@idx_server_component_data;
DROP INDEX versioned_attributes@idx_server_component_namespace;
DROP INDEX versioned_attributes@idx_server_component_id;

-- +goose StatementEnd
