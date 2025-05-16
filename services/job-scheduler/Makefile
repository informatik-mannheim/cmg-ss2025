.PHONY: run test fmt docker-up docker-down docker-build clean

run:
	WORKER_REGISTRY="http://localhost:8080" JOB_SERVICE="http://localhost:8080" CARBON_INTENSITY_PROVIDER="http://localhost:8080" go run main.go

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

test:
	go test ./... -cover