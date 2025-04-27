# TODO‑Scheduler

Простой веб‑сервис‑планировщик задач с повторениями и (опциональной) JWT‑аутентификацией.

## Описание проекта

Это REST‑сервис на Go + SQLite, реализующий функциональность TODO‑списка:

- Добавление/редактирование/удаление задач через API
- Повторения «d N» (каждые N дней) и «y» (раз в год), *опционально* — по неделям/месяцам
- Аутентификация по паролю (JWT + cookie), если задана переменная `TODO_PASSWORD`
- Фронтенд: статическая папка `web/`

## Задания со «звёздочкой»

- Расширенные правила повторения недель и месяцев (`w`, `m`)
- Поиск задач по строке и дате (`/api/tasks?search=…`)
- JWT‑аутентификация (`/api/signin`)
- Собираемый Docker‑образ

## Локальный запуск

1. Клонируем репозиторий и переходим в директорию проекта:
    ```bash
    git clone https://github.com/Haywnink/FinalProject.git
    cd FinalProject
    ```

2. (Опционально) Экспортируем переменные окружения:
    ```bash
    export TODO_PASSWORD="your_secret_password"
    export TODO_DBFILE="./scheduler.db"
    export TODO_PORT=7540
    ```

3. Запускаем приложение:
    ```bash
    go run main.go
    ```

   По умолчанию сервис будет доступен на `http://localhost:7540/`.

## Сборка и запуск в Docker

1. Собираем образ (порт и путь к БД можно изменить с помощью аргументов сборки):
    ```bash
    docker build \
      --build-arg TODO_PORT=8080 \
      --build-arg TODO_DBFILE=/data/scheduler.db \
      -t my-scheduler .
    ```

2. Запускаем контейнер:
    ```bash
    docker run -d \
      --name scheduler \
      -e TODO_PORT=8080 \
      -e TODO_DBFILE=/data/scheduler.db \
      -v /path/on/host/scheduler.db:/data/scheduler.db \
      -p 8080:8080 \
      my-scheduler
    ```

   После этого сервис будет доступен на `http://localhost:8080/`.
