.DEFAULT_GOAL := help

_YELLOW=\033[0;33m
_NC=\033[0m

.PHONY: help
help: ## prints this help
	@grep -hE '^[\.a-zA-Z/_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "${_YELLOW}%-16s${_NC} %s\n", $$1, $$2}'

.PHONY: setup
setup: ## downloads dependencies
	go mod download
	go get github.com/mgechev/revive


.PHONY: generate
generate: ## generates code
	go generate ./...

.PHONY: lint
lint: ## runs the code linter
	revive -formatter stylish -config revive.toml ./...

.PHONY: test/unit
test/unit: ## runs unit tests
	go test -tags=unit -count=1 -race -cover -v ./...

.PHONY: test/redis
test/redis: ## runs redis integration tests
	docker-compose -f pubsub/redis/docker-compose.yml up -d
	go test -tags=redis -count=1 -race -cover -v ./pubsub/redis

