-- +goose Up
CREATE TABLE users (
    id integer primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    name text UNIQUE
);

-- +goose Down
DROP TABLE users;
