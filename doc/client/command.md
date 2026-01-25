# Команды проекта

___На текущий момент не реализованны команды update & delete___

Обработка команд происходит при помощи библиотеки `cobra`  

## Пользователь

### Регистрация пользователя

Команда: `registration`
```
registration -l=login-test -p=password -m=secret

  -l, --login string        Логин (Обязательный)   
  -m, --masterPass string   Пароль для шифрования (Обязательный. Не подлежит восстановлению!!!)    
  -p, --password string     Пароль (Обязательный)

```

### Авторизация пользователя


Команда: `login`
```
login -l=login-test -p=password

  -l, --login string        Логин (Обязательный)     
  -p, --password string     Пароль (Обязательный)
```

### Проверка мастер пароля
Проверяет корректность мастер пароля.   
Пользователь должен быть авторизованным в  системе

Команда: `master`
```
master -m=secret

  -m, --masterPass string   Пароль для шифрования (Обязательный. Не подлежит восстановлению!!!)   
```

## Учетные данные (Пары логин/пароль)

### Добавление учетной данных
Команда: `credential-add`
```
credential-add -l=login-test -p=password -m=secret

  -l, --login string        Логин (обязательный)
  -m, --masterPass string   Пароль для шифрования (Обязательный)
  -p, --password string     Пароль (обязательный)
```

### Получение списка учетных данных
Команда: `credential-get`
```
credential-get -m=secret

  -m, --masterPass string   Пароль для шифрования (Обязательный)
```


### Удаление элемента из учетных данных
Команда: `credential-del`
```
credential-get -m=secret -l=login-test

  -m, --masterPass string   Пароль для шифрования (Обязательный)
  -l, --login string        Логин (обязательный)
```

## Данные банковских карт

### Добавление банковской карты
Команда: `card-add`
```
card-add --number "1234 1234 1234 1234" --expiry "12/23" --card_holder "Имя и фамилия" --cvv "123" -t=test-card -m=secret 

  -o, --card_holder string   Имя держателя карты (обязательное)
  -c, --cvv string           cvv (обязательное)
  -e, --expiry string        Срок действия карты (обязательное)
  -m, --masterPass string    Пароль для шифрования (Обязательный)
  -n, --number string        Номер карты (обязательное)
  -t, --title string         Произвольное название карты (обязательное)
```

### Получение списка банковских карт
Команда: `card-get`
```
card-get -m=secret

  -m, --masterPass string   Пароль для шифрования (Обязательный)
```

## Файлы

### Добавление файла

Команда: `file-add`

```
file-add -m=secret -f=../../test_file/1GB_file

  -f, --file string         путь к файлу (обязательный) 
  -m, --masterPass string   Пароль для шифрования (Обязательный)
```

### Получение списка файлов
Команда: `file-list`

```
file-add 

Без параметров
```

### Скачивание файла
Команда: `file-get`

```
file-get -f=1GB_file

-f, --file string   наименование файла для скачивания с сервера (Обязательный)
```