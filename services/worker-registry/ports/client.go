package ports

import (
	"context"
)

type ZoneClient interface {
	IsValidZone(code string, ctx context.Context) bool
	GetZones(ctx context.Context) (ZoneResponse, error)
}
