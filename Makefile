all: lint test
PHONY: test coverage lint golint clean vendor local-dev-databases docker-up docker-down integration-test unit-test
GOOS=linux
DB_STRING=host=localhost port=26257 user=root sslmode=disable
DEV_DB=${DB_STRING} dbname=hollow_dev
TEST_DB=${DB_STRING} dbname=hollow_test


test: | unit-test integration-test

integration-test: docker-up test-database
	@echo Running integration tests...
	@HOLLOW_TEST_DB="${TEST_DB}" go test -cover -tags testtools,integration ./... -p 1

unit-test: | lint
	@echo Running unit tests...
	@go test -cover -short -tags testtools ./...

coverage: | docker-up test-database
	@echo Generating coverage report...
	@go test ./... -race -coverprofile=.dev-data/coverage.out -covermode=atomic -tags testtools -p 1
	@go tool cover -func=.dev-data/coverage.out
	@go tool cover -html=.dev-data/coverage.out

lint: golint

golint: | vendor
	@echo Linting Go files...
	@golangci-lint run

clean: docker-down
	@echo Cleaning...
	@rm -rf ./dist/
	@rm -rf ./.dev-data/coverage.out
	@rm -rf ./.dev-data/compose/db/*
	@go clean -testcache

vendor:
	@go mod download

docker-up:
	@docker-compose up -d db

docker-down:
	@docker-compose down

dev-database:
	@docker exec hollow_db_1 cockroach sql --insecure -e "drop database if exists hollow_dev"
	@docker exec hollow_db_1 cockroach sql --insecure -e "create database hollow_dev"
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="${DEV_DB}" goose -dir=db/migrations up

test-database:
	@docker exec hollow_db_1 cockroach sql --insecure -e "drop database if exists hollow_test"
	@docker exec hollow_db_1 cockroach sql --insecure -e "create database hollow_test"
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="${TEST_DB}" goose -dir=db/migrations up
