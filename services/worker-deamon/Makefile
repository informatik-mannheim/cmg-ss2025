BINARY_NAME=worker-deamon-binary
DOCKER_IMAGE=worker-deamon

.PHONY: all build test integrationcheck clean

integrationcheck:
	@$(MAKE) build && \
	$(MAKE) test ; \
	$(MAKE) clean

build:
	go build -o $(BINARY_NAME) ./cmd/worker-deamon

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)

