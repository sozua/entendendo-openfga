const { getClient, waitForOpenFGA } = require('./client');
const { CHECK_ITERATIONS, DOCUMENT_COUNT, TEST_RUNS } = require('./config');
const { runCheckBenchmark } = require('./benchmarks/check');
const { generateCharts } = require('./utils/chart_model_complexity');
const scenarios = require('./scenarios/model_complexity');

const main = async () => {
  console.log('Waiting for OpenFGA...');
  await waitForOpenFGA();
  console.log('OpenFGA is ready.');

  const fgaClient = getClient();

  const allResults = {};

  for (const scenario of scenarios) {
    const storeName = `Bench ${scenario.id} ${Math.floor(Date.now() / 1000)}`;
    const createStoreResp = await fgaClient.createStore({ name: storeName });
    fgaClient.storeId = createStoreResp.id;

    const modelId = await scenario.generator(fgaClient);

    await scenario.seeder(fgaClient, DOCUMENT_COUNT);
    
    allResults[scenario.name] = {
      explicit: [],
      implicit: []
    };

    await runCheckBenchmark(fgaClient, modelId, scenario.checkTuple, 20); 

    for (let i = 1; i <= TEST_RUNS; i++) {
      const explicitLat = await runCheckBenchmark(fgaClient, modelId, scenario.checkTuple, CHECK_ITERATIONS);
      process.stdout.write(`Explicit Run ${i} for ${scenario.name}: ${explicitLat.toFixed(3)}ms\n`);
      allResults[scenario.name].explicit.push(explicitLat);

      const implicitLat = await runCheckBenchmark(fgaClient, undefined, scenario.checkTuple, CHECK_ITERATIONS);
      process.stdout.write(`Implicit Run ${i} for ${scenario.name}: ${implicitLat.toFixed(3)}ms\n`);
      allResults[scenario.name].implicit.push(implicitLat);
    }
    allResults[scenario.name].id = scenario.id;
  }

  console.log('\nAll Benchmarks Completed. Generating Charts...');
  await generateCharts(allResults);
};

main().catch(console.error);
