all: lint test
PHONY: test coverage lint golint clean vendor
GOOS=linux


test: | lint
	@echo Testing...
	@go test -cover ./...

coverage: | vendor
	@echo Generating coverage report...
	@go test -coverprofile=.dev-data/coverage.out  ./...
	@go tool cover -func=.dev-data/coverage.out
	@go tool cover -html=.dev-data/coverage.out

lint: golint

golint: | vendor
	@echo Linting Go files...
	@golangci-lint run

clean:
	@echo Cleaning...
	@rm -rf ./out/

vendor:
	@go mod download
