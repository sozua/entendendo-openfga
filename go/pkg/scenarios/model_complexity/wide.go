package model_complexity

import (
	"context"
	"fmt"

	"github.com/openfga/go-sdk/client"
	pkgclient "github.com/sozua/entendendo-openfga/pkg/client"
	"github.com/sozua/entendendo-openfga/pkg/utils"
)

func GetWideScenario() Scenario {
	return Scenario{
		ID:   "wide",
		Name: "Modelo Amplo (Muitos Tipos, Baixa Profundidade)",
		Generator: func(ctx context.Context, fgaClient *client.OpenFgaClient) (string, error) {
			typeDefinitions := []interface{}{
				map[string]interface{}{
					"type": "user",
				},
			}

			for i := 0; i < 80; i++ {
				typeDefinitions = append(typeDefinitions, map[string]interface{}{
					"type": fmt.Sprintf("type_%02d", i),
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
				})
			}

			body := map[string]interface{}{
				"schema_version":   "1.1",
				"type_definitions": typeDefinitions,
			}

			return pkgclient.WriteAuthorizationModelJSON(ctx, fgaClient, body)
		},
		Seeder: func(ctx context.Context, fgaClient *client.OpenFgaClient, count, batchSize, concurrency int) error {
			var tuples []client.ClientTupleKey
			for i := 0; i < count; i++ {
				typeIdx := i % 80
				tuples = append(tuples, client.ClientTupleKey{
					User:     "user:benchmark_user",
					Relation: "viewer",
					Object:   fmt.Sprintf("type_%02d:%d", typeIdx, i),
				})
			}
			return utils.WriteBatches(ctx, fgaClient, tuples, batchSize, concurrency)
		},
		CheckTuple: client.ClientCheckRequest{
			User:     "user:benchmark_user",
			Relation: "viewer",
			Object:   "type_00:1",
		},
	}
}
