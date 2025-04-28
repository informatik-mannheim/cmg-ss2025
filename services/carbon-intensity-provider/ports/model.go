package ports

type CarbonIntensityProvider struct {
	Zone             string `json:"zone"`
	CaarbonIntensity int64  `json:"carbon_intensity"`
}
