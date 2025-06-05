# Logging Concept for Go Microservices

This document describes the structured logging concept for Go-based microservices.

---

## 🔐 Goal

* Unified logging for all services
* Structured JSON output for cloud analysis (Azure)
* Minimal setup for developers
* Traceability via context

---

## ⚙️ Log Level

| Level | Description                         | Example                            |
| ----- | ----------------------------------- | ---------------------------------- |
| debug | Developer details, for debugging    | Payload content, other variables   |
| warn  | Unexpected behavior, no crash       | "Invalid secret provided"          |
| error | Failed operations                   | "Connection failed"                |

---

## 🔎 Logging Format (JSON)

### Fields:

* `time`: Timestamp
* `level`: Log level ("DEBUG", "WARN", "ERROR")
* `msg`: Main log message
* `service`: Service name
* `jobID`, `workerID`, ... : Context fields (optional)

---

Each log entry is output as JSON with the following structure:

```json
{
  "time":"2025-06-04T23:27:03.135205854Z",
  "level":"DEBUG",
  "msg":"Heartbeat received",
  "service":"worker-gateway",
  "workerID":"worker123",
  "status":"AVAILABLE"
}
```


## 📂 Logging Setup (`pkg/logging`)

Logging is implemented using Go's `log/slog` package (Go 1.21+).

### Initialization:

```go
import "github.com/informatik-mannheim/cmg-ss2025/pkg/logging"


func main() {
    logging.Init("worker-gateway")
}
```

### Basic Usage:

```go
logging.Debug("Response received", "status", 200)
logging.Warn("Result not sent", "jodID", id)
logging.Error("Server error", "error", err)
```

### Context Support:

To trace logs through the system use context.

In a Core function:

```go
func (s *WorkerGatewayService) Heartbeat(ctx context.Context, req ports.HeartbeatRequest) ([]ports.Job, error) {
	logging.From(ctx).Debug("Heartbeat received", "workerID", req.WorkerID, "status", req.Status)
}
```

In a HTTP handler:

```go
// POST /worker/heartbeat
func (h *Handler) HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	var req ports.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
    logging.From(r.Context()).Error("Invalid HeartbeatRequest", "error", err)
		return
	}
}
```

---

## Docker Compose:

```yaml
services:
  worker-gateway:
    environment:
      - LOG_LEVEL=debug
```

## Azure Logging:

* Logs are emitted via `stdout`
* Azure Container Apps capture this automatically
* Analysis via Azure Monitor / Log Analytics:

```kusto
ContainerLog
| where LogEntry contains "worker-gateway"
| where request_id == "abc123"
```
