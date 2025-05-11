# Step 1: Build the Go binary
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

# Copy go mod and sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the binary
RUN go build -o user-management .

# Step 2: Create a minimal runtime image
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/user-management /app/user-management

EXPOSE 8080 

CMD ["/app/user-management"]