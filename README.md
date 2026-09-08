# Subscription Aggregation Service API

REST-сервис для агрегации данных об онлайн-подписках.

## Стек технологий

- Go 1.25
- Gin
- PostgreSQL 15
- golang-migrate
- Zap
- Swagger

## Запуск

```bash
docker compose up --build
```

Сервер будет доступен на `http://localhost:8080`.

## Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `SERVER_ADDRESS` | Адрес сервера | `0.0.0.0` |
| `SERVER_PORT` | Порт сервера | `8080` |
| `DATABASE_URL` | Строка подключения к PostgreSQL | `postgres://user:password@localhost:5432/subscriptions?sslmode=disable` |
| `LOG_LEVEL` | Уровень логирования | `info` |
| `MIGRATIONS_PATH` | Путь к файлам миграций | `./migrations` |

## API Endpoints

- `GET /health` — проверка работоспособности
- `POST /subscriptions` — создать подписку
- `GET /subscriptions` — получить список подписок
- `GET /subscriptions/{id}` — получить подписку по ID
- `PUT /subscriptions/{id}` — обновить подписку
- `DELETE /subscriptions/{id}` — удалить подписку
- `GET /subscriptions/total` — посчитать общую стоимость подписок за период

## Миграции

Миграции применяются автоматически при старте приложения через `golang-migrate`. SQL-файлы находятся в директории `migrations/`.

## Swagger

Документация доступна по адресу `http://localhost:8080/swagger/index.html` при запущенном сервисе.
