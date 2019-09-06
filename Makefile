.DEFAULT_GOAL := help

_YELLOW=\033[0;33m
_NC=\033[0m

PKG_DIRS = $(shell go list -f '{{.Dir}}' ./...)

.PHONY: help setup doc # generic commands
help: ## prints this help
	@grep -hE '^[\.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "${_YELLOW}%-16s${_NC} %s\n", $$1, $$2}'

setup: ## downloads dependencies
	go get -u golang.org/x/lint/golint
	go get -u github.com/robertkrimen/godocdown/godocdown
	go mod tidy

doc: ## generates markdown documentation
	for d in ${PKG_DIRS}; do godocdown -o $$d/README.md $$d; done


.PHONY: autogenerate lint unit-test # go commands
autogenerate: ## autogenerates code
	go generate ./...

lint: ## runs the code linter
	go list ./... | grep -v /vendor/ | xargs golint -set_exit_status

unit-test: ## runs unit tests
	go test -count=1 -cover -v ./...


.PHONY: cicd # build pipeline commands
cicd: ## runs the CI/CD pipeline
	docker build -f Dockerfile.lint .
	docker build -f Dockerfile.unit-test .

