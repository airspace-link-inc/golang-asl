.PHONY: fmt clean deep-clean test test-watch help update

TEST := CONFIG_ENV=test go test ./...

fmt: ## Run gofmt over all go files
	find . -iname "*.go" ! -path *.git* -exec gofmt -w -s {} \;

update: ## Update dependencies
	go get -u -d ./...

clean: fmt ## gofmt, mod tidy, update packages, go generate
	go mod tidy
	go generate ./...

deep-clean: clean ## Run clean, then purge cache
	go clean -modcache -cache -i -r -x

test: fmt ## Run gofmt, run go vet, then run tests
	go vet ./...
	$(TEST)

test-update: ## test and update snapshots
	UPDATE_SNAPSHOTS=true $(TEST)

test-race: fmt ## Run gofmt, run go vet, then run tests
	go vet ./...
	CONFIG_ENV=test go test -race ./...

test-watch: ## Run tests with watchexec. Will re-run tests as you change files
	watchexec -c -i .git --exts go,geojson,json make test

help: ## Print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
