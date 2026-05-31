CREATE SCHEMA chat;

CREATE TABLE chat.users(
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    full_name VARCHAR(100) NOT NULL,
    login VARCHAR(25) NOT NULL,
    password VARCHAR(100) NOT NULL,

    CONSTRAINT chk_full_name_length CHECK (char_length(full_name) BETWEEN 3 AND 100),
    CONSTRAINT chk_login_length CHECK (char_length(full_name) BETWEEN 3 AND 25),
    CONSTRAINT chk_password_length CHECK (char_length(password) > 0),

    CONSTRAINT uq_login UNIQUE (login)
);

CREATE TABLE chat.messages(
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    user_id INT,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES chat.users(id),
    text VARCHAR(2048) NOT NULL CHECK (char_length(text) BETWEEN 1 AND 2048),
    sent_at TIMESTAMP
);