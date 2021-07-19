all: lint test
PHONY: test coverage lint golint clean vendor local-dev-databases docker-up docker-down integration-test unit-test
GOOS=linux


test: | unit-test integration-test

integration-test: docker-up local-dev-databases
	@echo Running integration tests...
	@go test -cover -tags testtools,integration ./... -p 1

unit-test: | lint
	@echo Running unit tests...
	@go test -cover -short -tags testtools ./...

coverage: | docker-up local-dev-databases
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

vendor:
	@go mod download

docker-up:
	@docker-compose up -d db

docker-down:
	@docker-compose down

local-dev-databases:
	@docker exec hollow_db_1 cockroach sql --insecure -e "create database if not exists hollow_dev"
	@docker exec hollow_db_1 cockroach sql --insecure -e "create database if not exists hollow_test"
