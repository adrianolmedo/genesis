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
