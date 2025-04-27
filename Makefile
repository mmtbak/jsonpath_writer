SHELL := /bin/bash
GO := go
GOTEST ?= go test
.PHONY: help

all: test

.PHONY: lint
lint:
	golangci-lint run --allow-parallel-runners


.PHONY: test
test:
	@echo "start unittest........"
	${GOTEST} -v ./...
