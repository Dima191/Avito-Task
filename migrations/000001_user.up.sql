CREATE TYPE user_role AS ENUM ('client','moderator');
CREATE TABLE users (
    user_id       BIGINT PRIMARY KEY,
    role          user_role    NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    hash_password VARCHAR(255) NOT NULL
);