-- +goose Up
-- +goose StatementBegin
-- Добавляем ограничение уникальности к полю name
ALTER TABLE users ADD CONSTRAINT unique_name UNIQUE (name);

-- Добавляем ограничение уникальности к полю email
ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаляем ограничение уникальности с поля name
ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_name;

-- Удаляем ограничение уникальности с поля email
ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_email;
-- +goose StatementEnd
