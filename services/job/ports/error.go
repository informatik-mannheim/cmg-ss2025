package ports

import "errors"

var (
	ErrNotExistingID         = errors.New("job ID must be provided")
	ErrInvalidIDFormat       = errors.New("job ID must be a valid UUID")
	ErrJobNotFound           = errors.New("job not found")
	ErrNotExistingJobName    = errors.New("job name must be provided")
	ErrNotExistingStatus     = errors.New("job status must be provided")
	ErrNotExistingZone       = errors.New("creation zone must be provided")
	ErrNotExistingImageName  = errors.New("image name must be provided")
	ErrNotExistingWorkerID   = errors.New("worker ID must be provided")
	ErrImageVersionIsInvalid = errors.New("image version format is invalid")
	ErrParamKeyValueEmpty    = errors.New("parameters cannot have empty keys or values")
	ErrErrorMessageEmpty     = errors.New("error message must be provided for failed jobs")
	ErrCarbonIsNegative      = errors.New("carbon intensity must be non-negative")
)

var (
	HTTPErr400MissId           = `{"error": "Bad Request","message": "The job ID must be provided"}`
	HTTPErr400InvalidId        = `{"error": "Bad Request","message": "The job ID format is invalid. Expected a UUID format."}`
	HTTPErr400JobNotFound      = `{"error": "Not Found","message": "A job with the specified ID does not exist. Please verify the ID."}`
	HTTPErr400FieldEmpty       = `{"error": "Bad Request","message": "jobname, creatzone and imagename must not be empty"}`
	HTTPErr400StatusEmpty      = `{"error": "Bad Request","message": "job status must not be empty"}`
	HTTPErr400InvalidInputData = `{"error": "Bad Request","message": "Invalid input data"}`
	HTTPErr500                 = `{"error": "Internal Server Error","message": "The server encountered an unexpected condition"}`
)
