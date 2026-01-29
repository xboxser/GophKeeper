-- удаление таблицы users_counters
BEGIN;
ALTER TABLE counters DROP CONSTRAINT unique_user_type_counter;
DROP TABLE counters;

ALTER TABLE bank_cards DROP COLUMN inc_id;
ALTER TABLE credentials DROP COLUMN inc_id;
COMMIT;