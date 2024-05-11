-- +goose Up
CREATE TABLE projects_users(
    project_id BIGINT REFERENCES projects(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, user_id)
);

-- +goose Down
DROP TABLE projects_users;