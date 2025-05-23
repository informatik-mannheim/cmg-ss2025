# Step 1: Build stage
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build the binary
RUN go build -o carbon-intensity-provider main.go

# Step 2: Final minimal image
FROM alpine:3.18

WORKDIR /app

# Add CA certs (for HTTPS)
RUN apk add --no-cache ca-certificates

# Copy the binary
COPY --from=builder /app/carbon-intensity-provider .

# Expose the service port
EXPOSE 8080

CMD ["./carbon-intensity-provider"]