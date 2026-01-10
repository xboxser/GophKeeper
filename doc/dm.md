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
		datetime created_at  ""  
		datetime uploaded_at  ""  
	}
    bank_cards {
        int id PK ""  
        int user_id FK ""  
        title string  ""
        number_enc string  "номер карты"  
        expiry_enc string  "срок действия"
        cvv_enc string  "cvv"
        card_holder_name string  "имя держателя карты"
        last4 string  "последние 4 цифры"
    }
    file {
        int id PK ""
        int user_id FK ""
        string path "путь до файла на сервере"
    }
	credentials}|--||users:"  "
    bank_cards}|--||users:"  "
    file}|--||users:"  "

```