SHELL := /bin/bash
PATH := $(GOPATH)/bin:$(PATH)

lint:
	golangci-lint run --fix

test:
	go test ./...

bench: # includes tests
	go test -bench=. ./...

clean:
	go clean -testcache

all: clean lint bench
