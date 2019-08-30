
.PHONY: test
test: lint ## Run tests with coverage
	go test -race -short -cover -coverprofile coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html

.PHONY: lint
lint: ## Run the code linter
	revive ./...

.PHONY: clean
clean: ## Clean all test files
	rm -rf coverage*

.PHONY: help
help:
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) |sort \
		|awk 'BEGIN{FS=":.*?## "};{printf "\033[36m%-30s\033[0m %s\n",$$1,$$2}'
