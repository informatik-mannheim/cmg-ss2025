# Worker Daemon Service

The **Worker Daemon Service** is responsible for:
- Registering itself with the Worker Gateway
- Sending regular heartbeats to indicate availability or computing status
- Receiving and executing jobs
- Sending job results back to the Gateway

It acts as a compute node that periodically contacts the central system via HTTP and reacts based on job assignments.

> **WARNING**  
> This is an early-stage implementation. Many features are minimal or simulated. Use with caution.

---

## Building
```bash
go build ./cmd/worker-deamon
```

## Testing

`go test ./...`

## Containerizing

`docker build -t worker-deamon .`

**WARNING**: Does not work inside the dev container

## Running

Inside the dev container: `go run ./cmd/worker-deamon`

Make sure a valid config.json file is present in the working directory.

As a container: 

`docker run -v $(pwd)/config.json:/config/config.json worker-deamon`

## Config
The daemon is configured using a config.json file:
```json
{
  "id": "worker123",
  "key": "12345678",
  "zone": "DE",
  "gateway_url": "http://localhost:8080",
  "heartbeat_interval_seconds": 10
}
```