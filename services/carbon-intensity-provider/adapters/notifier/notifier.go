package notifier

import (
	"context"
	"log"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/model"
)

type Notifier struct{}

func New() *Notifier {
	return &Notifier{}
}

func (d *Notifier) CarbonIntensityProviderChanged(data model.CarbonIntensityData, ctx context.Context) {
	log.Printf("[Notification] Notifier triggered for zone: %s, intensity: %.2f", data.Zone, data.CarbonIntensity)
}
