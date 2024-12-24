# GoAuthService

**GoAuthService** — это пример реализации части сервиса аутентификации, разработанного на языке Go. Сервис предоставляет два основных REST-метода для управления токенами аутентификации.

## Основные возможности

- **Access Token**:
  - Формат: JWT (HS512)
  - Не хранится в базе данных
  - Payload содержит: `user_id`, `ip_address`, время истечения

- **Refresh Token**:
  - Формат: base64-строка
  - В базе хранится только bcrypt-хеш
  - Одноразовый: после использования помечается как использованный
  - Связан с Access Token, можно обновлять только с тем Refresh Token, который был выдан вместе с ним
  - Payload также содержит `ip_address`

- **Безопасность**:
  - При изменении IP-адреса клиента при обновлении токена отправляется предупреждающее email-сообщение (мок-имитация отправки)

- **Инфраструктура**:
  - Используется PostgreSQL для хранения информации о пользователях и Refresh токенах
  - Реализовано Docker-окружение для быстрого запуска сервиса и базы данных
  - Миграции базы данных с помощью инструмента `goose` (встроенные миграции через `embed`)
  - Структурированное логирование через `zerolog`
  - Примеры unit-тестов (можно расширять)

## Используемый стек технологий

- **Язык**: Go (1.20+)
- **Веб-фреймворк**: [gorilla/mux](https://github.com/gorilla/mux)
- **БД**: PostgreSQL
- **JWT**: [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)
- **Хэширование**: [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **Миграции**: [goose](https://github.com/pressly/goose) + `embed`
- **Логирование**: [zerolog](https://github.com/rs/zerolog)
- **Контейнеризация**: Docker, Docker Compose
- **Тестирование**: Стандартный `testing` пакет Go + [testify/mock](https://github.com/stretchr/testify)

## Локальный запуск

### Требования

- [Docker](https://www.docker.com/get-started) и Docker Compose
- Go 1.20+ (если хотите запускать без Docker)
- `make` (опционально, для упрощения команд)

### Запуск через Docker Compose

1. **Склонируйте репозиторий:**

    ```bash
    git clone https://github.com/n-mukhin/GoAuthService.git
    cd GoAuthService
    ```

2. **Запустите сервис:**

    ```bash
    docker-compose up --build
    ```

    Docker Compose поднимет:

    - Контейнер с PostgreSQL
    - Контейнер с вашим Go-приложением `AuthService`

3. **После успешного запуска сервис будет доступен на порту `:8080`.**

## Открытые эндпоинты

### Получение пары токенов

- **Endpoint:**

    ```http
    GET http://localhost:8080/auth/token?user_id=11111111-1111-1111-1111-111111111111
    ```

- **Ответ:**

    ```json
    {
      "access_token": "<JWT>",
      "refresh_token": "<BASE64>"
    }
    ```

### Обновление токенов (Refresh Tokens)

- **Endpoint:**

    ```http
    POST http://localhost:8080/auth/refresh
    Content-Type: application/json
    ```

- **Тело запроса:**

    ```json
    {
      "access_token": "<OLD_JWT>",
      "refresh_token": "<OLD_BASE64>"
    }
    ```

- **Ответ:**

    ```json
    {
      "access_token": "<NEW_JWT>",
      "refresh_token": "<NEW_BASE64>"
    }
    ```

## Изменение настроек

Все настройки задаются через переменные окружения. Изменяются в файле `docker-compose.yml` или в файле `.env`.

### Основные переменные окружения:

- `DB_HOST`: Хост базы данных (по умолчанию `db`)
- `DB_PORT`: Порт базы данных (по умолчанию `5432`)
- `DB_USER`: Пользователь базы данных (по умолчанию `user`)
- `DB_PASSWORD`: Пароль базы данных (по умолчанию `pass`)
- `DB_NAME`: Имя базы данных (по умолчанию `authdb`)
- `JWT_SECRET`: Секрет для подписи JWT (например, `supersecretkey`)
- `SERVER_ADDR`: Адрес и порт сервера (по умолчанию `:8080`)
- `EMAIL_SENDER`: Адрес отправителя email (мок-отправка, например, `mock@example.com`)

