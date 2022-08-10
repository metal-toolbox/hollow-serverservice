-- +goose Up
-- +goose StatementBegin

ALTER TABLE server_secret_types RENAME TO server_credential_types;
ALTER TABLE server_secrets RENAME TO server_credentials;
ALTER TABLE server_credentials ADD COLUMN username string NOT NULL;
ALTER TABLE server_credentials RENAME value TO password;
ALTER TABLE server_credentials RENAME server_secret_type_id TO server_credential_type_id;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE server_credentials;
DROP TABLE server_credential_types;

-- +goose StatementEnd