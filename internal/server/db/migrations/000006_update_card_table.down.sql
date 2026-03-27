-- Откатываем типы колонок к VARCHAR и устанавливаем пустые строки
BEGIN;

ALTER TABLE bank_cards
  ALTER COLUMN number_enc TYPE VARCHAR(19) USING '',
  ALTER COLUMN expiry_enc TYPE VARCHAR(10) USING '',
  ALTER COLUMN cvv_enc TYPE VARCHAR(4) USING '',
  ALTER COLUMN card_holder_name TYPE VARCHAR(255) USING '';
  
COMMIT;