# Job Microservice

A microservice for managing consumer jobs in a queue system. This service provides operations to create, retrieve, and update jobs, facilitating operations from both scheduler and worker perspectives.

---

## Features

- **Job Management**: Create, retrieve, and update jobs within the system.
- **Status Filtering**: Retrieve jobs based on their status.
- **Scheduler and Worker Integration**: Update specific fields relevant to job schedulers and workers.
- **Distributed Tracing**: OpenTelemetry integration for request tracing across services.
- **Structured Logging**: Comprehensive logging with configurable log levels.
- **Database Flexibility**: Support for both PostgreSQL and in-memory storage.
- **Graceful Shutdown**: Proper signal handling for clean service termination.
- **Environment-based Configuration**: Full configuration via environment variables.

---

## Technology Stack

- **Language:** Go
- **API:** REST (JSON)
- **Database:** PostgreSQL with pgx driver
- **Containerization:** Docker, Docker Compose
- **Testing:** Go test framework, in-memory repository for fast tests
- **Tracing:** OpenTelemetry with Jaeger
- **Logging:** Structured logging with configurable levels
- **HTTP Router:** Gorilla Mux

---

## API Endpoints

### Get Jobs
Retrieve one or more jobs, filtered by their status.  
**Endpoint**: `GET /jobs`  
**Parameters**:  
- `status` (optional): Filter jobs by status as a comma-separated list (e.g., `queued,scheduled`).

### Create Job
Create a new job in the queue.  
**Endpoint**: `POST /jobs`  
**Payload**: See `JobCreate` schema.

### Get Job by ID
Retrieve a specific job using its unique ID.  
**Endpoint**: `GET /jobs/{id}`

### Get Job Outcome
Retrieve the outcome/result of a job by its unique ID.  
**Endpoint**: `GET /jobs/{id}/outcome`

### Update Job (Scheduler Perspective)
Update scheduler-related fields of a job.  
**Endpoint**: `PATCH /jobs/{id}/update-scheduler`

### Update Job (Worker Perspective)
Update worker-related fields of a job.  
**Endpoint**: `PATCH /jobs/{id}/update-workerdaemon`

---

## Environment Variables

The service supports the following environment variables:

### Core Configuration
- `PORT`: HTTP server port (default: `8080`)
- `JOB_REPO_TYPE`: Repository type (`inmemory` or `postgres`, default: `inmemory`)

### Database Configuration (PostgreSQL)
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: PostgreSQL username
- `DB_PASSWORD`: PostgreSQL password
- `DB_NAME`: PostgreSQL database name
- `SSL_MODE`: SSL mode (`true` or `false`)

### Observability
- `OTEL_EXPORTER_OTLP_ENDPOINT`: OpenTelemetry collector endpoint (e.g., `http://jaeger:4318/`)
- `LOG_LEVEL`: Logging level (`debug`, `warn`, `error`)

---

## Data Schemas

**The OpenAPI specification (`api.yaml`) shows the possible endpoints and the data schemas for API requests and responses.**

**Database schema:**  
The SQL schema for the PostgreSQL database (including the `jobs` table and its fields) is located in the [`database`](../../database) directory.  
You can find the table definitions and initialization scripts in files such as `job-init.sql`.

---

## Start Service without Docker Compose

You can also run the job microservice directly from the command line without Docker Compose.  
This is useful for local development, debugging, or running the service with custom environment variables.

### 1. In-Memory Mode (no database, for development/testing)

```sh
JOB_REPO_TYPE=inmemory SSL_MODE=false OTEL_EXPORTER_OTLP_ENDPOINT="http://jaeger:4318/" LOG_LEVEL=debug go run .
```

- This starts the service using in-memory storage. All data will be lost when the service stops.

### 2. PostgreSQL Mode (connect to a running Postgres instance)

