
## Stack

Go 1.25+ • Gin • PostgreSQL 16 • JWT • Docker • pgx/v5

## Features

- JWT аутентификация (bcrypt)
- Управление корзиной (add, remove, checkout, instant buy)
- CRUD адресов
- Каталог товаров + поиск

## Quick Start

```bash
cp .env.example .env
# Измените SECRET_KEY в .env
docker-compose up -d
```

**Endpoints:**
- API: http://localhost:8000
- PostgreSQL: localhost:5432
- pgAdmin: http://localhost:5050 (admin@admin.com / admin)

## API

### Public
```
POST   /users/signup          # Регистрация
POST   /users/login           # Вход
GET    /users/productview     # Все товары
GET    /users/search?name=    # Поиск
POST   /admin/addproduct      # Добавить товар
```

### Protected (Bearer token)
```
GET    /addtocart?id=         # В корзину
GET    /removeitem?id=        # Из корзины
GET    /listcart              # Просмотр корзины
GET    /cartcheckout          # Оформить заказ
GET    /instantbuy?id=        # Мгновенная покупка
```

## Structure

```
controllers/   # HTTP handlers
database/      # SQL queries
middleware/    # JWT auth
models/        # Data models
routes/        # Route definitions
tokens/        # JWT generation
migrations/    # DB schema
```

## Database

6 таблиц: users, products, cart, addresses, orders, order_items

Миграции выполняются автоматически при первом запуске.

## Development

```bash
docker-compose up -d              # Запуск
docker-compose logs -f app        # Логи
docker-compose down -v            # Очистка
docker-compose build app          # Пересборка

# Локально
docker-compose stop app
go run main.go
```
