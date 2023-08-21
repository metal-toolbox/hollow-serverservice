-- +goose Up
-- +goose StatementBegin

CREATE TABLE bom_info (
  serial_num STRING PRIMARY KEY NOT NULL,
  aoc_mac_address STRING NULL,
  bmc_mac_address STRING NULL,
  num_defi_pmi STRING NULL,
  num_def_pwd STRING NULL,
  metro STRING NULL
);

CREATE TABLE aoc_mac_address (
  aoc_mac_address STRING PRIMARY KEY NOT NULL,
  serial_num STRING NOT NULL REFERENCES bom_info(serial_num) ON DELETE CASCADE
);

CREATE TABLE bmc_mac_address (
  bmc_mac_address STRING PRIMARY KEY NOT NULL,
  serial_num STRING NOT NULL REFERENCES bom_info(serial_num) ON DELETE CASCADE
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE bom_info;
DROP TABLE aoc_mac_address;
DROP TABLE bmc_mac_address;

-- +goose StatementEnd