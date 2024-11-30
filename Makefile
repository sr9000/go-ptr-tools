SHELL := /bin/bash
PATH := $(GOPATH)/bin:$(PATH)

lint:
	golangci-lint run --fix
