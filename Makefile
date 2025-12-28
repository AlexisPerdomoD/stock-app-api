include .env
export
# COMMANDS
DOCKER_COMPOSE := docker compose
DOCKER := docker

# ENVS
CR_MIGRATION_DIR := ./internal/infrastructure/db/migrations/cockroachdb
CR_DB_URL := cockroach://$(CR_USER):$(CR_PASSWORD)@$(CR_HOST):$(CR_PORT)/$(CR_DB)?sslmode=$(CR_SSL)

# FILES
DB_COMPOSE_FILE := docker-compose-db.yml
TEST_DB_COMPOSE_FILE := docker-compose-db.test.yml

################################################################################
# LOCAL COMMANDS
################################################################################

local-up-db:
	$(DOCKER_COMPOSE) -f $(DB_COMPOSE_FILE) up -d
	@sleep 1
	@echo "👉 check if the database is ready"

local-down-db:
	$(DOCKER_COMPOSE) -f $(DB_COMPOSE_FILE) down

local-migrate-up: check-migrate local-up-db
	@migrate -path $(CR_MIGRATION_DIR) -database "cockroach://root@localhost:26257/defaultdb?sslmode=disable" up

local-migrate-down: check-migrate 
	@migrate -path $(CR_MIGRATION_DIR) -database "cockroach://root@localhost:26257/defaultdb?sslmode=disable" down
	
local-start: local-up-db  local-migrate-up
	@go run ./cmd/server

################################################################################
# TEST COMMANDS
################################################################################

test-local-up-db:
	$(DOCKER_COMPOSE) -f $(TEST_DB_COMPOSE_FILE) up -d
	@sleep 1
	@echo "👉 check if the database is ready"

test-local-down-db:
	$(DOCKER_COMPOSE) -f $(TEST_DB_COMPOSE_FILE) down

test-local-migrate-up: check-migrate test-local-up-db
	@migrate -path $(CR_MIGRATION_DIR) -database "cockroach://root@localhost:26257/defaultdb?sslmode=disable" up

test-local-migrate-down: check-migrate
	@migrate -path $(CR_MIGRATION_DIR) -database "cockroach://root@localhost:26257/defaultdb?sslmode=disable" down
test: test-local-up-db test-local-migrate-up
	@clear
	@echo "running tests"
	go test ./... -v | grep -v "^?"
	$(DOCKER_COMPOSE) -f $(TEST_DB_COMPOSE_FILE) down

################################################################################
# DOCKER COMMANDS
################################################################################

check-migrate:
	@command -v migrate >/dev/null 2>&1 || { \
		echo "❌ migrate not installed"; \
		echo "👉 install: https://github.com/golang-migrate/migrate"; \
		exit 1; \
	}

migrate-up: check-migrate
	migrate -path $(CR_MIGRATION_DIR) -database "$(CR_DB_URL)" up

migrate-down: check-migrate
	migrate -path $(CR_MIGRATION_DIR) -database "$(CR_DB_URL)" down

populate-db: migrate-up
	@go run ./cmd/populatedb

start: migrate-up
	@go run ./cmd/server
