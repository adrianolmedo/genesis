CREATE TABLE IF NOT EXISTS product (
    id BIGSERIAL,
    uuid UUID NOT NULL,
    name TEXT NOT NULL,
    observations TEXT NOT NULL DEFAULT '',
    price BIGINT NOT NULL DEFAULT 0 CHECK (price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT product_id_pk PRIMARY KEY (id)
);
