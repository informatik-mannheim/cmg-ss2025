BINARY_NAME=carbon-intensity-provider-binary
DOCKER_IMAGE=carbon-intensity-provider

.PHONY: all build test containerize clean integrationcheck deployment

integrationcheck:
	@$(MAKE) build && \
	$(MAKE) test ; \
	$(MAKE) clean
	
deployment:
	@$(MAKE) build && \
	$(MAKE) test && \
	$(MAKE) containerize ; \
	$(MAKE) clean

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

containerize:
	docker build -t $(DOCKER_IMAGE) .

clean:
	rm -f $(BINARY_NAME)