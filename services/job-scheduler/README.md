# Job Scheduler Service

This service Schedules Jobs from the `JobService` and assignes them workers from the `WorkerRegistry` while using the carbon intensity data from the `CarbonIntensityProvider`.

This service has no API.

> **WARNING**
> The implementation is in an early stage. Some functionality may be missing or subject to change.

---

## Architecture

- `adapter/`: Handles HTTP Requests and contains the repository implementation for the in-memory-database.
- `core/`: Implements the logic for the `JobScheduler`
- `ports/`: Defines the `JobScheduler` interface as well as the `Adapter` interfaces and all required models.
- `utils/`: Contains some utility functions to make life easier
- `main.go`: Wires everything and starts the server

---

# Environment Variables

This services uses the following environmentvariables:

| NAME                        | Required | Type   |
| --------------------------- | -------- | ------ |
| JOB_SCHEDULER_INTERVAL      | false    | Number |
| WORKER_REGISTRY             | true     | URL    |
| JOB_SERVICE                 | true     | URL    |
| CARBON_INTENSITY_PROVIDER   | true     | URL    |
| USER_MANAGEMENT_URL         | true     | URL    |
| AUTH_TOKEN                  | true     | String |
| OTEL_EXPORTER_OTLP_ENDPOINT | true     | URL    |

---

## Run Locally

```bash
make run       # go run main.go
make test      # run all unit tests with coverage
```

---

## Run with Docker

```bash
make docker-up    # start application in docker detached
make docker-down
```
