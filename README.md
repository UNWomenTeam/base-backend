# Go Restful API Boilerplate agroup07

## Подготовка и запуск

### PostgreSQL в Docker

- Нам нужно сделать эту БД доступной извне
  `docker run --name backend-pg -p 5432:5432 -e POSTGRES_PASSWORD=agroup -d postgres`
  `docker rm -f backend-pg`
  `docker start backend-pg`
- подключиться к контейнеру и уже оттуда запустить, например, psql:
  `psql --username=postgres --dbname=postgres`

### Pgadmin в Docker

`docker run --name admin4-pg -p 89:80 -e 'PGADMIN_DEFAULT_EMAIL=user@domain.com' -e 'PGADMIN_DEFAULT_PASSWORD=1234' -d dpage/pgadmin4`
`docker rm -f admin4-pg`

- http://admin4-pg/browser/
- user@domain.com
- 1234

### Компиляция и запуск

  `go build .` и `./base-backend serve`

### Быстрый запуск

- Клонировать этот репозиторий
- Создайте базу данных postgres и установите переменные среды для вашей базы данных соответственно, если они не используются по умолчанию.
- Запустите приложение, чтобы увидеть доступные команды: `go run main.go`
- Сначала инициализируйте базу данных, запустив сразу все миграции, найденные в ./database/migrate, с помощью команды _
  migrate_: `запустите main.go migrate`
- Запустите приложение с помощью команды _serve_: `go run main.go serve`

[![GoDoc Badge]][godoc] [![GoReportCard Badge]][goreportcard]

Легко расширяемый шаблон RESTful API, нацеленный на использование идиоматического подхода и лучших практик.

Цель иметь прочный и структурированный фундамент, на котором можно строить бэкенд.

Любые отзывы и запросы на вытягивание приветствуются и высоко ценятся. Не стесняйтесь открывать вопросы только для комментариев и
обсуждения.

## Features

The following feature set is a minimal selection of typical Web API requirements:

