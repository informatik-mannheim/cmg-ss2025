BINARY_NAME=cli-binary

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
	go build -o $(BINARY_NAME) .




run:
	docker	run	--rm -it	consumer-cli

run-local:
	docker run --rm -it -e GATEWAY_URL=http://localhost:8080/ consumer-cli

build-image:
	docker build -t consumer-cli .

