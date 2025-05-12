package notifier

import (
	"context"
	"log"

	"github.com/informatik-mannheim/cmg-ss2025/services/carbon-intensity-provider/ports"
)

type Notifier struct{}

func New() ports.Notifier {
	return &Notifier{}
}

func (d *Notifier) CarbonIntensityProviderChanged(data ports.CarbonIntensityData, ctx context.Context) {
	log.Printf("[Notification] Notifier triggered for zone: %s, intensity: %.2f", data.Zone, data.CarbonIntensity)
}

func (d *Notifier) Event(msg string) {
	log.Printf("[Notifier] %s", msg)
}

var _ ports.Notifier = (*Notifier)(nil)
