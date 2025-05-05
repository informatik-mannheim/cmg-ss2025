# Step 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the binary
RUN go build -o user-management .

# Step 2: Minimal runtime image
FROM scratch

# App listens on 8081
EXPOSE 8081

# Copy the statically compiled binary
COPY --from=builder /app/user-management /app/user-management

# Optional: copy static files like .env or users.json if you want to bundle them
# COPY --from=builder /app/users.json /app/users.json
# COPY --from=builder /app/.env /app/.env

CMD ["/app/user-management"]