package benchmarks

import (
	"context"
	"time"

	"github.com/openfga/go-sdk/client"
)

func RunListObjectsBenchmark(ctx context.Context, fgaClient *client.OpenFgaClient, modelID string, user, relation, objectType string, iterations int) (float64, error) {
	start := time.Now()

	options := client.ClientListObjectsOptions{
		AuthorizationModelId: &modelID,
	}

	body := client.ClientListObjectsRequest{
		User:     user,
		Relation: relation,
		Type:     objectType,
	}

	for i := 0; i < iterations; i++ {
		_, err := fgaClient.ListObjects(ctx).Body(body).Options(options).Execute()
		if err != nil {
			return 0, err
		}
	}

	elapsed := time.Since(start)
	return elapsed.Seconds() * 1000 / float64(iterations), nil
}
