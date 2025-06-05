BINARY_NAME=carbon-intensity-provider-binary
DOCKER_IMAGE=carbon-intensity-provider

.PHONY: all build test containerize

integrationcheck: build test 
deployment: build test containerize

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

containerize:
	docker build -t $(DOCKER_IMAGE) .