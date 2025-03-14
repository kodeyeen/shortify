![image](https://github.com/user-attachments/assets/f8283937-70cd-451f-8468-120c8ac5e9a8)
# Shortify

Shortify это сервис сокращатель ссылок, сделанный в рамках отбора на стажировку в компанию Ozon.

## Требования

- [Task](https://taskfile.dev/installation/) - для удобного запуска заготовленных задач (команд);
- Docker - все задачи выполняются в Docker контейнерах, что избавляет от необходимости устанавливать какие-либо зависимости.

## Быстрый старт

### Подготовьте файл с переменными окружения

Можно просто взять файл `.env.example` и убрать окончание `.example`.

### Запуск всего приложения

```
task start
```

Команда запустит базу данных PostgreSQL, применит все миграции и запустит HTTP сервер.  
По адресу http://localhost:8080/swagger/index.html можно будет открыть Swagger документацию.

![Снимок экрана 2025-03-14 155635](https://github.com/user-attachments/assets/a89a772c-769c-40f0-b098-080a8f538ada)


### Запуск юнит-тестов

```
task unit-test
```

### Особенности

TODO

### Структура проекта

TODO
