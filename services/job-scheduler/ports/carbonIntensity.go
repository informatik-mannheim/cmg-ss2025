package ports

type CarbonIntensityData struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
}

// CarbonIntensityResponse is the response from the carbon intensity provider
type CarbonIntensityResponse []CarbonIntensityData

type CarbonIntensityAdapter interface {
	GetCarbonIntensities(zones []string) (CarbonIntensityResponse, error)
}
