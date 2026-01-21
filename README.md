# Paraklit Shop API (учебный режим)

Backend API для учебного проекта интернет-магазина. Использует in-memory хранилище (данные не сохраняются между перезапусками).

## Требования

- **Go**: 1.24.4+

## Быстрый старт

### 1. Переменные окружения

Минимальный набор для запуска:

```powershell
$env:JWT_SECRET = "dev-secret-key-change-in-production"
$env:AUTH_BUYER_PASSWORD = "buyer123"
$env:AUTH_SELLER_PASSWORD = "seller123"
```

**Обязательные переменные:**
- `JWT_SECRET` — секрет для подписи JWT токенов
- `AUTH_BUYER_PASSWORD` — пароль для покупателя
- `AUTH_SELLER_PASSWORD` — пароль для продавца

### 2. Конфигурация

По умолчанию:
- Читается файл `config/local.yaml` (если существует)
- Поверх него применяются переменные окружения (имеют приоритет)

Путь к конфигу можно переопределить через `CONFIG_PATH`:
```powershell
$env:CONFIG_PATH = "config/production.yaml"
```

### 3. Запуск сервера

```powershell
go run ./cmd/api
```

Сервер запускается на `http://localhost:8080` (порт настраивается через `HTTP_PORT` или `config/local.yaml`).

## Аутентификация

### Тестовые пользователи

**Покупатель (buyer):**
- Email: `buyer@test.com`
- Пароль: значение переменной `AUTH_BUYER_PASSWORD`
- UserID: `1`
- Роль: `buyer`

**Продавец (seller):**
- Email: `seller@test.com`
- Пароль: значение переменной `AUTH_SELLER_PASSWORD`
- UserID: `2`
- Роль: `seller`

### Получение токена

```bash
POST /login
Content-Type: application/json

{
  "username": "buyer@test.com",
  "password": "buyer123"
}
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Используйте токен в заголовке для защищенных эндпоинтов:
```
Authorization: Bearer <token>
```

## API Эндпоинты

### Публичные маршруты

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/health` | Health check |
| `POST` | `/login` | Получить JWT токен |
| `GET` | `/products` | Список товаров |
| `GET` | `/swagger/*` | Swagger UI документация |

### Защищенные маршруты (требуют JWT)

Все маршруты под `/api/*` требуют валидный JWT токен в заголовке `Authorization`.

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/cart` | Просмотр корзины текущего пользователя |
| `POST` | `/api/cart/add/:productId/:qty` | Добавить товар в корзину |
| `DELETE` | `/api/cart/clear` | Очистить корзину |
| `DELETE` | `/api/cart/remove/:productId` | Удалить товар из корзины |
| `POST` | `/api/orders` | Создать заказ из корзины |
| `GET` | `/api/secret` | Тестовый защищенный эндпоинт |

### Ролевые маршруты

#### Для покупателя (`/api/buyer/*`)

Требуют роль `buyer` в JWT токене.

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/buyer/cart` | Просмотр корзины |
| `POST` | `/api/buyer/cart/add/:productId/:qty` | Добавить товар в корзину |
| `POST` | `/api/buyer/orders` | Создать заказ |

#### Для продавца (`/api/seller/*`)

Требуют роль `seller` в JWT токене.

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/seller/cart` | Просмотр корзины |
| `POST` | `/api/seller/cart/add/:productId/:qty` | Добавить товар в корзину |
| `POST` | `/api/seller/orders` | Создать заказ |
| `GET` | `/api/seller/products` | Список товаров (для продавца) |

## Примеры использования

### 1. Получение токена покупателя

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"buyer@test.com","password":"buyer123"}'
```

### 2. Просмотр товаров

```bash
curl http://localhost:8080/products
```

### 3. Добавление товара в корзину

```bash
curl -X POST http://localhost:8080/api/cart/add/1/2 \
  -H "Authorization: Bearer <your-token>"
```

### 4. Просмотр корзины

```bash
curl http://localhost:8080/api/cart \
  -H "Authorization: Bearer <your-token>"
```

### 5. Создание заказа

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer <your-token>"
```

## Структура проекта

```
paraklitshop/
├── cmd/api/              # Точка входа приложения
├── config/               # Конфигурационные файлы
├── internal/
│   ├── auth/             # JWT утилиты (генерация/парсинг токенов)
│   ├── config/           # Загрузка конфигурации
│   ├── handler/          # HTTP обработчики
│   ├── logger/           # Настройка логгера
│   ├── middleware/       # HTTP middleware (JWT, логирование, роли)
│   ├── model/            # Модели данных
│   ├── repository/       # Репозитории (in-memory для учебного режима)
│   ├── server/           # Настройка HTTP сервера
│   └── service/          # Бизнес-логика
└── docs/                 # Swagger документация
```

## Важные замечания

⚠️ **Учебный режим:**
- Данные хранятся в памяти (in-memory)
- Все данные теряются при перезапуске сервера
- PostgreSQL и Redis не используются (конфигурация есть, но подключения нет)

⚠️ **Безопасность:**
- Используйте сильные секреты в production (`JWT_SECRET`)
- Не коммитьте секреты в Git
- Переменные окружения имеют приоритет над конфигурационными файлами

## Swagger документация

После запуска сервера документация доступна по адресу:
- `http://localhost:8080/swagger/index.html`