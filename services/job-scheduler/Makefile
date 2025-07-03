BINARY_NAME=job-scheduler-binary
DOCKER_IMAGE=job-scheduler

.PHONY: all build containerize test run docker-up docker-down clean integrationcheck deployment junit coverage

integrationcheck: clean test build

deployment: clean test build containerize

clean:
	rm -f $(BINARY_NAME)
	rm -f report.xml
	rm -f coverage.out
	rm -f coverage.xml

test: junit coverage

junit:
	go test -v ./... | go-junit-report > report.xml

coverage:
	go test -coverprofile=coverage.out ./...
	gocover-cobertura < coverage.out > coverage.xml

build:
	go build -o $(BINARY_NAME) .

containerize:
	docker build -t $(DOCKER_IMAGE) .


run:
	WORKER_REGISTRY="http://localhost:8080" JOB_SERVICE="http://localhost:8080" CARBON_INTENSITY_PROVIDER="http://localhost:8080" go run main.go

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down










