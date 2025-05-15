# Job Scheduler Service

This service Schedules Jobs from the `JobService` and assignes them workers from the `WorkerRegistry` while using the carbon intensity data from the `CarbonIntensityProvider`.

This service has no API.

> **WARNING**
> The implementation is in an early stage. Some functionality may be missing or subject to change.

---

## Architecture

- `adapter/`: Handles HTTP Requests and contains the repository implementation for the in-memory-database.
- `core/`: Implements the logic for the `JobScheduler`
- `model/`: Contains all the data structures needed for `JobScheduler`
- `ports/`: Defines the `JobScheduler` interface as well as the `Adapter` interfaces
- `utils/`: Contains some utility functions to make life easier
- `main.go`: Wires everything and starts the server

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
