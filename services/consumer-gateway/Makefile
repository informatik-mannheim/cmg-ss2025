BINARY_NAME=consumer-gateway-binary
DOCKER_IMAGE=consumer-gateway

.PHONY: all build test containerize

integrationcheck: build test 
deployment: build test containerize

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

containerize:
	docker build -t $(DOCKER_IMAGE) .