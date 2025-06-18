-- +goose Up
CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    title TEXT,
    language TEXT,
    description TEXT
);

ALTER TABLE items
ADD COLUMN channel_id INTEGER REFERENCES channels(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE items
DROP COLUMN channel_id;

DROP TABLE channels;
