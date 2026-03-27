-- создание таблицы users_counters
BEGIN;

CREATE TABLE counters (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  count INTEGER NOT NULL,
  type VARCHAR(4) NOT NULL,

  CONSTRAINT fk_counters_users 
    FOREIGN KEY (user_id) 
    REFERENCES users(id)
    ON UPDATE CASCADE
);

ALTER TABLE credentials ADD COLUMN inc_id BIGINT NOT NULL;
ALTER TABLE bank_cards ADD COLUMN inc_id BIGINT NOT NULL;

-- Добавление уникального ограничения
ALTER TABLE counters ADD CONSTRAINT unique_user_type_counter UNIQUE (user_id, type);

COMMIT;