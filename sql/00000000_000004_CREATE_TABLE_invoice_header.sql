CREATE TABLE IF NOT EXISTS invoice_header (
    id BIGSERIAL,
    uuid UUID NOT NULL,
    cliente_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,

    CONSTRAINT invoice_header_id_pk PRIMARY KEY (id)
);
