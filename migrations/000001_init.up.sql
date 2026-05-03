CREATE SCHEMA chat;

CREATE TABLE chat.users(
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    full_name VARCHAR(100) NOT NULL CHECK (char_length(full_name) BETWEEN 3 AND 100),
    login VARCHAR(25) NOT NULL UNIQUE CHECK (char_length(login) BETWEEN 3 AND 25),
    password VARCHAR(100) NOT NULL CHECK (char_length(password) > 0)
);