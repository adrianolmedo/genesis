CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    uuid char(36) NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT users_id_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL NOT NULL,
    name VARCHAR(25) NOT NULL,
    observations VARCHAR(100),
    price INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    CONSTRAINT products_id_pk PRIMARY KEY (id)
);