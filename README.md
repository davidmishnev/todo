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
go run .
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
go test -v .

# Только тесты хранилища
go test -v ./tasks
```

## Требования

- Go 1.25.1 или выше
- Без внешних зависимостей (только стандартная библиотека)

