
include .env.example
export $(shell sed 's/=.*//' .env.example)

ENV_FILE=.env
COMPOSE=docker compose --env-file $(ENV_FILE)

.PHONY: init run stop refresh-db logs db_bash redis_bash app_bash swagger migrate migrate_down seed help set-env test coverage

set-env:
	@if [ ! -f .env ]; then cp .env.example .env; else echo ".env already exists"; fi

init:
	$(MAKE) set-env
	$(COMPOSE) build --no-cache
	$(COMPOSE) up -d postgres redis
	$(COMPOSE) up -d app
	$(MAKE) migrate
	$(MAKE) seed
	
	@echo "Initialization complete."
	@echo "Everything works! You can now run 'make logs' to see the application logs."
	@echo "Use 'make help' to see available commands."
	@echo "Happy coding!"

run:
	$(COMPOSE) up -d

stop:
	$(COMPOSE) down

refresh_db:
	$(MAKE) migrate_down
	$(MAKE) migrate
	$(MAKE) seed

logs:
	$(COMPOSE) logs -f

db_bash:
	$(COMPOSE) exec postgres psql -U "$(DB_USER)" -d "$(DB_NAME)"

redis_bash:
	$(COMPOSE) exec redis redis-cli

app_bash:
	$(COMPOSE) exec app sh

swagger:
	swag init -g cmd/api/main.go -o docs/swagger

migrate:
	$(COMPOSE) exec app migrate -path /app/migrations -database "$(DATABASE_URL)" up

migrate_down:
	$(COMPOSE) exec app migrate -path /app/migrations -database "$(DATABASE_URL)" down

seed:
	$(COMPOSE) exec app /app/seed

test:
	go test -v ./internal/...

coverage:
	go test ./internal/... -coverprofile=/tmp/coverage.out
	go tool cover -func=/tmp/coverage.out | grep total
	@echo "For HTML report, run: go tool cover -html=/tmp/coverage.out"

help:
	@echo "init | run | stop | refresh_db | logs | db_bash | redis_bash | app_bash | swagger | migrate | migrate_down | seed | test | coverage"