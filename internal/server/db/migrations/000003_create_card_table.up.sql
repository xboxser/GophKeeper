-- Создание таблицы для хранения банковских карт пользователей
BEGIN;

CREATE TABLE bank_cards (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  title VARCHAR(255) NOT NULL,
  number_enc VARCHAR(19) NOT NULL,
  expiry_enc VARCHAR(10) NOT NULL,
  cvv_enc VARCHAR(4) NOT NULL,
  card_holder_name VARCHAR(255) NOT NULL,
  last4 VARCHAR(4),

  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  uploaded_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_bank_cards_users 
    FOREIGN KEY (user_id) 
    REFERENCES users(id)
    ON UPDATE CASCADE
);

COMMIT;