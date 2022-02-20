THIS_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

deps:
	go get -u
	go mod tidy

build:
	go build -o out/ $(THIS_DIR)

run:
	go run $(THIS_DIR)

test:
	go test $(THIS_DIR)
