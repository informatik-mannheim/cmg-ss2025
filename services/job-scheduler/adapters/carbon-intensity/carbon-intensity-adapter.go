package carbonintensity

import (
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

type CarbonIntensityAdapter struct {
	environments model.Environments
}

var _ ports.CarbonIntensityAdapter = (*CarbonIntensityAdapter)(nil)

func NewCarbonIntensityAdapter(environments model.Environments) *CarbonIntensityAdapter {
	return &CarbonIntensityAdapter{
		environments: environments,
	}
}

func (adapter *CarbonIntensityAdapter) GetCarbonIntensities(zones []string) (model.CarbonIntensityResponse, error) {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	responses := make([]model.CarbonIntensityData, len(zones))

	for _, zone := range zones {
		endpoint := model.GetCarbonEndpoint(adapter.environments.CarbonIntensityProviderUrl, zone)

		data, err := utils.GetRequest[model.CarbonIntensityData](endpoint)
		if err != nil {
			return nil, err
		}

		responses = append(responses, data)
	}

	return responses, nil
}
