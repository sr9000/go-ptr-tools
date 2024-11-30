SHELL := /bin/bash
PATH := /home/sr9000/go/bin:$(PATH)

lint:
	golangci-lint run --fix
