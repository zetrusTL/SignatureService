# Signature Service

Модуль: `github.com/zetrus/signature-service`.

Сервис на Go для создания устройств подписи (RSA / ECDSA на P-256) и подписи транзакций по описанному в задании формату `secured_data_to_be_signed`.

## Запуск

```bash
go mod tidy
go run ./cmd/server
```

По умолчанию слушает `ADDR` (если не задано — `:8080`).

## HTTP API

| Метод | Путь | Описание |
|--------|------|----------|
| `POST` | `/v1/signature-devices` | Создание устройства |
| `GET` | `/v1/signature-devices` | Список устройств |
| `GET` | `/v1/signature-devices/{id}` | Устройство по UUID |
| `POST` | `/v1/signature-devices/{id}/signatures` | Подпись данных |

Создание (`POST /v1/signature-devices`), тело JSON:

```json
{
  "id": "",
  "algorithm": "RSA",
  "label": "опционально"
}
```

Поле `id` можно опустить или передать пустую строку — будет сгенерирован UUID. Непустой `id` должен быть корректным UUID.

Подпись (`POST .../signatures`):

```json
{ "data": "данные для подписи" }
```

Ответ:

```json
{
  "signature": "<base64>",
  "signed_data": "<counter>_<data>_<third_base64>"
}
```

## Тесты

```bash
go test ./...
go test -race ./...
```

