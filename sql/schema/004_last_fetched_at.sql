-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
INSERT INTO feeds DROP COLUMN last_fetched_at;

