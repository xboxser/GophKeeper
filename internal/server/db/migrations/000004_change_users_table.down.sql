-- Удаляем колонку code из таблицы users
BEGIN;

ALTER TABLE users 
DROP COLUMN code;

COMMIT;