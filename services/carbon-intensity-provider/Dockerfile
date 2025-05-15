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

# Copy persistent storage files (optional, if you want to seed local files)
COPY zones.json zones.json
COPY zones_metadata.json zones_metadata.json

# Expose the service port
EXPOSE 8080

# Set environment variables via docker-compose
CMD ["./carbon-intensity-provider"]