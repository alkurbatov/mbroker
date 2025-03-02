MAKEFLAGS = --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

SERVICE_NAME = mbroker
BUILD_DIR = dist
SEED := on

# NB (alkurbatov): Although this template has small coverage threshold
# production-ready services must have >= 80% coverage.
TEST_COVERAGE_THRESHOLD = 0

.DEFAULT_GOAL := help
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: download
download: ## Download go.mod dependencies
	@echo Downloading go.mod dependencies
	go mod download -x

.PHONY: update
update: ## Update all Golang modules at once
	go get -u ./...
	go mod tidy

.PHONY: build
build:
	./scripts/build cmd/$(SERVICE_NAME) $(BUILD_DIR)/$(SERVICE_NAME)

.PHONY: run
run: build ## Run the project
	$(BUILD_DIR)/$(SERVICE_NAME)

.PHONY: clean
clean: stop
	rm -rf $(BUILD_FOLDER)

.PHONY: lint-golang
lint-golang: ## Lint Golang source code
	golangci-lint run
	go run golang.org/x/tools/cmd/deadcode@latest -test ./... | tee deadcode.out && [ ! -s deadcode.out ]

.PHONY: lint-shell
lint-shell: ## Lint shell scripts
	shellcheck --severity=warning ./scripts/*

.PHONY: lint
lint: lint-golang lint-shell ## Lint project source

.PHONY: fmt
fmt: ## Format the source code
	go run mvdan.cc/gofumpt@latest -l -w -extra .
	go run golang.org/x/tools/cmd/goimports@latest -l -w .
	go run github.com/daixiang0/gci@latest write \
		--skip-generated \
		--custom-order \
		-s standard \
		-s default \
		-s prefix\(github.com/alkurbatov/mbroker\) \
		-s blank \
		-s dot \
		.

.PHONY: unit-tests
unit-tests: ## Run unit tests
	go test -v -race -shuffle=$(SEED) ./{internal,pkg}/... -coverprofile=coverage.out -covermode atomic
	@grep -v -E "(_mock|.pb).go" coverage.out > coverage.out.tmp
	@mv coverage.out.tmp coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
	@./scripts/check-coverage coverage.out $(TEST_COVERAGE_THRESHOLD)

.PHONY: docs
docs: ## View project documentation
	@echo "Project and packages documentation available at:"
	@echo -e "\thttp://127.0.0.1:3000/pkg/"
	@go run golang.org/x/tools/cmd/godoc@v0.30.0 -http=:3000 -index > /dev/null
