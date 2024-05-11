-- +goose Up
ALTER TABLE tasks ADD COLUMN description TEXT NOT NULL DEFAULT 'New Task';

-- +goose Down
ALTER TABLE tasks DROP COLUMN description;