package zone_validator

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type ZoneValidator struct {
	zones []ports.Zone
}

var _ ports.ZoneValidator = (*ZoneValidator)(nil)

func NewZoneValidator() *ZoneValidator {
	return &ZoneValidator{
		zones: []ports.Zone{
			{Code: "DE", Name: "Germany"},
			{Code: "EN", Name: "England"},
			{Code: "FR", Name: "France"},
		},
	}
}

func (z *ZoneValidator) IsValidZone(code string, ctx context.Context) bool {
	for _, zone := range z.zones {
		if zone.Code == code {
			return true
		}
	}
	return false
}
