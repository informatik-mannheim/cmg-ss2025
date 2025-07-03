BINARY_NAME=cli-binary

.PHONY: all build test clean integrationcheck

integrationcheck:
	@$(MAKE) build && \
	$(MAKE) test ; \
	$(MAKE) clean

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)

run:
	docker	run	--rm -it	consumer-cli

run-local:
	docker run --rm -it -e GATEWAY_URL=http://localhost:8080/ consumer-cli

build-image:
	docker build -t consumer-cli .

