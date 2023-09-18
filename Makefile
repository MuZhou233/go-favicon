SHELL := /bin/bash

VERSION ?= $(shell git describe --tags --always --abbrev=10)
BUILD_DATE ?= $(shell date)
LDFLAGS ?= -X "main.version=$(VERSION)" -X "main.buildDate=$(BUILD_DATE)"

GOFILES := $(shell find . -type f -name '*.go' )
GOFMT ?= gofmt -s -w

.PHONY: init
init:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: format
format:
	$(GOFMT) $(GOFILES)

.PHONY: golint
golint:
	golint -set_exit_status ./...

.PHONY: golangci
golangci:
	golangci-lint run

.PHONY: lint
lint: golint golangci

.PHONY: test
test:
	go test -json -cover -coverprofile cover.out | tparse

.PHONY: test-coverage
test-coverage:
	go tool cover -html=cover.out -o cover.html

.PHONY: build
build:
	go build -ldflags '$(LDFLAGS)' -v ./cmd/favicon

.PHONY: clean
clean:
	go clean
	rm -rf ./favicon

.PHONY: run
run:
	go run ./cmd/favicon/main.go

.PHONY: all
all: tidy format lint test build

.PHONY: checks
checks: tidy format lint
