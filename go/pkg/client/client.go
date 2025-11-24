package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openfga/go-sdk/client"
)

func GetClient(apiURL string) (*client.OpenFgaClient, error) {
	return client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl: apiURL,
	})
}

func WaitForOpenFGA(apiURL string) error {
	for i := 0; i < 60; i++ {
		resp, err := http.Get(apiURL + "/healthz")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("OpenFGA failed to start")
}

func CreateStore(ctx context.Context, fgaClient *client.OpenFgaClient, name string) (string, error) {
	store, err := fgaClient.CreateStore(ctx).Body(client.ClientCreateStoreRequest{
		Name: name,
	}).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to create store: %w", err)
	}
	return store.Id, nil
}

func WriteAuthorizationModelJSON(ctx context.Context, fgaClient *client.OpenFgaClient, bodyJSON map[string]interface{}) (string, error) {
	storeID, err := fgaClient.GetStoreId()
	if err != nil {
		return "", fmt.Errorf("failed to get store ID: %w", err)
	}
	apiURL := fgaClient.GetConfig().ApiUrl

	jsonData, err := json.Marshal(bodyJSON)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %w", err)
	}

	url := fmt.Sprintf("%s/stores/%s/authorization-models", apiURL, storeID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		AuthorizationModelID string `json:"authorization_model_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.AuthorizationModelID, nil
}
