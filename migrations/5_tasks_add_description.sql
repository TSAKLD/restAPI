-- +goose Up
ALTER TABLE tasks ADD COLUMN description TEXT NOT NULL DEFAULT 'New Task';

-- +goose Down
DROP TABLE projects_users;