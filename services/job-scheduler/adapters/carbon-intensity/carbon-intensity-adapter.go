package carbonintensity

import (
	"fmt"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

func GetCarbonEndpoint(base, zone string) string {
	return fmt.Sprintf("%s/carbon-intensity/%s", base, zone)
}

type CarbonIntensityAdapter struct {
	baseUrl string
}

var _ ports.CarbonIntensityAdapter = (*CarbonIntensityAdapter)(nil)

func NewCarbonIntensityAdapter(baseUrl string) *CarbonIntensityAdapter {
	return &CarbonIntensityAdapter{
		baseUrl: baseUrl,
	}
}

func (adapter *CarbonIntensityAdapter) GetCarbonIntensities(zones []string) (ports.CarbonIntensityResponse, error) {
	// For now its kept simple and return an error as soon as it gets one, changes in Phase 3
	responses := make([]ports.CarbonIntensityData, len(zones))

	for _, zone := range zones {
		endpoint := GetCarbonEndpoint(adapter.baseUrl, zone)

		// StatusCode is not relevant yet
		data, _, err := utils.GetRequest[ports.CarbonIntensityData](endpoint)
		if err != nil {
			return nil, err
		}

		responses = append(responses, data)
	}

	return responses, nil
}
