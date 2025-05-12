package carbonintensity

import (
	"fmt"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"
	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

type CarbonIntensityAdapterMock struct {
	shouldGetCarbonsFail  bool
	shouldGetCarbonsEmpty bool
}

var _ ports.CarbonIntensityAdapter = (*CarbonIntensityAdapterMock)(nil)

func CreateCarbonIntensityAdapterMock(shouldGetCarbonsFail, shouldGetCarbonsEmpty bool) CarbonIntensityAdapterMock {
	return CarbonIntensityAdapterMock{
		shouldGetCarbonsFail:  shouldGetCarbonsFail,
		shouldGetCarbonsEmpty: shouldGetCarbonsEmpty,
	}
}

func (adapter *CarbonIntensityAdapterMock) GetCarbonIntensities(zones []string) (model.CarbonIntensityResponse, error) {
	if adapter.shouldGetCarbonsFail {
		return nil, fmt.Errorf("some carbon get error")
	}
	if adapter.shouldGetCarbonsEmpty {
		return model.CarbonIntensityResponse{}, nil
	}
	response := utils.Filter(MockCarbons, func(carbon model.CarbonIntensityData) bool {
		return utils.Some(zones, func(zone string) bool {
			return zone == carbon.Zone
		})
	})

	return response, nil
}
