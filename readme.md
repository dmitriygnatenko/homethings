### Формат файла конфигурации приложения

config/.env

```
APP_PORT=8080

DB_HOST=localhost
DB_PORT=3306
DB_NAME=
DB_USER=
DB_PASSWORD=

CORS_ALLOW_ORIGING=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE

AUTH_USER=
AUTH_PASSWORD=
```

### Команды

- make run (запуск приложения)
- make test (запуск тестов)
- make test-cover (статистика по покрытию тестами)
- make swag (генерация документации)
- make lint (запуск линтера)
- make migration-status (статус миграций)
- make migration-up (применение миграций)
- make migration-down (откат миграций)
- make docker-build (сборка контейнеров)
- make docker-up (запуск контейнеров)
- make docker-down (остановка контейнеров)
- make install-deps (установка зависимостей)
- make app-build (компиляция приложения)


### Документация по методам API

/docs/