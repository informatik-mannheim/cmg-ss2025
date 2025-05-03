package ports

import (
	"context"
)

type Notifier interface {
	CarbonIntensityProviderChanged(carbonIntensityProvider CarbonIntensityProvider, ctx context.Context)
}
