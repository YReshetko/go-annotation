-include .env
export

.PHONY: help
help: ## Display the Makefile helper
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[35m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@printf "\033[31m'*' - The targets are not runnable in make container\n \033[0m"

.PHONY: setup
setup: ## Downloads and installs all libraries and dependencies for the project
	@go mod tidy
	@go get -d -v ./...
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

.PHONY: test
test: ## Runs unit tests against the codebase
	go clean -testcache
	go test ./...

.PHONY: lint
lint: ## Runs linter against the service codebase
	@golangci-lint run --config golangci-lint.yaml
	@printf "[✔ ] \033[32mLinter passed\033[0m\n"
