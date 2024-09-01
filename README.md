# Накопительная система лояльности «Гофермарт»

Получение информации о расчёте начислений

Шаблон из: https://github.com/yandex-praktikum/go-musthave-diploma-tpl


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