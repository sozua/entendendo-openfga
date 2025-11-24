package utils

import (
	"context"
	"fmt"
	"sync"

	"github.com/openfga/go-sdk/client"
)

func WriteBatches(ctx context.Context, fgaClient *client.OpenFgaClient, tuples []client.ClientTupleKey, batchSize, concurrency int) error {
	var batches [][]client.ClientTupleKey
	for i := 0; i < len(tuples); i += batchSize {
		end := i + batchSize
		if end > len(tuples) {
			end = len(tuples)
		}
		batches = append(batches, tuples[i:end])
	}

	queue := make(chan []client.ClientTupleKey, len(batches))
	for _, batch := range batches {
		queue <- batch
	}
	close(queue)

	var wg sync.WaitGroup
	errChan := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range queue {
				body := client.ClientWriteRequest{
					Writes: batch,
				}
				_, err := fgaClient.Write(ctx).Body(body).Execute()
				if err != nil {
					errChan <- fmt.Errorf("error seeding batch: %w", err)
				}
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
