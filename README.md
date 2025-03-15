![image](https://github.com/user-attachments/assets/a89a772c-769c-40f0-b098-080a8f538ada)

# Shortify

Shortify это сервис сокращатель ссылок, сделанный в рамках отбора на стажировку в компанию Ozon.

## Требования

- [Task](https://taskfile.dev/installation/) - для удобного запуска заготовленных задач (команд);
- Docker - все задачи выполняются в Docker контейнерах, что избавляет от необходимости устанавливать какие-либо зависимости.

## Быстрый старт

### Подготовьте файл с переменными окружения

Можно просто взять файл `.env.example` и убрать окончание `.example`.

```shell
mv .env.example .env
```

Параметр `PERSISTENCE_TYPE` отвечает за тип хранилища ссылок.  
Доступно `inmemory` и `postgres`.

### Запуск всего приложения

```shell
task start
```

Команда запустит базу данных PostgreSQL, применит все миграции и запустит HTTP сервер.  
По адресу http://localhost:8080/swagger/index.html можно будет открыть Swagger документацию.

### Запуск юнит-тестов

```shell
task unit-test
```

## Структура проекта

Сервис разработан согласно принципам SOLID и чистой архитектуры для большей поддерживаемости и масштабируемости.
```
.
├── cmd
│   └── api-server            # команда, запускающая API сервер
├── configs
├── docs                      # Swagger документация
├── internal
│   ├── config
│   ├── delivery              # способы доставки данных в наше приложение будь то http, cli или kafka
│   │   └── http              # REST API
│   │       └── v1            # версионирование REST API
│   ├── domain                # доменный слой, который содержит всего одну сущность - URL
│   ├── dto                   # DTO сервисов для общения со слоем контроллеров.
│   └── generation            # реализация различных схем предоставления коротких ссылок
│       └── rand              # генерация на основе пакета crypto/rand
│       └── kgs               # здесь же могла бы быть реализация, обращающаяся к какому-то внешнему сервису (Key Generation Service)
│   └── persistence           # реализации различных схем хранения данных
│       └── inmemory          # в памяти
│       └── postgres          # в базе данных
│   └── url                   # сервисный слой
├── migrations
├── v1                        # DTO http контроллеров
│   └── url.go                # и одновременно это пакет для других Go'шных сервисов. Здесь же можно предоставить HTTP клиент
├── .env.example
├── .mockery.yaml             # конфигурация mockery для генерации моков
├── Dockerfile
├── docker-compose.yaml
├── go.mod
├── go.sum
```
