# Job Microservice

A microservice for managing consumer jobs in a queue system. This service provides operations to create, retrieve, and update jobs, facilitating operations from both scheduler and worker perspectives.

---

## Features

- **Job Management**: Create, retrieve, and update jobs within the system.
- **Status Filtering**: Retrieve jobs based on their status.
- **Scheduler and Worker Integration**: Update specific fields relevant to job schedulers and workers.
- **Comprehensive Testing**: Detailed unit tests covering a wide range of edge cases.

---

## Technology Stack

- **Language:** Go
- **API:** REST (JSON)
- **Containerization:** Docker
- **Testing:** Go test framework, in-memory repository for fast tests

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

## Data Schemas

**The OpenAPI specification (`api.yaml`) shows the possible endpoints and the data schemas for API requests and responses.**

**Database schema:**  
The SQL schema for the PostgreSQL database (including the `jobs` table and its fields) is located in the [`database`](../../database) directory.  
You can find the table definitions and initialization scripts in files such as `job-init.sql`.

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

### 2. Stopping the Services

To stop the services, press `Ctrl+C` or run:

```sh
docker-compose down
```

### 3. Data Persistence

- PostgreSQL data is persisted in `<select_a_path_on_your_local_system>` on the host.
- The database schema is automatically initialized from the `database` directory via Docker Compose.

### 4. Environment Variables

The following environment variables are set automatically in `docker-compose.yaml`:

- `JOB_REPO_TYPE=postgres`
- `PG_HOST=postgres`
- `PG_PORT=5432`
- `PG_USER=jobuser`
- `PG_PASS=jobpass`
- `PG_DB=jobdb`

---

## Example cURL Commands

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

- `memory` (default): Use in-memory storage (for development/testing)
- `postgres`: Use PostgreSQL database

If `postgres` is selected, set the following variables:

- `PG_HOST`
- `PG_PORT`
- `PG_USER`
- `PG_PASS`
- `PG_DB`

---

## Development & Testing

- To run tests:
  ```sh
  go test ./...
  ```

- To build locally:
  ```sh
  go build -o job-service
  ```

---

## Notes

- The service expects all IDs to be valid UUIDs.
- Status values must be one of: `queued`, `scheduled`, `running`, `completed`, `failed`, `cancelled`.
- For more details, see the [api.yaml](api.yaml) and the code in the [core](core/), [adapters/handler-http](adapters/handler-http/), and [ports](ports/) directories.

---