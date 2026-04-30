-- +goose Up
CREATE TABLE users (
  id TEXT PRIMARY KEY NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  password TEXT NOT NULL,
  profile_picture_url TEXT
);

-- +goose Down
DROP TABLE users;
