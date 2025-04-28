package model

// CarbonIntensityData represents carbon intensity information for a single zone.
type CarbonIntensityData struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
}

// AvailableZonesResponse is used to respond with a list of available zones.
type AvailableZonesResponse struct {
	Zones []string `json:"zones"`
}
