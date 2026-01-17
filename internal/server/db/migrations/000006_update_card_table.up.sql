-- Обновление тип данных для хранения банковских карт пользователей
BEGIN;

ALTER TABLE bank_cards
  ALTER COLUMN number_enc TYPE BYTEA USING ''::BYTEA,
  ALTER COLUMN expiry_enc TYPE BYTEA USING ''::BYTEA,
  ALTER COLUMN cvv_enc TYPE BYTEA USING ''::BYTEA,
  ALTER COLUMN card_holder_name TYPE BYTEA USING ''::BYTEA;

COMMIT;