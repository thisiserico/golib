.DEFAULT_GOAL := help

_YELLOW=\033[0;33m
_NC=\033[0m

PKG_DIRS = $(shell go list -f '{{.Dir}}' ./...)

.PHONY: help setup doc # generic commands
help: ## prints this help
	@grep -hE '^[\.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "${_YELLOW}%-16s${_NC} %s\n", $$1, $$2}'

setup: ## downloads dependencies
	GO111MODULE=off go get golang.org/x/lint/golint
	GO111MODULE=off go get github.com/robertkrimen/godocdown/godocdown

doc: ## generates markdown documentation
	for d in ${PKG_DIRS}; do godocdown -o $$d/README.md $$d; done


.PHONY: generate lint test # go commands
generate: ## generates code
	go generate ./...

lint: ## runs the code linter
	go list ./... | xargs golint -set_exit_status

test: ## runs tests
	go test -count=1 -cover -v ./...

