CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL NOT NULL,
    uuid UUID NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT user_id_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS customer (
    id SERIAL NOT NULL,
    uuid UUID NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT customer_id_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS product (
    id SERIAL NOT NULL,
    uuid UUID NOT NULL,
    name VARCHAR(25) NOT NULL,
    observations VARCHAR(100),
    price INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    CONSTRAINT product_id_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS invoice_header (
    id SERIAL NOT NULL,
    uuid UUID NOT NULL,
    client_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    CONSTRAINT invoice_header_id_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS invoice_item (
    id SERIAL NOT NULL,
    uuid UUID NOT NULL,
    invoice_header_id INT NOT NULL,
    product_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    CONSTRAINT invoice_item_id_pk PRIMARY KEY (id),
    CONSTRAINT invoice_item_invoice_header_id_pk FOREIGN KEY (invoice_header_id) REFERENCES invoice_header (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT invoice_item_product_id_pk FOREIGN KEY (product_id) REFERENCES product (id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
