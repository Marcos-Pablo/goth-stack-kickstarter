-- +goose Up
CREATE TABLE sessions (
  token TEXT PRIMARY KEY NOT NULL,
  data BLOB NOT NULL,
  expiry TIMESTAMP NOT NULL
);

CREATE INDEX idx_sessions_expiry ON sessions (expiry);

-- +goose Down
DROP INDEX idx_sessions_expiry;

DROP TABLE sessions;
