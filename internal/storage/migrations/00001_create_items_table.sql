-- +goose Up
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    pub_date TIMESTAMP
);

-- +goose Down
DROP TABLE items;
