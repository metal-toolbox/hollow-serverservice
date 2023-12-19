all: lint test
PHONY: test coverage lint golint clean vendor local-dev-databases docker-up docker-down integration-test unit-test
GOOS=linux
DB_STRING=host=localhost port=26257 user=root sslmode=disable
DEV_DB=${DB_STRING} dbname=serverservice
TEST_DB=${DB_STRING} dbname=serverservice_test
DOCKER_IMAGE  := "ghcr.io/metal-toolbox/fleetdb"

## run all tests
test: | unit-test integration-test

## run integration tests
integration-test: test-database
	@echo Running integration tests...
	@SERVERSERVICE_CRDB_URI="${TEST_DB}" go test -cover -tags testtools,integration -p 1 ./...

## run lint and unit tests
unit-test: | lint
	@echo Running unit tests...
	@SERVERSERVICE_CRDB_URI="${TEST_DB}" go test -cover -short -tags testtools ./...

## check test coverage
coverage: | test-database
	@echo Generating coverage report...
	@SERVERSERVICE_CRDB_URI="${TEST_DB}" go test ./... -race -coverprofile=coverage.out -covermode=atomic -tags testtools -p 1
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

## lint
lint: golint

golint: | vendor
	@echo Linting Go files...
	@golangci-lint run

## clean docker files
clean: docker-clean
	@echo Cleaning...
	@rm -rf ./dist/
	@rm -rf coverage.out
	@go clean -testcache

vendor:
	@go mod download
	@go mod tidy

## setup docker compose test env
docker-up:
	@docker-compose -f quickstart.yml up -d crdb

## stop docker compose test env
docker-down:
	@docker-compose -f quickstart.yml down

## clean docker volumes
docker-clean:
	@docker-compose -f quickstart.yml down --volumes

## setup devel database
dev-database: | vendor
	@cockroach sql --insecure -e "drop database if exists serverservice"
	@cockroach sql --insecure -e "create database serverservice"
	@SERVERSERVICE_CRDB_URI="${DEV_DB}" go run main.go migrate up

## setup test database
test-database: | vendor
	@cockroach sql --insecure -e "drop database if exists serverservice_test"
	@cockroach sql --insecure -e "create database serverservice_test"
	@SERVERSERVICE_CRDB_URI="${TEST_DB}" go run main.go migrate up
	@cockroach sql --insecure -e "use serverservice_test; ALTER TABLE attributes DROP CONSTRAINT check_server_id_server_component_id; ALTER TABLE versioned_attributes DROP CONSTRAINT check_server_id_server_component_id;"


## Build linux bin
build-linux:
	GOOS=linux GOARCH=amd64 go build -o fleetdb

## build docker image and tag as ghcr.io/metal-toolbox/fleetdb:latest
build-image: build-linux
	docker build --rm=true -f Dockerfile -t ${DOCKER_IMAGE}:latest  . \
							 --label org.label-schema.schema-version=1.0 \
							 --label org.label-schema.vcs-ref=$(GIT_COMMIT_FULL) \
							 --label org.label-schema.vcs-url=$(REPO)

## build and push devel docker image to KIND image repo used by the sandbox - https://github.com/metal-toolbox/sandbox
push-image-devel: build-image
	docker tag ${DOCKER_IMAGE}:latest localhost:5001/fleetdb:latest
	docker push localhost:5001/fleetdb:latest
	kind load docker-image localhost:5001/fleetdb:latest


# https://gist.github.com/prwhite/8168133
# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


TARGET_MAX_CHAR_NUM=20
## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
