#step 1: build the service in a builder container
FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go get ./...
RUN go build -o /worker-registry .

#step 2: copy the binary (and nothing else) into the run container
FROM alpine:3.21

# Add CA certs (for HTTPS)
RUN apk add --no-cache ca-certificates

EXPOSE 8080
COPY --from=builder /worker-registry /app/worker-registry
CMD ["/app/worker-registry"]
