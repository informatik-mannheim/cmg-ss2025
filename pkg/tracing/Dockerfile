FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# binary
RUN go build -o tracing .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/tracing .

# ports
ENV PORT=8080
EXPOSE 8080

# start
ENTRYPOINT ["./tracing"]