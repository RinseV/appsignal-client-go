.DEFAULT_GOAL := help

GO ?= go
ENV_FILE ?= .env

.PHONY: help
help: ## Show this help
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Compile the package
	$(GO) build ./...

.PHONY: test
test: ## Run the unit tests
	$(GO) test -v ./...

.PHONY: test-e2e
test-e2e: ## Run the E2E tests (needs APPSIGNAL_API_TOKEN, read from .env if present)
	@set -a; [ -f $(ENV_FILE) ] && . ./$(ENV_FILE); set +a; \
		$(GO) test -tags=e2e -v ./...

.PHONY: cover
cover: ## Run the unit tests and write coverage.out
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

.PHONY: fmt
fmt: ## Format the code
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet, including the e2e-tagged files
	$(GO) vet ./...
	$(GO) vet -tags=e2e ./...

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	$(GO) mod tidy

.PHONY: check
check: fmt vet test ## Format, vet and run the unit tests
