-- +goose Up
-- +goose StatementBegin

CREATE TABLE public.servers (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NULL,
    facility_code STRING NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    INDEX idx_facility (facility_code ASC),
    FAMILY "primary" (id, name, facility_code, created_at, updated_at, deleted_at)
);

CREATE TABLE public.server_component_types (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    slug STRING NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE INDEX idx_name (name ASC),
    UNIQUE INDEX server_component_types_slug_key (slug ASC),
    FAMILY "primary" (id, name, created_at, updated_at, slug)
);

CREATE TABLE public.server_components (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NULL,
    vendor STRING NULL,
    model STRING NULL,
    serial STRING NULL,
    server_component_type_id UUID NOT NULL,
    server_id UUID NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    INDEX idx_server_component_type_id (server_component_type_id ASC),
    INDEX idx_server_id (server_id ASC),
    UNIQUE INDEX idx_server_components (server_id ASC, serial ASC, server_component_type_id ASC) WHERE ((server_id IS NOT NULL) AND (serial IS NOT NULL)) AND (server_component_type_id IS NOT NULL),
    FAMILY "primary" (id, name, vendor, model, serial, server_component_type_id, server_id, created_at, updated_at)
);

CREATE TABLE public.versioned_attributes (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    server_id UUID NULL,
    namespace STRING NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    -- ensure exactly one relationship is set
    CHECK (
        (
            (server_id is not null)::integer +
            (server_component_id is not null)::integer
        ) = 1
    ),
    server_component_id UUID NULL,
    tally INT8 NOT NULL DEFAULT 0:::INT8,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    INDEX idx_server_id (server_id ASC) WHERE server_id IS NOT NULL,
    INDEX idx_server_namespace (server_id ASC, namespace ASC, created_at ASC) WHERE server_id IS NOT NULL,
    INVERTED INDEX idx_server_data (server_id, namespace, data) WHERE server_id IS NOT NULL,
    INDEX idx_server_component_id (server_component_id ASC) WHERE server_component_id IS NOT NULL,
    INDEX idx_server_component_namespace (server_component_id ASC, namespace ASC, created_at ASC) WHERE server_component_id IS NOT NULL,
    INVERTED INDEX idx_server_component_data (server_component_id, namespace, data) WHERE server_component_id IS NOT NULL,
    FAMILY "primary" (id, server_id, namespace, data, created_at, updated_at, server_component_id, tally)
);

CREATE TABLE public.attributes (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    server_id UUID NULL,
    server_component_id UUID NULL,
    namespace STRING NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    -- ensure exactly one relationship is set
    CHECK (
        (
            (server_id is not null)::integer +
            (server_component_id is not null)::integer
        ) = 1
    ),
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    INDEX idx_server_id (server_id ASC) WHERE server_id IS NOT NULL,
    UNIQUE INDEX idx_server_namespace (server_id ASC, namespace ASC) WHERE server_id IS NOT NULL,
    INVERTED INDEX idx_server_data (server_id, namespace, data) WHERE server_id IS NOT NULL,
    INDEX idx_server_component_id (server_component_id ASC) WHERE server_component_id IS NOT NULL,
    UNIQUE INDEX idx_server_component_namespace (server_component_id ASC, namespace ASC) WHERE server_component_id IS NOT NULL,
    INVERTED INDEX idx_server_component_data (server_component_id, namespace, data) WHERE server_component_id IS NOT NULL,
    FAMILY "primary" (id, server_id, server_component_id, namespace, data, created_at, updated_at)
);

CREATE TABLE public.component_firmware_version (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    component STRING NOT NULL,
    vendor STRING NOT NULL,
    model STRING[] NOT NULL,
    filename STRING NOT NULL,
    version STRING NOT NULL,
    checksum STRING NOT NULL,
    upstream_url STRING NOT NULL,
    repository_url STRING NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE INDEX vendor_component_version_filename_unique (vendor ASC, component ASC, version ASC, filename ASC),
    FAMILY "primary" (id, component, vendor, model, filename, version, checksum, upstream_url, repository_url, created_at, updated_at)
);

CREATE TABLE public.server_credential_types (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    slug STRING NOT NULL,
    builtin BOOL NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE INDEX server_secret_types_slug_key (slug ASC),
    FAMILY "primary" (id, name, slug, builtin, created_at, updated_at)
);

CREATE TABLE public.server_credentials (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    server_id UUID NOT NULL,
    server_credential_type_id UUID NOT NULL,
    password STRING NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    username STRING NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE INDEX idx_server_secrets_by_type (server_id ASC, server_credential_type_id ASC),
    FAMILY "primary" (id, server_id, server_credential_type_id, password, created_at, updated_at, username)
);

CREATE TABLE public.component_firmware_set (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE INDEX idx_name (name ASC),
    FAMILY "primary" (id, name, created_at, updated_at)
);

