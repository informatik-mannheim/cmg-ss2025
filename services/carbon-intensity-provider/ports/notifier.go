package ports

import (
	"context"
)

type Notifier interface {
	CarbonIntensityProviderChanged(data CarbonIntensityData, ctx context.Context)
}
