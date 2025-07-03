BINARY_NAME=worker-deamon-binary

.PHONY: all build test containerize clean integrationcheck junit coverage

integrationcheck: clean test build

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
	go build -o $(BINARY_NAME) ./cmd/worker-deamon



