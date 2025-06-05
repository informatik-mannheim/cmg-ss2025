BINARY_NAME=worker-registry-binary
DOCKER_IMAGE=worker-registry

.PHONY: all build test containerize run fmt docker-up docker-down clean integrationcheck deployment

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
	go run main.go

docker-up:
	docker compose up --build -d 

docker-down:
	docker compose down

clean:
	rm -f $(BINARY_NAME)
