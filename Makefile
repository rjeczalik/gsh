all: build cmd test

.PHONY: test
test: ## Run tests
	@go test -race -v ./...

.PHONY: build
build: ## Build package
	@go build -race -v ./...

.PHONY: cmd
cmd: ## Build package
	@go build -o gossh ./cmd/gossh

.PHONY: help
help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[36m%-25s\033[0m %s\n", $$1, $$NF}' $(MAKEFILE_LIST)
