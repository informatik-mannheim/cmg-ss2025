FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY . .

WORKDIR /app/services/worker-gateway
RUN go mod download
RUN go build -o worker-gateway

FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/services/worker-gateway/worker-gateway .

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["./worker-gateway"]