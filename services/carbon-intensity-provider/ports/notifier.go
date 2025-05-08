package ports

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
)

type Notifier interface {
	CarbonIntensityProviderChanged(data model.CarbonIntensityData, ctx context.Context)
}
