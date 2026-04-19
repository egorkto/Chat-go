CREATE SCHEMA chat;

CREATE TABLE chat.users(
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    full_name VARCHAR(20) NOT NULL UNIQUE CHECK (char_length(full_name) BETWEEN 3 AND 20),
    password VARCHAR(100) NOT NULL
);