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
- Docker‑образ

## Локальный запуск

```bash
git clone https://github.com/Haywnink/FinalProject.git
cd FinalProject
# Опционально:
export TODO_PASSWORD=random
export TODO_DBFILE=./scheduler.db
export TODO_PORT=7540
go run main.go
