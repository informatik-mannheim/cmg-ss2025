.PHONY: run test fmt docker-down clean

run:
	docker compose up --build

test:
	go test ./... -cover

fmt:
	gofmt -w .

docker-down:
	docker compose down

clean:
	rm -f user-management