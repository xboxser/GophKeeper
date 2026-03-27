-- Создание таблицы для хранения информации по файлам
BEGIN;

CREATE TABLE files (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  uploaded_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_files_users 
    FOREIGN KEY (user_id) 
    REFERENCES users(id)
    ON UPDATE CASCADE
);


COMMIT;