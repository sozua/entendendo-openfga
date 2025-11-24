package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ChartDataset struct {
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	BorderColor     string    `json:"borderColor"`
	BackgroundColor string    `json:"backgroundColor"`
	Fill            bool      `json:"fill"`
	Tension         float64   `json:"tension"`
	PointRadius     int       `json:"pointRadius"`
	PointHoverRadius int      `json:"pointHoverRadius"`
}

type ChartConfig struct {
	Type string `json:"type"`
	Data struct {
		Labels   []string       `json:"labels"`
		Datasets []ChartDataset `json:"datasets"`
	} `json:"data"`
	Options map[string]interface{} `json:"options"`
}

type QuickChartRequest struct {
	Width           int         `json:"width"`
	Height          int         `json:"height"`
	BackgroundColor string      `json:"backgroundColor"`
	Chart           ChartConfig `json:"chart"`
}

var scenarioColors = map[string]map[string]string{
	"Base (Direct Assignment)": {
		"borderColor":     "rgba(74, 222, 128, 1)",
		"backgroundColor": "rgba(74, 222, 128, 0.1)",
	},
	"Intersection (A AND B)": {
		"borderColor":     "rgba(100, 180, 255, 1)",
		"backgroundColor": "rgba(100, 180, 255, 0.1)",
	},
	"Exclusion (A BUT NOT B)": {
		"borderColor":     "rgba(255, 100, 100, 1)",
		"backgroundColor": "rgba(255, 100, 100, 0.1)",
	},
	"Both ((A AND B) BUT NOT C)": {
		"borderColor":     "rgba(255, 220, 100, 1)",
		"backgroundColor": "rgba(255, 220, 100, 0.1)",
	},
}

func GenerateComparisonCharts(results map[string][]float64, testRuns int) error {
	runLabels := make([]string, testRuns)
	for i := 0; i < testRuns; i++ {
		runLabels[i] = fmt.Sprintf("Run %d", i+1)
	}

	var datasets []ChartDataset
	for name, data := range results {
		colors := scenarioColors[name]
		datasets = append(datasets, ChartDataset{
			Label:            name,
			Data:             data,
			BorderColor:      colors["borderColor"],
			BackgroundColor:  colors["backgroundColor"],
			Fill:             true,
			Tension:          0.3,
			PointRadius:      6,
			PointHoverRadius: 8,
		})
	}

	config := ChartConfig{
		Type: "line",
	}
	config.Data.Labels = runLabels
	config.Data.Datasets = datasets
	config.Options = map[string]interface{}{
		"title": map[string]interface{}{
			"display":    true,
			"text":       "ListObjects Performance por operação (Menor é melhor)",
			"fontColor":  "#ffffff",
			"fontSize":   20,
			"fontFamily": "Overused Grotesk",
		},
		"legend": map[string]interface{}{
			"display": true,
			"labels": map[string]interface{}{
				"fontColor":  "#cdd6f4",
				"fontFamily": "Overused Grotesk",
			},
		},
		"scales": map[string]interface{}{
			"yAxes": []map[string]interface{}{
				{
					"ticks": map[string]interface{}{
						"beginAtZero": true,
						"fontColor":   "#a6adc8",
						"fontFamily":  "Overused Grotesk",
					},
					"gridLines": map[string]interface{}{
						"color": "#333344",
					},
					"scaleLabel": map[string]interface{}{
						"display":    true,
						"labelString": "Latency (ms)",
						"fontColor":  "#6c7086",
						"fontFamily": "Overused Grotesk",
					},
				},
			},
			"xAxes": []map[string]interface{}{
				{
					"ticks": map[string]interface{}{
						"fontColor":  "#a6adc8",
						"fontFamily": "Overused Grotesk",
					},
					"gridLines": map[string]interface{}{
						"display": false,
					},
				},
			},
		},
	}

	return generateChart(config, "output/chart_operators_list_objects.png")
}

func generateChart(config ChartConfig, filename string) error {
	req := QuickChartRequest{
		Width:           800,
		Height:          500,
		BackgroundColor: "#1e1e2e",
		Chart:           config,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal chart config: %w", err)
	}

	resp, err := http.Post("https://quickchart.io/chart/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request to QuickChart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("QuickChart API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode QuickChart response: %w", err)
	}

	imgResp, err := http.Get(result.URL)
	if err != nil {
		return fmt.Errorf("failed to download chart image: %w", err)
	}
	defer imgResp.Body.Close()

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, imgResp.Body); err != nil {
		return fmt.Errorf("failed to write chart image: %w", err)
	}

	fmt.Printf("Generated: %s\n", filename)
	return nil
}
