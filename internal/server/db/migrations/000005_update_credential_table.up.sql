-- Очистка данных и изменение типа на BYTEA для хранения зашифрованных паролей пользователя
BEGIN;
ALTER TABLE credentials ALTER COLUMN password DROP NOT NULL;
UPDATE credentials SET password = NULL;

ALTER TABLE credentials ALTER COLUMN password TYPE BYTEA USING NULL;


COMMIT;