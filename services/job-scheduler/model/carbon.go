package model

import "fmt"

func GetCarbonEndpoint(zone string) string {
	// FIXME: Add base
	// FIXME: change string to UUID
	return fmt.Sprintf("TODO:ADDBASE/carbon-intensity/%s", zone)
}

// CarbonIntensityResponse is the response from the carbon intensity provider
type CarbonIntensityResponse struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
}
