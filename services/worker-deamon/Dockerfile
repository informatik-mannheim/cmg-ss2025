FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o worker-daemon ./cmd/worker-deamon

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/worker-daemon .

CMD ["./worker-daemon"]