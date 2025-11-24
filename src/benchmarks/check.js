const runCheckBenchmark = async (fgaClient, modelId, tuple, iterations) => {
  const start = Date.now();
  const options = { authorizationModelId: modelId };
  for (let i = 0; i < iterations; i++) {
    await fgaClient.check(tuple, options);
  }
  return (Date.now() - start) / iterations;
};

module.exports = { runCheckBenchmark };