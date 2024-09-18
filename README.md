# Накопительная система лояльности «Гофермарт»

Получение информации о расчёте начислений

Шаблон из: https://github.com/yandex-praktikum/go-musthave-diploma-tpl
ТЗ: https://github.com/yandex-praktikum/go-musthave-group-diploma-tpl/blob/master/SPECIFICATION.md

### Алгоритм Луна онлайн: https://planetcalc.ru/2464/
Последний разряд контрольной суммы: 0 - заказ валидный для http://localhost:8080/api/orders


### Accrual
 - Регистрация нового совершённого заказа
    POST http://localhost:8080/api/orders
```json
{
    "order": "12345678908513",
    "goods": [
        {
            "description": "Чайник Bork",
            "price": 200
        },
        {
            "description": "Ботинок левыый",
            "price": 100
        }
    ]
} 
```

- Получение информации о расчёте начислений
    GET http://localhost:8080/api/orders/12345678908513
```json

```

- Регистрация информации о вознаграждении за товар
    POST http://localhost:8080/api/goods

```json
{
    "match": "Bork",
    "reward": 10,
    "reward_type": "%"
} 
```

### Gose (Migrations) https://github.com/pressly/goose
```bash
git clone https://github.com/pressly/goose
cd goose
go mod tidy
go build -o goose ./cmd/goose

./goose --version
# goose version:(devel)
```

### Запуск, тестирование

1. Настроить build/package/docker/postgres/.env и собрать build/package/docker/postgres/docker-compose.yml docker
2. Запустить БД tools/start_docker_postgres.sh
3. Запустить автотесты tools/run_gophermarttest.sh
4. Для запуска приложения(без использования автотестов) дополнительно нужно запустить сервис расчёта баллов tools/accrual_start.sh
5. После этого можно запускать main функцию cmd/gophermart/main.go

### Конфигурация
**Переменные окружениия:**
 - RUN_ADDRESS Адрес сервера и порт
 - ACCRUAL_SYSTEM_ADDRESS Внешняя система расчёта бонусов
 - DATABASE_URI Строка подключения к БД
 - LOG_LEVEL Уроверь логирования

**Параметры командной строки:**
 - a адрес и порт запуска сервиса
 - d адрес подключения к базе данных
 - r адрес системы расчёта начислений
 - l уровень логирования