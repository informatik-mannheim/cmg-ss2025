package clients

import (
	"context"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
)

type MockZoneClient struct{}

func (f MockZoneClient) GetZones(ctx context.Context) (ports.ZoneResponse, error) {
	return ports.ZoneResponse{
		Zones: []ports.Zone{
			{Code: "DE", Name: "Germany"},
			{Code: "EN", Name: "England"},
		},
	}, nil
}

func (f MockZoneClient) IsValidZone(code string, ctx context.Context) bool {
	return code == "DE" || code == "EN"
}
