# Job Microservice

A microservice for managing consumer jobs in a queue system. This service provides operations to create, retrieve, and update jobs, facilitating operations from both scheduler and worker perspectives.

## Features

- **Job Management**: Create, retrieve, and update jobs within the system.
- **Status Filtering**: Retrieve jobs based on their status.
- **Scheduler and Worker Integration**: Update specific fields relevant to job schedulers and workers.
- **Comprehensive Testing**: Detailed unit tests covering a wide range of edge cases.

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

### Update Job (Scheduler Perspective)
Update scheduler-related fields of a job.  
**Endpoint**: `PATCH /jobs/{id}/update-scheduler`

### Update Job (Worker Perspective)
Update worker-related fields of a job.  
**Endpoint**: `PATCH /jobs/{id}/update-workerdaemon`


