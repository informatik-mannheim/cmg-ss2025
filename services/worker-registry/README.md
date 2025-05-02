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

## API Overview
### `GET /workers`

Returns the list of all registered workers.

#### Example Command

```bash
curl -X 'GET' \
  'http://localhost:8080/workers' \
  -H 'accept: application/json'
```
#### Sample Response

```json
[
  {
    "id": "101a8fa4-dce6-4d10-a65a-e88e6edbcdf7",
    "status": "RUNNING",
    "zone": "DE"
  },
  {
    "id": "06842e74-1121-4700-b6d5-6558c8af6199",
    "status": "AVAILABLE",
    "zone": "EN"
  }
]
```
### `GET /workers?status=AVAILABLE&zone=DE`

Returns the list of all registered workers in the zone `DE` that are `AVAILABLE`.

#### Example Command

```bash
curl -X 'GET' \
  'http://localhost:8080/workers?status=AVAILABLE&zone=DE' \
  -H 'accept: application/json'
```
#### Sample Response

```json
[
  {
    "id": "101a8fa4-dce6-4d10-a65a-e88e6edbcdf7",
    "status": "AVAILABLE",
    "zone": "DE"
  },
  {
    "id": "42842c54-1121-4700-b6d5-6331c8af616a",
    "status": "AVAILABLE",
    "zone": "DE"
  }
]
```

### `GET /workers/{id}`

Returns a worker with the specified `id`.

#### Example Command

```bash
curl -X 'GET' \
  'http://localhost:8080/workers/101a8fa4-dce6-4d10-a65a-e88e6edbcdf7' \
  -H 'accept: application/json'
```

#### Example Response

```json
[
  {
    "id": "101a8fa4-dce6-4d10-a65a-e88e6edbcdf7",
    "status": "AVAILABLE",
    "zone": "DE"
  }
]
```

### `POST /workers`

Creates a worker from given `zone`.

#### Example Command
```bash
curl -X 'POST' \
  'http://localhost:8080/workers?zone=EN' \
  -H 'accept: application/json' \
  -d ''
```
#### Example Response
```json
[
  {
  "id": "90bb1e74-22f1-4b91-bf0b-fd17e542cb3e",
  "status": "AVAILABLE",
  "zone": "EN"
  }
]
```

### `PUT /workers/{id}/status`

Updates the `status` of a specific worker (`AVAILABLE` or `RUNNING`).

#### Example Command
```bash
curl -X 'PUT' \
  'http://localhost:8080/workers/101a8fa4-dce6-4d10-a65a-e88e6edbcdf7/status' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "workerStatus": "RUNNING"
}'
```
#### Example Response
```json
[
  {
    "id": "101a8fa4-dce6-4d10-a65a-e88e6edbcdf7",
    "status": "RUNNING",
    "zone": "DE"
  }
]
```

