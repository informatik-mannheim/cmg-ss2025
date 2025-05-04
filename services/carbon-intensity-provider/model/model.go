package model

type CarbonIntensityData struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
}

type Zone struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type AvailableZonesResponse struct {
	Zones []Zone `json:"zones"`
}
