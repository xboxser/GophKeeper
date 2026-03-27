-- Создание таблицы для хранения учетных данных пользователей
BEGIN;

CREATE TABLE credentials (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  login VARCHAR(255) NOT NULL,
  password VARCHAR(60) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  uploaded_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_credentials_users 
    FOREIGN KEY (user_id) 
    REFERENCES users(id)
    ON UPDATE CASCADE
);

COMMIT;