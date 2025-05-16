package ports

import "github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/model"

type CarbonIntensityAdapter interface {
	GetCarbonIntensities(zones []string) (model.CarbonIntensityResponse, error)
}
