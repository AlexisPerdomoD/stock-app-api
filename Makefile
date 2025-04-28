DB_COMPOSE_FILE := docker-compose-db.yml
TEST_DB_COMPOSE_FILE := docker-compose-db.test.yml

up-db:
	docker-compose -f $(DB_COMPOSE_FILE) up -d

down-db:
	docker-compose -f $(DB_COMPOSE_FILE) down

populate-db:
	@go run ./cmd/populatedb

server:
	@go run ./cmd/server

start-local:
	@clear
	docker-compose -f $(DB_COMPOSE_FILE) up -d
	@sleep 1
	@go run ./cmd/populatedb
	@clear
	@go run ./cmd/server

up-test-db:
	docker-compose -f $(TEST_DB_COMPOSE_FILE) up -d

down-test-db:
	docker-compose -f $(TEST_DB_COMPOSE_FILE) down

test:
	@clear
	@echo "running tests"
	docker-compose -f $(TEST_DB_COMPOSE_FILE) up -d
	@sleep 1
	go test ./... -v | grep -v "^?"
	docker-compose -f $(TEST_DB_COMPOSE_FILE) down
