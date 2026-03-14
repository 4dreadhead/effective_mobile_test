# Сервис агрегации подписок

Сервис предоставляет CRUD-операции по подпискам через `/subscriptions` и расчёт суммы
за период через `/subscriptions/total`. Для всех запросов требуется заголовок
`Authorization` со значением `API_KEY`, поля `from`/`to` принимаются в формате
`MM.YYYY`.

Структура кода: точка входа — `cmd/server`, бизнес-логика и слои данных - в
`internal/subscriptions` (model/usecase/presenter/controller/view), инфраструктура -
в `internal/platform` (auth, db, http). Миграции находятся в `migrations/`.

## Сборка образа / приложения

- Docker-образ:
  `docker build -t subscription-service .`
- Локальный бинарник:
  `go build -o bin/server ./cmd/server`

## Настройка окружения

1. Скопируйте пример файла окружения:
   `cp .env.example .env`
2. Отредактируйте `.env` и задайте:
   - `DATABASE_DSN`
   - `API_KEY`
   - Опционально: `HTTP_ADDRESS`
3. Если используете Docker Compose, скопируйте:
   `cp docker-compose.yml.example docker-compose.yml`

## Запуск проекта

- Docker Compose:
  `docker compose up --build`
- Локальный бинарник:
  `./bin/server`

## Документация API

- Swagger UI: `http://localhost:8080/swagger/index.html`
- Все эндпоинты требуют заголовок `Authorization` со значением `API_KEY`.

## Тесты

Для use cases написаны тесты:

`go test ./internal/subscriptions/usecase/...`
