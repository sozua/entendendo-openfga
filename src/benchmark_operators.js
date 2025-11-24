const { getClient, waitForOpenFGA } = require('./client');
const { LIST_OBJECTS_ITERATIONS, DOCUMENT_COUNT, TEST_RUNS } = require('./config');
const { runListObjectsBenchmark } = require('./benchmarks/listObjects');
const { generateComparisonCharts } = require('./utils/chart');
const scenarios = require('./scenarios/operators');

const main = async () => {
  console.log('Waiting for OpenFGA...');
  await waitForOpenFGA();
  console.log('OpenFGA is ready.');

  const fgaClient = getClient();

  const results = {
    list_objects: {}
  };

  for (const scenario of scenarios) {
    const storeName = `OpBench ${scenario.id} ${Math.floor(Date.now() / 1000)}`;
    const createStoreResp = await fgaClient.createStore({ name: storeName });
    fgaClient.storeId = createStoreResp.id;

    const modelId = await scenario.generator(fgaClient);

    await scenario.seeder(fgaClient, DOCUMENT_COUNT);
    
    results.list_objects[scenario.name] = [];
    await runListObjectsBenchmark(fgaClient, modelId, 'user:benchmark_user', 'viewer', 'document', 10);

    for (let i = 1; i <= TEST_RUNS; i++) {
      const lat = await runListObjectsBenchmark(fgaClient, modelId, 'user:benchmark_user', 'viewer', 'document', LIST_OBJECTS_ITERATIONS);
      process.stdout.write(`ListObjects Run ${i} for ${scenario.name}: ${lat.toFixed(3)}ms\n`);
      results.list_objects[scenario.name].push(lat);
    }
  }

  console.log('\nAll Benchmarks Completed. Generating Charts...');
  await generateComparisonCharts(results);
};

main().catch(console.error);
