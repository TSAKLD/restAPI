-- +goose Up
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    password   TEXT        NOT NULL,
    email      TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE sessions(
    id uuid PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL
);

CREATE TABLE projects(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL
);

-- +goose Down
DROP TABLE projects;
DROP TABLE sessions;
DROP TABLE users;