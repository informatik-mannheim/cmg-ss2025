# Consumer Gateway Service

The consumer-gateway service is a forwarding service that routes all client requests to the appropriate service. 

> The implementation is in an early stage. Many things are still missing. Use with care.
---
## Usage

> You must first start the respective service by running their `main.go` file.
---
### Jobs
**Endpoints:** `/jobs` , `/jobs/{id}` , `/jobs/{id}/outcome `

**Create New Job:**
```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "jobName": "my-job",
    "creationZone": "DE",
    "image": {
      "name": "img-123"
    },
    "parameters": {
      "flag": "-a"
    }
  }'
``` 
This will return a long and unfiltered (for now) Response.

**Get job outcome:** <br>
To get a specific job, you must copy the id attribute from the create job response.
```bash
curl -X GET http://localhost:8080/jobs/{id}/outcome -H "Content-Type: application/json" 

```

**Get job:** <br>
To get a specific job, you must copy the id attribute from the create job response.
```bash
curl -X GET http://localhost:8080/jobs/{id} -H "Content-Type: application/json" 

```
---

### Login
**Endpoints:** `/auth/login`

**Login:**<br>
Returns token with correct secret.
```bash
curl -X POST http://localhost:8080/auth/login -H  "Content-Type: application/json" -d '{"secret": "this is so secret 123"}'
````