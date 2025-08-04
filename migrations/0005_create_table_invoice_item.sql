-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS invoice_item (
    id BIGSERIAL,
    invoice_header_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT invoice_item_id_pk PRIMARY KEY (id),

    CONSTRAINT invoice_item_invoice_header_id_pk FOREIGN KEY (invoice_header_id)
        REFERENCES invoice_header (id) ON UPDATE RESTRICT ON DELETE RESTRICT,

    CONSTRAINT invoice_item_product_id_pk FOREIGN KEY (product_id)
        REFERENCES product (id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoice_item DROP CONSTRAINT IF EXISTS invoice_item_id_pk;
ALTER TABLE invoice_item DROP CONSTRAINT IF EXISTS invoice_item_invoice_header_id_pk;
ALTER TABLE invoice_item DROP CONSTRAINT IF EXISTS invoice_item_product_id_pk;
DROP TABLE IF EXISTS invoice_item;
-- +goose StatementEnd