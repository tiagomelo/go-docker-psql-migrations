# go-docker-psql-migrations

A sample project that demonstrates how to manage [PostgreSQL](https://www.postgresql.org/) database migrations using [Go](https://go.dev), [golang-migrate](https://github.com/golang-migrate/migrate) and [Docker](https://docker.com).

## Prerequisites

- [Go](https://go.dev)
- [Docker](https://docker.com)
- [docker-compose](https://docs.docker.com/compose/)

## running it

```
make run
```

it will display all existing tables along with their column names.

## unit tests

```
make test
```

## available Make targets

```
$ make help
Usage: make [target]

  help                    shows this help message
  start-psql              starts psql instance
  stop-psql               stops psql instance
  psql-console            opens psql terminal
  create-migration        creates a migration file
  migrate-up              runs migrations up to N version (optional)
  migrate-down            runs migrations down to N version (optional)
  migrate-version         shows current migration version number
  migrate-force-version   forces migrations to version V
  test                    run unit tests
  run                     runs the app
```