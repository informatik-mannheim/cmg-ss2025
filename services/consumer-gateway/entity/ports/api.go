package ports

import (
	"context"
	"errors"
)

var ErrConsumerNotFound = errors.New("Consumer Gateway not found")

type Api interface {
	Set(consumer Consumer, ctx context.Context) error
	Get(id string, ctx context.Context) (Consumer, error)
}
