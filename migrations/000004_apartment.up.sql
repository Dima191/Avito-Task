CREATE TYPE moderation_status AS ENUM ('created', 'approved', 'declined', 'on moderation');
CREATE TABLE apartments (
    apartment_id      BIGINT PRIMARY KEY,
    apartment_number  INT NOT NULL,
    house_id          BIGINT NOT NULL REFERENCES houses (house_id),
    price             INT NOT NULL,
    number_of_rooms   INT NOT NULL,
    moderation_status moderation_status
);

CREATE INDEX ON apartments(moderation_status);

CREATE OR REPLACE FUNCTION update_last_apartment_added_at()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE houses SET last_apartment_added_at = NOW() WHERE house_id = NEW.house_id;
    RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER apartment_added
    AFTER INSERT
    ON apartments
    FOR EACH ROW
    EXECUTE FUNCTION update_last_apartment_added_at();

