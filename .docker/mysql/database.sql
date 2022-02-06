CREATE DATABASE IF NOT EXISTS go_practice_restapi;

USE go_practice_restapi;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT,
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
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    name VARCHAR(25) NOT NULL,
    observations VARCHAR(100),
    price INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS invoice_headers (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    client_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS invoice_items (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    invoice_header_id INT NOT NULL,
    product_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    CONSTRAINT invoice_items_invoice_header_id_pk FOREIGN KEY (invoice_header_id) REFERENCES invoice_headers (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT invoice_items_product_id_pk FOREIGN KEY (product_id) REFERENCES products (id) ON UPDATE RESTRICT ON DELETE RESTRICT
);