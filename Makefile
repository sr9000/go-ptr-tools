SHELL := /bin/bash
PATH := $(GOPATH)/bin:$(PATH)

lint:
	golangci-lint run --fix

test:
	go test -race ./... -skip="^.+_rc$$"

	# race condition test
	go test -race -run=Example_wrong_unprotectedConcurrentAccess_rc ./... >/dev/null 2>&1; \
		test $$? -eq 1 || echo "expected race condition to happen"

bench: # includes tests
	go test -bench=. ./... -run=^$$

clean:
	go clean -testcache

all: clean lint test bench
