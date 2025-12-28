# TODO API

HTTP-сервер на Go для управления задачами.

## Возможности

- Создание, чтение, обновление и удаление задач
- Управление статусами задач (Assigned, InProgress, Completed, Dropped)
- Валидация данных
- Логирование запросов
- Использование контекста для таймаутов
- Потокобезопасное хранилище в памяти

## Структура задачи

```json
{
  "TaskID": 123,
  "Header": "Название задачи",
  "Description": "Описание",
  "Status": 0
}
```

**Статусы:** 0 - Assigned, 1 - InProgress, 2 - Completed, 3 - Dropped

## API

| Метод | Путь | Описание |
|-------|------|----------|
| POST | /todos | Создать задачу |
| GET | /todos | Получить все задачи |
| GET | /todos/{id} | Получить задачу по ID |
| PUT | /todos/{id} | Обновить задачу |
| DELETE | /todos/{id} | Удалить задачу |

## Запуск

```bash
go run ./cmd/server
```

Сервер запустится на `http://localhost:8080`

## Примеры использования

### Создать задачу

**PowerShell:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/todos -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"Header":"Buy milk","Description":"At the store","Status":0}'
```

**Bash/curl:**
```bash
curl -X POST http://localhost:8080/todos -H "Content-Type: application/json" -d '{"Header":"Buy milk","Description":"At the store","Status":0}'
```

### Получить все задачи

**PowerShell:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/todos
```

**Bash/curl:**
```bash
curl http://localhost:8080/todos
```

### Обновить задачу

**PowerShell:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/todos/1 -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"Header":"Buy milk","Description":"Completed","Status":2}'
```

**Bash/curl:**
```bash
curl -X PUT http://localhost:8080/todos/1 -H "Content-Type: application/json" -d '{"Header":"Buy milk","Description":"Completed","Status":2}'
```

### Удалить задачу

**PowerShell:**
```powershell
Invoke-WebRequest -Uri http://localhost:8080/todos/1 -Method DELETE
```

**Bash/curl:**
```bash
curl -X DELETE http://localhost:8080/todos/1
```

## Тестирование

```bash
# Запустить все тесты
go test ./...

# С подробным выводом
go test -v ./...

# Только тесты сервера
go test -v ./internal/server

# Только тесты хранилища
go test -v ./internal/storage
```

## Структура проекта

```
todo/
├── cmd/
│   └── server/          # Точка входа приложения
│       └── main.go
├── internal/
│   ├── server/          # HTTP-обработчики
│   │   ├── server.go
│   │   └── server_test.go
│   └── storage/         # Хранилище данных
│       ├── task.go
│       ├── errors.go
│       ├── storage.go
│       └── storage_test.go
├── Dockerfile
├── go.mod
└── README.md
```

## Docker

```bash
# Собрать образ
docker build -t todo-server .

# Запустить контейнер
docker run -p 8080:8080 todo-server
```

## Требования

- Go 1.23 или выше
- Без внешних зависимостей (только стандартная библиотека)

