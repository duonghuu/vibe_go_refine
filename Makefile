.PHONY: build up down logs restart clean migrate-up migrate-down migrate-docker

# Project variables
COMPOSE_FILE = docker-compose.yml
MIGRATION_PATH = apps/backend/database/migrations
DB_URL = mysql://root:root@tcp(127.0.0.1:3306)/vibe_go_refine

build:
	docker compose -f $(COMPOSE_FILE) build

up:
	docker compose -f $(COMPOSE_FILE) up -d

down:
	docker compose -f $(COMPOSE_FILE) down

logs:
	docker compose -f $(COMPOSE_FILE) logs -f

restart:
	docker compose -f $(COMPOSE_FILE) restart

clean:
	docker compose -f $(COMPOSE_FILE) down -v --remove-orphans

migrate-up:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down

migrate-docker:
	docker run --rm -v $(PWD)/$(MIGRATION_PATH):/migrations --network host migrate/migrate -path=/migrations -database "$(DB_URL)" up
