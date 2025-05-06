# Step 1: Build the Go binary
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

# Copy go mod and sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the binary
RUN go build -o user-management .

# Step 2: Minimal runtime image
FROM alpine:3.18

# add ca-certificates to the image
RUN apk add --no-cache ca-certificates

# Create a non-root user
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/user-management /app/user-management

# Expose the port the app runs on
EXPOSE 8081

# Command to run the binary
CMD ["/app/user-management"]
