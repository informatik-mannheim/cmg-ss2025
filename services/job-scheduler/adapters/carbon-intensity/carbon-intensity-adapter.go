package carbonintensity

import (
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
)

type CarbonIntensityAdapter struct{}

var _ ports.CarbonIntensityAdapter = (*CarbonIntensityAdapter)(nil)

func CreateCarbonIntensityAdapter() CarbonIntensityAdapter {
	return CarbonIntensityAdapter{}
}

func (adapter *CarbonIntensityAdapter) GetCarbonIntensities(zones []string) (model.CarbonIntensityResponse, error) {
	// FIXME: implement
	panic("unimplemented")
}
