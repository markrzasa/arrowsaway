THIS_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

deps:
	go get -u
	go mod tidy -compat=1.17

build: deps
	go build -o out/ $(THIS_DIR)

run:
	go run $(THIS_DIR)

test:
	go test $(THIS_DIR)