```sh
JOB_REPO_TYPE=postgres DB_HOST=postgres DB_PORT=5432 DB_USER=jobuser DB_PASSWORD=jobpass DB_NAME=jobdb SSL_MODE=false OTEL_EXPORTER_OTLP_ENDPOINT="http://jaeger:4318/" LOG_LEVEL=debug go run .
```

- This starts the service using a PostgreSQL database.  
- Make sure a compatible Postgres instance is running and accessible with the provided credentials.

---

## Running the Service with Docker Compose

You can run the job microservice together with a PostgreSQL database using Docker Compose.

### 1. Build and Start the Services

From the `/workspaces/cmg-ss2025/services/job` directory, run:

```sh
docker-compose up --build
```

This will:
- Build the jobservice Docker image.
- Start both the jobservice and a PostgreSQL database.
- Expose the jobservice on [http://localhost:8080](http://localhost:8080).
- Automatically initialize the database schema from the `database` directory.

### 2. Stopping the Services

To stop the services, press `Ctrl+C` or run:

```sh
docker-compose down
```

### 3. Data Persistence

- PostgreSQL data is persisted in a configurable host directory.
- The database schema is automatically initialized from the `database` directory via Docker Compose.

### 4. Environment Variables

The following environment variables are set automatically in `docker-compose.yaml`:

- `JOB_REPO_TYPE=postgres`
- `DB_HOST=postgres`
- `DB_PORT=5432`
- `DB_USER=jobuser`
- `DB_PASSWORD=jobpass`
- `DB_NAME=jobdb`
- `SSL_MODE=false`
- `OTEL_EXPORTER_OTLP_ENDPOINT="http://jaeger:4318/"`
- `LOG_LEVEL=debug`

---

## Service Features

### Database Connection Resilience
- **Automatic Retry**: Up to 10 connection attempts with 3-second delays
- **Fallback Strategy**: Automatic fallback to in-memory storage if PostgreSQL is unavailable
- **SSL Support**: Configurable SSL mode for secure database connections

### Observability
- **Distributed Tracing**: Full OpenTelemetry integration with Jaeger
- **Structured Logging**: Configurable log levels with context-aware logging
- **Request Tracing**: All HTTP requests are automatically traced
- **Error Tracking**: Comprehensive error logging with context

### Graceful Shutdown
- **Signal Handling**: Responds to SIGINT and SIGTERM signals
- **Clean Shutdown**: Proper resource cleanup and connection termination
- **Tracing Cleanup**: Automatic trace data flushing on shutdown

### Validation & Error Handling
- **UUID Validation**: All job IDs must be valid UUIDs
- **Status Validation**: Job status must be one of: `queued`, `scheduled`, `running`, `completed`, `failed`, `cancelled`
- **Input Validation**: Comprehensive validation for all API endpoints
- **Error Responses**: Structured error responses with appropriate HTTP status codes

---

## Example cURL Commands

### `For the Projekt-Docker-Compose(Root) Port 8089 is exposed`

### 1. GET `/jobs`: Successful retrieval of jobs (including scenarios without filters and with invalid status filters)

Retrieve all jobs (no status filters): 
```sh
curl -X GET "http://localhost:8080/jobs"
```

Retrieving jobs with a valid status filter: 
```sh
curl -G "http://localhost:8080/jobs" --data-urlencode "status=queued,completed"
```

Attempt to retrieve jobs with an invalid status filter: 
```sh
curl -G "http://localhost:8080/jobs" --data-urlencode "status=invalidStatus"
```

---

### 2. POST `/jobs`: Creating a new job (including scenarios with invalid data)

Successful creation of a job:
```sh
curl -X POST "http://localhost:8080/jobs" -H "Content-Type: application/json" -d '{
  "jobName": "Example Job",
  "creationZone": "DE",
  "image": {
    "name": "exampleApp",
    "version": "1.0"
  },
  "parameters": {
    "param1": "value1"
  }
}'
```

Incorrect creation due to invalid data (missing fields):
```sh
curl -X POST "http://localhost:8080/jobs" -H "Content-Type: application/json" -d '{
  "jobName": "",
  "creationZone": "DE"
}'
```

---

### 3. GET `/jobs/{id}`: Retrieve a specific job by ID (including scenarios with an invalid ID).

Successful retrieval of a specific job:
```sh
curl -X GET "http://localhost:8080/jobs/{id}"
```

Incorrect retrieval due to an invalid ID:
```sh
curl -X GET "http://localhost:8080/jobs/invalid-id"
```

---

### 4. GET `/jobs/{id}/outcome`: Retrieving the result of a job by ID.

Successful retrieval of the job result:
```sh
curl -X GET "http://localhost:8080/jobs/{id}/outcome"
```

Incorrect retrieval due to an invalid ID:
```sh
curl -X GET "http://localhost:8080/jobs/invalid-id/outcome"
```

---

### 5. PATCH `/jobs/{id}/update-scheduler`: Updating job properties from the scheduler's perspective.

Successful updating of a job:
```sh
curl -X PATCH "http://localhost:8080/jobs/{id}/update-scheduler" -H "Content-Type: application/json" -d '{
  "workerId": "123e4567-e89b-12d3-a456-426614174000",
  "computeZone": "DE",
  "carbonIntensity": 80,
  "carbonSavings": 20,
  "status": "scheduled"
}'
```

Incorrect update due to invalid data:
```sh
curl -X PATCH "http://localhost:8080/jobs/{id}/update-scheduler" -H "Content-Type: application/json" -d '{
  "workerId": "",
  "computeZone": "",
  "carbonIntensity": -1,
  "carbonSavings": -1,
  "status": "unknown"
}'
```

---

### 6. PATCH `/jobs/{id}/update-workerdaemon`: Updating job properties from the perspective of a worker daemon.

Successful updating of a job:
```sh
curl -X PATCH "http://localhost:8080/jobs/{id}/update-workerdaemon" -H "Content-Type: application/json" -d '{
  "status": "completed",
  "result": "Execution successful",
  "errorMessage": ""
}'
```

Incorrect update due to invalid data (missing status):
```sh
curl -X PATCH "http://localhost:8080/jobs/{id}/update-workerdaemon" -H "Content-Type: application/json" -d '{
  "result": "Some result",
  "errorMessage": "Some error"
}'
```

---

## Repository Selection

Set the environment variable `JOB_REPO_TYPE` to select the repository:

- `inmemory` (default): Use in-memory storage (for development/testing)
- `postgres`: Use PostgreSQL database

If `postgres` is selected, set the following variables:

- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `SSL_MODE`

---

## Development & Testing

### Running Tests
```sh
go test ./...
```

### Building Locally (binary)
```sh
go build -o job-service
```

### Docker Build
```sh
docker build -t job-service .
```

---

## Architecture

The service follows Clean Architecture principles with:

- **Ports**: Interfaces and data models (`ports/`)
- **Adapters**: External integrations (`adapters/`)
  - HTTP handlers (`adapters/handler-http/`)
  - Database repository (`adapters/repo-database/`)
  - In-memory repository (`adapters/repo-in-memory/`)
- **Core**: Business logic (`core/`)
- **Utils**: Shared utilities (`utils/`)

---

## Monitoring & Observability

### Tracing
- All HTTP requests are automatically traced
- Database operations are traced
- Custom spans can be added for business logic
- Traces are exported to Jaeger

### Logging
- Structured logging with JSON format
- Configurable log levels
- Request ID correlation
- Error context preservation

---

## Notes

- The service expects all IDs to be valid UUIDs.
- Status values must be one of: `queued`, `scheduled`, `running`, `completed`, `failed`, `cancelled`.
- Image version validation ensures proper semver or tag format.
- Failed jobs require an error message when updating status.
- The service automatically handles database connection failures with in-memory fallback.
- All timestamps are in UTC format.
- For more details, see the [api.yaml](api.yaml) and the code in the [core](core/), [adapters/handler-http](adapters/handler-http/), and [ports](ports/) directories.

---
