package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sozua/entendendo-openfga/pkg/benchmarks"
	pkgclient "github.com/sozua/entendendo-openfga/pkg/client"
	"github.com/sozua/entendendo-openfga/pkg/config"
	"github.com/sozua/entendendo-openfga/pkg/scenarios/operators"
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

	results := make(map[string][]float64)
	scenarios := operators.GetAllScenarios()

	for _, scenario := range scenarios {
		storeName := fmt.Sprintf("OpBench %s %d", scenario.ID, time.Now().Unix())
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

		results[scenario.Name] = make([]float64, 0, config.TestRuns)

		_, err = benchmarks.RunListObjectsBenchmark(ctx, fgaClient, modelID, "user:benchmark_user", "viewer", "document", 10)
		if err != nil {
			log.Fatalf("Warmup failed for %s: %v", scenario.Name, err)
		}

		for i := 1; i <= config.TestRuns; i++ {
			lat, err := benchmarks.RunListObjectsBenchmark(ctx, fgaClient, modelID, "user:benchmark_user", "viewer", "document", config.ListObjectsIterations)
			if err != nil {
				log.Fatalf("Benchmark failed for %s: %v", scenario.Name, err)
			}
			fmt.Printf("ListObjects Run %d for %s: %.3fms\n", i, scenario.Name, lat)
			results[scenario.Name] = append(results[scenario.Name], lat)
		}
	}

	fmt.Println("\nAll Benchmarks Completed. Generating Charts...")
	if err := utils.GenerateComparisonCharts(results, config.TestRuns); err != nil {
		log.Fatalf("Failed to generate charts: %v", err)
	}
}
