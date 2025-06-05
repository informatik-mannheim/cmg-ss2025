BINARY_NAME=worker-deamon-binary
DOCKER_IMAGE=worker-deamon

.PHONY: all build test 

integrationcheck: build test 

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...


