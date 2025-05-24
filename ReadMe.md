# Go Offline Test - Сервис цитат
## Описание проекта
Проект представляет собой сервис для работы с цитатами, реализованный на Go. Сервис предоставляет REST API для добавления, получения, удаления цитат и работы с авторами.

## Структура проекта
```
go-offline-test/
├── app/                  # Основное приложение
├── internal/             # Внутренние пакеты
│   ├── controllers/      # HTTP контроллеры
│   ├── repository/       # Репозиторий для хранения данных
│   ├── services/         # Бизнес-логика
│   └── shared/           # Общие структуры
│       └── dto/          # Data Transfer Objects
├── tests/                # Тесты
├── .dockerignore
├── .gitignore
├── docker-compose.yml    # Конфигурация Docker
├── Dockerfile            # Сборка образа
└── go.mod               # Зависимости Go
```
## Настройка окружения
1. Убедитесь, что вы находитесь в корневой папке проекта - `go-offline-test`.
2. Создайте новый файл с расширением `.env`:
### Для Linux:
``` bash
touch .env
```
### Для Windows
``` bash
type nul >> .env
```
3. Откройте только что созданый файл в корневой папке проекта с расширением `.env` и вставьте туда следующее:
``` .env
ADDR_CONFIG=http://localhost
PORT_CONFIG=:8080
```
## Установка и запуск
### Требования
Go 1.21+
Docker (опционально)

### Запуск без Docker
1. Клонируйте репозиторий:
``` bash
git clone https://github.com/yourusername/go-offline-test.git
cd go-offline-test
```
2. Установите зависимости:
``` bash
go mod download
```
3. Раскомментируйте в app/app.go следующие строки:
``` go
	// projectRoot, err := filepath.Abs("../")
	// if err != nil {
	// 	log.Fatal("Не удалось получить projectRoot:", err)
	// }
	//
	// envPath := filepath.Join(projectRoot, ".env")
	// if err := godotenv.Load(envPath); err != nil {
	//	log.Fatal("Не удалось загрузить .env файл")
	// }
	// log.Printf("INFO: .env файл успешно загружен")
```
4. Запустите сервер:
``` bash
go run app/app.go
```
### Запуск с Docker
``` bash
docker-compose up --build
```
## API Endpoints
### Цитаты
`GET /quotes` - Получить все цитаты

`GET /quotes/random` - Получить случайную цитату

`POST /quotes` - Добавить новую цитату

`DELETE /quotes/{id}` - Удалить цитату по ID

### Цитаты по авторам
`GET /quotes?author={name}` - Получить цитаты автора

## Примеры запросов
### Добавление цитаты
``` bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"text": "Пример цитаты", "authorName": "Пример Автора"}'
```
### Получение случайной цитаты
``` bash
curl http://localhost:8080/quotes/random
```
### Получение цитат автора
``` bash
curl http://localhost:8080/quotes?author=Пример%20Автора
```
## Интерфейсы:
### Сервис
``` go
type IQuoteService interface {
    AddQuote(ctx context.Context, quote *dto.Quote) error
    ListQuotes(ctx context.Context) ([]*dto.Quote, error)
    RandomQuote(ctx context.Context) (*dto.Quote, error)
    QuotesByAuthor(ctx context.Context, authorName string) ([]*dto.Quote, error)
    DeleteQuote(ctx context.Context, quoteID int) error
    ValidateData(text, authorName, mode string) error
}
```
### Репозиторий:
``` go
type IQuoteRepository interface {
    AddQuote(ctx context.Context, quote *dto.Quote) error
    Quotes(ctx context.Context) ([]*dto.Quote, error)
    RandomQuote(ctx context.Context) (*dto.Quote, error)
    QuotesByAuthor(ctx context.Context, authorName string) ([]*dto.Quote, error)
    DeleteQuote(ctx context.Context, idQuote int) error
}
```
## Валидация данных
### Сервис выполняет строгую валидацию:

* Цитата: 1-500 символов, не пустая

* Имя автора: 2-100 символов, только буквы, пробелы и дефисы

* Не может начинаться/заканчиваться дефисом

## Обработка ошибок
### Сервис возвращает детализированные ошибки с HTTP статусами:

* 400 - Невалидные данные

* 404 - Цитата/автор не найден

* 500 - Внутренняя ошибка сервера

## Тестирование
### Для запуска тестов:

``` bash
cd tests
go test -v ./...
```
## Лицензия
MIT License