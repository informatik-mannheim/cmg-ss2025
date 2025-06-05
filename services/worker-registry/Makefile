BINARY_NAME=worker-registry-binary
DOCKER_IMAGE=worker-registry

.PHONY: all build test containerize run fmt docker-up docker-down

integrationcheck: build test 
deployment: build test containerize

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

containerize:
	docker build -t $(DOCKER_IMAGE) .

run:
	go run main.go

docker-up:
	docker compose up --build -d 

docker-down:
	docker compose down
