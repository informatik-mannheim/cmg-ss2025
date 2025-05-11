.PHONY: build run test fmt docker-up docker-down docker-build clean

build:
	go build -o user-management main.go

run:
	go run main.go

test:
	go test ./... -cover

fmt:
	gofmt -w .

docker-build:
	docker build -t user-management .

docker-up:
	docker compose up --build

docker-down:
	docker compose down

clean:
	rm -f user-management