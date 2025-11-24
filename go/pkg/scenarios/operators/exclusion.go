package operators

import (
	"context"
	"fmt"

	"github.com/openfga/go-sdk/client"
	pkgclient "github.com/sozua/entendendo-openfga/pkg/client"
	"github.com/sozua/entendendo-openfga/pkg/utils"
)

func GetExclusionScenario() Scenario {
	return Scenario{
		ID:   "but_not",
		Name: "Exclusion (A BUT NOT B)",
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
							"a": map[string]interface{}{
								"this": map[string]interface{}{},
							},
							"b": map[string]interface{}{
								"this": map[string]interface{}{},
							},
							"viewer": map[string]interface{}{
								"difference": map[string]interface{}{
									"base": map[string]interface{}{
										"computedUserset": map[string]interface{}{
											"relation": "a",
										},
									},
									"subtract": map[string]interface{}{
										"computedUserset": map[string]interface{}{
											"relation": "b",
										},
									},
								},
							},
						},
						"metadata": map[string]interface{}{
							"relations": map[string]interface{}{
								"a": map[string]interface{}{
									"directly_related_user_types": []interface{}{
										map[string]interface{}{"type": "user"},
									},
								},
								"b": map[string]interface{}{
									"directly_related_user_types": []interface{}{
										map[string]interface{}{"type": "user"},
									},
								},
								"viewer": map[string]interface{}{
									"directly_related_user_types": []interface{}{},
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
			for i := 1; i <= count; i++ {
				tuples = append(tuples, client.ClientTupleKey{
					User:     "user:benchmark_user",
					Relation: "a",
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
