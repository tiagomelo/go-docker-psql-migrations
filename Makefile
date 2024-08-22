include .env
export

# ==============================================================================
# Useful variables

# Version - this is optionally used on goto command
V?=

# Number of migrations - this is optionally used on up and down commands
N?=

# PSQL domain source name string
PSQL_DSN ?= $(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST)/$(POSTGRES_DB)

.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]\n"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==============================================================================
# DB

.PHONY: start-psql
## start-psql: starts psql instance
start-psql:
	@ docker-compose up -d
	@ echo "Waiting for Postgres to start..."
	@ until docker exec $(POSTGRES_DATABASE_CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)  -c "SELECT 1;" >/dev/null 2>&1; do \
		echo "Postgres not ready, sleeping for 5 seconds..."; \
		sleep 5; \
	done
	@ echo "Postgres is up and running."

.PHONY: stop-psql
## stop-psql: stops psql instance
stop-psql:
	@ docker-compose down

.PHONY: psql-console
## psql-console: opens psql terminal
psql-console: export PGPASSWORD=$(POSTGRES_PASSWORD)
psql-console:
	@ docker exec -it $(POSTGRES_DATABASE_CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

# ==============================================================================
# DB migrations

.PHONY: create-migration
## create-migration: creates a migration file
create-migration:
	@ if [ -z "$(NAME)" ]; then echo >&2 "please set the name of the migration via the variable NAME"; exit 2; fi
	@ docker run --rm -v `pwd`/db/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(NAME)

.PHONY: migrate-up
## migrate-up: runs migrations up to N version (optional)
migrate-up: start-psql
	@ docker run --rm --network $(POSTGRES_DATABASE_CONTAINER_NETWORK_NAME) -v `pwd`/db/migrations:/migrations migrate/migrate -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_DATABASE_CONTAINER_NAME):5432/$(POSTGRES_DB)?sslmode=disable' -path /migrations up $(N)

.PHONY: migrate-down
## migrate-down: runs migrations down to N version (optional)
migrate-down:
	@ if [ -z "$(N)" ]; then \
		docker run --rm --network $(POSTGRES_DATABASE_CONTAINER_NETWORK_NAME) -v `pwd`/db/migrations:/migrations migrate/migrate -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_DATABASE_CONTAINER_NAME):5432/$(POSTGRES_DB)?sslmode=disable' -path /migrations down -all; \
	else \
		docker run --rm --network $(POSTGRES_DATABASE_CONTAINER_NETWORK_NAME) -v `pwd`/db/migrations:/migrations migrate/migrate -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_DATABASE_CONTAINER_NAME):5432/$(POSTGRES_DB)?sslmode=disable' -path /migrations down $(N); \
	fi

.PHONY: migrate-version
## migrate-version: shows current migration version number
migrate-version:
	@ docker run --rm --network $(POSTGRES_DATABASE_CONTAINER_NETWORK_NAME) -v `pwd`/db/migrations:/migrations migrate/migrate -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_DATABASE_CONTAINER_NAME):5432/$(POSTGRES_DB)?sslmode=disable' -path /migrations version

.PHONY: migrate-force-version
## migrate-force-version: forces migrations to version V
migrate-force-version:
	@ if [ -z "$(V)" ]; then echo >&2 please set version via variable V; exit 2; fi
	@ docker run --rm --network $(POSTGRES_DATABASE_CONTAINER_NETWORK_NAME) -v `pwd`/db/migrations:/migrations migrate/migrate -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_DATABASE_CONTAINER_NAME):5432/$(POSTGRES_DB)?sslmode=disable' -path /migrations force $(V)

# ==============================================================================
# Unit tests

.PHONY: test
## test: run unit tests
test:
	@ go test -v ./... -count=1

# ==============================================================================
# App execution

.PHONY: run
## run: runs the app
run: migrate-up
	@ go run cmd/main.go
