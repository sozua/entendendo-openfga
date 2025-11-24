package operators

import (
	"context"

	"github.com/openfga/go-sdk/client"
)

type Scenario struct {
	ID         string
	Name       string
	Generator  func(context.Context, *client.OpenFgaClient) (string, error)
	Seeder     func(context.Context, *client.OpenFgaClient, int, int, int) error
	CheckTuple client.ClientCheckRequest
}

func GetAllScenarios() []Scenario {
	return []Scenario{
		GetBaseScenario(),
		GetIntersectionScenario(),
		GetExclusionScenario(),
		GetBothScenario(),
	}
}
