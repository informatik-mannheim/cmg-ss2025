BINARY_NAME=job-service-binary
DOCKER_IMAGE=job-service

.PHONY: all build test containerize clean

all: build test clean

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./... -cover

containerize:
	docker build -t $(DOCKER_IMAGE) .

clean:
	rm -f $(BINARY_NAME)
	-docker image rm $(DOCKER_IMAGE):latest || true