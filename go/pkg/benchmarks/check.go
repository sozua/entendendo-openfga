package benchmarks

import (
	"context"
	"time"

	"github.com/openfga/go-sdk/client"
)

func RunCheckBenchmark(ctx context.Context, fgaClient *client.OpenFgaClient, modelID *string, tuple client.ClientCheckRequest, iterations int) (float64, error) {
	start := time.Now()

	options := client.ClientCheckOptions{}
	if modelID != nil {
		options.AuthorizationModelId = modelID
	}

	for i := 0; i < iterations; i++ {
		_, err := fgaClient.Check(ctx).Body(tuple).Options(options).Execute()
		if err != nil {
			return 0, err
		}
	}

	elapsed := time.Since(start)
	return elapsed.Seconds() * 1000 / float64(iterations), nil
}
