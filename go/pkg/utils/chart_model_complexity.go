package utils

import (
	"fmt"
)

type ModelComplexityResult struct {
	Explicit []float64
	Implicit []float64
	ID       string
	Name     string
}

var modelComplexityColors = map[string]map[string]string{
	"Modelo Curto (Poucos Tipos, Baixa Profundidade)": {
		"borderColor":     "rgba(74, 222, 128, 1)",
		"backgroundColor": "rgba(74, 222, 128, 0.1)",
	},
	"Modelo Amplo (Muitos Tipos, Baixa Profundidade)": {
		"borderColor":     "rgba(100, 180, 255, 1)",
		"backgroundColor": "rgba(100, 180, 255, 0.1)",
	},
	"Modelo Profundo (Poucos Tipos, Alta Profundidade 15)": {
		"borderColor":     "rgba(255, 100, 100, 1)",
		"backgroundColor": "rgba(255, 100, 100, 0.1)",
	},
	"Modelo Amplo e Profundo (Muitos Tipos + Alta Profundidade)": {
		"borderColor":     "rgba(255, 220, 100, 1)",
		"backgroundColor": "rgba(255, 220, 100, 0.1)",
	},
}

func GenerateModelComplexityCharts(results map[string]ModelComplexityResult, testRuns int) error {
	runLabels := make([]string, testRuns)
	for i := 0; i < testRuns; i++ {
		runLabels[i] = fmt.Sprintf("Run %d", i+1)
	}

	for scenarioName, data := range results {
		filename := fmt.Sprintf("output/chart_model_complexity_%s.png", data.ID)

		datasets := []ChartDataset{
			{
				Label:            "Implicit Model ID",
				Data:             data.Implicit,
				BorderColor:      "rgba(255, 100, 100, 1)",
				BackgroundColor:  "rgba(255, 100, 100, 0.1)",
				Fill:             true,
				Tension:          0.3,
				PointRadius:      6,
				PointHoverRadius: 8,
			},
			{
				Label:            "Explicit Model ID",
				Data:             data.Explicit,
				BorderColor:      "rgba(100, 180, 255, 1)",
				BackgroundColor:  "rgba(100, 180, 255, 0.1)",
				Fill:             true,
				Tension:          0.3,
				PointRadius:      6,
				PointHoverRadius: 8,
			},
		}

		config := ChartConfig{
			Type: "line",
		}
		config.Data.Labels = runLabels
		config.Data.Datasets = datasets
		config.Options = map[string]interface{}{
			"title": map[string]interface{}{
				"display":    true,
				"text":       fmt.Sprintf("%s (Menor é melhor)", scenarioName),
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
							"display":     true,
							"labelString": "Latency (ms)",
							"fontColor":   "#6c7086",
							"fontFamily":  "Overused Grotesk",
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

		if err := generateChart(config, filename); err != nil {
			return fmt.Errorf("failed to generate chart for %s: %w", scenarioName, err)
		}
	}

	return nil
}
