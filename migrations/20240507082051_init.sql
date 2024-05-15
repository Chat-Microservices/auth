-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access (
  id SERIAL PRIMARY KEY,
  role_id INTEGER REFERENCES roles(id),
  path VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access;
-- +goose StatementEnd
