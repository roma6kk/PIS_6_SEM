# GO04_01 — REST API сервер для коллекции Celebrities

## Структура проекта

```
GO04_01/
├── main.go            # Исходный код сервера
├── go.mod             # Модуль Go
├── go.sum             # Хэши зависимостей
├── Celebrities.json   # Файл данных коллекции
└── README.md          # Инструкция
```

## Структура элемента коллекции (Celebrity)

```json
{
  "Id":         1,
  "Name":       "Elon Musk",
  "Profession": "Entrepreneur",
  "Country":    "USA",
  "BirthYear":  1971
}
```

## REST API (порт 3000)

| Метод  | URL                  | Описание                                              |
|--------|----------------------|-------------------------------------------------------|
| GET    | /Celebrities/All     | Вся коллекция                                         |
| GET    | /Celebrities/{id}    | Элемент с Id = id                                     |
| POST   | /Celebrities         | Добавить элемент (409 Conflict при дублировании Id)   |
| PUT    | /Celebrities/{id}    | Изменить элемент (404 Not Found если отсутствует)     |
| DELETE | /Celebrities/{id}    | Удалить элемент (404 Not Found если отсутствует)      |

---

## Сборка и запуск (Задание 1)

### 1. Загрузить зависимости

```bash
go mod tidy
```

### 2. Запустить напрямую (для разработки)

```bash
go run main.go
```

### 3. Собрать исполняемый файл (Windows)

```bash
go build -o GO04_01.exe .
```

### 4. Перенести и запустить

```bash
# Создать директорию и скопировать файлы
mkdir C:\GO04_01_dist
copy GO04_01.exe C:\GO04_01_dist\
copy Celebrities.json C:\GO04_01_dist\

# Запустить
cd C:\GO04_01_dist
GO04_01.exe
```

---

## Сборка под Linux (Задание 2)

### На Windows — кросс-компиляция:

```bash
set GOOS=linux
set GOARCH=amd64
t```

### Или на самом Linux:

```bash
go build -o GO04_01 .
./GO04_01
```

### Перенос на Linux:

```bash
# Скопировать файлы на Linux-машину
scp GO04_01_linux user@linux-host:/home/user/GO04_01/
scp Celebrities.json user@linux-host:/home/user/GO04_01/

# На Linux:
chmod +x GO04_01_linux
./GO04_01_linux
```

---

## Примеры запросов (curl)

### GET — вся коллекция

```bash
curl http://localhost:3000/Celebrities/All
```

### GET — элемент по Id

```bash
curl http://localhost:3000/Celebrities/1
```

### POST — добавить элемент

```bash
curl -X POST http://localhost:3000/Celebrities \
  -H "Content-Type: application/json" \
  -d '{"Id":6,"Name":"Lionel Messi","Profession":"Football Player","Country":"Argentina","BirthYear":1987}'
```

### PUT — изменить элемент

```bash
curl -X PUT http://localhost:3000/Celebrities/1 \
  -H "Content-Type: application/json" \
  -d '{"Name":"Elon Musk","Profession":"CEO of Tesla & SpaceX","Country":"USA","BirthYear":1971}'
```

### DELETE — удалить элемент

```bash
curl -X DELETE http://localhost:3000/Celebrities/3
```

---

## Статус-коды ответов

| Ситуация                          | HTTP Status             |
|-----------------------------------|-------------------------|
| Успешный GET / DELETE             | 200 OK / 204 No Content |
| Успешный POST                     | 201 Created             |
| Элемент не найден (GET/PUT/DEL)   | 404 Not Found           |
| Дублирование Id (POST)            | 409 Conflict            |
| Некорректный JSON                 | 400 Bad Request         |
| Ошибка файловой системы           | 500 Internal Server Error |
