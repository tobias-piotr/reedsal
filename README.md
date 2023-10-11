# Reedsal

Simple live chat app that uses:

- Email + password for a simple authentication
- JWT for authorization
- Websockets for communication
- PostgreSQL as a database
- sqlx for the queries
- goose for migrations
- Redis for messaging and pubsub

Is it overengineered? For sure. But it was just for fun.

I got too tired and bored with this, so there are some parts that will remain unfinished, like tests or lack of settings.

## Migrations

To create a new migration, use this command and move the file to `db/migrations/` directory.

```zsh
go install github.com/pressly/goose/v3/cmd/goose@latest
goose create <migration> sql
```
