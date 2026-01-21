## Paraklit Shop (учебный режим)

### 1. Требования
- **Go**: 1.24.4+

### 2. Переменные окружения (обязательно)
Минимальный набор для запуска:

- `JWT_SECRET` — секрет для подписи JWT  
- `AUTH_BUYER_PASSWORD` — пароль покупателя `buyer@test.com`  
- `AUTH_SELLER_PASSWORD` — пароль продавца `seller@test.com`

Пример (PowerShell):

```powershell
$env:JWT_SECRET = "dev-secret"
$env:AUTH_BUYER_PASSWORD = "buyer123"
$env:AUTH_SELLER_PASSWORD = "seller123"
```

### 3. Конфигурация
По умолчанию:

- Читается файл `config/local.yaml` (если есть)  
- Поверх него применяются переменные окружения  

Путь к конфигу можно переопределить:

- `CONFIG_PATH=path/to/config.yaml`

### 4. Запуск сервера

Из корня проекта:

```powershell
go run ./cmd/api
```

Сервер по умолчанию слушает `http://localhost:8080`.

### 5. Тестовый логин и роли

- Покупатель:  
  - `username`: `buyer@test.com`  
  - `password`: значение `AUTH_BUYER_PASSWORD`

- Продавец:  
  - `username`: `seller@test.com`  
  - `password`: значение `AUTH_SELLER_PASSWORD`

### 6. Основные эндпоинты (учебный сценарий)

- `POST /login` — получить JWT токен  
- `GET /products` — список товаров (in-memory)  
- `GET /swagger/index.html` — Swagger UI (если включён)

Защищённые маршруты (требуют `Authorization: Bearer <token>`):

- `/api/cart`, `/api/cart/add/:productId/:qty`, `/api/cart/clear`, `/api/cart/remove/:productId`
- `/api/orders`

