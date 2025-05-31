# 💳 Wallets Service

Сервис для управления криптовалютными кошельками с поддержкой переводов между пользователями.

## 📋 Описание

Сервис предоставляет API для:
- Создания кошельков
- Просмотра баланса
- Переводов между кошельками
- Депозитов и снятий средств
- Мониторинга транзакций

## 🏗️ Архитектура

- **Язык**: Go 1.24.3
- **База данных**: PostgreSQL
- **Кэш**: Redis опционально
- **Мониторинг**: Prometheus + Grafana
- **Документация API**: Swagger
- **Миграции**: Goose

## 🚀 Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.24.3+ (для локальной разработки)
- Make

### Запуск в Docker (Production)

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/Sanchir01/wallets.git
cd wallets
```

2. **Запустите все сервисы:**
```bash
make compose-prod
```

Это запустит:
- Приложение на порту `8080`
- PostgreSQL на порту `5439`
- Redis на порту `6380`
- Prometheus на порту `9090`
- Grafana на порту `3000`

3. **Выполните миграции базы данных:**
```bash
make migrations-up-prod
```

### Запуск для разработки

1. **Запустите только инфраструктуру:**
```bash
make docker
```

2. **Выполните миграции:**
```bash
make migrations-up
```

3. **Запустите приложение локально:**
```bash
make run
```

## 🔧 Конфигурация

### Переменные окружения

#### Development (.env.dev)
```env
CONFIG_PATH="config/dev.yaml"
DB_PASSWORD="test"
```

#### Production (.env.prod)
```env
CONFIG_PATH="config/prod.yaml"
DB_USER_PROD="postgres"
DB_PASSWORD="postgres"
DB_HOST_PROD="localhost"
DB_PORT_PROD="5432"
DB_NAME_PROD="postgres"
```

### Файлы конфигурации

- `config/dev.yaml` - настройки для разработки
- `config/prod.yaml` - настройки для продакшена

## 📊 API Endpoints

### Кошельки

- `GET /api/v1/wallets` - Получить все кошельки
- `GET /api/v1/wallets/{id}` - Получить баланс кошелька
- `POST /api/v1/wallets/create` - Создать новый кошелек

### Переводы

- `POST /api/v1/wallets/transfer` - Перевод между кошельками

### Документация API

После запуска приложения Swagger документация доступна по адресу:
```
http://localhost:8080/swagger/index.html
```

## 🗄️ База данных

### Схема

#### Таблица `wallets`
```sql
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    balance DECIMAL(10, 2) DEFAULT 0,
    currency TEXT DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Таблица `transactions`
```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID REFERENCES wallets (id),
    sender_wallet_id UUID REFERENCES wallets (id),
    amount DECIMAL(15, 2) NOT NULL,
    type operation_type NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Enum `operation_type`
```sql
CREATE TYPE operation_type AS ENUM ('DEPOSIT', 'WITHDRAW', 'TRANSFER');
```

### Миграции

```bash
# Создать новую миграцию
make migrations-new MIGRATION_NAME=add_new_feature

# Применить миграции (dev)
make migrations-up

# Откатить миграции (dev)
make migrations-down

# Статус миграций (dev)
make migrations-status

# Применить миграции (prod)
make migrations-up-prod

# Откатить миграции (prod)
make migrations-down-prod

# Статус миграций (prod)
make migrations-status-prod
```

## 📈 Мониторинг

### Prometheus

Метрики доступны по адресу: `http://localhost:9090`

### Grafana

- URL: `http://localhost:3000`
- Логин: `admin`
- Пароль: `admin`

## 🛠️ Команды Make

```bash
# Сборка приложения
make build

# Запуск приложения
make run

# Генерация Swagger документации
make swag

# Сборка Docker образа
make docker-build

# Запуск инфраструктуры для разработки
make docker

# Полный запуск в Docker
make docker-app

# Запуск в production режиме
make compose-prod
```

## 🐳 Docker

### Development
```bash
# Запуск только инфраструктуры (PostgreSQL, Redis, Prometheus, Grafana)
docker-compose up -d
```

### Production
```bash
# Запуск всех сервисов включая приложение
docker-compose -f docker-compose.prod.yaml up --build -d
```

## 🔍 Логирование

Приложение использует структурированное логирование с библиотекой `slog`:
- Логи выводятся в JSON формате в production
- Цветной вывод в development режиме

## 🧪 Тестирование

```bash
# Запуск тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...
```

## 📦 Зависимости

### Основные
- `github.com/go-chi/chi/v5` - HTTP роутер
- `github.com/jackc/pgx/v5` - PostgreSQL драйвер
- `github.com/redis/go-redis/v9` - Redis клиент
- `github.com/Masterminds/squirrel` - SQL билдер
- `github.com/google/uuid` - UUID генератор

### Мониторинг
- `github.com/prometheus/client_golang` - Prometheus метрики

### Документация
- `github.com/swaggo/http-swagger/v2` - Swagger UI

### Валидация
- `github.com/go-playground/validator/v10` - Валидация данных

## 🚨 Устранение неполадок

### Проблемы с базой данных

1. **Ошибка подключения к PostgreSQL:**
```bash
# Проверьте, что PostgreSQL запущен
docker-compose ps

# Проверьте логи
docker-compose logs db
```

2. **Ошибки миграций:**
```bash
# Проверьте статус миграций
make migrations-status

# Пересоздайте базу данных
docker-compose down -v
docker-compose up -d db
make migrations-up
```

### Проблемы с enum типами

Если возникает ошибка `invalid input value for enum operation_type`, убедитесь что используете правильные значения:
- `DEPOSIT` (не `deposit`)
- `WITHDRAW` (не `withdraw`)
- `TRANSFER` (не `transfer`)

### Проблемы с портами

Если порты заняты, измените их в `docker-compose.yaml`:
```yaml
ports:
  - "5440:5432"  # вместо 5439:5432
```

## 📄 Лицензия

MIT License

## 👥 Разработка

### Структура проекта

```
wallets/
├── cmd/main/           # Точка входа приложения
├── internal/
│   ├── app/           # Инициализация приложения
│   ├── feature/       # Бизнес-логика
│   └── server/        # HTTP сервер
├── pkg/               # Общие пакеты
├── migrations/        # Миграции базы данных
├── config/           # Конфигурационные файлы
└── docs/             # Документация
```

**Автор**: [Sanchir01](https://github.com/Sanchir01)
