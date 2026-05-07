# GO05_01 — Лабораторная работа 5

REST-сервер на Go с PostgreSQL (`database/sql` + `gorilla/mux`).

---

## Структура Celebrity

```json
{
  "id": 1,
  "name": "Elon Musk",
  "age": 52,
  "country": "USA",
  "profession": "Entrepreneur"
}
```

---

## Установка и запуск

### 1. Подготовить PostgreSQL

Создайте базу данных `celebrities` и при необходимости пользователя:

```sql
CREATE DATABASE celebrities;
```

По умолчанию приложение подключается к:

```bash
postgres://postgres:postgres@localhost:5432/celebrities?sslmode=disable
```

Если нужны другие параметры, задайте переменную окружения `DATABASE_URL`.

### 2. Установить зависимости

```bash
cd GO05_01
go mod tidy
```

### 3. Собрать исполняемый файл

```bash
# Linux / macOS
go build -o GO05_01 .

# Windows
go build -o GO05_01.exe .
```

### 4. Запустить приложение

```bash
# Linux / macOS
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/celebrities?sslmode=disable"
./GO05_01

# Windows PowerShell
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/celebrities?sslmode=disable"
./GO05_01.exe
```

При старте приложение само создаёт таблицу `Celebrities`, если её ещё нет.

---

## Тестирование API

### GET — все элементы
```
GET http://localhost:3000/Celebrities/All
```

### GET — по Id
```
GET http://localhost:3000/Celebrities/1
```

### POST — добавить элемент
```
POST http://localhost:3000/Celebrities
Content-Type: application/json

{
  "id": 1,
  "name": "Elon Musk",
  "age": 52,
  "country": "USA",
  "profession": "Entrepreneur"
}
```
- При дублировании Id → **409 Conflict**

### PUT — изменить элемент
```
PUT http://localhost:3000/Celebrities/1
Content-Type: application/json

{
  "name": "Elon Musk",
  "age": 53,
  "country": "USA",
  "profession": "CEO"
}
```
- При отсутствии Id → **404 Not Found**

### DELETE — удалить элемент
```
DELETE http://localhost:3000/Celebrities/1
```
- При отсутствии Id → **404 Not Found**

---

## Тестирование через curl

```bash
# Добавить
curl -X POST http://localhost:3000/Celebrities \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"Elon Musk","age":52,"country":"USA","profession":"Entrepreneur"}'

# Получить все
curl http://localhost:3000/Celebrities/All

# Получить по Id
curl http://localhost:3000/Celebrities/1

# Изменить
curl -X PUT http://localhost:3000/Celebrities/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Elon Musk","age":53,"country":"USA","profession":"CEO"}'

# Удалить
curl -X DELETE http://localhost:3000/Celebrities/1
```
