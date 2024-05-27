# OZON
Тестовое задание

## :open_file_folder: Структура проекта

- `cmd/main.go` - главный файл
- `config` - папка с конфигом
- `db` - папка с миграциями
- `internal` - папка с кодом проекта
    - `app` - папка с точкой запуска приложения
    - `auth` - логика авторизации
    - `core` - основная логика
        - `delivery` - уровень delivery с реализацией graphql
        - `model` - модели проекта
        - `repository` - уровень repository c реализацией postgres и redis 
        - `usecase` - уровень usecase с основной бизнес логикой
    - `middleware` - папка с middleware
- `pkg`
- `docker-compose.yml` - docker-compose для запуска контейнеров
- `Dockerfile` - Dockerfile для запуска сервиса в докер-контейнере
- `...`

## О проекте

IN_MEMORY_MODE = bool - env переменная отвечающая за способ хранения данных(Redis/PostgresQL)

http://localhost:19090/ - GraphQL Playground

### Технологии:
 - UberFx для построения Dependency injection
 - Uber.Zap - логгер
 - Paseto токены для реализации авторизации с access и refresh токенами
 - Valkey для удобной реализации логики с токенами авторизации
 - Redis в качестве in-memory базы данных
 - Jaeger для трассировка запросов(http://localhost:16686/)
 - Golangci-lint - линтер
 - Gin - http фреймворк
 - Goose - система миграции с автоматизацией через docker-compose

Для решения проблым n + 1 запросов были использованы dataloader`ы
Пагинация была реализована методом курсора(Cursor Based Pagination). 
```
query GetContent {
  content(postId: 8) {
    id,
    content,
    commentsEnabled,
    comments(first: 10, cursor: "String", deep: false) {
      edges {
        cursor,
        hasSubComments,
        node {
          id,
          postId,
          parentId,
          authorUuid,
          content,
          createdAt
        }
      },
      pageInfo {
        startCursor,
        endCursor,
        hasNextPage
      }
    }
  }
}
```
Где first - это кол-во возвращаемых элементов, cursor - шифрованная относительная позиция, deep - флаг для получения вложенных комментариев.
При depp = true, будут получены вложенные комментарии от позиции курсора. От курсоров первого уровня можно получить комментарии первого уровня или с deep = true
элементы второго уровня. У каждого элемента есть hasSubComments, который указывает о наличии вложенности.

## (*) GraphQL Subscriptions

Был реализован функционал подписки на новые комментария под постом
```
subscription CommentSub {
  notification(postID: 1) {
    id,
	postId,
    parentId,
    authorUuid,
	content,
    createdAt
  }
}
```

## :hammer: Запуск

docker compose up --build

## Тесты

go run ./...
Тестами была покрыта основная логика сервиса.
Для тестирования Repository был использован testcontainers, с образами postgres и redis
