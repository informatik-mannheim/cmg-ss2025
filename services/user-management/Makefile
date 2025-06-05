BINARY_NAME=user-management-binary
DOCKER_IMAGE=user-management

.PHONY: all build test containerize run fmt docker-down clean integrationcheck deployment

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

run:
	docker compose up --build

fmt:
	gofmt -w .

docker-down:
	docker compose down

clean:
	rm -f $(BINARY_NAME)