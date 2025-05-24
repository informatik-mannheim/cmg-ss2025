# Worker Gateway Service

The **Worker Gateway Service** manages the communication between:
- **Worker Daemon(s)** – execute compute jobs
- **Job Service** – manages job assignments and statuses
- **Worker Registry** – tracks worker availability

The gateway provides worker with jobs, handles job results, receives heartbeats from workers, and sends state updates across services.


> **WARNING**
> The implementation is in an early stage. Many things are still missing. Use with care.


## Building

`go build .`

## Testing

`go test ./...`

## Containerizing

`docker build -t worker-gateway .`

**WARNING**: Does not work inside the dev container

## Running

Inside the dev container: `go run .`

As a container: 

`docker run -e PORT=8080 -p 8080:8080 worker-gateway`

## Usage
### Register a Worker
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "id": "worker123",
  "key": "12345678",
  "zone": "DE"
}' http://localhost:8080/register
```

### Send Heartbeat
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "workerId": "worker123",
  "status": "AVAILABLE"
}' http://localhost:8080/worker/heartbeat
```

### Submit Job Result
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "jobId": "job456",
  "status": "completed",
  "result": "Job Computed.",
  "errorMessage": ""
}' http://localhost:8080/result
```

