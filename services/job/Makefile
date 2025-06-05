BINARY_NAME=job-service-binary
DOCKER_IMAGE=job-service

.PHONY: all build test containerize

all: build test clean

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./... -cover

containerize:
	docker build -t $(DOCKER_IMAGE) .

