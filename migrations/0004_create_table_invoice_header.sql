-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS invoice_header (
    id BIGSERIAL,
    uuid UUID NOT NULL,
    client_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT invoice_header_id_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoice_header DROP CONSTRAINT IF EXISTS invoice_header_id_pk;
DROP TABLE IF EXISTS invoice_header;
-- +goose StatementEnd