-- Удаление поля в таблице files
BEGIN;

ALTER TABLE files DROP COLUMN size;

COMMIT;