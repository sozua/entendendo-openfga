package operators

import (
	"context"
	"fmt"

	"github.com/openfga/go-sdk/client"
	pkgclient "github.com/sozua/entendendo-openfga/pkg/client"
	"github.com/sozua/entendendo-openfga/pkg/utils"
)

func GetBothScenario() Scenario {
	return Scenario{
		ID:   "and_but_not",
		Name: "Both ((A AND B) BUT NOT C)",
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
							"c": map[string]interface{}{
								"this": map[string]interface{}{},
							},
							"temp_and": map[string]interface{}{
								"intersection": map[string]interface{}{
									"child": []interface{}{
										map[string]interface{}{
											"computedUserset": map[string]interface{}{
												"relation": "a",
											},
										},
										map[string]interface{}{
											"computedUserset": map[string]interface{}{
												"relation": "b",
											},
										},
									},
								},
							},
							"viewer": map[string]interface{}{
								"difference": map[string]interface{}{
									"base": map[string]interface{}{
										"computedUserset": map[string]interface{}{
											"relation": "temp_and",
										},
									},
									"subtract": map[string]interface{}{
										"computedUserset": map[string]interface{}{
											"relation": "c",
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
								"c": map[string]interface{}{
									"directly_related_user_types": []interface{}{
										map[string]interface{}{"type": "user"},
									},
								},
								"temp_and": map[string]interface{}{
									"directly_related_user_types": []interface{}{},
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
				tuples = append(tuples, client.ClientTupleKey{
					User:     "user:benchmark_user",
					Relation: "b",
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
