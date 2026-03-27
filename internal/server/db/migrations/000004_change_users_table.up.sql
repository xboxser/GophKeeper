-- Добавление нового поля в таблицу users
BEGIN;

ALTER TABLE users 
ADD COLUMN code BYTEA;

COMMIT;