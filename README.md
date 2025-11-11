## Получение токена Spirit

Отправить запрос на получение SMS кода (номер - 10 цифр, начиная с 9)

```
curl --location 'https://app.spiritfit.ru/Fitness/hs/mobile/login/check' \
--header 'Content-Type: application/json' \
--data '{
    "phone": "..."
}'
```

Полученный код использовать в следующем запросе

```
curl --location 'https://app.spiritfit.ru/Fitness/hs/mobile/login/code' \
--header 'Content-Type: application/json' \
--data '{
    "phone": "...",
    "code": "..."
}'
```

Полученный токен можно использовать при получении данных

```
curl --location 'https://app.spiritfit.ru/Fitness/hs/mobile/clubs/01' \
--header 'Authorization: ...'
```
