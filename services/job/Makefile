BINARY_NAME=job-binary
DOCKER_IMAGE=job

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

