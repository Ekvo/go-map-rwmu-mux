# Quotation book (go-map-rwmu-mux)

---
Мини-сервис “Цитатник”  

Реализует **REST API**-сервис на **Go** для хранения и управления цитатами.

**Основные принципы:** DTO, SOLID  
**Технологический стек:** Go 1.24.1  

#### Возможности
* сохранять цитату
* читать случайную цитату
* читать список всех цитат
* читать список цитат по автору
* удалять цитаты по ID
___

#### Directory structure
```
.
├── cmd/app
│       └──── main.go
├── init
│   └──── .env  // сервис не коммерческий .env не в .gitignore
├── internal
|   ├── app 
|   │   └── app.go      // инициализация, запуск и остановка сервиса
|   ├── config 
|   │   └── config.go   // получение данных из .env
|   ├── db 
|   │   ├── db.go       // описание базы и методов 
|   │   ├── requests.go // реализация запросов в базу      
|   │   └── requests_test.go 
|   ├── model 
|   │   └──── quote.go     
|   ├── server  
|   │   └──── server.go   
|   ├── servises
|   │   ├── create_quote.go      
|   │   ├── delete_quote.go             
|   │   ├── deserializer.go    // обработка запроса      
|   │   ├── read_quote_list.go
|   │   ├── read_quotes_by_author.go     
|   │   ├── read_random_quote.go
|   │   ├── serializer.go      // создание ответа
|   │   └── services.go        // бизнес логика     
|   └── transport   
|       ├── router_test.go     
|       ├── route.go      // реализация запросов
|       └── transport.go  // маршрутизация 
└── pkg/utils
    └──── utils.go        // вспомогательные функции

Dockerfile
go.mod
README.md
```
---

## Запуск приложения
**Важно**:  
Проверьте, свободен ли у вас порт `8080`. Например, если у вас установлен локально `postgresql`, нужно освободить его.  
Возможное решение: [link](https://stackoverflow.com/questions/47026506/edb-postgres-server-from-local-host-apache-server-is-up-and-running-the-default "https://stackoverflow.com/questions/47026506/edb-postgres-server-from-local-host-apache-server-is-up-and-running-the-default") 

Для запуска приложения должен быть установлен **Go компилятор**, и среда для разработки на языке **Go**, а также **Git**:

* [Шаги установки Go и VSCode](https://learn.microsoft.com/ru-ru/azure/developer/go/configure-visual-studio-code "https://learn.microsoft.com/ru-ru/azure/developer/go/configure-visual-studio-code")
* [Шаги установки Git](https://git-scm.com/book/ru/v2/%D0%92%D0%B2%D0%B5%D0%B4%D0%B5%D0%BD%D0%B8%D0%B5-%D0%A3%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0-Git "https://git-scm.com/book/ru/v2/%D0%92%D0%B2%D0%B5%D0%B4%D0%B5%D0%BD%D0%B8%D0%B5-%D0%A3%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0-Git") 

#### Первый способ:
- открыть **git bash**
- склонировать к себе на локальную машину 
```bash
git clone https://github.com/Ekvo/go-map-rwmu-mux.git
```
- перейти в директорию `go-map-rwmu-mux` и запустить приложение
```bash
cd go-map-rwmu-mux && 
go run cmd/app/main.go
```

#### Второй способ:
Для запуска данным способом необходима установка **Docker**

* [установка Docker для Mac OS](https://docs.docker.com/desktop/setup/install/mac-install/ "https://docs.docker.com/desktop/setup/install/mac-install/") 
* [видео установки для Mac](https://www.youtube.com/watch?v=S2kvJw58504 "https://www.youtube.com/watch?v=S2kvJw58504")
* [установка Docker для Windows](https://docs.docker.com/desktop/setup/install/windows-install/ "https://docs.docker.com/desktop/setup/install/windows-install/")
* [видео установки для Windows](https://www.youtube.com/watch?v=xQDh6dJWTf8 "https://www.youtube.com/watch?v=xQDh6dJWTf8")
* [установка Docker для Linux](https://docs.docker.com/desktop/setup/install/linux/ "https://docs.docker.com/desktop/setup/install/linux/")
* [видео установки для Linux Ubuntu](https://www.youtube.com/watch?v=ozEXL4JnedE "https://www.youtube.com/watch?v=ozEXL4JnedE")

- Склонировать к себе на локальную машину
```bash
git clone https://github.com/Ekvo/go-map-rwmu-mux.git
```
- перейти в директорию `go-map-rwmu-mux` и создать образ
```bash
cd go-map-rwmu-mux &&
docker build -t quotebook:v1.0.0 .
```
- запустить 
```bash
docker run -p 8080:8080 quotebook:v1.0.0
```
----

#### Curl
Можно покидать запросы:

* создание цитаты
```http request
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```
* получение списка цитат
```http request
curl http://localhost:8080/quotes 
```
* случайная цитата
```http request
curl http://localhost:8080/quotes/random
```
* список цитат по автору
```http request
curl http://localhost:8080/quotes?author=Confucius
```
* удаление цитаты по ID
```http request
curl -X DELETE http://localhost:8080/quotes/1
```
---

#### Tests
Для запуска тестов, после того как сервис склонирован, с `github.com/Ekvo/go-map-rwmu-mux` и вы находитесь в директории `go-map-rwmu-mux` 
```bash
go test ./...
```

Также можно посмотреть покрытие `coverage` тестами:
```bash
go test ./internal/db -coverprofile=coverage.put
```
Вместо `./internal/db` можно также использовать `./internal/transport` и `./pkg/utils`.  
Для просмотра результатов
```bash
go tool cover -html=coverage
```
### coverage
| direct                                          | percent `%` |
|:------------------------------------------------|----------:|
| go-map-rwmu-mux/internal/db/db.go               |      88.9 |
| go-map-rwmu-mux/internal/db/requests.go         |      81.1 |
|                                                 |           |
| go-map-rwmu-mux/internal/transport/transport.go |     100.0 |
| go-map-rwmu-mux/internal/transport/router.go    |      86.3 |
|                                                 |           |
| go-map-rwmu-mux/pkg/utils/utils.go              |      82.1 |


ps. Thank you for your time:)
