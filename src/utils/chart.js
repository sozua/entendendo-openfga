const QuickChart = require("quickchart-js");
const { TEST_RUNS } = require("../config");
const path = require("path");

const scenarioColors = {
  "Base (Direct Assignment)": {
    borderColor: "rgba(74, 222, 128, 1)",
    backgroundColor: "rgba(74, 222, 128, 0.1)",
  },
  "Intersection (A AND B)": {
    borderColor: "rgba(100, 180, 255, 1)",
    backgroundColor: "rgba(100, 180, 255, 0.1)",
  },
  "Exclusion (A BUT NOT B)": {
    borderColor: "rgba(255, 100, 100, 1)",
    backgroundColor: "rgba(255, 100, 100, 0.1)",
  },
  "Both ((A AND B) BUT NOT C)": {
    borderColor: "rgba(255, 220, 100, 1)",
    backgroundColor: "rgba(255, 220, 100, 0.1)",
  },
};

const generateComparisonCharts = async (results) => {
  const runLabels = Array.from({ length: TEST_RUNS }, (_, i) => `Run ${i + 1}`);
  const scenarioNames = Object.keys(results.list_objects);

  const createChart = async (title, dataKey, filename) => {
    const datasets = scenarioNames.map((name) => ({
      label: name,
      data: results[dataKey][name],
      borderColor: scenarioColors[name].borderColor,
      backgroundColor: scenarioColors[name].backgroundColor,
      fill: true,
      tension: 0.3,
      pointRadius: 6,
      pointHoverRadius: 8,
    }));

    const chart = new QuickChart();
    chart.setWidth(800);
    chart.setHeight(500);
    chart.setBackgroundColor("#1e1e2e");
    chart.setConfig({
      type: "line",
      data: {
        labels: runLabels,
        datasets: datasets,
      },
      options: {
        title: {
          display: true,
          text: title,
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
    await chart.toFile(path.join(__dirname, "../../output", filename));
    console.log(`Generated: output/${filename}`);
  };

  await createChart(
    "ListObjects Performance por operação (Menor é melhor)",
    "list_objects",
    "chart_operators_list_objects.png",
  );
};

module.exports = { generateComparisonCharts };

