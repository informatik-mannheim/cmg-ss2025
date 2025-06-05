BINARY_NAME=cli-binary

.PHONY: all build test

integrationcheck: build test 

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