CREATE TABLE public.component_firmware_set_map (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    firmware_set_id UUID NOT NULL,
    firmware_id UUID NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE INDEX component_firmware_set_map_firmware_set_id_firmware_id_key (firmware_set_id ASC, firmware_id ASC),
    FAMILY "primary" (id, firmware_set_id, firmware_id)
);

CREATE TABLE public.attributes_firmware_set (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    firmware_set_id UUID NULL,
    namespace STRING NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    INDEX idx_firmware_set_id (firmware_set_id ASC) WHERE firmware_set_id IS NOT NULL,
    INVERTED INDEX idx_firmware_set_data (firmware_set_id, namespace, data) WHERE firmware_set_id IS NOT NULL,
    UNIQUE INDEX idx_firmware_set_namespace (firmware_set_id ASC, namespace ASC) WHERE firmware_set_id IS NOT NULL,
    FAMILY "primary" (id, firmware_set_id, namespace, data, created_at, updated_at)
);

CREATE TABLE public.bom_info (
    serial_num STRING NOT NULL,
    aoc_mac_address STRING NULL,
    bmc_mac_address STRING NULL,
    num_defi_pmi STRING NULL,
    num_def_pwd STRING NULL,
    metro STRING NULL,
    CONSTRAINT "primary" PRIMARY KEY (serial_num ASC),
    FAMILY "primary" (serial_num, aoc_mac_address, bmc_mac_address, num_defi_pmi, num_def_pwd, metro)
);

CREATE TABLE public.aoc_mac_address (
    aoc_mac_address STRING NOT NULL,
    serial_num STRING NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (aoc_mac_address ASC),
    FAMILY "primary" (aoc_mac_address, serial_num)
);

CREATE TABLE public.bmc_mac_address (
    bmc_mac_address STRING NOT NULL,
    serial_num STRING NOT NULL,
    CONSTRAINT "primary" PRIMARY KEY (bmc_mac_address ASC),
    FAMILY "primary" (bmc_mac_address, serial_num)
);

INSERT INTO server_credential_types(name, slug, builtin, created_at, updated_at)
  VALUES ('BMC', 'bmc', true, now(), now());

ALTER TABLE public.server_components ADD CONSTRAINT fk_server_component_type_id_ref_server_component_types FOREIGN KEY (server_component_type_id) REFERENCES public.server_component_types(id);
ALTER TABLE public.server_components ADD CONSTRAINT fk_server_id_ref_servers FOREIGN KEY (server_id) REFERENCES public.servers(id) ON DELETE CASCADE;
ALTER TABLE public.versioned_attributes ADD CONSTRAINT fk_server_id_ref_servers FOREIGN KEY (server_id) REFERENCES public.servers(id) ON DELETE CASCADE;
ALTER TABLE public.versioned_attributes ADD CONSTRAINT fk_server_component_id_ref_server_components FOREIGN KEY (server_component_id) REFERENCES public.server_components(id) ON DELETE CASCADE;
ALTER TABLE public.attributes ADD CONSTRAINT fk_server_id_ref_servers FOREIGN KEY (server_id) REFERENCES public.servers(id) ON DELETE CASCADE;
ALTER TABLE public.attributes ADD CONSTRAINT fk_server_component_id_ref_server_components FOREIGN KEY (server_component_id) REFERENCES public.server_components(id) ON DELETE CASCADE;
ALTER TABLE public.server_credentials ADD CONSTRAINT fk_server_id_ref_servers FOREIGN KEY (server_id) REFERENCES public.servers(id) ON DELETE CASCADE;
ALTER TABLE public.server_credentials ADD CONSTRAINT fk_server_secret_type_id_ref_server_secret_types FOREIGN KEY (server_credential_type_id) REFERENCES public.server_credential_types(id);
ALTER TABLE public.component_firmware_set_map ADD CONSTRAINT fk_firmware_set_id_ref_component_firmware_set FOREIGN KEY (firmware_set_id) REFERENCES public.component_firmware_set(id) ON DELETE CASCADE;
ALTER TABLE public.component_firmware_set_map ADD CONSTRAINT fk_firmware_id_ref_component_firmware_version FOREIGN KEY (firmware_id) REFERENCES public.component_firmware_version(id) ON DELETE RESTRICT;
ALTER TABLE public.attributes_firmware_set ADD CONSTRAINT fk_firmware_set_id_ref_component_firmware_set FOREIGN KEY (firmware_set_id) REFERENCES public.component_firmware_set(id) ON DELETE CASCADE;
ALTER TABLE public.aoc_mac_address ADD CONSTRAINT fk_serial_num_ref_bom_info FOREIGN KEY (serial_num) REFERENCES public.bom_info(serial_num) ON DELETE CASCADE;
ALTER TABLE public.bmc_mac_address ADD CONSTRAINT fk_serial_num_ref_bom_info FOREIGN KEY (serial_num) REFERENCES public.bom_info(serial_num) ON DELETE CASCADE;

-- +goose StatementEnd
-- +goose Down
