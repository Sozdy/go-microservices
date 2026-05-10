# microservices
![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/Sozdy/d00f582c6b62437ef4e5002ef99e5982/raw/coverage.json)

Учебный проект по микросервисам на Go. Содержит три сервиса — `order`, `inventory`, `payment` — и общий модуль `shared` с proto-контрактами и OpenAPI-спецификациями.

## Требования

- Go 1.26
- [Task](https://taskfile.dev)

## Быстрый старт

```bash
task setup      # установить инструменты (buf, protoc-плагины, linter и т.д.)
task gen        # сгенерировать код из .proto и OpenAPI
task lint       # прогнать golangci-lint
task test:api   # запустить API-тесты
```

Полный список команд — `task --list`.

## Структура

```
order/       # сервис заказов
inventory/   # сервис склада
payment/     # сервис оплат
shared/      # общие proto, OpenAPI, сгенерированный код
```
