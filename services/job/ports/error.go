package ports

import "errors"

var (
	ErrNotExistingID         = errors.New("job ID must be provided")
	ErrInvalidIDFormat       = errors.New("job ID must be a valid UUID")
	ErrJobNotFound           = errors.New("job not found")
	ErrNotExistingJobName    = errors.New("job name must be provided")
	ErrNotExistingStatus     = errors.New("job status must be provided")
	ErrNotExistingImageName  = errors.New("image name must be provided")
	ErrNotExistingWorkerID   = errors.New("worker ID must be provided")
	ErrImageVersionIsInvalid = errors.New("image version format is invalid")
	ErrParamKeyValueEmpty    = errors.New("parameters cannot have empty keys or values")
	ErrErrorMessageEmpty     = errors.New("error message must be provided for failed jobs")
	ErrCarbonIsNegative      = errors.New("carbon intensity must be non-negative")
)
