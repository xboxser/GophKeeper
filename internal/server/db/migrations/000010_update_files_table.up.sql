-- Добавление поля status в таблицу files
BEGIN;

ALTER TABLE files ADD COLUMN status VARCHAR(4);

COMMIT;