- Configuration using [viper](https://github.com/spf13/viper)
- CLI features using [cobra](https://github.com/spf13/cobra)
- PostgreSQL support including migrations using [go-pg](https://github.com/go-pg/pg)
- Structured logging with [zap](https://github.com/zap)
- Routing with [chi router](https://github.com/go-chi/chi) and middleware
- JWT Authentication using [lestrrat-go/jwx](https://github.com/lestrrat-go/jwx) with example passwordless email
  authentication
- Request data validation using [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
- HTML emails with [go-mail](https://github.com/go-mail/mail)


### Migration

- Миграции исполняются автоматически при старте ПО  
- Автоматически добавляется два пользователя

| Email                       | Login | Password             | Grant                       |
|-----------------------------|-------|----------------------|-----------------------------|
| admin@agroup07.ru          | root  | agroup              | (has access to admin panel) |
| user@agroup07.ru           | user  | agroup07            |                             |



Проверьте [routes.md](routes.md) для сгенерированного обзора предоставленных маршрутов API.

### API Authentication

For passwordless login following routes are available:

| Path          | Method | Required JSON            | Request JSON                    | Header                                | Description             |
| ------------- | ------ |--------------------------|---------------------------------| ------------------------------------- |-------------------------|
| /auth/login   | POST   | access && refrash token  | {login: "логин", pwd: "пароль"} |                                       | токены доступа          |
| /auth/refresh | POST   |                          |                                 | Authorization: "Bearer refresh_token" | refresh JWTs            |
| /auth/logout  | POST   |                          |                                 | Authorizaiton: "Bearer refresh_token" | logout from this device |

### Example API

Помимо /auth/_ API предоставляет два основных маршрута /api/_ и /admin/\*, например, для разделения приложения и
административный контекст. Для последнего требуется войти в систему как администратор, предоставив соответствующий JWT в
Заголовок авторизации.

| Path                        | Method       | Required JSON                            | Header                                                 | Description                    |
|-----------------------------|--------------|------------------------------------------|--------------------------------------------------------|--------------------------------|
| /admin                      | GET          | Hello Admin                              |                                                        | пинг прав администратора       |
| /admin/accounts             | GET          | [docs/accounts.json](docs/accounts.json) |                                                        | список пользователей           |
| /admin/accounts             | POST         | [docs/account.json](docs/account.json)   | [docs/account_payload.json](docs/account_payload.json) | создать пользователя           |
| /admin/accounts/{accountID} | GET          | [docs/account.json](docs/account.json)   |                                                        | получить пользователя          |
| /admin/accounts/{accountID} | PUT          |                                          | [docs/account_payload.json](docs/account_payload.json) | обновить пользователя          |
| /admin/accounts/{accountID} | DELETE       |                                          |                                                        | удалить пользователя           |
| /api/account                | GET          | [docs/account.json](docs/account.json)   |                                                        | получить текущего пользователя |
| /api/account                | PUT          | [docs/account.json](docs/account.json)   | [docs/account_payload.json](docs/account_payload.json) | обновить текущего пользователя |
| /api/account                | DELETE       |                                          |                                                        | удалить текущего пользователя  |
| /api/profile                | GET          | <Profile>                                |                                                        |                                |
| /api/profile                | PUT          |                                          | <Profile>                                              |                                |
| /ping                       | GET          | pong                                     |                                                        | пинг                           |


### Client API Access and CORS

Сервер настроен для обслуживания клиента Progressive Web App (PWA) из папки _./public_ (этот репозиторий обслуживает только
пример index.html, см. ниже демо-клиент PWA, который можно разместить здесь). В этом случае включение CORS не требуется, т.к.
клиент обслуживается с того же хоста, что и API.

Если вы хотите получить доступ к API от клиента, который находится на сервере с другого хоста, в том числе, например. разработка в прямом эфире
перезагружая сервер с приведенным ниже демо-клиентом, вы должны сначала включить CORS на сервере, установив переменную среды _ENABLE_CORS=true_ на сервере,
чтобы разрешить соединения API от клиентов, серверы которых находятся на других хостах.

#### Demo client application

Для демонстрации функций входа в систему и управления учетной записью этот API служит демо [Vue.js](https://vuejs.org) PWA.
Исходный код клиента можно найти [здесь](https://github.com/UNWomenTeam/base-backend-vue). Соберите и поместите его в
api в папке _./public_ или используйте сервер активной разработки (требуется включенный CORS).

### Environment Variables

По умолчанию viper ищет файл конфигурации в $WORK_DIR/config.yml. Установка вашей конфигурации в качестве переменных среды.

| Name                    | Type          | Default                     | Description                                                                                                                                                                                               |
| ----------------------- | ------------- | --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| PORT                    | string        | localhost:3000              | http address (accepts also port number only for heroku compability)                                                                                                                                       |
| LOG_LEVEL               | string        | debug                       | log level                                                                                                                                                                                                 |
| LOG_TEXTLOGGING         | bool          | false                       | defaults to json logging                                                                                                                                                                                  |
| DB_NETWORK              | string        | tcp                         | database 'tcp' or 'unix' connection                                                                                                                                                                       |
| DB_ADDR                 | string        | localhost:5432              | database tcp address or unix socket                                                                                                                                                                       |
| DB_USER                 | string        | postgres                    | database user name                                                                                                                                                                                        |
| DB_PASSWORD             | string        | postgres                    | database user password                                                                                                                                                                                    |
| DB_DATABASE             | string        | postgres                    | database shema name                                                                                                                                                                                       |
| AUTH_JWT_SECRET         | string        | random                      | jwt sign and verify key - value "random" creates random 32 char secret at startup (and automatically invalidates existing tokens on app restarts, so during dev you might want to set a fixed value here) |
| AUTH_JWT_EXPIRY         | time.Duration | 15m                         | jwt access token expiry                                                                                                                                                                                   |
| AUTH_JWT_REFRESH_EXPIRY | time.Duration | 1h                          | jwt refresh token expiry                                                                                                                                                                                  |
| ENABLE_CORS             | bool          | false                       | enable CORS requests                                                                                                                                                                                      |

### Testing

Пакет auth/pwdless содержит примеры тестов API с использованием фиктивной базы данных.

[godoc]: https://godoc.org/github.com/UNWomenTeam/base-backend

[godoc badge]: https://godoc.org/github.com/UNWomenTeam/base-backend?status.svg

[goreportcard]: https://goreportcard.com/report/github.com/UNWomenTeam/base-backend

[goreportcard badge]: https://goreportcard.com/badge/github.com/UNWomenTeam/base-backend
