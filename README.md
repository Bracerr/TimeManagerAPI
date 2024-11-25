# Настройка переменных окружения

## Для корректной работы приложения необходимо установить следующие переменные окружения в корневой папке:

## Обязательные переменные

- ### `MONGO_URI`: URL для подключения к базе данных.
- ### `MONGO_DB_NAME`: Название БД в Монго.
- ### `JWT_SIGNING_KEY`: Ключ для шифровки Jwt Token.
- ### `JWT_EXPIRATION_HOUR`: Время истечения токена в часах.
## Пример .env для Docker-compose
```plaintext
MONGO_URI='mongodb://mongo:27017'
MONGO_DB_NAME='JWT_BASE'

JWT_SIGNING_KEY="C9kT7hzVbeSlth5avv5ihjSziZ5v23W0f+bCk+hTMOA=="
JWT_EXPIRATION_HOUR="24"
```

# Запуск
```bash
docker-compose up --build -d
```