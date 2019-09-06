.DEFAULT_GOAL := help

_YELLOW=\033[0;33m
_NC=\033[0m

PKG_DIRS = $(shell go list -f '{{.Dir}}' ./...)

.PHONY: help setup # generic commands
help: ## prints this help
	@grep -hE '^[\.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "${_YELLOW}%-16s${_NC} %s\n", $$1, $$2}'

setup: ## downloads dependencies
	go get -u golang.org/x/lint/golint
	go get -u github.com/robertkrimen/godocdown/godocdown
	go mod tidy

.PHONY: doc autogenerate lint unit-test # go commands
doc: ## generates markdown documentation
	for d in ${PKG_DIRS}; do godocdown -o $$d/README.md $$d; done

autogenerate: ## autogenerates code
	go generate ./...

lint: ## runs the code linter
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

unit-test: ## runs unit tests
	go test -cover -v ./...

.PHONY: cicd # pipeline commands
cicd: ## runs the CI/CD pipeline
	@make lint unit-test

