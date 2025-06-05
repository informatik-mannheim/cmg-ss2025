#step 1: build the service in a builder container
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app
COPY . .

WORKDIR /app/services/consumer-gateway
RUN go mod download
RUN go build -o consumer-gateway

#step 2: copy the binary (and nothing else) into the run container
FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/services/consumer-gateway/consumer-gateway .

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["./consumer-gateway"]
