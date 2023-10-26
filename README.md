## Получение учетных данных Google Cloud

1. Создать Google Cloud проект
1. Создать сервисный аккаунт в Google Cloud проекте, сгенерировать для него ключи (они автоматически скачаются в браузере)
1. Скачанный файл положить в этот проект и переименовать в `credentials.json`
1. Создать Google Sheet, добавить сервисный аккаунт как редактора в таблицу по почте сервисного аккаунта

## Получение токена Spirit

Отправить запрос на получение SMS кода

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
