package ports

import (
	"context"
)

type ZoneValidator interface {
	IsValidZone(code string, ctx context.Context) bool
}
