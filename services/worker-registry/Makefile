.PHONY: run test fmt docker-up docker-down docker-build clean

run:
	go run main.go

docker-up:
	docker compose up --build -d 

docker-down:
	docker compose down

test:
	go test ./... -cover