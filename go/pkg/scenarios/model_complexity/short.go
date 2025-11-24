package model_complexity

import (
	"context"
	"fmt"

	"github.com/openfga/go-sdk/client"
	pkgclient "github.com/sozua/entendendo-openfga/pkg/client"
	"github.com/sozua/entendendo-openfga/pkg/utils"
)

func GetShortScenario() Scenario {
	return Scenario{
		ID:   "short",
		Name: "Modelo Curto (Poucos Tipos, Baixa Profundidade)",
		Generator: func(ctx context.Context, fgaClient *client.OpenFgaClient) (string, error) {
			body := map[string]interface{}{
				"schema_version": "1.1",
				"type_definitions": []interface{}{
					map[string]interface{}{
						"type": "user",
					},
					map[string]interface{}{
						"type": "document",
						"relations": map[string]interface{}{
							"viewer": map[string]interface{}{
								"this": map[string]interface{}{},
							},
						},
						"metadata": map[string]interface{}{
							"relations": map[string]interface{}{
								"viewer": map[string]interface{}{
									"directly_related_user_types": []interface{}{
										map[string]interface{}{"type": "user"},
									},
								},
							},
						},
					},
				},
			}

			return pkgclient.WriteAuthorizationModelJSON(ctx, fgaClient, body)
		},
		Seeder: func(ctx context.Context, fgaClient *client.OpenFgaClient, count, batchSize, concurrency int) error {
			var tuples []client.ClientTupleKey
			for i := 0; i < count; i++ {
				tuples = append(tuples, client.ClientTupleKey{
					User:     "user:benchmark_user",
					Relation: "viewer",
					Object:   fmt.Sprintf("document:%d", i),
				})
			}
			return utils.WriteBatches(ctx, fgaClient, tuples, batchSize, concurrency)
		},
		CheckTuple: client.ClientCheckRequest{
			User:     "user:benchmark_user",
			Relation: "viewer",
			Object:   "document:1",
		},
	}
}
