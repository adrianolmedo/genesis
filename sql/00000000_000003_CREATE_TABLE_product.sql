CREATE TABLE IF NOT EXISTS product (
    id BIGSERIAL,
    uuid UUID NOT NULL,
    name VARCHAR(25) NOT NULL,
    observations VARCHAR(100),
    price INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,

    CONSTRAINT product_id_pk PRIMARY KEY (id)
);
