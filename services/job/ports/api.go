package ports

import (
	"context"
	"errors"
)

var ErrJobNotFound = errors.New("job not found")

type Api interface {
	Set(job Job, ctx context.Context) error
	Get(id string, ctx context.Context) (Job, error)
}
