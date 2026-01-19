-- Добавление поля в таблицу files
BEGIN;


ALTER TABLE files ADD COLUMN size BIGINT NOT NULL DEFAULT 0;

COMMIT;