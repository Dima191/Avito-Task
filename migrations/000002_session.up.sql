CREATE TABLE sessions (
    session_id         BIGINT PRIMARY KEY,
    user_id            BIGINT REFERENCES users (user_id),
    hash_refresh_token VARCHAR(255) NOT NULL UNIQUE,
    expires_at         TIMESTAMP    NOT NULL
);