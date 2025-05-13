package model

import "fmt"

type CarbonIntensityData struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
}

// -------------------------- Endpoints --------------------------

func GetCarbonEndpoint(base, zone string) string {
	// FIXME: change string to UUID
	return fmt.Sprintf("%s/carbon-intensity/%s", base, zone)
}

// -------------------------- Response & Request --------------------------

// CarbonIntensityResponse is the response from the carbon intensity provider
type CarbonIntensityResponse []CarbonIntensityData
