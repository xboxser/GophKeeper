# Схема БД

```mermaid
---
config:
  theme: neutral
---
erDiagram
	direction TB
	credentials {
		int id PK ""  
        int user_id FK ""  
		string password  ""  
        string login  ""  
		datetime created_at  ""  
		datetime uploaded_at  ""  
	}
	users {
		int id PK ""  
		string login  ""  
		string password  "hash BCRYPT"  
        code BYTEA  ""  
		datetime created_at  ""  
		datetime uploaded_at  ""  
	}
    bank_cards {
        int id PK ""  
        int user_id FK ""  
        title string  ""
        number_enc BYTEA  "номер карты"  
        expiry_enc BYTEA  "срок действия"
        cvv_enc BYTEA  "cvv"
        card_holder_name BYTEA  "имя держателя карты"
        last4 string  "последние 4 цифры"
    }
    files {
        int id PK ""
        int user_id FK ""
        string name "имя файла"
        int size "размер файла"
        datetime created_at  ""  
		datetime uploaded_at  ""  
    }
	credentials}|--||users:"  "
    bank_cards}|--||users:"  "
    files}|--||users:"  "

```