APP_CONTAINER=task_manager_app
DB_CONTAINER=task_manager_postgres
REDIS_CONTAINER=task_manager_redis

COMPOSE=docker compose
ENV_FILE=.env

.PHONY: help init run stop logs refresh_db db_bash redis_bash app_bash swagger

help:
	@echo "Available commands:"
	@echo "  make init         Build images, start containers"
	@echo "  make run          Start containers"
	@echo "  make stop         Stop containers"
	@echo "  make logs         Show logs"
	@echo "  make refresh_db   Recreate database"
	@echo "  make db_bash      Enter postgres container"
	@echo "  make redis_bash   Enter redis container"
	@echo "  make app_bash     Enter app container"
	@echo "  make swagger      Rebuild swagger"

init:
	$(COMPOSE) --env-file $(ENV_FILE) up -d --build

run:
	$(COMPOSE) --env-file $(ENV_FILE) up -d

stop:
	$(COMPOSE) --env-file $(ENV_FILE) down

logs:
	$(COMPOSE) logs -f

refresh_db:
	$(COMPOSE) stop postgres
	$(COMPOSE) rm -f postgres
	$(COMPOSE) up -d postgres

db_bash:
	docker exec -it $(DB_CONTAINER) sh

redis_bash:
	docker exec -it $(REDIS_CONTAINER) sh

app_bash:
	docker exec -it $(APP_CONTAINER) sh

swagger:
	@echo "Swagger will be generated in next phases"
