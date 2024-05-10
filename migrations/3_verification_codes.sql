-- +goose Up
CREATE TABLE verification_codes(
    code uuid PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id)  ON DELETE CASCADE
);

ALTER TABLE  users ADD COLUMN is_verified BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP COLUMN is_verified;
DROP TABLE verification_codes;