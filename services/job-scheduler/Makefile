BINARY_NAME=job-scheduler-binary
DOCKER_IMAGE=job-scheduler

.PHONY: all build containerize test run docker-up docker-down  

integrationcheck: build test 
deployment: build test containerize

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

containerize:
	docker build -t $(DOCKER_IMAGE) .

run:
	WORKER_REGISTRY="http://localhost:8080" JOB_SERVICE="http://localhost:8080" CARBON_INTENSITY_PROVIDER="http://localhost:8080" go run main.go

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down








