# step 1: build the service in a builder container
FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app

# Copy only necessary files
COPY services/consumer-gateway/go.mod .
COPY services/consumer-gateway/go.sum .
RUN go mod download

COPY services/consumer-gateway/ .

RUN go build -o /consumer-gateway .

# step 2: create lightweight runtime image
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
EXPOSE 8080
COPY --from=builder /consumer-gateway /app/consumer-gateway
CMD ["/app/consumer-gateway"]