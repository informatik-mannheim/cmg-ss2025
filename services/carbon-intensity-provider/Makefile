BINARY_NAME=carbon-intensity-provider-binary
DOCKER_IMAGE=carbon-intensity-provider

.PHONY: all build test containerize clean integrationcheck deployment junit coverage

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
