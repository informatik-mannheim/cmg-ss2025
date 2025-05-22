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

**openAPI specified api.yaml shows the possible endpoints and the data schemas**

---

## Running the Service with Docker

1. **Build the Docker image:**
   ```sh
   docker build -t job-service .
   ```

2. **Run the container:**
   ```sh
   docker run -p 8080:8080 job-service
   ```

   The service will be available at `http://localhost:8080`.

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