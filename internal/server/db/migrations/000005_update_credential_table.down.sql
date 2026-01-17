-- Откат правок для таблицы  для учетных данных пользователей
BEGIN;

ALTER TABLE credentials ALTER COLUMN password TYPE VARCHAR(60) USING NULL;

COMMIT;