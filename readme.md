### Формат файла конфигурации

.env

```
APP_PORT=8080

DB_HOST=localhost
DB_PORT=3306
DB_NAME=homethings
DB_USER=user
DB_PASSWORD=pass
DB_MAX_OPEN_CONNS=0
DB_MAX_IDLE_CONNS=5
DB_MAX_CONN_LIFETIME=0
DB_MAX_IDLE_CONN_LIFETIME=60

CORS_ALLOW_ORIGING=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE

JWT_SECRET_KEY=secret
JWT_LIFETIME=21600

SMTP_HOST=smtp.example.ru
SMTP_PORT=2525
SMTP_USER=user@example.ru
SMTP_PASSWORD=pass

ERRORS_EMAIL=user@example.ru
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