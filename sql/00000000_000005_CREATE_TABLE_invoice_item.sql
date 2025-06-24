CREATE TABLE IF NOT EXISTS invoice_item (
    id BIGSERIAL,
    invoice_header_id INT NOT NULL,
    product_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT invoice_item_id_pk PRIMARY KEY (id),

    CONSTRAINT invoice_item_invoice_header_id_pk FOREIGN KEY (invoice_header_id)
        REFERENCES invoice_header (id) ON UPDATE RESTRICT ON DELETE RESTRICT,

    CONSTRAINT invoice_item_product_id_pk FOREIGN KEY (product_id)
        REFERENCES product (id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
