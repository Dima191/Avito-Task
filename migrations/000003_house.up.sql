CREATE TABLE houses (
    house_id                BIGINT PRIMARY KEY,
    address                 varchar(255) NOT NULL UNIQUE,
    year                    SMALLINT     NOT NULL,
    developer               varchar(255),
    created_at              TIMESTAMP DEFAULT NOW(),
    last_apartment_added_at TIMESTAMP
);