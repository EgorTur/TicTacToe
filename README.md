# Крестики-нолики (Tic-Tac-Toe Backend)
Учебный проект School 21 (Sber) - бэкенд-сервис для игры в крестики-нолики с ИИ и многопользовательским режимом.

## Возможности
- Игра с компьютером (алгоритм minimax) и режим «игрок против игрока».
- Валидация ходов и проверка завершения игры.
- Одновременная работа с несколькими игровыми сессиями.
- Регистрация и аутентификация пользователей (JWT: access + refresh токены).
- История завершённых игр и таблица лидеров (топ N по проценту побед).
- Чистая архитектура (Domain, Application, Datasource, Web, DI).
- Внедрение зависимостей через uber/fx.
- Хранение данных в PostgreSQL (драйвер jackc/pgx).
- Полная контейнеризация с Docker и Docker Compose.
- Написаны unit-тесты для ключевых компонентов и интеграционные тесты для репозиториев.

## Технологический стек
- Язык: Go
- База данных: PostgreSQL
- Драйвер БД: jackc/pgx
- DI-контейнер: uber/fx
- Аутентификация: JWT
- HTTP-сервер: net/http
- Контейнеризация: Docker, Docker Compose
- Тестирование: стандартный пакет testing

## Архитектура
Проект построен по принципам чистой архитектуры:

- domain/entity – доменные модели (игровое поле, игра, пользователь, лидерборд).
- domain/repository – интерфейсы репозиториев.
- application/service – бизнес-логика (игровой сервис, сервис пользователей).
- application/auth – сервис аутентификации и JWT-провайдер.
- datasource – реализация репозиториев, мапперы, SQL-запросы.
- db – подключение к PostgreSQL, миграции.
- web – HTTP-обработчики, middleware, маппинг моделей.
- di – конфигурация графа зависимостей uber/fx.

<details>
<summary><b> Полная структура проекта (кликните чтобы развернуть)</b></summary>
  
```text
.
├── README.md
└── src
    ├── Dockerfile
    ├── cmd
    │   └── app
    │       └── main.go
    ├── docker-compose.yml
    ├── go.mod
    ├── go.sum
    └── internal
        ├── application
        │   ├── auth
        │   │   ├── interface.go
        │   │   ├── jwt_provider.go
        │   │   ├── jwt_provider_test.go
        │   │   ├── model.go
        │   │   ├── service.go
        │   │   └── service_test.go
        │   └── service
        │       ├── game_service.go
        │       ├── service.go
        │       ├── service_test.go
        │       ├── user_interface.go
        │       ├── user_service.go
        │       └── user_service_test.go
        ├── datasource
        │   ├── game_repo_db.go
        │   ├── mapper.go
        │   ├── mapper_test.go
        │   ├── models.go
        │   ├── sql
        │   │   ├── 001_create_users.sql
        │   │   ├── 002_create_games.sql
        │   │   ├── game_wiew.sql
        │   │   ├── games_comleted_by_user.sql
        │   │   ├── games_create.sql
        │   │   ├── games_get_by_id.sql
        │   │   ├── games_list_available.sql
        │   │   ├── games_update.sql
        │   │   ├── get_id.sql
        │   │   ├── leaderboard.sql
        │   │   ├── user_create.sql
        │   │   └── users_get_by_login.sql
        │   ├── sql_queries.go
        │   ├── user_models.go
        │   └── user_repo_db.go
        ├── db
        │   ├── migrate.go
        │   └── postgres.go
        ├── di
        │   └── fx.go
        ├── domain
        │   ├── entity
        │   │   ├── board.go
        │   │   ├── board_test.go
        │   │   ├── game.go
        │   │   ├── game_test.go
        │   │   ├── leaderboard.go
        │   │   └── user.go
        │   └── repository
        │       ├── game_repository.go
        │       ├── mocks
        │       │   ├── game_repository_mock.go
        │       │   └── user_repository_mock.go
        │       └── user_repository.go
        └── web
            ├── auth_handler.go
            ├── auth_handler_test.go
            ├── auth_middleware.go
            ├── auth_middleware_test.go
            ├── handler.go
            ├── handler_test.go
            ├── mapper.go
            ├── mapper_test.go
            ├── mock_gameservice_test.go
            ├── model.go
            ├── server.go
            └── user_handler.go
```
</details> 

## API Endpoints

### Аутентификация
| Метод | Путь | Описание | Доступ |
|-------|------|----------|--------|
| POST | /sign-up | Регистрация нового пользователя | Без авторизации |
| POST | /sign-in | Получение JWT (access+refresh) | Без авторизации |
| POST | /refresh | Обновление access токена | Без авторизации |
| GET | /me | Информация о текущем пользователе | Авторизован |

### Игровые эндпоинты
| Метод | Путь | Описание | Доступ |
|-------|------|----------|--------|
| POST | /games | Создать новую игру (против игрока или бота) | Авторизован |
| GET | /games | Список доступных игр, ожидающих игроков | Авторизован |
| POST | /games/{id}/join | Присоединиться к игре | Авторизован |
| POST | /game/{id} | Сделать ход (обновить поле) | Авторизован |
| GET | /game/{id} | Получить информацию о текущей игре | Авторизован |
| GET | /games/history | История завершённых игр пользователя | Авторизован |
| GET | /leaderboard?top=N | Топ-N игроков по проценту побед | Авторизован |

## Установка и запуск

### Локальный запуск (без Docker)
1. Клонируйте репозиторий.
2. Настройте переменные окружения для подключения к PostgreSQL (см. .env.example).
3. Примените миграции из директории internal/datasource/sql/ (или используйте автоматическую миграцию при старте).
4. Запустите приложение:
```bash
go run ./cmd/app/main.go
``` 

### Запуск через Docker
Проект полностью контейнеризирован. Для запуска выполните команду:
```bash
docker-compose up -d --build
```

Будет запущено два контейнера:
- API (Go-приложение) – доступен на порту 8080
- PostgreSQL – доступен на порту 5432

Все переменные окружения задаются в файле .env:
- POSTGRES_USER – пользователь БД
- POSTGRES_PASSWORD – пароль БД
- POSTGRES_DB – название БД
- DATABASE_URL – строка подключения к БД
- JWT_ACCESS_SECRET – секрет для access-токенов
- JWT_REFRESH_SECRET – секрет для refresh-токенов

При первом запуске миграции применяются автоматически.

## Тестирование

В проекте реализованы unit-тесты для:
- игровой логики (проверка корректности ходов, определение победителя, работа ИИ);
- сервисного слоя (моки репозиториев через интерфейсы);
- JWT-провайдера и middleware авторизации.

## Этапы разработки
1. Базовый слой – игровая логика, in-memory хранилище, REST API, DI-контейнер.
2. База данных и авторизация – миграция на PostgreSQL, базовая аутентификация, многопользовательский режим.
3. JWT, история и лидерборд – замена базовой авторизации на JWT, история завершённых игр, рейтинг игроков.
4. Контейнеризация и тестирование – добавлены Dockerfile и docker-compose, написаны unit‑test.