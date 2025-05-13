.PHONY: run test fmt docker-up docker-down docker-build clean

docker-build:
	docker build -t worker-registry .

docker-up:
	docker run -d -p 8080:8080 --name worker-registry worker-registry

docker-down:
	docker stop worker-registry && docker rm worker-registry

run:
	go run main.go

test:
	go test ./... -cover

fmt:
	gofmt -w .

clean:
	rm -f worker-registry