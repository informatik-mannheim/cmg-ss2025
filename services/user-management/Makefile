BINARY_NAME=user-management-binary
DOCKER_IMAGE=user-management

.PHONY: all build test containerize run fmt docker-down clean integrationcheck deployment junit coverage

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
	docker compose up --build

fmt:
	gofmt -w .

docker-down:
	docker compose down
