package handler_http

var (
	HTTPErr400MissId           = `{"error": "Bad Request","message": "The job ID must be provided"}`
	HTTPErr400InvalidId        = `{"error": "Bad Request","message": "The job ID format is invalid. Expected a UUID format."}`
	HTTPErr400JobNotFound      = `{"error": "Not Found","message": "A job with the specified ID does not exist. Please verify the ID."}`
	HTTPErr400FieldEmpty       = `{"error": "Bad Request","message": "jobname and imagename must not be empty"}`
	HTTPErr400StatusEmpty      = `{"error": "Bad Request","message": "job status must not be empty"}`
	HTTPErr400InvalidInputData = `{"error": "Bad Request","message": "Invalid input data"}`
	HTTPErr500                 = `{"error": "Internal Server Error","message": "The server encountered an unexpected condition"}`
)
