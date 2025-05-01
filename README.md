# EffectiveMobile Test Task

## Описание
Тестовое задание для EffectiveMobile, реализующее REST API сервис для работы с данными о людях. Сервис позволяет получать, создавать, обновлять и удалять информацию о людях, а также получать статистику по возрасту и полу.

## Структура проекта
```
server/
├── db/              # Работа с базой данных
├── handlers/        # Обработчики HTTP запросов
│   └── api/         # API обработчики
├── models/          # Модели данных
├── routes/          # Маршрутизация
├── services/        # Сервисы
└── requests-manager # Менеджер запросов
```

## Требования
- Go 1.21 или выше
- PostgreSQL

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/Hirogava/EffectiveMobileTest.git
cd EffectiveMobileTest
```

2. Создайте файл .env на основе .env_example:
```bash
cp server/.env_example server/.env
```

3. Настройте переменные окружения в файле .env:
```
DB_CONNECTION_STRING="user=your_user password=your_password dbname=effective sslmode=disable"
SERVER_PORT=8080
```

4. Запустите сервер:
```bash
cd server
go run main.go
```

## API Endpoints

### Люди (Peoples)
- `GET /api/peoples` - Получить список всех людей
- `GET /api/peoples/{id}` - Получить информацию о человеке по ID
- `POST /api/peoples` - Создать новую запись о человеке
- `PUT /api/peoples/{id}` - Обновить информацию о человеке
- `DELETE /api/peoples/{id}` - Удалить запись о человеке

## Модели данных

### People
```json
{
    "id": 1,
    "name": "Ivan",
    "surname": "Ivanov",
    "patronymic": "Ivanovich",
    "age": 30,
    "gender": "male",
    "countries": [
        {
            "country_id": "RU",
            "probability": 0.95
        }
    ]
}
```

## Особенности
- Автоматическое определение возраста и пола по имени
- Определение национальности с вероятностью
- Поддержка отчества (опционально)
- Статистика по возрасту и полу

## Разработка
Для разработки рекомендуется использовать:
- VS Code с расширением Go
- Postman для тестирования API
- pgAdmin для работы с базой данных