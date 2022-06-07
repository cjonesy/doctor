TAG?=""

default: help ;

.PHONY: help
help: ## This help dialog
	@echo "Please use 'make <target>' where <target> is one of:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	 awk 'BEGIN {FS = ":.*?## "}; \
	 {printf "%-15s %s\n", $$1, $$2}'
	@echo "\nCheck the Makefile to know exactly what each target is doing."

.PHONY: test
test: lint test-unit go-mod-tidy ## Run all tests and linters

.PHONY: test-unit
test-unit: ## Run unit tests
	go test -v -race ./...

.PHONY: go-mod-tidy
go-mod-tidy: ## Clean go.mod
	go mod tidy
	git diff --exit-code go.sum

.PHONY: lint
lint: ## Run linters - staticcheck, vet, gofmt
	staticcheck ./...
	go vet ./...
	test -z "$(shell gofmt -l .)"

.PHONY: test-release
test-release: ## Run a test release with goreleaser
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: clean
clean: ## Clean up any cruft left over from old builds
	rm -rf doctor dist/

.PHONY: build
build: clean ## Build the application
	CGO_ENABLED=0 go build ./cmd/doctor

.PHONY: tag
tag: ## Create and push a git tag
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)

# Requires GITHUB_TOKEN environment variable to be set
.PHONY: release
release: clean ## Create a new release with goreleaser
	goreleaser
