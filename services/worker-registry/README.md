# Worker Registry Service

This service provides an API to store and retrieve worker data, following the Ports & Adapters (Hexagonal) architecture.

> **WARNING**
> The implementation is in an early stage. Some functionality may be missing or subject to change.

---

## Architecture
- `adapter/`: Handles HTTP Requests and contains the repository implementation for the in-memory-database.
- `api/`: contains the OpenApi specification for the `WorkerRegistry` endpoints.
- `core/`: Implements the logic for the `WorkerRegistry`
- `main.go`: Wires everything and starts the server
- `model/`: Contains the data structure for `Worker`
- `ports/`: Defines the `WorkerRegistry` interface

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

---


## API Overview
### `GET /workers`

Returns the list of all registered workers.

#### Example Command

```bash
curl 'localhost:8080/workers'
```
#### Sample Response

```json
[
  {
    "id": "5fda654b-3343-42ae-bab2-0faeffb78f2e",
    "status": "RUNNING",
    "zone": "DE"
  },
  {
    "id": "a35c739a-ee51-4cdc-9393-f75a8f1d76ad",
    "status": "AVAILABLE",
    "zone": "EN"
  }
]
```
### `GET /workers?status=AVAILABLE&zone=DE`

Returns the list of all registered workers in the zone `DE` that are `AVAILABLE`.

#### Example Command

```bash
curl 'localhost:8080/workers?status=AVAILABLE&zone=DE'
```
#### Sample Response

```json
[
  {
    "id": "5fda654b-3343-42ae-bab2-0faeffb78f2e",
    "status": "AVAILABLE",
    "zone": "DE"
  },
  {
    "id": "a35c739a-ee51-4cdc-9393-f75a8f1d76ad",
    "status": "AVAILABLE",
    "zone": "DE"
  }
]
```

### `GET /workers/{id}`

Returns a worker with the specified `id`.

#### Example Command

```bash
curl 'localhost:8080/workers/a35c739a-ee51-4cdc-9393-f75a8f1d76ad'
```

#### Example Response

```json
[
  {
    "id": "a35c739a-ee51-4cdc-9393-f75a8f1d76ad",
    "status": "AVAILABLE",
    "zone": "DE"
  }
]
```

### `POST /workers`

Creates a worker from given `zone`.

#### Example Command
```bash
curl -X 'POST' 'localhost:8080/workers?zone=EN'
```
#### Example Response
```json
[
  {
  "id": "5fda654b-3343-42ae-bab2-0faeffb78f2e",
  "status": "AVAILABLE",
  "zone": "EN"
  }
]
```

### `PUT /workers/{id}/status`

Updates the `status` of a specific worker (`AVAILABLE` or `RUNNING`).

#### Example Command
```bash
curl -X 'PUT' 'localhost:8080/workers/5fda654b-3343-42ae-bab2-0faeffb78f2e/status' -d '{"status": "RUNNING"}'
```
#### Example Response
```json
[
  {
    "id": "5fda654b-3343-42ae-bab2-0faeffb78f2e",
    "status": "RUNNING",
    "zone": "DE"
  }
]
```

