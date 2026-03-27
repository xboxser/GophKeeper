-- Создание таблицы для хранения пользователей
BEGIN;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  login VARCHAR(255) NOT NULL,
  password VARCHAR(60) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  uploaded_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;