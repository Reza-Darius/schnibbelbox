-- +goose Up
CREATE TABLE IF NOT EXISTS test (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT
);

-- +goose Down
SELECT 'down SQL query';
