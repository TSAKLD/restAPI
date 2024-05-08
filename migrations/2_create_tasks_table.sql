-- +goose Up
CREATE TABLE tasks(
                      id BIGSERIAL PRIMARY KEY,
                      name TEXT NOT NULL,
                      project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                      user_id BIGINT NOT NULL  REFERENCES users(id) ON DELETE CASCADE,
                      created_at timestamptz NOT NULL
);

-- +goose Down
DROP TABLE tasks;