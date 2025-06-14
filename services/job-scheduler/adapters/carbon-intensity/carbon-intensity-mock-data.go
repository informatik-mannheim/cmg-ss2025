package carbonintensity

import "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/ports"

var MockCarbons = []ports.CarbonIntensityData{
	{
		Zone:            "DE",
		CarbonIntensity: 100,
	},
	{
		Zone:            "US",
		CarbonIntensity: 10,
	},
	{
		Zone:            "JP",
		CarbonIntensity: 50,
	},
	{
		Zone:            "FR",
		CarbonIntensity: 20,
	},
	{
		Zone:            "CH",
		CarbonIntensity: 5,
	},
}
