all: lint test
PHONY: test coverage lint golint clean vendor local-dev-databases docker-up docker-down integration-test unit-test
GOOS=linux
DB_STRING=host=localhost port=26257 user=root sslmode=disable
DEV_DB=${DB_STRING} dbname=serverservice
TEST_DB=${DB_STRING} dbname=serverservice_test

test: | unit-test integration-test

integration-test: test-database
	@echo Running integration tests...
	@SERVERSERVICE_CRDB_URI="${TEST_DB}" go test -cover -tags testtools,integration -p 1 ./...

unit-test: | lint
	@echo Running unit tests...
	@go test -cover -short -tags testtools ./...

coverage: | test-database
	@echo Generating coverage report...
	@SERVERSERVICE_CRDB_URI="${TEST_DB}" go test ./... -race -coverprofile=coverage.out -covermode=atomic -tags testtools -p 1
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

lint: golint

golint: | vendor
	@echo Linting Go files...
	@golangci-lint run

clean: docker-clean
	@echo Cleaning...
	@rm -rf ./dist/
	@rm -rf coverage.out
	@go clean -testcache

vendor:
	@go mod download
	@go mod tidy

docker-up:
	@docker-compose -f quickstart.yml up -d crdb

docker-down:
	@docker-compose -f quickstart.yml down

docker-clean:
	@docker-compose -f quickstart.yml down --volumes

dev-database: | vendor
	@cockroach sql --insecure -e "drop database if exists serverservice"
	@cockroach sql --insecure -e "create database serverservice"
	@SERVERSERVICE_CRDB_URI="${DEV_DB}" go run main.go migrate up

test-database: | vendor
	@cockroach sql --insecure -e "drop database if exists serverservice_test"
	@cockroach sql --insecure -e "create database serverservice_test"
	@SERVERSERVICE_CRDB_URI="${TEST_DB}" go run main.go migrate up
	@cockroach sql --insecure -e "use serverservice_test; ALTER TABLE attributes DROP CONSTRAINT check_server_id_server_component_id; ALTER TABLE versioned_attributes DROP CONSTRAINT check_server_id_server_component_id;"
