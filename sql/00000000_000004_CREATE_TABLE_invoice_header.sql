CREATE TABLE IF NOT EXISTS invoice_header (
    id BIGSERIAL,
    uuid UUID NOT NULL,
    client_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT invoice_header_id_pk PRIMARY KEY (id)
);
