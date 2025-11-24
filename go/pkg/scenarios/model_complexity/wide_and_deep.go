package model_complexity

import (
	"context"
	"fmt"

	"github.com/openfga/go-sdk/client"
	pkgclient "github.com/sozua/entendendo-openfga/pkg/client"
	"github.com/sozua/entendendo-openfga/pkg/utils"
)

func GetWideAndDeepScenario() Scenario {
	const depth = 15

	return Scenario{
		ID:   "wide_deep",
		Name: "Modelo Amplo e Profundo (Muitos Tipos + Alta Profundidade)",
		Generator: func(ctx context.Context, fgaClient *client.OpenFgaClient) (string, error) {
			typeDefinitions := []interface{}{
				map[string]interface{}{
					"type": "user",
				},
			}

			for i := 0; i < 65; i++ {
				typeDefinitions = append(typeDefinitions, map[string]interface{}{
					"type": fmt.Sprintf("junk_%d", i),
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

			for i := 0; i < depth; i++ {
				currentType := fmt.Sprintf("node_%02d", i)
				var nextType string
				if i == depth-1 {
					nextType = "user"
				} else {
					nextType = fmt.Sprintf("node_%02d", i+1)
				}

				typeDefinitions = append(typeDefinitions, map[string]interface{}{
					"type": currentType,
					"relations": map[string]interface{}{
						"next": map[string]interface{}{
							"this": map[string]interface{}{},
						},
					},
					"metadata": map[string]interface{}{
						"relations": map[string]interface{}{
							"next": map[string]interface{}{
								"directly_related_user_types": []interface{}{
									map[string]interface{}{"type": nextType},
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

			for i := 0; i < depth; i++ {
				currentObj := fmt.Sprintf("node_%02d:1", i)
				var nextObj string
				if i == depth-1 {
					nextObj = "user:benchmark_user"
				} else {
					nextObj = fmt.Sprintf("node_%02d:1", i+1)
				}
				tuples = append(tuples, client.ClientTupleKey{
					User:     nextObj,
					Relation: "next",
					Object:   currentObj,
				})
			}

			for i := 0; i < count; i++ {
				tuples = append(tuples, client.ClientTupleKey{
					User:     "user:random",
					Relation: "viewer",
					Object:   fmt.Sprintf("junk_%d:%d", i%65, i),
				})
			}

			return utils.WriteBatches(ctx, fgaClient, tuples, batchSize, concurrency)
		},
		CheckTuple: client.ClientCheckRequest{
			User:     "user:benchmark_user",
			Relation: "next",
			Object:   "node_00:1",
		},
	}
}
