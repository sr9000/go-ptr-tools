SHELL := /bin/bash
PATH := $(GOPATH)/bin:$(PATH)

.PHONY: all clean lint test cover bench

lint:
	golangci-lint run --fix

test:
	go test -race ./... -skip="^.+_rc$$"

	# race condition test
	go test -race -run=Example_wrong_unprotectedConcurrentAccess_rc ./... >/dev/null 2>&1; \
		test $$? -eq 1 || echo "expected race condition to happen"

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

bench: # includes tests
	go test -bench=. ./... -run=^$$

clean:
	go clean -testcache

all: clean lint test bench
