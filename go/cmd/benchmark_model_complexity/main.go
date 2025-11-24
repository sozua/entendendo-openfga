package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sozua/entendendo-openfga/pkg/benchmarks"
	pkgclient "github.com/sozua/entendendo-openfga/pkg/client"
	"github.com/sozua/entendendo-openfga/pkg/config"
	"github.com/sozua/entendendo-openfga/pkg/scenarios/model_complexity"
	"github.com/sozua/entendendo-openfga/pkg/utils"
)

func main() {
	ctx := context.Background()

	fmt.Println("Waiting for OpenFGA...")
	if err := pkgclient.WaitForOpenFGA(config.OpenFGAURL); err != nil {
		log.Fatal(err)
	}
	fmt.Println("OpenFGA is ready.")

	fgaClient, err := pkgclient.GetClient(config.OpenFGAURL)
	if err != nil {
		log.Fatal(err)
	}

	allResults := make(map[string]utils.ModelComplexityResult)
	scenarios := model_complexity.GetAllScenarios()

	for _, scenario := range scenarios {
		storeName := fmt.Sprintf("Bench %s %d", scenario.ID, time.Now().Unix())
		storeID, err := pkgclient.CreateStore(ctx, fgaClient, storeName)
		if err != nil {
			log.Fatalf("Failed to create store: %v", err)
		}
		fgaClient.SetStoreId(storeID)

		modelID, err := scenario.Generator(ctx, fgaClient)
		if err != nil {
			log.Fatalf("Failed to generate model for %s: %v", scenario.Name, err)
		}

		if err := scenario.Seeder(ctx, fgaClient, config.DocumentCount, config.BatchSize, config.Concurrency); err != nil {
			log.Fatalf("Failed to seed data for %s: %v", scenario.Name, err)
		}

		allResults[scenario.Name] = utils.ModelComplexityResult{
			Explicit: make([]float64, 0, config.TestRuns),
			Implicit: make([]float64, 0, config.TestRuns),
			ID:       scenario.ID,
			Name:     scenario.Name,
		}

		_, err = benchmarks.RunCheckBenchmark(ctx, fgaClient, &modelID, scenario.CheckTuple, 20)
		if err != nil {
			log.Fatalf("Warmup failed for %s: %v", scenario.Name, err)
		}

		result := allResults[scenario.Name]
		for i := 1; i <= config.TestRuns; i++ {
			explicitLat, err := benchmarks.RunCheckBenchmark(ctx, fgaClient, &modelID, scenario.CheckTuple, config.CheckIterations)
			if err != nil {
				log.Fatalf("Explicit benchmark failed for %s: %v", scenario.Name, err)
			}
			fmt.Printf("Explicit Run %d for %s: %.3fms\n", i, scenario.Name, explicitLat)
			result.Explicit = append(result.Explicit, explicitLat)

			implicitLat, err := benchmarks.RunCheckBenchmark(ctx, fgaClient, nil, scenario.CheckTuple, config.CheckIterations)
			if err != nil {
				log.Fatalf("Implicit benchmark failed for %s: %v", scenario.Name, err)
			}
			fmt.Printf("Implicit Run %d for %s: %.3fms\n", i, scenario.Name, implicitLat)
			result.Implicit = append(result.Implicit, implicitLat)
		}
		allResults[scenario.Name] = result
	}

	fmt.Println("\nAll Benchmarks Completed. Generating Charts...")
	if err := utils.GenerateModelComplexityCharts(allResults, config.TestRuns); err != nil {
		log.Fatalf("Failed to generate charts: %v", err)
	}
}
