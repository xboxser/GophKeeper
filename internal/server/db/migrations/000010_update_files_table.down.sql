-- Удаление поля status в таблице files
BEGIN;

ALTER TABLE files DROP COLUMN status;

COMMIT;