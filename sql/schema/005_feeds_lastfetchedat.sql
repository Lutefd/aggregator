-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLUMN last_fetched_at;
