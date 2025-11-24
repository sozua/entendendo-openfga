const QuickChart = require("quickchart-js");
const { TEST_RUNS } = require("../config");
const path = require("path");

const scenarioColors = {
  "Short Model (Low Types, Low Depth)": {
    borderColor: "rgba(74, 222, 128, 1)",
    backgroundColor: "rgba(74, 222, 128, 0.1)",
  },
  "Wide Model (Many Types, Low Depth)": {
    borderColor: "rgba(100, 180, 255, 1)",
    backgroundColor: "rgba(100, 180, 255, 0.1)",
  },
  "Deep Model (Few Types, High Depth 15)": {
    borderColor: "rgba(255, 100, 100, 1)",
    backgroundColor: "rgba(255, 100, 100, 0.1)",
  },
  "Wide & Deep Model (Many Types + High Depth)": {
    borderColor: "rgba(255, 220, 100, 1)",
    backgroundColor: "rgba(255, 220, 100, 0.1)",
  },
};

const generateCharts = async (allResults) => {
  const runLabels = Array.from({ length: TEST_RUNS }, (_, i) => `Run ${i + 1}`);
  const scenarioNames = Object.keys(allResults);

  for (const scenarioName of scenarioNames) {
    const data = allResults[scenarioName];
    const filename = `chart_model_complexity_${data.id}.png`;
    const outputPath = path.join(__dirname, "../../output", filename);

    const chart = new QuickChart();
    chart.setWidth(800);
    chart.setHeight(500);
    chart.setBackgroundColor("#1e1e2e");

    chart.setConfig({
      type: "line",
      data: {
        labels: runLabels,
        datasets: [
          {
            label: "Implicit Model ID",
            data: data.implicit,
            borderColor: "rgba(255, 100, 100, 1)",
            backgroundColor: "rgba(255, 100, 100, 0.1)",
            fill: true,
            tension: 0.3,
            pointRadius: 6,
            pointHoverRadius: 8,
          },
          {
            label: "Explicit Model ID",
            data: data.explicit,
            borderColor: "rgba(100, 180, 255, 1)",
            backgroundColor: "rgba(100, 180, 255, 0.1)",
            fill: true,
            tension: 0.3,
            pointRadius: 6,
            pointHoverRadius: 8,
          },
        ],
      },
      options: {
        title: {
          display: true,
          text: `${scenarioName} (Menor é melhor)`,
          fontColor: "#ffffff",
          fontSize: 20,
          fontFamily: "Overused Grotesk",
        },
        legend: {
          display: true,
          labels: { fontColor: "#cdd6f4", fontFamily: "Overused Grotesk" },
        },
        scales: {
          yAxes: [
            {
              ticks: {
                beginAtZero: true,
                fontColor: "#a6adc8",
                fontFamily: "Overused Grotesk",
              },
              gridLines: { color: "#333344" },
              scaleLabel: {
                display: true,
                labelString: "Latency (ms)",
                fontColor: "#6c7086",
                fontFamily: "Overused Grotesk",
              },
            },
          ],
          xAxes: [
            {
              ticks: { fontColor: "#a6adc8", fontFamily: "Overused Grotesk" },
              gridLines: { display: false },
            },
          ],
        },
      },
    });
    try {
      await chart.toFile(outputPath);
      console.log(`Generated: ${outputPath}`);
    } catch (err) {
      console.error(`Failed to generate chart for ${scenarioName}:`, err);
    }
  }
};

module.exports = { generateCharts };

